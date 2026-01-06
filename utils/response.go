package utils

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{Error: message})
}

func RespondWithSuccess(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func RespondWithMessage(c *gin.Context, code int, message string) {
	c.JSON(code, SuccessResponse{Message: message})
}

func HandleServiceError(c *gin.Context, err error) {
	switch err.Error() {
	case "user with this email already exists":
		RespondWithError(c, http.StatusConflict, err.Error())
	case "invalid credentials":
		RespondWithError(c, http.StatusUnauthorized, err.Error())
	case "user not found":
		RespondWithError(c, http.StatusNotFound, err.Error())
	case "invalid refresh token":
		RespondWithError(c, http.StatusUnauthorized, err.Error())
	case "admin must be assigned to an organisation to import users":
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case "invalid CSV format, expected: firstName,lastName,email":
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case "CSV file is empty":
		RespondWithError(c, http.StatusBadRequest, err.Error())
	case "user does not belong to admin's organisation":
		RespondWithError(c, http.StatusForbidden, err.Error())
	case "availability configuration not found":
		RespondWithError(c, http.StatusNotFound, err.Error())
	case "availability configuration already exists for this user":
		RespondWithError(c, http.StatusConflict, err.Error())
	case "at least one availability slot must be selected":
		RespondWithError(c, http.StatusBadRequest, err.Error())
	default:
		log.Printf("Unhandled service error: %v", err)
		if strings.Contains(err.Error(), "sendgrid") {
			RespondWithError(c, http.StatusInternalServerError, "Failed to send email: "+err.Error())
		} else {
			RespondWithError(c, http.StatusInternalServerError, "Internal server error")
		}
	}
}
