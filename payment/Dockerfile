FROM golang:latest AS builder
WORKDIR /app
COPY . .
WORKDIR /app/microservices/payment
RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go

FROM scratch
COPY --from=builder /app/microservices/payment/payment ./payment
CMD ["./payment"]
