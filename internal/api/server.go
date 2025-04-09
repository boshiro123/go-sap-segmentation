package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-test/internal/handlers"
	"go-test/internal/repository"
	"go-test/internal/sap"
	"go-test/pkg/config"
)

type Server struct {
	router              *gin.Engine
	logger              *slog.Logger
	cfg                 *config.Config
	segmentationHandler *handlers.SegmentationHandler
	healthHandler       *handlers.HealthHandler
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

	// Инициализация обработчиков
	segmentationHandler := handlers.NewSegmentationHandler(logger, sapClient, segmentationRepo)
	healthHandler := handlers.NewHealthHandler(logger)

	server := &Server{
		router:              router,
		logger:              logger,
		cfg:                 cfg,
		segmentationHandler: segmentationHandler,
		healthHandler:       healthHandler,
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
			segmentation.GET("/", s.segmentationHandler.GetAll)
			segmentation.GET("/:id", s.segmentationHandler.GetByID)
			segmentation.POST("/import", s.segmentationHandler.Import)
		}

		api.GET("/health", s.healthHandler.Check)
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
