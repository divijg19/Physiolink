# Build stage
FROM golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -trimpath/-buildvcs=false keep the binary reproducible and free of local paths.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -buildvcs=false -ldflags="-s -w" -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -buildvcs=false -ldflags="-s -w" -o worker ./cmd/worker

# Run stage
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./api"]
