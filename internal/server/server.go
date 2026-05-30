package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"albums-api/internal/config"
	"albums-api/internal/db"
	"albums-api/internal/handlers"
	"albums-api/internal/repository"

	_ "albums-api/docs" // ← Этот импорт должен быть именно так

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	router *gin.Engine
	db     *db.DB
	config *config.Config
}

func New(cfg *config.Config) (*Server, error) {
	config.InitLogger(cfg.Env)

	database, err := db.NewConnection(cfg)
	if err != nil {
		return nil, err
	}

	repo := repository.NewAlbumRepository(database)
	albumHandler := handlers.NewAlbumHandler(repo)
	healthHandler := handlers.NewHealthHandler()

	router := gin.Default()

	// CORS Middleware
	router.Use(corsMiddleware())

	// Logging Middleware
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		config.Log.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  c.Writer.Status(),
			"latency": time.Since(start),
			"client":  c.ClientIP(),
		}).Info("request completed")
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	api := router.Group("/api/v1")
	{
		albums := api.Group("/albums")
		{
			albums.GET("", albumHandler.GetAll)
			albums.GET("/:id", albumHandler.GetByID)
			albums.POST("", albumHandler.Create)
			albums.PUT("/:id", albumHandler.Update)
			albums.DELETE("/:id", albumHandler.Delete)
		}

		api.GET("/health", healthHandler.Check)
	}

	return &Server{
		router: router,
		db:     database,
		config: cfg,
	}, nil
}

// corsMiddleware ... (оставь как было)
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:    ":" + s.config.Server.Port,
		Handler: s.router,
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("🚀 Сервер запущен на http://localhost:%s\n", s.config.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Ожидаем сигнал остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Получен сигнал завершения...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("❌ Ошибка graceful shutdown: %v\n", err)
	}

	s.db.Close()
	fmt.Println("👋 Сервер остановлен gracefully")
	return nil
}
