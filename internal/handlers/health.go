package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler обрабатывает запросы для проверки работоспособности системы
type HealthHandler struct {
	logger *slog.Logger
}

// NewHealthHandler создает новый обработчик для проверки здоровья
func NewHealthHandler(logger *slog.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// Check эндпоинт для проверки работоспособности сервера
// @Summary Проверка работоспособности
// @Description Проверяет работоспособность API сервера
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	h.logger.Debug("health check requested")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
