# Go-React Starter

A full-stack application starter template featuring a Go (Gin) backend and a React (Vite) frontend, complete with Docker support.

## Tech Stack

- **Backend:** Go (Golang), Gin Framework, GORM
- **Frontend:** React, Vite, Tailwind CSS
- **Database:** MySQL
- **Caching:** Redis
- **Infrastructure:** Docker, Docker Compose, Nginx

## Prerequisites

- **Docker** and **Docker Compose** (Recommended)
- **Go 1.25+** (For local backend development)
- **Node.js 20+** (For local frontend development)
- **Make** (Optional, for using Makefile shortcuts)

## Getting Started

### Using Docker (Recommended)

The easiest way to run the application is using Docker Compose.

1.  **Start the application:**
    ```bash
    make docker-up
    # OR
    docker-compose up -d --build
    ```

2.  **Access the application:**
    - Frontend: [http://localhost:3000](http://localhost:3000)
    - Backend API: [http://localhost:8080](http://localhost:8080)
    - MySQL: Port `3307`
    - Redis: Port `6380`

3.  **View Logs:**
    ```bash
    make docker-logs
    # OR
    docker-compose logs -f
    ```

4.  **Stop the application:**
    ```bash
    make docker-down
    # OR
    docker-compose down
    ```

5.  **Run Migrations & Seeding:**
    - To run migrations:
      ```bash
      make docker-migrate
      # OR
      docker-compose exec backend ./migrate
      ```
    - To run seeds:
      ```bash
      make docker-seed
      # OR
      docker-compose exec backend ./seeder
      ```

### Manual Setup

If you prefer to run services individually without Docker:

#### Backend

1.  Navigate to the backend directory:
    ```bash
    cd backend
    ```
2.  Install dependencies:
    ```bash
    go mod tidy
    ```
3.  Configure environment variables:
    - Copy `.env-example` to `.env` (if not already done).
    - Update database and redis credentials to match your local services.
4.  Run the server:
    ```bash
    go run ./cmd/api/main.go
    ```

#### Frontend

1.  Navigate to the frontend directory:
    ```bash
    cd frontend
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the development server:
    ```bash
    npm run dev
    ```

## Project Structure

```
├── backend/            # Go backend application
│   ├── cmd/            # Application entrypoints
│   ├── internal/       # Private application code
│   └── ...
├── frontend/           # React frontend application
│   ├── src/            # Source code
│   ├── public/         # Static assets
│   └── ...
├── docker-compose.yml  # Docker services configuration
├── Makefile            # Shortcuts for common commands
└── README.md           # Project documentation
```

## Environment Variables

Key environment variables used in `backend/.env`:

| Variable     | Description                        | Default          |
| :----------- | :--------------------------------- | :--------------- |
| `APP_PORT`   | Port for the backend server        | `8080`           |
| `DB_HOST`    | Database host (use `db` in Docker) | `localhost`      |
| `DB_PORT`    | Database port                      | `3306`           |
| `REDIS_HOST` | Redis host (use `redis` in Docker) | `localhost`      |
| `JWT_SECRET` | Secret key for JWT tokens          | `secret-jwt-key` |

## Makefile Commands

- `make run-backend`: Run backend locally
- `make run-frontend`: Run frontend locally
- `make docker-up`: Start Docker services
- `make docker-down`: Stop Docker services
- `make docker-build`: Rebuild Docker images
- `make docker-logs`: View container logs
- `make docker-migrate`: Run database migrations (Docker)
- `make docker-seed`: Run database seeder (Docker)


## BUGS
1. login with invalid password will crash the query OK
2. Request id empty on system logs OK
3. email input background color only on login page
4. permission ubah ke bitwise operator


## TODO
1. Implement reset password OK
2. implement remember me OK
3. Implement 2fa by email anf HOTP
4. Storage with one time link and user accesspermission
5. web description untuk set title, favicon, seo
6. reset 2fa secret key
7. module api key
8. native event for notifications