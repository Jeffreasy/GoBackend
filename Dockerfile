FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM scratch
COPY --from=builder /app/main /main
COPY --from=builder /app/migrations /migrations
EXPOSE 8080
CMD ["/main"]
