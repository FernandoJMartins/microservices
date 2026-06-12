package main

import (
	"log"

	"github.com/FernandoJMartins/microservices/order/config"
	"github.com/FernandoJMartins/microservices/order/internal/adapters/db"
	payment_adapter "github.com/FernandoJMartins/microservices/order/internal/adapters/payment"

	"github.com/FernandoJMartins/microservices/order/internal/adapters/grpc"
	"github.com/FernandoJMartins/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf(" Failed to connect to database . Error : %v", err)
	}

	paymentAdapter, err := payment_adapter.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf(" Failed to initialize payment stub . Error : %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
