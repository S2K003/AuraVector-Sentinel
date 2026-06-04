# Dockerfile
# Stage 1: Build the optimized Go binary
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
# We disable CGO and build a statically linked binary for pure performance
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auravector cmd/auravector/main.go

# Stage 2: Create a minimal runtime container
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auravector .
EXPOSE 8080
CMD ["./auravector"]