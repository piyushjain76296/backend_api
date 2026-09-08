package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"ticket-system/internal/api/handlers"
	"ticket-system/internal/api/middleware"
	"ticket-system/internal/core/models"
	"ticket-system/internal/core/services"
	"ticket-system/internal/infrastructure/database"
	"ticket-system/internal/infrastructure/repositories"
)

func setupTestServer(t *testing.T) *http.ServeMux {
	// Use in-memory SQLite for tests
	err := database.InitDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}

	// Set JWT secret
	os.Setenv("JWT_SECRET", "test-secret")

	userRepo := repositories.NewUserRepository(database.DB)
	ticketRepo := repositories.NewTicketRepository(database.DB)

	authService := services.NewAuthService(userRepo)
	ticketService := services.NewTicketService(ticketRepo)

	authHandler := handlers.NewAuthHandler(authService)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.HealthHandler)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("POST /tickets", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.CreateTicket)))
	mux.Handle("GET /tickets", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.ListTickets)))
	mux.Handle("GET /tickets/{id}", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.GetTicket)))
	mux.Handle("PATCH /tickets/{id}/status", middleware.AuthMiddleware(http.HandlerFunc(ticketHandler.UpdateStatus)))

	return mux
}

func doRequest(mux *http.ServeMux, method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}

	req := httptest.NewRequest(method, url, &buf)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func TestCompleteFlow(t *testing.T) {
	mux := setupTestServer(t)

	// 1. Register User A
	registerBody := handlers.AuthRequest{Email: "usera@example.com", Password: "password123"}
	rr := doRequest(mux, "POST", "/auth/register", registerBody, "")
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %v", rr.Code)
	}

	// 2. Duplicate Registration
	rr = doRequest(mux, "POST", "/auth/register", registerBody, "")
	if rr.Code != http.StatusConflict {
		t.Fatalf("Expected status 409 for duplicate, got %v", rr.Code)
	}

	// 3. Login User A
	rr = doRequest(mux, "POST", "/auth/login", registerBody, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for login, got %v", rr.Code)
	}
	var loginResp handlers.LoginResponse
	json.NewDecoder(rr.Body).Decode(&loginResp)
	tokenA := loginResp.Token

	// 4. Register and Login User B
	registerBodyB := handlers.AuthRequest{Email: "userb@example.com", Password: "password123"}
	doRequest(mux, "POST", "/auth/register", registerBodyB, "")
	rr = doRequest(mux, "POST", "/auth/login", registerBodyB, "")
	json.NewDecoder(rr.Body).Decode(&loginResp)
	tokenB := loginResp.Token

	// 5. Create Ticket A
	ticketReq := handlers.CreateTicketRequest{Title: "Ticket A", Description: "Desc A"}
	rr = doRequest(mux, "POST", "/tickets", ticketReq, tokenA)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected status 201 for ticket creation, got %v", rr.Code)
	}
	var ticketA models.Ticket
	json.NewDecoder(rr.Body).Decode(&ticketA)
	if ticketA.Status != "open" {
		t.Fatalf("Expected ticket to be open, got %v", ticketA.Status)
	}

	// 6. List tickets as User A
	rr = doRequest(mux, "GET", "/tickets", nil, tokenA)
	var listResp handlers.TicketsListResponse
	json.NewDecoder(rr.Body).Decode(&listResp)
	if len(listResp.Tickets) != 1 {
		t.Fatalf("Expected 1 ticket for User A, got %v", len(listResp.Tickets))
	}

	// 7. Attempt to access Ticket A using User B token
	rr = doRequest(mux, "GET", "/tickets/1", nil, tokenB)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 for User B accessing User A ticket, got %v", rr.Code)
	}

	// 8. Update Ticket A as User A
	updateReq := handlers.UpdateStatusRequest{Status: "in_progress"}
	rr = doRequest(mux, "PATCH", "/tickets/1/status", updateReq, tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 for status update, got %v", rr.Code)
	}

	// 9. Close Ticket A
	updateReq = handlers.UpdateStatusRequest{Status: "closed"}
	rr = doRequest(mux, "PATCH", "/tickets/1/status", updateReq, tokenA)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 for closing ticket, got %v", rr.Code)
	}

	// 10. Attempt to reopen Ticket A
	updateReq = handlers.UpdateStatusRequest{Status: "open"}
	rr = doRequest(mux, "PATCH", "/tickets/1/status", updateReq, tokenA)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 for invalid transition, got %v", rr.Code)
	}
}

func TestHealthCheck(t *testing.T) {
	mux := setupTestServer(t)
	rr := doRequest(mux, "GET", "/health", nil, "")
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for health check, got %v", rr.Code)
	}
}
