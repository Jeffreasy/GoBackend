# Dockerfile

# Deze Dockerfile bouwt een productiebare build van de Go applicatie. 
# Eerst wordt een builder stage gebruikt om de binary te bouwen, vervolgens wordt een minimalistische scratch image gebruikt voor runtime.

FROM golang:1.20 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM scratch
COPY --from=builder /app/main /main
EXPOSE 8080
CMD ["/main"]
