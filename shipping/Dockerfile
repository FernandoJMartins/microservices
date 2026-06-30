FROM golang:latest AS builder
WORKDIR /app
COPY microservices-proto/ ./microservices-proto/
COPY microservices/shipping/ ./microservices/shipping/
WORKDIR /app/microservices/shipping
RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shipping ./cmd/main.go

FROM scratch
COPY --from=builder /app/microservices/shipping/shipping ./shipping
CMD ["./shipping"]
