package main

import (
	"log"

	"balance/config"
	"balance/internal/database"
	"balance/internal/handlers"
	"balance/internal/repository"
	"balance/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	balanceService := service.NewBalanceService(userRepo, transactionRepo)
	balanceHandler := handlers.NewBalanceHandler(balanceService)

	router := gin.Default()

	router.GET("/user/:userId/balance", balanceHandler.GetBalance)
	router.POST("/user/:userId/transaction", balanceHandler.ProcessTransaction)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
