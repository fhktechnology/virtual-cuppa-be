package handlers

import (
	"net/http"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/services"
	"virtual-cuppa-be/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) ImportCSV(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "CSV file is required")
		return
	}
	defer file.Close()

	count, err := h.userService.ImportUsersFromCSV(userID.(uint), file)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "Users imported successfully",
		"count":   count,
	})
}

func (h *UserHandler) ConfirmUser(c *gin.Context) {
	adminID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.ConfirmUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.userService.ConfirmUser(adminID.(uint), input.UserID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "User confirmed successfully",
	})
}

func (h *UserHandler) GetOrganisationUsers(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organisation := c.Query("organisation")
	if organisation == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "organisation query parameter is required")
		return
	}

	users, err := h.userService.GetUsersByOrganisation(organisation)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}
