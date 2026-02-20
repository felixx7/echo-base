package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"echo-base/config"
	"echo-base/database"
	"echo-base/domain/repository"
	"echo-base/domain/usecase"
	"echo-base/http/handler"
	"echo-base/http/middleware"
	"echo-base/http/routes"
)

func main() {
	// Load config
	cfg := config.Load()
	dbCfg := config.LoadDatabaseConfig()

	// Initialize database
	db, err := database.Connect(dbCfg)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("error running migrations: %v", err)
	}

	// Initialize Echo instance
	e := echo.New()
	e.HideBanner = true

	// Initialize repositories (using PostgreSQL)
	userRepo := repository.NewUserRepository(db)

	// Initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)

	// Register global middleware
	e.Use(middleware.LoggerMiddleware())
	e.Use(middleware.RecoverMiddleware())
	e.Use(middleware.CORSMiddleware())

	// Register routes (moved to http/routes)
	routes.RegisterRoutes(e, userHandler)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("[%s] Server running on %s\n", cfg.AppName, addr)
	if err := e.Start(addr); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
