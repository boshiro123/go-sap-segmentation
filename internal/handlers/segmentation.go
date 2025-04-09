package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-test/internal/repository"
	"go-test/internal/sap"
)

// SegmentationHandler обрабатывает запросы, связанные с сегментацией
type SegmentationHandler struct {
	logger           *slog.Logger
	sapClient        *sap.Client
	segmentationRepo *repository.SegmentationRepository
}

// NewSegmentationHandler создает новый обработчик для сегментации
func NewSegmentationHandler(
	logger *slog.Logger,
	sapClient *sap.Client,
	segmentationRepo *repository.SegmentationRepository,
) *SegmentationHandler {
	return &SegmentationHandler{
		logger:           logger,
		sapClient:        sapClient,
		segmentationRepo: segmentationRepo,
	}
}

// GetAll возвращает все сегменты
// @Summary Получить все сегменты
// @Description Возвращает список всех сегментов из базы данных
// @Tags segmentation
// @Accept json
// @Produce json
// @Success 200 {array} model.Segmentation
// @Failure 500 {object} map[string]string
// @Router /api/segmentation [get]
func (h *SegmentationHandler) GetAll(c *gin.Context) {
	segments, err := h.segmentationRepo.GetAll()
	if err != nil {
		h.logger.Error("failed to get all segments", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get segments"})
		return
	}

	c.JSON(http.StatusOK, segments)
}

// GetByID возвращает сегмент по ID
// @Summary Получить сегмент по ID
// @Description Возвращает сегмент с указанным SAP ID
// @Tags segmentation
// @Accept json
// @Produce json
// @Param id path string true "SAP ID сегмента"
// @Success 200 {object} model.Segmentation
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/segmentation/{id} [get]
func (h *SegmentationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	segment, err := h.segmentationRepo.GetByAddressSapID(id)
	if err != nil {
		h.logger.Error("failed to get segment by ID", "error", err.Error(), "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "segment not found"})
		return
	}

	c.JSON(http.StatusOK, segment)
}

// Import запускает импорт сегментации из SAP API
// @Summary Импортировать сегментацию
// @Description Запускает процесс импорта данных из SAP API в базу данных
// @Tags segmentation
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/segmentation/import [post]
func (h *SegmentationHandler) Import(c *gin.Context) {
	h.logger.Info("starting segmentation import")

	segments, err := h.sapClient.FetchSegmentation()
	if err != nil {
		h.logger.Error("failed to fetch segmentation", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch segmentation data"})
		return
	}

	if len(segments) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no segmentation data to import"})
		return
	}

	if err := h.segmentationRepo.InsertOrUpdate(segments); err != nil {
		h.logger.Error("failed to save segmentation data", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save segmentation data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "import completed successfully",
		"count":   len(segments),
	})
}
