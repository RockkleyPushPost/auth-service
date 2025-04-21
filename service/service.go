package service

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
	"pushpost/internal/config"
	"pushpost/pkg/di"
)

const serviceName string = "AuthService"

type Option func(*service)

type service struct {
	Di     *di.DI
	server *fiber.App
	logger *log.Logger
	config *config.Config
}

func NewService(opts ...Option) (Service, error) {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}

	return s, s.validate()
}

//func WithContainer(container *di.Container) Option {
//	return func(s *service) {
//		s.container = container
//	}
//}

func WithLogger(logger *log.Logger) Option {
	return func(s *service) {
		s.logger = logger
	}
}

func WithServer(server *fiber.App) Option {
	return func(s *service) {
		s.server = server
	}
}

func WithConfig(config *config.Config) Option {
	return func(s *service) {
		s.config = config
	}
}

func WithDI(container *di.DI) Option {
	return func(s *service) {
		s.Di = container

	}
}

func (s *service) validate() error {

	//if s.container == nil {
	//	return errors.New("missing container")
	//}
	if s.logger == nil {

		return errors.New("missing logger")
	}
	if s.config == nil {

		return errors.New("missing config")
	}

	if s.server == nil {

		return errors.New("missing server")
	}

	return nil
}

func (s *service) Run(ctx context.Context) error {
	s.logger.Printf("starting %s\n", s.Name())

	err := s.server.Listen(s.config.Server.Host + ":" + s.config.Server.Port)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *service) Shutdown(ctx context.Context) error {
	s.logger.Printf("shutting down %s...\n", serviceName)
	err := s.server.Shutdown()

	if err != nil {

		return err
	}

	ctx.Done()

	return nil
}

func (s *service) Name() string {
	return serviceName
}
