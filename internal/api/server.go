package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-test/internal/sap"
	"go-test/pkg/config"
	"go-test/pkg/repository"
)

type Server struct {
	router           *gin.Engine
	logger           *slog.Logger
	cfg              *config.Config
	sapClient        *sap.Client
	segmentationRepo *repository.SegmentationRepository
}

func NewServer(
	cfg *config.Config,
	logger *slog.Logger,
	sapClient *sap.Client,
	segmentationRepo *repository.SegmentationRepository,
) *Server {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(loggerMiddleware(logger))

	server := &Server{
		router:           router,
		logger:           logger,
		cfg:              cfg,
		sapClient:        sapClient,
		segmentationRepo: segmentationRepo,
	}

	server.initRoutes()

	return server
}

func loggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"ip", c.ClientIP(),
		)

		c.Next()

		logger.Info("request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
		)
	}
}

func (s *Server) initRoutes() {
	api := s.router.Group("/api")
	{
		segmentation := api.Group("/segmentation")
		{
			segmentation.GET("/", s.getAllSegments)
			segmentation.GET("/:id", s.getSegmentByID)
			segmentation.POST("/import", s.importSegmentation)
		}

		api.GET("/health", s.healthCheck)
	}

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (s *Server) Run(addr string) error {
	s.logger.Info("starting API server", "address", addr)
	return s.router.Run(addr)
}

// getAllSegments возвращает все сегменты
// @Summary Получить все сегменты
// @Description Возвращает список всех сегментов из базы данных
// @Tags segmentation
// @Accept json
// @Produce json
// @Success 200 {array} model.Segmentation
// @Failure 500 {object} map[string]string
// @Router /api/segmentation [get]
func (s *Server) getAllSegments(c *gin.Context) {
	segments, err := s.segmentationRepo.GetAll()
	if err != nil {
		s.logger.Error("failed to get all segments", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get segments"})
		return
	}

	c.JSON(http.StatusOK, segments)
}

// getSegmentByID возвращает сегмент по ID
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
func (s *Server) getSegmentByID(c *gin.Context) {
	id := c.Param("id")

	segment, err := s.segmentationRepo.GetByAddressSapID(id)
	if err != nil {
		s.logger.Error("failed to get segment by ID", "error", err.Error(), "id", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "segment not found"})
		return
	}

	c.JSON(http.StatusOK, segment)
}

// importSegmentation запускает импорт сегментации из SAP API
// @Summary Импортировать сегментацию
// @Description Запускает процесс импорта данных из SAP API в базу данных
// @Tags segmentation
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/segmentation/import [post]
func (s *Server) importSegmentation(c *gin.Context) {
	s.logger.Info("starting segmentation import")

	segments, err := s.sapClient.FetchSegmentation()
	if err != nil {
		s.logger.Error("failed to fetch segmentation", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch segmentation data"})
		return
	}

	if len(segments) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no segmentation data to import"})
		return
	}

	if err := s.segmentationRepo.InsertOrUpdate(segments); err != nil {
		s.logger.Error("failed to save segmentation data", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save segmentation data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "import completed successfully",
		"count":   len(segments),
	})
}

// healthCheck эндпоинт для проверки работоспособности сервера
// @Summary Проверка работоспособности
// @Description Проверяет работоспособность API сервера
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/health [get]
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
