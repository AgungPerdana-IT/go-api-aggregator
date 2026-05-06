# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o app cmd/server/main.go

# Run stage
FROM scratch
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 3000
CMD ["./app"]