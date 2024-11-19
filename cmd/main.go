package main

import (
	"api-rest-test/internal/di"
	"api-rest-test/internal/shutdown"
	"api-rest-test/pkg/http"
	"api-rest-test/pkg/logger"
)

func main() {
	err := di.GetContainer().Invoke(func(server *http.Server, shutdown *shutdown.ShutdownManager) {
		go shutdown.EnableSignalHandling()

		server.Start()
	})
	if err != nil {
		logger.GetLogger().Fatal("main", "main", err)
	}

}
