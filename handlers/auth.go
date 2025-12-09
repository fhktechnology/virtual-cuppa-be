package handlers

import (
	"net/http"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/services"
	"virtual-cuppa-be/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Register(&input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.Login(&input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input models.RefreshTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.authService.RefreshToken(input.RefreshToken)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, response)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"user": user})
}
