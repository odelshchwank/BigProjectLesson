# Golang ToDo API

This repository contains a full-featured ToDo list application built with Go. It serves as a comprehensive example of building a RESTful API following clean architecture principles, complete with a database, containerization, and a simple web frontend.

The application allows users to manage tasks, with capabilities for creating, viewing, updating, and deleting both users and their associated tasks. It also provides an endpoint for task statistics.

## Features

*   **User Management**: Full CRUD operations for users (`/users`).
*   **Task Management**: Full CRUD operations for tasks, linked to a user (`/tasks`).
*   **Task Statistics**: Endpoint to retrieve statistics like tasks created, completed, completion rate, and average completion time (`/statistics`).
*   **Pagination & Filtering**: List endpoints for users and tasks support pagination (`limit`, `offset`) and filtering (e.g., tasks by `user_id`).
*   **Partial Updates**: `PATCH` endpoints utilize a three-state logic (value set, value unset, value set to null) for flexible updates.
*   **API Documentation**: Automatically generated Swagger 2.0 documentation available via a UI.
*   **Web UI**: A single-page HTML interface (`/public/index.html`) for interacting with the API.

## Project Structure

The project is structured using clean architecture principles to promote separation of concerns, testability, and maintainability.

```
.
├── cmd/                # Application entry points
├── docs/               # Swagger API documentation files
├── internal/
│   ├── core/           # Core application components (server, logger, db pool)
│   └── features/       # Business logic for distinct features (users, tasks)
│       ├── repository/ # Data access layer (e.g., Postgres implementation)
│       ├── service/    # Business logic layer
│       └── transport/  # API layer (HTTP handlers, DTOs)
├── migrations/         # Database migration files
├── public/             # Static frontend assets (HTML, CSS, JS)
├── docker-compose.yaml # Docker Compose configuration
└── Makefile            # Helper commands for development and deployment
```

*   **`cmd/todoapp`**: The main application executable. It initializes all components (config, logger, database, services, and HTTP server) and wires them together.
*   **`internal/core`**: Contains foundational code shared across the application, such as the HTTP server setup, middleware, database connection pool, and configuration loading.
*   **`internal/features`**: Each sub-directory represents a distinct feature (e.g., `users`, `tasks`). Within each feature, the code is layered into:
    *   **`transport/http`**: Defines HTTP handlers, request/response DTOs, and routing. This is the entry point for API requests.
    *   **`service`**: Implements the core business logic for the feature, orchestrating operations and validations.
    *   **`repository/postgres`**: Manages data persistence and communication with the PostgreSQL database.
*   **`migrations`**: SQL scripts for creating and updating the database schema, managed via `migrate`.

## Tech Stack

*   **Language**: Go
*   **Database**: PostgreSQL
*   **API**: Standard `net/http` library
*   **Database Driver**: `pgx/v5`
*   **Containerization**: Docker & Docker Compose
*   **Logging**: `uber-go/zap`
*   **Configuration**: `kelseyhightower/envconfig`
*   **API Documentation**: `swaggo/swag`

## Getting Started

### Prerequisites

*   Docker
*   Docker Compose
*   Make (optional, for convenience)

### Configuration

1.  Create a `.env` file by copying the example:
    ```sh
    cp .env.example .env
    ```
2.  Edit the `.env` file and provide values for your PostgreSQL instance:
    ```env
    POSTGRES_USER=your_user
    POSTGRES_PASSWORD=your_password
    POSTGRES_DB=your_db
    ```

### Running with Docker Compose

This is the recommended way to run the application and its database.

1.  **Start the database container:**
    ```sh
    make env-up
    ```
    This will start a PostgreSQL instance in a Docker container.

2.  **Apply database migrations:**
    ```sh
    make migrate-up
    ```
    This command runs the SQL scripts in the `migrations` folder to set up the necessary tables.

3.  **Build and deploy the application:**
    ```sh
    make todoapp-deploy
    ```
    This builds the Go application Docker image and starts the `todoapp` container.

The application will be accessible at `http://localhost:5050`.

### Running Locally (for Development)

If you prefer to run the Go application directly on your host machine for development:

1.  **Start the database container:**
    ```sh
    make env-up
    ```

2.  **Forward the database port to your local machine:**
    ```sh
    make env-port-forwarder
    ```
    This allows the local Go application to connect to the PostgreSQL instance running in Docker.

3.  **Apply database migrations:**
    ```sh
    make migrate-up
    ```

4.  **Run the application:**
    ```sh
    make todoapp-run
    ```
    This will compile and run the `main.go` file. The server will start on `localhost:5050`.

## API Documentation

Once the application is running, you can access the interactive Swagger UI for detailed API documentation and testing:

**[http://80.249.146.85:5050/swagger/index.html](http://80.249.146.85:5050/swagger/index.html)**

To regenerate the Swagger documentation after making code changes, run:
```sh
make swagger-gen
```

## Available Commands

The `Makefile` provides several useful commands to streamline development:

*   `make env-up`: Starts the PostgreSQL service using Docker Compose.
*   `make env-down`: Stops the PostgreSQL service.
*   `make env-cleanup`: Stops and removes all Docker volumes, effectively wiping the database. **Use with caution.**
*   `make env-port-forwarder`: Starts a container to forward the database port to `localhost:5432`.
*   `make migrate-up`: Applies all pending database migrations.
*   `make migrate-down`: Rolls back the last applied database migration.
*   `make todoapp-deploy`: Builds and deploys the application as a Docker container.
*   `make todoapp-undeploy`: Stops the application container.
*   `make todoapp-run`: Runs the Go application locally (requires a running database).
*   `make swagger-gen`: Regenerates the Swagger documentation files from code annotations.
*   `make ps`: Shows the status of the running Docker containers for this project.