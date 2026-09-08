# Backend Intern Ticket System

A production-quality REST API in Golang for a simple ticket management system. Built for a Backend Intern Assignment, this system prioritizes correctness, exact API contract compliance, and clean Go code.

## Features

* User Registration & Login
* JWT-based Authentication
* Secure Password Hashing (bcrypt)
* Create, Read, and Update Tickets
* Strict Status Lifecycle Enforcement (`open` -> `in_progress` -> `closed`)
* Data Isolation & Ownership Authorization (Users can only access their own tickets)
* Complete API Contract Compliance
* Docker & Docker Compose Ready
* In-memory SQLite for testing, disk-backed SQLite for production

## Tech Stack

* **Language**: Go 1.22
* **Routing**: Standard Library `net/http` (with 1.22 enhanced routing)
* **Database**: SQLite (`github.com/mattn/go-sqlite3`)
* **Security**: `golang.org/x/crypto/bcrypt`, `github.com/golang-jwt/jwt/v5`
* **Deployment**: Docker

## Project Structure

```
ticket-system/
│
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
│
├── internal/
│   ├── api/
│   │   ├── handlers/         # HTTP handlers mapping endpoints to services
│   │   └── middleware/       # JWT Authentication middleware
│   │
│   ├── core/
│   │   ├── models/           # Domain entities (User, Ticket)
│   │   └── services/         # Business logic and validation
│   │
│   └── infrastructure/
│       ├── database/         # SQLite connection setup
│       └── repositories/     # Data access layer
│
├── test/
│   └── api_test.go           # End-to-end integration tests
│
├── .env.example              # Example environment variables
├── Dockerfile                # Multi-stage production build
└── README.md
```

## Setup & Run Locally

### Prerequisites
* Go 1.22 or higher
* GCC (for CGO / SQLite compilation)

### Steps

1. **Clone the repository:**
   ```bash
   git clone <REPOSITORY_URL>
   cd ticket-system
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   ```

3. **Set environment variables (optional):**
   Copy the example environment file and customize it.
   ```bash
   cp .env.example .env
   ```
   *Note: If no environment variables are set, the app defaults to port 8080, an insecure test JWT secret, and `./tickets.db` for the database.*

4. **Run the server:**
   ```bash
   go run ./cmd/server
   ```

5. **Run the tests:**
   ```bash
   go test ./... -v
   ```

## Docker

The application is containerized using a multi-stage Docker build to ensure a small final image footprint.

### Build the Image
```bash
docker build -t ticket-system .
```

### Run the Container
```bash
docker run -p 8080:8080 ticket-system
```

## API Documentation

All API responses are in JSON. In the event of an error, the API will respond with `{"error": "description"}` and an appropriate HTTP status code.

### Public Endpoints

#### 1. Health Check
* **Method**: `GET`
* **URL**: `/health`
* **Response**: `200 OK`
  ```json
  {
    "status": "ok"
  }
  ```

#### 2. Register
* **Method**: `POST`
* **URL**: `/auth/register`
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
* **Response**: `201 Created`
  ```json
  {
    "message": "user registered successfully"
  }
  ```

#### 3. Login
* **Method**: `POST`
* **URL**: `/auth/login`
* **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
* **Response**: `200 OK`
  ```json
  {
    "token": "<jwt-token>"
  }
  ```

### Protected Endpoints
**Authentication Required**: Header `Authorization: Bearer <token>` must be present.

#### 4. Create Ticket
* **Method**: `POST`
* **URL**: `/tickets`
* **Request Body**:
  ```json
  {
    "title": "Cannot login",
    "description": "I get a 500 error when clicking login."
  }
  ```
* **Response**: `201 Created`
  ```json
  {
    "id": 1,
    "user_id": 1,
    "title": "Cannot login",
    "description": "I get a 500 error when clicking login.",
    "status": "open",
    "created_at": "2026-09-08T12:00:00Z",
    "updated_at": "2026-09-08T12:00:00Z"
  }
  ```

#### 5. List Own Tickets
* **Method**: `GET`
* **URL**: `/tickets`
* **Response**: `200 OK`
  ```json
  {
    "tickets": [
      {
        "id": 1,
        "user_id": 1,
        "title": "Cannot login",
        "description": "...",
        "status": "open",
        "created_at": "...",
        "updated_at": "..."
      }
    ]
  }
  ```

#### 6. Get Ticket By ID
* **Method**: `GET`
* **URL**: `/tickets/{id}`
* **Response**: `200 OK` (Returns the single ticket object if owned by the authenticated user)

#### 7. Update Ticket Status
* **Method**: `PATCH`
* **URL**: `/tickets/{id}/status`
* **Request Body**:
  ```json
  {
    "status": "in_progress"
  }
  ```
* **Response**: `200 OK`
  ```json
  {
    "message": "status updated successfully"
  }
  ```

## Ticket Status Flow

The system enforces a strict, one-way state machine for ticket statuses:
```text
open -> in_progress -> closed
```

* New tickets are always forced to the `open` state upon creation, regardless of user input.
* `open` can transition to `in_progress`.
* `in_progress` can transition to `closed`.
* `closed` is a terminal state. Closed tickets can **never** be reopened or changed to `in_progress`.
* Invalid transitions will result in a `400 Bad Request`.

## Deployment Instructions

To deploy this API for free, you can follow these steps using **Render.com**:

1. Push this repository to GitHub.
2. Sign up / log in to [Render](https://render.com/).
3. Click **New** -> **Web Service**.
4. Connect your GitHub account and select the repository.
5. In the configuration:
   * **Name**: `ticket-system-api`
   * **Language**: `Docker`
   * **Region**: `Oregon (US West)` (or closest)
   * **Branch**: `main`
6. Click **Advanced** and set the Environment Variables:
   * `PORT`: `8080`
   * `JWT_SECRET`: Generate a secure random string and paste it here.
   * `DATABASE_PATH`: `./data/tickets.db`
7. Click **Create Web Service**. Render will automatically build the `Dockerfile` and deploy the service.

Once deployed, update the section below with your live URLs.

```text
GitHub Repository:
<YOUR_GITHUB_URL>

Deployed Application:
<YOUR_RENDER_URL>

Health Check:
<YOUR_RENDER_URL>/health
```

## Assumptions & Design Decisions
* **Database**: SQLite was chosen to simplify the project setup process and avoid requiring external dependencies like a PostgreSQL container.
* **Architecture**: A layered architecture (Handlers -> Services -> Repositories) was adopted to ensure separation of concerns and testability, without going into over-engineered hexagonal or clean architecture patterns.
* **Routing**: Adopted Go 1.22's native `net/http` method routing (`GET /path`) to completely eliminate the need for third-party routers like `chi` or `gorilla/mux`.
* **Security**: Not leaking data. If a user tries to access another user's ticket by ID, the system returns a `404 Not Found` rather than a `403 Forbidden` to prevent leaking the existence of other users' tickets.
