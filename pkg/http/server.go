package http

import (
	"api-rest-test/internal/app/controller"
	"api-rest-test/internal/config"
	"api-rest-test/pkg/logger"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	configuration  *config.Configuration
	userController *controller.UserController
	logger         *logger.Log
	app            *fiber.App
}

func NewServer(user *controller.UserController, configuration *config.Configuration, logger *logger.Log) *Server {
	return &Server{
		configuration:  configuration,
		logger:         logger,
		userController: user,
		app:            fiber.New(),
	}
}

func (s *Server) Start() {
	s.app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: false,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	s.setupRoutes()
	port := fmt.Sprintf(":%d", s.configuration.Port)
	s.logger.Infof("Http server started on port: %d", s.configuration.Port)

	if err := s.app.Listen(port); err != nil {
		s.logger.Errorf("Start", "Listen", "Error starting server: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Infof("Shutting down http server...")
	return s.app.Shutdown()
}

func (s *Server) setupRoutes() {
	s.app.Get("/user/health", s.userController.HealthCheck)

	s.app.Post("/user/register/", s.userController.Register)
	s.app.Post("/user/login/", s.userController.Login)
}
