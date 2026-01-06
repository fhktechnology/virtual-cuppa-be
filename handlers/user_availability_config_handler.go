package handlers

import (
	"net/http"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/services"
	"virtual-cuppa-be/utils"

	"github.com/gin-gonic/gin"
)

type UserAvailabilityConfigHandler struct {
	configService services.UserAvailabilityConfigService
	matchService  services.MatchService
}

func NewUserAvailabilityConfigHandler(configService services.UserAvailabilityConfigService, matchService services.MatchService) *UserAvailabilityConfigHandler {
	return &UserAvailabilityConfigHandler{
		configService: configService,
		matchService:  matchService,
	}
}

// CreateConfig creates availability configuration for the authenticated user
func (h *UserAvailabilityConfigHandler) CreateConfig(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.CreateAvailabilityConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	config, err := h.configService.CreateConfig(userID.(uint), input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Try to generate a match for the user now that they have availability config
	_ = h.matchService.TryGenerateMatchForUser(userID.(uint))

	utils.RespondWithSuccess(c, http.StatusCreated, config)
}

// GetConfig retrieves availability configuration for the authenticated user
func (h *UserAvailabilityConfigHandler) GetConfig(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	config, err := h.configService.GetConfig(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, config)
}

// UpdateConfig updates availability configuration for the authenticated user
func (h *UserAvailabilityConfigHandler) UpdateConfig(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.UpdateAvailabilityConfigInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	config, err := h.configService.UpdateConfig(userID.(uint), input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Try to generate a match for the user after updating availability config
	_ = h.matchService.TryGenerateMatchForUser(userID.(uint))

	utils.RespondWithSuccess(c, http.StatusOK, config)
}

// DeleteConfig deletes availability configuration for the authenticated user
func (h *UserAvailabilityConfigHandler) DeleteConfig(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := h.configService.DeleteConfig(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Availability configuration deleted successfully",
	})
}

// HasConfig checks if the authenticated user has availability configuration
func (h *UserAvailabilityConfigHandler) HasConfig(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	hasConfig, err := h.configService.HasConfig(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"hasConfig": hasConfig,
	})
}
