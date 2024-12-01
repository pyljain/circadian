package routes

import (
	"net/http"

	"circadian/internal/db"
	"circadian/internal/models"

	"github.com/labstack/echo/v4"
)

type HealthCheckResultHandler struct {
	db *db.Db
}

func NewHealthCheckResultHandler(database *db.Db) *HealthCheckResultHandler {
	return &HealthCheckResultHandler{db: database}
}

func (h *HealthCheckResultHandler) CreateHealthCheckResult(c echo.Context) error {
	hcr := new(models.HealthCheckResult)
	if err := c.Bind(hcr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.db.InsertHealthCheckResult(hcr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, hcr)
}

func (h *HealthCheckResultHandler) GetHealthCheckResult(c echo.Context) error {
	id := c.Param("id")

	hcr, err := h.db.GetHealthCheckResult(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, hcr)
}

func (h *HealthCheckResultHandler) UpdateHealthCheckResult(c echo.Context) error {
	hcr := new(models.HealthCheckResult)
	if err := c.Bind(hcr); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err := h.db.UpdateHealthCheckResult(hcr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, hcr)
}

func (h *HealthCheckResultHandler) DeleteHealthCheckResult(c echo.Context) error {
	id := c.Param("id")

	err := h.db.DeleteHealthCheckResult(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Health check result deleted"})
}

func (h *HealthCheckResultHandler) ListHealthCheckResults(c echo.Context) error {

	results, err := h.db.ListHealthCheckResults()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}

func RegisterHealthCheckResultRoutes(e *echo.Echo, database *db.Db) {
	handler := NewHealthCheckResultHandler(database)

	e.POST("/health-check-results", handler.CreateHealthCheckResult)
	e.GET("/health-check-results/:id", handler.GetHealthCheckResult)
	e.PUT("/health-check-results", handler.UpdateHealthCheckResult)
	e.DELETE("/health-check-results/:id", handler.DeleteHealthCheckResult)
	e.GET("/health-check-results", handler.ListHealthCheckResults)
}
