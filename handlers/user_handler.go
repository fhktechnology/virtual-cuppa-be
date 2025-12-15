package handlers

import (
	"fmt"
	"log"
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
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("GetOrganisationUsers: userID not found in context")
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	log.Printf("GetOrganisationUsers: userID=%v", userID)

	// Get current user to retrieve organisationID
	currentUser, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		log.Printf("GetOrganisationUsers: error getting user: %v", err)
		utils.HandleServiceError(c, err)
		return
	}

	log.Printf("GetOrganisationUsers: currentUser=%+v, organisationID=%v", currentUser.Email, currentUser.OrganisationID)

	if currentUser.OrganisationID == nil {
		log.Println("GetOrganisationUsers: user has no organisationID")
		utils.RespondWithError(c, http.StatusBadRequest, "User not assigned to any organisation")
		return
	}

	users, err := h.userService.GetUsersByOrganisation(*currentUser.OrganisationID)
	if err != nil {
		log.Printf("GetOrganisationUsers: error getting users: %v", err)
		utils.HandleServiceError(c, err)
		return
	}

	log.Printf("GetOrganisationUsers: found %d users", len(users))
	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"users": users,
		"count": len(users),
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.CreateUser(userID.(uint), &input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusCreated, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	targetUserID := c.Param("id")
	if targetUserID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var id uint
	if _, err := fmt.Sscanf(targetUserID, "%d", &id); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(userID.(uint), id); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) UpdateTags(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	targetUserID := c.Param("userId")
	if targetUserID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	var id uint
	if _, err := fmt.Sscanf(targetUserID, "%d", &id); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input models.UpdateTagsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.UpdateUserTags(userID.(uint), id, input.Tags); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "Tags updated successfully"})
}
