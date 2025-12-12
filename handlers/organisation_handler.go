package handlers

import (
	"net/http"

	"virtual-cuppa-be/models"
	"virtual-cuppa-be/services"
	"virtual-cuppa-be/utils"

	"github.com/gin-gonic/gin"
)

type OrganisationHandler struct {
	organisationService services.OrganisationService
	userService         services.UserService
}

func NewOrganisationHandler(organisationService services.OrganisationService, userService services.UserService) *OrganisationHandler {
	return &OrganisationHandler{
		organisationService: organisationService,
		userService:         userService,
	}
}

func (h *OrganisationHandler) UpsertOrganisation(c *gin.Context) {
	var input models.UpsertOrganisationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	organisation, err := h.organisationService.UpsertOrganisation(&input)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	// Assign organisation to the current user
	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	user.OrganisationID = &organisation.ID
	if err := h.userService.UpdateUser(user); err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, organisation)
}

func (h *OrganisationHandler) GetOrganisation(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	if user.OrganisationID == nil {
		utils.RespondWithError(c, http.StatusNotFound, "User not assigned to any organisation")
		return
	}

	organisation, err := h.organisationService.GetOrganisationByID(*user.OrganisationID)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, organisation)
}
