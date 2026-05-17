# 🎬 Lenara's Showtime Service – Complete Build Guide
## Overview
This microservice manages movie showtimes and seat availability.  
It follows **Clean Architecture**, uses **PostgreSQL + Redis**, and exposes **6 gRPC endpoints**.
## Prerequisites
- Go 1.25+, Docker & Docker Compose, protoc with Go plugins, golang-migrate, grpcurl
## Project Structure
showtime-service/
├── cmd/
│ ├── main.go # entry point
│ ├── seed/main.go # populates 200 showtimes with seats
│ └── initdb/main.go # creates database & user
├── internal/
│ ├── server/grpc.go # gRPC handler registration
│ ├── service/ # business logic (interfaces)
│ ├── repository/ # PostgreSQL queries
│ └── cache/ # Redis cache
├── proto/showtime.proto # 6 RPC definitions
├── migrations/ # SQL up/down migrations
├── Dockerfile
├── docker-compose.yml
└── go.mod
## 6 gRPC Endpoints
1. `GetShowtime` – fetch a showtime by ID
2. `ListShowtimesByMovie` – all showtimes for a movie
3. `GetSeatAvailability` – seats for a specific showtime
4. `CreateShowtime` – add a new showtime
5. `UpdateShowtime` – modify existing showtime
6. `DeleteShowtime` – remove a showtime
## Database & Cache
- **PostgreSQL**: tables `showtimes` and `seats`, migrations handled by `golang-migrate`
- **Redis**: caches `GetShowtime` results with 5min TTL
## Testing
- **Unit tests** (`internal/service/showtime_test.go`): mock repository & cache
- **Integration tests** (`internal/repository/postgres_test.go`): uses testcontainers-go to spin up a real PostgreSQL
## Step-by-step Terminal Commands (Demo)
### 1. Start infrastructure
```bash
cd showtime-service
docker-compose up -d postgres-showtime redis
2. Apply migrations
bash
migrate -path migrations -database "postgres://showtime:showtime@localhost:5433/showtimedb?sslmode=disable" up
3. Run the service
bash
DATABASE_URL="postgres://showtime:showtime@localhost:5433/showtimedb?sslmode=disable" \
REDIS_ADDR=localhost:6379 \
GRPC_PORT=50052 \
go run ./cmd/main.go
4. Test gRPC endpoints
bash
grpcurl -plaintext localhost:50052 list
grpcurl -plaintext -d '{"movie_id":"m1","hall":"Hall A","start_time":"2026-06-01T12:00:00Z","available_seats":50}' localhost:50052 showtime.ShowtimeService/CreateShowtime
grpcurl -plaintext -d '{"movie_id":"m1"}' localhost:50052 showtime.ShowtimeService/ListShowtimesByMovie
5. Run tests
bash
go test ./internal/service/ -v
INTEGRATION=1 go test ./internal/repository/ -v
Clean Architecture Layers
Handler: gRPC server delegates to service
Service: business logic, depends on interfaces
Repository: PostgreSQL implementation
Cache: Redis implementation
