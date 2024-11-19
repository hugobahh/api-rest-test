package di

import (
	controllerRegister "api-rest-test/internal/app/controller"
	repositoryRegister "api-rest-test/internal/app/repository"
	serviceRegister "api-rest-test/internal/app/service"
	"api-rest-test/pkg/logger"
	"sync"

	"api-rest-test/internal/config"
	"api-rest-test/internal/shutdown"
	"api-rest-test/pkg/http"

	"api-rest-test/pkg/database/mysql"

	"go.uber.org/dig"
)

var (
	container *dig.Container
	once      sync.Once
)

func GetContainer() *dig.Container {
	once.Do(func() {
		container = buildContainer()
	})
	return container
}

func buildContainer() *dig.Container {
	container := dig.New()
	container.Provide(logger.NewLog)
	container.Provide(config.NewConfiguration)
	container.Provide(mysql.NewMySQLConnector)
	container.Provide(repositoryRegister.NewUserRepository)
	container.Provide(controllerRegister.NewUserController)
	container.Provide(serviceRegister.NewUserService)
	container.Provide(shutdown.NewShutdownManager)
	container.Provide(http.NewServer)

	return container
}
