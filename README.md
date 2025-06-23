# todo-api

A simple RESTful API for managing users and todo items, built with Go.

## Features

- **User Registration:** Register new users via `/register`.
- **User Login:** Authenticate existing users via `/login`.
- **Todo Management:** Add, update, fetch, and delete todo items via `/todos` endpoints.
- **Rate Limiting:** Middleware integration for security and rate limiting.
- **Docker Support:** Includes Dockerfile and `docker-compose.yml` for easy containerization.

## Endpoints

| Method | Endpoint            | Description                  |
|--------|---------------------|------------------------------|
| POST   | `/register`         | Register a new user          |
| POST   | `/login`            | User login                   |
| POST   | `/todos`            | Add a new todo item          |
| PUT    | `/todos/{id}`       | Update a todo item           |
| DELETE | `/todos/{id}`       | Delete a todo item           |
| GET    | `/todos`            | Fetch all todo items         |

## Project Structure

- `main.go` - Application entrypoint, routes, and server setup.
- `data/` - User and todo data handling.
- `encrypt/` - Middleware and security logic.
- `schema/` - (Likely) database schema or models.
- `Dockerfile` & `docker-compose.yml` - Containerization support.

## Running Locally

### Prerequisites

- Go (latest version recommended)
- Docker (optional, for containerization)

### Clone & Run

```bash
git clone https://github.com/Prateesh94/todo-api.git
cd todo-api
go run main.go
```

The server will be running at `http://localhost:8080`.

### Using Docker

```bash
docker-compose up --build
```

## Testing

Tests are provided in `main_test.go`. Run them using:

```bash
go test
```

## License

MIT

---

**Note:** This README is based on limited file visibility. Please update with additional details as needed.  
See more files or details in the [GitHub UI](https://github.com/Prateesh94/todo-api/tree/main/).