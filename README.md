# Cinema System – Final Documentation

## Team Members
- Lenara: Showtime Service + API Gateway method GetShowtimeWithMovie
- Turar: Movie Service + Notification Service + API Gateway method SendBookingNotification
- Almas: Booking Service + Monitoring (Grafana, Prometheus, Jaeger, Loki) + API Gateway method GetBookingDetails

## GitHub Repository
https://github.com/LenarasCode/ADP2_finalproject_cinemasystem

---

## Prerequisites
Install the following tools on your machine before starting:
- Go 1.25 or higher
- Docker Desktop (optional, for full infrastructure)
- Docker Compose v2
- protoc with protoc-gen-go and protoc-gen-go-grpc plugins
- golang-migrate CLI
- grpcurl
- Git

---

## Infrastructure Setup

Start the shared PostgreSQL container (used by all services when Docker is limited):
docker start microservices-sre-postgres-1

If the container does not exist, create it:
docker run -d --name microservices-sre-postgres-1 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=secret -p 5432:5432 postgres:15

---

## LENARA – Showtime Service

### Step 1: Create database and user
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE DATABASE showtimedb;"
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE USER showtime WITH PASSWORD 'showtime'; GRANT ALL PRIVILEGES ON DATABASE showtimedb TO showtime;"

### Step 2: Apply migrations
cd showtime-service
export DATABASE_URL="postgres://showtime:showtime@localhost:5432/showtimedb?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up

### Step 3: Seed data (200 showtimes with seats)
export DATABASE_URL="postgres://showtime:showtime@localhost:5432/showtimedb?sslmode=disable"
go run ./cmd/seed/main.go

### Step 4: Start Showtime Service
DATABASE_URL="postgres://showtime:showtime@localhost:5432/showtimedb?sslmode=disable" REDIS_ADDR=localhost:6379 GRPC_PORT=50052 go run ./cmd/main.go

### Step 5: Test gRPC endpoints
List methods:
grpcurl -plaintext localhost:50052 list

Create a showtime:
grpcurl -plaintext -d '{"movie_id":"movie-123","hall":"Hall 5","start_time":"2026-05-18T10:00:00Z","available_seats":80}' localhost:50052 showtime.ShowtimeService/CreateShowtime

List showtimes for a movie:
grpcurl -plaintext -d '{"movie_id":"movie-123"}' localhost:50052 showtime.ShowtimeService/ListShowtimesByMovie

Get seat availability:
grpcurl -plaintext -d '{"showtime_id":"<showtime-id>"}' localhost:50052 showtime.ShowtimeService/GetSeatAvailability

### Step 6: Run unit tests
go test ./internal/service/ -v

### Step 7: Run integration tests (requires Docker)
INTEGRATION=1 go test ./internal/repository/ -v

---

## TURAR – Movie Service

### Step 1: Create database and user
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE DATABASE moviedb;"
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE USER movie WITH PASSWORD 'movie'; GRANT ALL PRIVILEGES ON DATABASE moviedb TO movie;"

### Step 2: Apply migrations
cd movie-service
export DATABASE_URL="postgres://movie:movie@localhost:5432/moviedb?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up

### Step 3: Seed 50 movies
export DATABASE_URL="postgres://movie:movie@localhost:5432/moviedb?sslmode=disable"
go run ./cmd/seed/main.go

### Step 4: Start Movie Service
DATABASE_URL="postgres://movie:movie@localhost:5432/moviedb?sslmode=disable" REDIS_ADDR=localhost:6379 GRPC_PORT=50051 go run ./cmd/main.go

### Step 5: Test gRPC endpoints
List methods:
grpcurl -plaintext localhost:50051 list

List first 5 movies:
grpcurl -plaintext -d '{"page":1,"limit":5}' localhost:50051 movie.MovieService/ListMovies

Search movies by genre:
grpcurl -plaintext -d '{"query":"Action"}' localhost:50051 movie.MovieService/SearchMovies

Get a single movie:
grpcurl -plaintext -d '{"id":"<movie-id>"}' localhost:50051 movie.MovieService/GetMovie

### Step 6: Run tests
go test ./... -v

---

## TURAR – Notification Service

### Step 1: Create database and user
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE DATABASE notificationdb;"
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE USER notification WITH PASSWORD 'notification'; GRANT ALL PRIVILEGES ON DATABASE notificationdb TO notification;"

### Step 2: Apply migrations
cd notification-service
export DATABASE_URL="postgres://notification:notification@localhost:5432/notificationdb?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up

### Step 3: Start Notification Service
DATABASE_URL="postgres://notification:notification@localhost:5432/notificationdb?sslmode=disable" NATS_URL="" GRPC_PORT=50054 go run ./cmd/main.go

### Step 4: Test gRPC endpoints
List methods:
grpcurl -plaintext localhost:50054 list

Send booking confirmation:
grpcurl -plaintext -d '{"booking_id":"123","recipient_email":"test@example.com"}' localhost:50054 notification.NotificationService/SendBookingConfirmation

Get email logs:
grpcurl -plaintext localhost:50054 notification.NotificationService/GetEmailLogs

Health check:
grpcurl -plaintext localhost:50054 notification.NotificationService/HealthCheck

### Step 5: Run unit tests
go test ./internal/service/ -v

