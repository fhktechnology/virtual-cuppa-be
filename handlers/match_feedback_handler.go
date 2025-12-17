package handlers

import (
	"net/http"
	"strconv"
	"virtual-cuppa-be/models"
	"virtual-cuppa-be/services"
	"virtual-cuppa-be/utils"

	"github.com/gin-gonic/gin"
)

type MatchFeedbackHandler struct {
	matchService services.MatchService
}

func NewMatchFeedbackHandler(matchService services.MatchService) *MatchFeedbackHandler {
	return &MatchFeedbackHandler{
		matchService: matchService,
	}
}

type SubmitFeedbackRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

func (h *MatchFeedbackHandler) SubmitFeedback(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid match ID")
		return
	}

	var req SubmitFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.matchService.SubmitFeedback(userID.(uint), uint(matchID), req.Rating, req.Comment)
	if err != nil {
		switch err {
		case services.ErrMatchNotFound:
			utils.RespondWithError(c, http.StatusNotFound, "Match not found")
		case services.ErrUnauthorizedMatch:
			utils.RespondWithError(c, http.StatusForbidden, "Unauthorized to provide feedback for this match")
		case services.ErrMatchNotAccepted:
			utils.RespondWithError(c, http.StatusBadRequest, "Can only provide feedback for accepted matches")
		case services.ErrFeedbackAlreadyExists:
			utils.RespondWithError(c, http.StatusConflict, "Feedback already submitted for this match")
		case services.ErrInvalidRating:
			utils.RespondWithError(c, http.StatusBadRequest, "Rating must be between 1 and 5")
		default:
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to submit feedback")
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Feedback submitted successfully",
	})
}

func (h *MatchFeedbackHandler) GetMatchFeedbacks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid match ID")
		return
	}

	feedbacks, err := h.matchService.GetMatchFeedbacks(userID.(uint), uint(matchID))
	if err != nil {
		switch err {
		case services.ErrMatchNotFound:
			utils.RespondWithError(c, http.StatusNotFound, "Match not found")
		case services.ErrUnauthorizedMatch:
			utils.RespondWithError(c, http.StatusForbidden, "Unauthorized to view feedbacks for this match")
		default:
			utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve feedbacks")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedbacks": feedbacks,
		"count":     len(feedbacks),
	})
}

func (h *MatchFeedbackHandler) GetPendingFeedback(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	matches, err := h.matchService.GetMatchesPendingFeedback(userID.(uint))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve pending feedback")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   len(matches),
	})
}

// Admin handler to view all feedbacks for a match
func (h *MatchFeedbackHandler) AdminGetMatchFeedbacks(c *gin.Context) {
	// Check if user is admin
	accountType, exists := c.Get("accountType")
	if !exists || accountType != models.AccountTypeAdmin {
		utils.RespondWithError(c, http.StatusForbidden, "Admin access required")
		return
	}

	matchID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid match ID")
		return
	}

	// For admin, we bypass the authorization check
	// Get match to verify it exists and belongs to admin's organisation
	userID, _ := c.Get("userID")
	match, err := h.matchService.GetMatchFeedbacks(userID.(uint), uint(matchID))
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve feedbacks")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"feedbacks": match,
		"count":     len(match),
	})
}
