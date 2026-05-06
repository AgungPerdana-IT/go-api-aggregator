# build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o app cmd/server/main.go

# run stage (super kecil)
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]
