package main

import (
	"context"
	"github.com/RockkleyPushPost/auth-service/service"
	"github.com/RockkleyPushPost/common/config"
	"github.com/RockkleyPushPost/common/di"
	lg "github.com/RockkleyPushPost/common/logger"
	"github.com/RockkleyPushPost/common/setup"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//kafkaBroker := os.Getenv("KAFKA_BROKER")
	//
	//usecase := usecase.AuthUseCase{kafkaBroker}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadYamlConfig(os.Getenv("AUTH_SERVICE_CONFIG_PATH"))

	if err != nil {

		log.Fatalf("failed to load config: %v", err)
	}
	srvLogger := lg.InitLogger(cfg.ServiceName)

	fiberConfig := fiber.Config{ // FIXME no hardcoded config here (move to config)
		AppName:                 cfg.ServiceName,
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

	server := setup.NewFiber(fiberConfig, corsConfig)
	db, err := setup.Database(cfg.Database)

	if err != nil {

		log.Fatal(err)
	}

	DI := di.NewDI(server, cfg.JwtSecret)

	server.Use(fiberLogger)

	err = service.Setup(DI, server, db, cfg)

	if err != nil {

		srvLogger.Fatal(err)
	}

	srv, err := service.NewService(
		service.WithConfig(cfg),
		service.WithDI(DI),
		service.WithLogger(srvLogger),
		service.WithServer(server),
	)

	if err != nil {

		srvLogger.Fatal(err)
	}

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