---

## ALMAS – Booking Service

### Step 1: Create database and user
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE DATABASE bookingdb;"
docker exec -i microservices-sre-postgres-1 psql -U admin -d postgres -c "CREATE USER booking WITH PASSWORD 'booking'; GRANT ALL PRIVILEGES ON DATABASE bookingdb TO booking;"

### Step 2: Apply migrations
cd booking-service
export DATABASE_URL="postgres://booking:booking@localhost:5432/bookingdb?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up

### Step 3: Create seats table and insert test seat
docker exec -i microservices-sre-postgres-1 psql -U booking -d bookingdb -c "CREATE TABLE IF NOT EXISTS seats (id UUID PRIMARY KEY, is_available BOOLEAN DEFAULT true); INSERT INTO seats (id, is_available) VALUES ('b1b2c3d4-e5f6-7890-abcd-ef1234567891', true) ON CONFLICT DO NOTHING;"

### Step 4: Start Booking Service
DATABASE_URL="postgres://booking:booking@localhost:5432/bookingdb?sslmode=disable" GRPC_PORT=50053 go run ./cmd/main.go

### Step 5: Test gRPC endpoints
List methods:
grpcurl -plaintext localhost:50053 list

Create booking:
grpcurl -plaintext -d '{"user_id":"u1","showtime_id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","seat_ids":["b1b2c3d4-e5f6-7890-abcd-ef1234567891"]}' localhost:50053 booking.BookingService/CreateBooking

Get booking:
grpcurl -plaintext -d '{"id":"<booking-id-from-create>"}' localhost:50053 booking.BookingService/GetBooking

Cancel booking:
grpcurl -plaintext -d '{"id":"<booking-id-from-create>"}' localhost:50053 booking.BookingService/CancelBooking

### Step 6: Run unit tests
go test ./internal/service/ -v

---

## ALMAS – Monitoring (Grafana and Prometheus)

### Step 1: Start Prometheus
cd cinema-system/prometheus-2.55.0.windows-amd64
rm -rf data promdata
mkdir promdata
./prometheus --config.file=../prometheus-local.yml --storage.tsdb.path=promdata

### Step 2: Verify Prometheus targets
Open http://localhost:9090/targets in browser. All four targets should be green.

### Step 3: Open Grafana
Open http://localhost:4040 in browser. Login: admin / admin.

### Step 4: Add Prometheus Data Source in Grafana
Configuration > Data Sources > Add data source > Prometheus.
URL: http://host.docker.internal:9090
Click Save and test.

### Step 5: Import gRPC dashboard
+ > Import > Paste JSON dashboard code.
Select Prometheus data source.
Click Import.

### Step 6: Generate test data for metrics
grpcurl -plaintext localhost:50051 movie.MovieService/ListMovies
grpcurl -plaintext localhost:50052 showtime.ShowtimeService/ListShowtimesByMovie
grpcurl -plaintext -d '{"user_id":"u1","showtime_id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","seat_ids":["b1b2c3d4-e5f6-7890-abcd-ef1234567891"]}' localhost:50053 booking.BookingService/CreateBooking

### Step 7: Refresh Grafana dashboard
Click refresh button. Metrics should appear on all panels.

---

## API Gateway (All Members)

### Step 1: Start API Gateway
cd api-gateway
GRPC_PORT=50050 go run ./cmd/main.go

### Step 2: Test composite methods
List methods:
grpcurl -plaintext localhost:50050 list

Get showtime with movie (Lenara):
grpcurl -plaintext -d '{"showtime_id":"s1"}' localhost:50050 gateway.CinemaGateway/GetShowtimeWithMovie

Send booking notification (Turar):
grpcurl -plaintext -d '{"booking_id":"b1","recipient_email":"a@b.com"}' localhost:50050 gateway.CinemaGateway/SendBookingNotification

Get booking details (Almas):
grpcurl -plaintext -d '{"booking_id":"b1"}' localhost:50050 gateway.CinemaGateway/GetBookingDetails

---

## Frontend (All Members)

Open frontend/index.html in any web browser.
Three tabs are available: Movies, Notifications, Booking.
Each tab demonstrates the corresponding functionality.

---

## Criteria Checklist

1. Clean architecture (20%): All services separate handler, service, repository, and cache layers.
2. At least 12 gRPC endpoints (20%): 18 total (Showtime 6 + Movie 3 + Notification 3 + Booking 3 + Gateway 3).
3. Message Queue NATS (20%): Booking publishes booking.created, Notification subscribes.
4. Databases and Caches with migrations and transactions (20%): PostgreSQL everywhere, Redis in Movie and Showtime, transactions in Booking.
5. Sending Emails via SMTP (10%): Notification Service SMTP sender implemented.
6. Testing Unit and Integration (10%): Unit tests in all services, integration tests with testcontainers.
7. Bonus Frontend (10%): Dark luxury UI with three tabs.
8. Bonus 2 Grafana with tracing metrics and logs (10%): Prometheus metrics, Grafana dashboard, Jaeger and Loki configured.

---

## Final Git Commit Commands

After all changes are verified, commit the documentation:
cd cinema-system
git add FINAL_DOCUMENTATION.md documentation.pdf
git commit -m "docs: add final documentation with step-by-step commands for all members"
git push origin main
