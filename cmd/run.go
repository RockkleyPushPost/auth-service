package main

import (
	"context"
	"log"
	"os"
	"pushpost/internal/config"
	lg "pushpost/pkg/logger"
)

const ServiceName = "auth-service"

func main() {
	//kafkaBroker := os.Getenv("KAFKA_BROKER")
	//
	//usecase := usecase.AuthUseCase{kafkaBroker}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srvLogger := lg.InitLogger(ServiceName)

	cfg, err := config.LoadYamlConfig(os.Getenv("AUTH_SERVICE_CONFIG_PATH"))

	if err != nil {

		log.Fatalf("failed to load gateway service config: %v", err)
	}
}
