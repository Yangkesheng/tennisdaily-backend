package main

import (
	"context"
	"log"
	"time"

	"tennisdaily-backend/internal/config"
	"tennisdaily-backend/internal/handler"
	"tennisdaily-backend/internal/logger"
	"tennisdaily-backend/internal/middleware"
	"tennisdaily-backend/internal/repository"
	"tennisdaily-backend/internal/response"
	"tennisdaily-backend/internal/service"
	"tennisdaily-backend/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
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

	racketRepo := repository.NewRacketRepository(db)

	contentSecurityService := service.NewContentSecurityService(cfg)
	authService := service.NewAuthService(cfg, userRepo, contentSecurityService)
	sessionService := service.NewSessionService(sessionRepo, userRepo, racketRepo, cfg.SessionCategoryResolver, contentSecurityService)
	statsService := service.NewStatsService(sessionRepo, racketRepo, cfg.SessionCategoryResolver)
	racketService := service.NewRacketService(racketRepo, userRepo, contentSecurityService, cfg.PolyesterStringHealth)
	homeService := service.NewHomeService(sessionRepo, racketRepo, cfg.SessionCategoryResolver)
	enumService := service.NewEnumService(cfg.SessionCategoryResolver)

	authHandler := handler.NewAuthHandler(authService)
	sessionHandler := handler.NewSessionHandler(sessionService)
	statsHandler := handler.NewStatsHandler(statsService)
	racketHandler := handler.NewRacketHandler(racketService)
	homeHandler := handler.NewHomeHandler(homeService)
	enumHandler := handler.NewEnumHandler(enumService)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		logger.Debug("GET /health start clientIP=%s", c.ClientIP())
		logger.Debug("GET /health success status=ok")
		response.OK(c, gin.H{"status": "ok"})
	})

	startRacketImageMigration(cfg, racketRepo)

	api := r.Group("/api")
	api.POST("/auth/wechat-login", authHandler.WechatLogin)
	api.POST("/auth/phone-login", authHandler.PhoneLogin)
	api.GET("/session-config", enumHandler.SessionConfig)

	authed := api.Group("")
	authed.Use(middleware.Auth(authService))
	{
		authed.GET("/auth/me", authHandler.Me)
		authed.POST("/auth/logout", authHandler.Logout)
		authed.PUT("/auth/profile", authHandler.UpdateProfile)

		authed.GET("/home/summary", homeHandler.Summary)

		authed.GET("/sessions", sessionHandler.List)
		authed.POST("/sessions", sessionHandler.Create)
		authed.GET("/sessions/latest", sessionHandler.Latest)
		authed.GET("/sessions/calendar", sessionHandler.Calendar)
		authed.GET("/sessions/:id", sessionHandler.Get)
		authed.PUT("/sessions/:id", sessionHandler.Update)
		authed.DELETE("/sessions/:id", sessionHandler.Delete)
		authed.GET("/stats/month", statsHandler.Month)
		authed.GET("/stats/charts", statsHandler.Charts)

		authed.GET("/racket-brands", racketHandler.Brands)
		authed.GET("/racket-series", racketHandler.Series)
		authed.GET("/racket-library", racketHandler.Library)
		authed.GET("/racket-library/stats", racketHandler.LibraryStats)
		authed.GET("/my-rackets", racketHandler.MyRackets)
		authed.GET("/my-rackets/primary", racketHandler.MyPrimaryRacket)
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
		authed.PUT("/rackets/:id/stringing-records/:recordId", racketHandler.UpdateStringingRecord)
		authed.DELETE("/rackets/:id/stringing-records/:recordId", racketHandler.DeleteStringingRecord)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

// startRacketImageMigration 配置开启时，注册球拍库图片每日迁移定时任务。
func startRacketImageMigration(cfg config.Config, racketRepo *repository.RacketRepository) {
	if !cfg.StorageEnabled {
		return
	}
	if cfg.StorageEnvID == "" || cfg.StorageBucket == "" || cfg.StorageRegion == "" || cfg.StorageFolder == "" || cfg.StorageSchedule == "" {
		log.Fatalf("storage enabled but envId/bucket/region/folder/schedule missing")
	}

	cloudStore := storage.New(storage.Config{
		EnvID:  cfg.StorageEnvID,
		Bucket: cfg.StorageBucket,
		Region: cfg.StorageRegion,
		Folder: cfg.StorageFolder,
	})
	migrationService := service.NewRacketImageMigration(racketRepo, cloudStore)

	runOnce := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
		defer cancel()
		if _, err := migrationService.Run(ctx); err != nil {
			logger.Error("racket image migration error: %v", err)
		}
	}

	if cfg.StorageRunOnStart {
		logger.Debug("racket image migration runOnStart triggered")
		go runOnce()
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	scheduler, err := gocron.NewScheduler(gocron.WithLocation(loc))
	if err != nil {
		log.Fatalf("create gocron scheduler: %v", err)
	}
	if _, err := scheduler.NewJob(
		gocron.CronJob(cfg.StorageSchedule, false),
		gocron.NewTask(runOnce),
	); err != nil {
		log.Fatalf("register racket image migration job: %v", err)
	}
	scheduler.Start()
	logger.Debug("racket image migration scheduler started schedule=%s", cfg.StorageSchedule)
}
