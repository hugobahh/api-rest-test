package shutdown

import (
	"api-rest-test/internal/config"
	"api-rest-test/pkg/http"
	"api-rest-test/pkg/logger"
	"context"
	"os/signal"
	"syscall"
	"time"
)

type ShutdownHandler interface {
	Shutdown(context.Context) error
}

type ShutdownManager struct {
	handlers []ShutdownHandler
	logger   *logger.Log
	timeout  time.Duration
}

func NewShutdownManager(server *http.Server,
	configuration *config.Configuration, logger *logger.Log) *ShutdownManager {
	return &ShutdownManager{
		handlers: []ShutdownHandler{
			server,
		},
		logger:  logger,
		timeout: configuration.ShutdownTimeout,
	}
}

func (s *ShutdownManager) EnableSignalHandling() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	<-ctx.Done()
	s.logger.Infof("Received signal. Shutting down...")
	s.gracefulShutdown(ctx)
}

func (s *ShutdownManager) gracefulShutdown(ctx context.Context) {
	context, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	for _, handler := range s.handlers {
		if err := handler.Shutdown(context); err != nil {
			s.logger.Error("ShutdownManager", "gracefulShutdown", err)
		}
	}
}
