package handlers

import (
	"github.com/fishmanDK/rutube-test-task/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

const (
	contextTimeResponse = 50 * time.Millisecond
)

type Handlers struct {
	service *service.Service
	logger  *slog.Logger
}

func MustHandlers(service *service.Service, logger *slog.Logger) *Handlers {
	return &Handlers{
		service: service,
		logger:  logger,
	}
}

func (h *Handlers) InitRouts() *gin.Engine {
	router := gin.New()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/update-access-token", h.apdateAccessToken)
	}

	api := router.Group("/api")
	api.Use(h.authMiddleware)
	{
		api.POST("/subscribe", h.subscribe)
		api.POST("/unsubscribe", h.unsubscribe)
		api.GET("/subscription-list", h.allSubs)
	}

	return router
}
