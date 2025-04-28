package service

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"pushpost/internal/config"
	"pushpost/internal/services/user_service/domain"
	"pushpost/internal/services/user_service/domain/usecase"
	transport2 "pushpost/internal/services/user_service/transport"
	"pushpost/internal/services/user_service/transport/handlers"
	"pushpost/internal/services/user_service/transport/routing"
	"pushpost/pkg/di"
)

func Setup(DI *di.DI, server *fiber.App, db *gorm.DB, cfg *config.Config) error {

	jwtSecret := cfg.JwtSecret

	// Auth
	var authUseCase domain.AuthUseCase = &usecase.AuthUseCase{JwtSecret: jwtSecret}
	var authHandler transport2.AuthHandler = &transport.AuthHandler{}

	if err := DI.Register(
		server, db, authUseCase, authHandler); err != nil {
		log.Fatalf("failed to register %v", err)

		return err
	}

	if err := DI.Bind(server, db, authUseCase, authHandler); err != nil {
		log.Fatalf("failed to bind %v", err)

		return err
	}

	authRoutes := routing.AuthRoutes{
		Register:       authHandler.RegisterUser,
		Login:          authHandler.Login,
		VerifyEmailOTP: authHandler.VerifyEmailOTP,
		SendNewOTP:     authHandler.SendNewOTP,
	}

	if err := DI.RegisterRoutes(authRoutes, "/auth"); err != nil {
		log.Fatalf("failed to register routes: %v", err)

		return err
	}

	return nil
}
