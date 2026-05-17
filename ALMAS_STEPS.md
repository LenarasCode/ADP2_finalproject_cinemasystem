# 📦 Almas – Booking Service & Grafana

## Overview
Almas implemented the Booking microservice with transactional booking, NATS event publishing, gRPC, and set up Grafana with Prometheus, Jaeger, Loki for full observability.

## Features
- **Booking Service**: CreateBooking, GetBooking, CancelBooking with PostgreSQL transactions
- **NATS**: publishes `booking.created` event
- **Grafana**: dashboards for gRPC metrics, tracing, logs
- **Prometheus** metrics endpoint on each service
- **Jaeger** distributed tracing
- **Loki** log aggregation
- **Frontend** booking panel (luxury dark UI)

## Demo script (English)

### 1. Start infrastructure
```bash
docker compose up -d postgres-booking nats prometheus jaeger loki grafana
2. Apply migrations
bash
cd booking-service
migrate -path migrations -database "postgres://booking:booking@localhost:5436/bookingdb?sslmode=disable" up
3. Run Booking Service
bash
DATABASE_URL="postgres://booking:booking@localhost:5436/bookingdb?sslmode=disable" NATS_URL=nats://localhost:4222 GRPC_PORT=50053 go run ./cmd/main.go
4. gRPC calls
bash
grpcurl -plaintext localhost:50053 list
grpcurl -plaintext -d '{"user_id":"u1","showtime_id":"s1","seat_ids":["seat1"]}' localhost:50053 booking.BookingService/CreateBooking
grpcurl -plaintext -d '{"id":"booking-id"}' localhost:50053 booking.BookingService/GetBooking
grpcurl -plaintext -d '{"id":"booking-id"}' localhost:50053 booking.BookingService/CancelBooking
5. Show Grafana
Open http://localhost:3000 (admin/admin). Go to Dashboards → gRPC metrics.

6. Show tracing in Jaeger
Open http://localhost:16686. Search traces for booking-service.

7. Run tests
bash
cd booking-service
go test ./internal/service/ -v
INTEGRATION=1 go test ./internal/repository/ -v
Russian demo script (русский)
То же самое, но по-русски (кратко):

bash
# BD and monitoring
docker compose up -d postgres-booking nats prometheus jaeger loki grafana
# Migrations
cd booking-service
migrate -path migrations -database "postgres://booking:booking@localhost:5436/bookingdb?sslmode=disable" up
# Service
DATABASE_URL="postgres://booking:booking@localhost:5436/bookingdb?sslmode=disable" NATS_URL=nats://localhost:4222 GRPC_PORT=50053 go run ./cmd/main.go &
# gRPC
grpcurl -plaintext localhost:50053 list
grpcurl -plaintext -d '{"user_id":"u1","showtime_id":"s1","seat_ids":["seat1"]}' localhost:50053 booking.BookingService/CreateBooking
# Open Grafana (localhost:3000) and Jaeger (localhost:16686)
