package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	http2 "github.com/zhosyaaa/RoommateTap/internal/delivery/http"
	"github.com/zhosyaaa/RoommateTap/internal/repository"
	"github.com/zhosyaaa/RoommateTap/internal/server"
	"github.com/zhosyaaa/RoommateTap/internal/service"
	"github.com/zhosyaaa/RoommateTap/pkg/auth"
	"github.com/zhosyaaa/RoommateTap/pkg/cache"
	"github.com/zhosyaaa/RoommateTap/pkg/database"
	"github.com/zhosyaaa/RoommateTap/pkg/database/redis"
	"github.com/zhosyaaa/RoommateTap/pkg/email/smtp"
	"github.com/zhosyaaa/RoommateTap/pkg/hash"
	"github.com/zhosyaaa/RoommateTap/pkg/logger"
	"github.com/zhosyaaa/RoommateTap/pkg/otp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	fmt.Println("---------RUN-------------")
	cfg, err := config.Init(configPath)
	if err != nil {
		log.Fatalf("Error getting configs: %v", err)
	}
	fmt.Println("cfg: ", cfg)

	db, err := database.GetDBInstance(*cfg)
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}

	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)

		return
	}
	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	rdb := redis.NewRedisClient(cfg.Redis)

	fmt.Println("-------------------")
	fmt.Println("rdb: ", rdb)
	fmt.Println("-------------------")

	repos := repository.NewRepositories(db)
	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Pass, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)

		return
	}

	otpGenerator := otp.NewGOTPGenerator()
	services := service.NewServices(service.Deps{
		Repos:                  repos,
		Cache:                  memCache,
		Hasher:                 hasher,
		TokenManager:           tokenManager,
		EmailSender:            emailSender,
		EmailConfig:            cfg.Email,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		CacheTTL:               int64(cfg.CacheTTL.Seconds()),
		OtpGenerator:           otpGenerator,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		Environment:            cfg.Environment,
		Domain:                 cfg.HTTP.Host,
		RedisClient:            rdb,
	})

	handlers := http2.NewHandler(services, tokenManager)
	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg))

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		fmt.Errorf("failed to stop server: %v", err)
	}

}
