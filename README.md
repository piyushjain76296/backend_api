# Ticket Management System API

This repository contains a REST API for a ticket management system, built as my submission for the Backend Engineering Intern assignment. My primary focus during this project was ensuring strict API contract compliance, robust security, and writing clean, idiomatic Go code.

## Key Features Implemented

* **User Authentication**: Secure registration and login using bcrypt for password hashing and JWT for stateless session management.
* **Ticket Management**: Full ticket creation and retrieval with strict status lifecycle enforcement (`open` -> `in_progress` -> `closed`).
* **Data Isolation**: Robust authorization checks ensuring users can only view and modify their own tickets.
* **Integrated Frontend UI**: A premium, glassmorphic single-page application built with Vanilla HTML/CSS/JS, served directly from the Go backend.
* **Containerization**: A multi-stage Docker build for a minimal and secure production image.

## Tech Stack

* **Language**: Go 1.22
* **Routing**: Native `net/http` (leveraging Go 1.22 routing features)
* **Database**: SQLite (`github.com/mattn/go-sqlite3`)
* **Deployment**: Docker, Render

## Project Architecture

I chose a lightweight, layered architecture to keep the codebase maintainable without over-engineering:

```text
ticket-system/
├── cmd/server/main.go          # Application entry point
├── frontend/                   # Vanilla HTML/CSS/JS static files for the UI
├── internal/
│   ├── api/                    # HTTP handlers and JWT middleware
│   ├── core/                   # Domain models and business logic (services)
│   └── infrastructure/         # SQLite database configuration and repositories
├── test/api_test.go            # Integration tests
└── Dockerfile                  # Multi-stage build configuration
```

## Running the Project

### Prerequisites
* Go 1.22+
* GCC (required for SQLite CGO compilation)

### Local Setup
1. Clone the repository and download dependencies:
   ```bash
   git clone https://github.com/piyushjain76296/backend_api.git
   cd backend_api
   go mod download
   ```
2. Start the server:
   ```bash
   go run ./cmd/server
   ```
3. Run the integration tests:
   ```bash
   go test ./... -v
   ```

### Docker Setup
To run the application using Docker:
```bash
docker build -t ticket-system .
docker run -p 8080:8080 ticket-system
```

## API Endpoints

The API consumes and produces `application/json`. Requests failing validation or authorization return clear JSON error messages and appropriate HTTP status codes.

### Public Routes
* `GET /` - Serves the frontend web interface.
* `GET /health` - Health check endpoint.
* `POST /auth/register` - Registers a new user. Requires `email` and `password`.
* `POST /auth/login` - Authenticates a user and returns a JWT.

### Protected Routes (Requires `Authorization: Bearer <token>`)
* `POST /tickets` - Creates a new ticket (automatically initialized to `open`).
* `GET /tickets` - Retrieves all tickets owned by the authenticated user.
* `GET /tickets/{id}` - Retrieves a specific ticket if owned by the user.
* `PATCH /tickets/{id}/status` - Updates ticket status following strict lifecycle rules.

## Design Decisions & Assumptions

* **Database Choice**: I chose SQLite because it perfectly fulfills the requirement for a simple persistent storage solution while making the setup process entirely frictionless for reviewers. It avoids the need to configure external database containers.
* **Routing**: Instead of relying on third-party routers like `chi` or `gorilla/mux`, I utilized the native HTTP method routing introduced in Go 1.22 (`GET /path`) to minimize external dependencies.
* **Security & Error Handling**: To prevent data leakage, if a user attempts to fetch a ticket ID belonging to another user, my API intentionally returns a `404 Not Found` rather than a `403 Forbidden`. This obscures the existence of other users' data.
* **Status Enforcement**: The transition logic (`open` -> `in_progress` -> `closed`) is strictly handled in the service layer. I ensured that closed tickets are in a terminal state and cannot be reopened.

## Live Deployment

I have containerized and deployed the application to Render. You can test the live endpoints here:

* **Deployed Application**: `https://ticket-system-api-dfid.onrender.com`
* **Health Check**: `https://ticket-system-api-dfid.onrender.com/health`
