package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"pushpost/internal/config"
	"pushpost/internal/setup"
	lg "pushpost/pkg/logger"
	"time"
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

	fiberConfig := fiber.Config{ // FIXME no hardcoded config here (move to config)
		AppName:                 ServiceName,
		ReadTimeout:             30 * time.Second,
		WriteTimeout:            30 * time.Second,
		IdleTimeout:             120 * time.Second,
		EnableTrustedProxyCheck: true,
		ProxyHeader:             fiber.HeaderXForwardedFor,
	}

	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Trace-ID",
	}

	fiberLogger := logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	})

	app := setup.NewFiber(fiberConfig, corsConfig)

	app.Use(fiberLogger)
}
