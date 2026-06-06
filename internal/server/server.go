package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"albums-api/internal/client"
	"albums-api/internal/config"
	"albums-api/internal/db"
	"albums-api/internal/handlers"
	"albums-api/internal/middleware"
	"albums-api/internal/repository"

	_ "albums-api/docs" // ← Этот импорт должен быть именно так

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer *http.Server
	router     *gin.Engine
	cfg        *config.Config
}

func New(cfg *config.Config) *Server {
	// Подключаемся к БД
	database, err := db.New(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	// Создаём репозиторий
	albumRepo := repository.NewAlbumRepository(database)

	// Создаём клиент Auth Service
	authClient := client.NewAuthClient(cfg.Auth.URL)

	// Создаём handler
	albumHandler := handlers.NewAlbumHandler(albumRepo)

	router := gin.New()

	// Global middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// ====================== SWAGGER ======================
	// Должен быть зарегистрирован ДО групп
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	healthHandler := handlers.NewHealthHandler()
	router.GET("/api/v1/health", healthHandler.Check)

	// API v1
	v1 := router.Group("/api/v1")
	{
		albums := v1.Group("/albums")
		{
			// Чтение — доступно всем
			albums.GET("", albumHandler.GetAll)
			albums.GET("/:id", albumHandler.GetByID)

			// Мутации — требуют авторизацию
			albums.POST("", middleware.JWTAuth(authClient), albumHandler.Create)
			albums.PUT("/:id", middleware.JWTAuth(authClient), albumHandler.Update)
			albums.DELETE("/:id", middleware.JWTAuth(authClient), albumHandler.Delete)
		}
	}

	return &Server{
		router: router,
		cfg:    cfg,
	}
}

// GetRouter возвращает роутер для использования в тестах
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

func (s *Server) Start() error {
	s.httpServer = &http.Server{
		Addr:    ":" + s.cfg.Server.Port,
		Handler: s.router,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	fmt.Printf("🚀 Сервер запущен на http://localhost:%s\n", s.cfg.Server.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Получен сигнал завершения...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	fmt.Println("✅ Сервер gracefully остановлен")
	return nil
}
