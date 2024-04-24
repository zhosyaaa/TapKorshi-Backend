package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhosyaaa/RoommateTap/internal/config"
	"github.com/zhosyaaa/RoommateTap/internal/server"
	"github.com/zhosyaaa/RoommateTap/pkg/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		return
	}

	db, err := database.GetDBInstance(*cfg)
	if err != nil {
		log.Fatalf("Error initializing database connection: %v", err)
	}

	// Check if the database connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// HTTP Server
	srv := server.NewServer(cfg)

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
