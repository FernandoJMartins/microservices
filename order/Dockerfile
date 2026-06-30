FROM golang:latest AS builder
WORKDIR /app
COPY microservices-proto/ ./microservices-proto/
COPY microservices/order/ ./microservices/order/
WORKDIR /app/microservices/order
RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order ./cmd/main.go

FROM scratch
COPY --from=builder /app/microservices/order/order ./order
CMD ["./order"]
