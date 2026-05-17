# 📧 Turar's Microservices – Complete Build Guide

## Overview
Turar handles two microservices:
1. **Notification Service** – email notifications via SMTP and NATS
2. **Movie Service** – movie catalog with 50 films

Both follow Clean Architecture, use PostgreSQL, Redis, gRPC, and are fully tested.

## Notification Service
### Features
- 3 gRPC methods: SendBookingConfirmation, GetEmailLogs, HealthCheck
- NATS subscriber for `booking.created` events
- SMTP email sending (Gmail / Microsoft)
- PostgreSQL for email log storage
- Unit and integration tests
- Dockerized

### How to test
```bash
cd notification-service
# 1. Start infrastructure
docker-compose up -d postgres-notification nats
# 2. Apply migrations
migrate -path migrations -database "postgres://notification:notification@localhost:5435/notificationdb?sslmode=disable" up
# 3. Run service (with real SMTP if available)
SMTP_USER=your@gmail.com SMTP_PASS=your_app_password \
DATABASE_URL="postgres://notification:notification@localhost:5435/notificationdb?sslmode=disable" \
NATS_URL=nats://localhost:4222 GRPC_PORT=50054 \
go run ./cmd/main.go
# 4. gRPC checks
grpcurl -plaintext localhost:50054 list
grpcurl -plaintext -d '{"booking_id":"123","recipient_email":"test@example.com"}' localhost:50054 notification.NotificationService/SendBookingConfirmation
grpcurl -plaintext localhost:50054 notification.NotificationService/GetEmailLogs
# 5. Simulate NATS event
nats pub booking.created "booking-456"
# 6. Tests
go test ./internal/service/ -v
INTEGRATION=1 go test ./internal/repository/ -v
# 7. Frontend: open frontend/index.html in browser
Movie Service (also Turar)
Features
3 gRPC methods: GetMovie, ListMovies, SearchMovies

PostgreSQL with 50 movies seed

Redis caching

Unit tests (included in service code)

Dockerized

How to test
bash
cd movie-service
# 1. Start database
docker-compose up -d postgres-movie redis
# 2. Apply migrations
migrate -path migrations -database "postgres://movie:movie@localhost:5434/moviedb?sslmode=disable" up
# 3. Seed 50 movies
go run ./cmd/seed/main.go
# 4. Run service
DATABASE_URL="postgres://movie:movie@localhost:5434/moviedb?sslmode=disable" \
REDIS_ADDR=localhost:6379 GRPC_PORT=50051 \
go run ./cmd/main.go
# 5. gRPC checks
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"id":"<movie-id>"}' localhost:50051 movie.MovieService/GetMovie
grpcurl -plaintext -d '{"page":1,"limit":5}' localhost:50051 movie.MovieService/ListMovies
grpcurl -plaintext -d '{"query":"Action"}' localhost:50051 movie.MovieService/SearchMovies
# 6. Tests
go test ./...
Clean Architecture
Handler (gRPC server) → Service (business logic) → Repository (PostgreSQL) + Cache (Redis)

Interfaces for testability (in Notification Service)

Migrations via golang-migrate
