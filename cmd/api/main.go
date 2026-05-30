package main

import (
	"log"

	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/handler"
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/middleware"
	"tennisdaily-backend/internal/repository"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	authService := service.NewAuthService(cfg, userRepo)
	sessionService := service.NewSessionService(sessionRepo)
	statsService := service.NewStatsService(sessionRepo)

	authHandler := handler.NewAuthHandler(authService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	statsHandler := handler.NewStatsHandler(statsService)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		logger.Debug("GET /health start clientIP=%s", c.ClientIP())
		logger.Debug("GET /health success status=ok")
		response.OK(c, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	api.POST("/auth/wechat-login", authHandler.WechatLogin)

	authed := api.Group("")
	authed.Use(middleware.Auth(authService))
	{
		authed.GET("/sessions", sessionHandler.List)
		authed.POST("/sessions", sessionHandler.Create)
		authed.GET("/sessions/latest", sessionHandler.Latest)
		authed.GET("/sessions/:id", sessionHandler.Get)
		authed.PUT("/sessions/:id", sessionHandler.Update)
		authed.DELETE("/sessions/:id", sessionHandler.Delete)
		authed.GET("/stats/month", statsHandler.Month)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
