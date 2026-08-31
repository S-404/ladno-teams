package main

import (
	"context"
	"ladno-teams/infrastructure"
	"ladno-teams/infrastructure/database"
	"ladno-teams/internal/config"
	"ladno-teams/internal/handler"
	"ladno-teams/internal/repository"
	"ladno-teams/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.InitDB(database.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSsl,
		User:     cfg.DBUser,
		Password: cfg.DBPass,
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	repositories := repository.NewRepository(db)
	services := service.NewService(*cfg, repositories)
	if err := services.User.EnsureAdmin(*cfg); err != nil {
		log.Fatalf("Failed to ensure admin: %v", err)
	}
	handlers := handler.NewHandler(*cfg, services)

	srv := infrastructure.NewHttpServer(cfg.HTTPPort, handlers.InitRoutes())
	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()
	log.Printf("HttpServer listening %s...", cfg.HTTPPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("App shutting down...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("error occurred on db connection close: %s", err.Error())
	}
}
