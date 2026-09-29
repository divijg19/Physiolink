# Build stage
# Floats to the latest 1.26.x patch, matching the `go-version` input of
# .github/actions/go-setup. The toolchain is not pinned, so the built binary is
# only byte-reproducible for a given base image; -trimpath and
# -buildvcs=false below keep it free of local paths and VCS stamps regardless.
FROM docker.io/library/golang:1.26 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
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
