# Thrive Track Workout 💪

A RESTful workout tracking API built with Go and PostgreSQL, designed to help fitness enthusiasts log and monitor their training progress.

[![Go Version](https://img.shields.io/badge/go-1.22-blue)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features ✨
- RESTful API endpoints for workout management
- PostgreSQL database integration
- Dockerized development environment
- Health check endpoint monitoring
- Chi router for efficient routing
- Structured logging
- Proper connection pooling configuration

## Tech Stack 🛠️
- **Language:** Go 1.22
- **Database:** PostgreSQL
- **ORM:** pgx (PostgreSQL driver)
- **Router:** Chi
- **Containerization:** Docker + Docker Compose

## Getting Started 🚀

### Prerequisites
- Go 1.22+
- PostgreSQL 15+
- Docker (optional)

### Installation

#### Using Docker (Recommended):
```bash
docker-compose up -d
``` 
#### Manual Setup:
1. Clone repository:

```bash
git clone https://github.com/nickemma/thrive-track-workout.git
cd thrive-track-workout
``` 
2. Install dependencies:
```
go mod download
```
3. Start PostgreSQL database:
```
docker run -d -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres postgres:15
```
4. Run application:
```
go run main.go
```
### Environment Variables
Create .env file:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=postgres
DB_NAME=postgres
```
### API Reference 📚
Endpoints

- GET /healthcheck - Service health monitoring
- GET /api/workouts - List all workouts
- POST /api/workouts - Create new workout
- GET /api/workouts/{id} - Get specific workout
- PUT /api/workouts/{id} - Update workout
- DELETE /api/workouts/{id} - Delete workout

### Running Tests 🧪
```
go test -v ./...
```
### Contributing 🤝
1. Fork the project
2. Create your feature branch (git checkout -b feature/AmazingFeature)
3. Commit changes (git commit -m 'Add AmazingFeature')
4. Push branch (git push origin feature/AmazingFeature)
5. Open Pull Request