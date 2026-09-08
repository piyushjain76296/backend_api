package main

import (
	"log"
	"net/http"
	"os"

	"ticket-system/internal/api/handlers"
	"ticket-system/internal/api/middleware"
	"ticket-system/internal/core/services"
	"ticket-system/internal/infrastructure/database"
	"ticket-system/internal/infrastructure/repositories"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./tickets.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	// Initialize Repositories
	userRepo := repositories.NewUserRepository(database.DB)
	ticketRepo := repositories.NewTicketRepository(database.DB)

	// Initialize Services
	authService := services.NewAuthService(userRepo)
	ticketService := services.NewTicketService(ticketRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	// Public routes
	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Protected routes (require JWT)
	mux.Handle("POST /tickets", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.CreateTicket)))
	mux.Handle("GET /tickets", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.ListTickets)))
	mux.Handle("GET /tickets/{id}", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.GetTicket)))
	mux.Handle("PATCH /tickets/{id}/status", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.UpdateStatus)))

	log.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
