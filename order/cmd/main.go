package main

import (
	"log"

	"github.com/FernandoJMartins/microservices/order/config"
	"github.com/FernandoJMartins/microservices/order/internal/adapters/db"
	"github.com/FernandoJMartins/microservices/order/internal/adapters/grpc"

	//" github .com/ FernandoJMartins / microservices / order / internal / adapters / rest "

	"github.com/FernandoJMartins/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf(" Failed to connect to database . Error : %v", err)
	}
	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
