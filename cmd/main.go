package main

import (
	"log"

	tm "github.com/gordejka179/test-manager"
	"github.com/gordejka179/test-manager/internal/api"
)

type Config struct {
	appPort string
}

func main() {
	conf := InitConfig()

	handlers := new(api.Handler)

	srv := new(tm.Server)

	if err := srv.Run(conf.appPort, handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}

func InitConfig() Config {
	return Config{
		appPort: "8080",
	}
}
