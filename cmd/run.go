package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"os/signal"
	"pushpost/internal/config"
	"pushpost/internal/services/auth_service/service"
	"pushpost/internal/setup"
	"pushpost/pkg/di"
	lg "pushpost/pkg/logger"
	"syscall"
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
	srv, err := service.NewService(service.WithServer(app), service.WithConfig(cfg), service.WithLogger(srvLogger))

	if err != nil {
		log.Fatal(err)
	}
	DI := di.NewDI(app, cfg.JwtSecret)

	err = service.Setup(DI, app, db, cfg)
	go handleShutdown(ctx, cancel, srv, srvLogger)

	srvLogger.Fatal(srv.Run(ctx))
}

func handleShutdown(ctx context.Context, cancel context.CancelFunc, srv service.Service, logger *log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logger.Printf("received signal: %v", sig)
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Printf("shutdown error: %v", err)
		}
	case <-ctx.Done():
		logger.Println("context cancelled")
	}
}
