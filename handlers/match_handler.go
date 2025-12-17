package handlers

import (
	"net/http"
	"strconv"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/scheduler"
	"virtual-cuppa-be/services"
	"virtual-cuppa-be/utils"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService   services.MatchService
	matchScheduler *scheduler.MatchScheduler
}

func NewMatchHandler(matchService services.MatchService, matchScheduler *scheduler.MatchScheduler) *MatchHandler {
	return &MatchHandler{
		matchService:   matchService,
		matchScheduler: matchScheduler,
	}
}

// GetCurrentMatch returns the current pending match for the user
func (h *MatchHandler) GetCurrentMatch(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	match, err := h.matchService.GetCurrentMatch(userID.(uint))
	if err != nil {
		if err == services.ErrMatchNotFound {
			c.Status(http.StatusNoContent)
			return
		}
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, match)
}

// GetMatchHistory returns all matches for the user
func (h *MatchHandler) GetMatchHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matches, err := h.matchService.GetMatchHistory(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// AcceptMatch accepts a pending match with availability
func (h *MatchHandler) AcceptMatch(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matchIDStr := c.Param("id")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid match ID")
		return
	}

	var req struct {
		Availability models.Availability `json:"availability" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body. Expected format: {\"availability\": {\"2025-02-18\": [\"09:30\", \"10:30\"]}}")
		return
	}

	match, err := h.matchService.AcceptMatchWithAvailability(userID.(uint), uint(matchID), req.Availability)
	if err != nil {
		if err == services.ErrMatchNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Match not found")
			return
		}
		if err == services.ErrUnauthorizedMatch {
			utils.RespondWithError(c, http.StatusForbidden, "Unauthorized to modify this match")
			return
		}
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Match accepted successfully",
		"match":   match,
	})
}

// RejectMatch rejects a pending match
func (h *MatchHandler) RejectMatch(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matchIDStr := c.Param("id")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid match ID")
		return
	}

	if err := h.matchService.RejectMatch(userID.(uint), uint(matchID)); err != nil {
		if err == services.ErrMatchNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Match not found")
			return
		}
		if err == services.ErrUnauthorizedMatch {
			utils.RespondWithError(c, http.StatusForbidden, "Unauthorized to modify this match")
			return
		}
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Match rejected successfully"})
}

// GenerateMatches generates matches for the admin's organisation (admin only)
func (h *MatchHandler) GenerateMatches(c *gin.Context) {
	organisationID, exists := c.Get("organisationID")
	if !exists || organisationID == nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Admin not assigned to any organisation")
		return
	}

	orgID, ok := organisationID.(uint)
	if !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid organisation ID")
		return
	}

	count, err := h.matchService.GenerateMatchesForOrganisation(orgID)
	if err != nil {
		if err == services.ErrNoUsersToMatch {
			utils.RespondWithError(c, http.StatusBadRequest, "Not enough users to create matches")
			return
		}
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Matches generated successfully",
		"count":   count,
	})
}

// GetOrganisationMatches returns all matches for the admin's organisation
func (h *MatchHandler) GetOrganisationMatches(c *gin.Context) {
	organisationID, exists := c.Get("organisationID")
	if !exists || organisationID == nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Admin not assigned to any organisation")
		return
	}

	orgID, ok := organisationID.(uint)
	if !ok {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid organisation ID")
		return
	}

	matches, err := h.matchService.GetOrganisationMatches(orgID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// TriggerScheduler manually triggers the scheduler to check and generate matches
func (h *MatchHandler) TriggerScheduler(c *gin.Context) {
	if h.matchScheduler == nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Scheduler not initialized")
		return
	}

	go h.matchScheduler.RunNow()
	
	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Scheduler triggered successfully - matches will be generated for all organisations",
	})
}

// GetMatchAvailabilities returns availabilities for a specific match
func (h *MatchHandler) GetMatchAvailabilities(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matchIDStr := c.Param("id")
	matchID, err := strconv.ParseUint(matchIDStr, 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid match ID")
		return
	}

	availabilities, err := h.matchService.GetMatchAvailabilities(userID.(uint), uint(matchID))
	if err != nil {
		if err == services.ErrMatchNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Match not found")
			return
		}
		if err == services.ErrUnauthorizedMatch {
			utils.RespondWithError(c, http.StatusForbidden, "Unauthorized to view this match")
			return
		}
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"availabilities": availabilities,
	})
}
