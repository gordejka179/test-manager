package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/service"
	"github.com/gordejka179/test-manager/internal/storage"
	"github.com/gordejka179/test-manager/internal/transport"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// tests
	TestRepository, err := storage.NewSQLiteStorage("tmp.db")
	if err != nil {
		log.Fatalf("Failed to create SQLite storage: %v", err)
	}

	Runner := storage.NewRunner(TestRepository.DB)

	Service := service.NewService(TestRepository, Runner)

	ServiceHandler := transport.NewServiceHandler(Service)

	router.GET("/tests", h.ListTests)

	router.POST("/tests", h.CreateTest)
	router.GET("/tests/:id", h.GetTest)
	router.PUT("/tests/:id", h.UpdateTest)
	router.DELETE("/tests/:id", h.DeleteTest)

	// configs
	router.POST("/tests/:id/configs", h.AddConfig)
	router.GET("/tests/:id/configs", h.ListConfigs)
	router.GET("/tests/:id/configs/:configId", h.GetConfig)
	router.DELETE("/tests/:id/configs/:configId", h.DeleteConfig)

	// test running
	router.POST("/tests/:id/run", h.RunTest)

	return router
}
