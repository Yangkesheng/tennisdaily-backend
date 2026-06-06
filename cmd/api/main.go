package main

import (
	"fmt"
	"log"
	"os"

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
	printCurrentDirectoryInfo()

	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	racketRepo := repository.NewRacketRepository(db)

	authService := service.NewAuthService(cfg, userRepo)
	sessionService := service.NewSessionService(sessionRepo)
	statsService := service.NewStatsService(sessionRepo)
	racketService := service.NewRacketService(racketRepo)

	authHandler := handler.NewAuthHandler(authService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	statsHandler := handler.NewStatsHandler(statsService)
	racketHandler := handler.NewRacketHandler(racketService)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		logger.Debug("GET /health start clientIP=%s", c.ClientIP())
		logger.Debug("GET /health success status=ok")
		response.OK(c, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	api.POST("/auth/wechat-login", authHandler.WechatLogin)
	api.POST("/auth/phone-login", authHandler.PhoneLogin)

	authed := api.Group("")
	authed.Use(middleware.Auth(authService))
	{
		authed.GET("/auth/me", authHandler.Me)
		authed.POST("/auth/logout", authHandler.Logout)
		authed.PUT("/auth/profile", authHandler.UpdateProfile)

		authed.GET("/sessions", sessionHandler.List)
		authed.POST("/sessions", sessionHandler.Create)
		authed.GET("/sessions/latest", sessionHandler.Latest)
		authed.GET("/sessions/:id", sessionHandler.Get)
		authed.PUT("/sessions/:id", sessionHandler.Update)
		authed.DELETE("/sessions/:id", sessionHandler.Delete)
		authed.GET("/stats/month", statsHandler.Month)

		authed.GET("/racket-library", racketHandler.Library)
		authed.GET("/my-rackets", racketHandler.MyRackets)
		authed.GET("/rackets/stats", racketHandler.Stats)
		authed.GET("/rackets", racketHandler.List)
		authed.POST("/rackets", racketHandler.Create)
		authed.GET("/rackets/selectable", racketHandler.Selectable)
		authed.GET("/rackets/:id", racketHandler.Detail)
		authed.PUT("/rackets/:id", racketHandler.Update)
		authed.DELETE("/rackets/:id", racketHandler.Delete)
		authed.POST("/rackets/:id/set-primary", racketHandler.SetPrimary)
		authed.POST("/rackets/:id/retire", racketHandler.Retire)
		authed.POST("/rackets/:id/stringing-records", racketHandler.CreateStringingRecord)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

func printCurrentDirectoryInfo() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("当前路径: 获取失败: %v\n", err)
	} else {
		fmt.Printf("当前路径: %s\n", wd)
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("读取当前目录失败: %v\n", err)
		return
	}

	fmt.Println("当前目录文件列表:")
	for _, entry := range entries {
		fmt.Println(entry.Name())
	}
}
