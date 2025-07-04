package main

import (
	test_manager "app"
	"log"
	"os"

	"github.com/gordejka179/test-manager/internal/api"
)

type Config struct {
	appPort       string
	dbUsername    string
	dbPassword    string
	dbName        string
	SessionSecret string
	dbHost        string
}

func main() {
	conf := InitConfig()

	handlers := new(api.Handler)

	srv := new(test_manager.Server)

	if err := srv.Run(conf.appPort, handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running http server: %s", err.Error())
	}

}

func InitConfig() Config {

	result := Config{
		appPort:       os.Getenv("APP_PORT"),
		dbUsername:    os.Getenv("DB_USER"),
		dbPassword:    os.Getenv("DB_PASSWORD"),
		dbName:        os.Getenv("DB_NAME"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		dbHost:        os.Getenv("DB_HOST"),
	}

	// default value
	if result.dbHost == "" {
		return Config{
			appPort:       "8080",
			dbUsername:    "postgres",
			dbPassword:    "qwerty",
			dbName:        "postgres",
			SessionSecret: "session",
			dbHost:        "localhost",
		}
	}

	return result
}
