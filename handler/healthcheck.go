package handler

import (
	"golang-api-crowdfunding/healthcheck"
	"golang-api-crowdfunding/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthcheckHandler struct {
	service healthcheck.Service
}

func NewHealthcheckHandler(service healthcheck.Service) *healthcheckHandler {
	return &healthcheckHandler{service}
}

func (h *healthcheckHandler) HealthcheckHandler(c *gin.Context) {
	// Perform the health check using the service
	health, err := h.service.HealthcheckService()
	if err != nil {
		response := helper.ApiResponse("Internal server error", http.StatusInternalServerError, "error", err)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.ApiResponse("Health check passed successfully", http.StatusOK, "success", health)

	// Return the health check result as JSON
	c.JSON(http.StatusOK, response)
}
