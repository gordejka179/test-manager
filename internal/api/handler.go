package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/service"
	"github.com/gordejka179/test-manager/internal/storage"
	"github.com/gordejka179/test-manager/internal/transport/rest/handler"
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

	Service := service.NewTestService(TestRepository)
	ServiceHandler := handler.NewTestServiceHandler(Service)

	//Runner := storage.NewRunner(TestRepository.DB)

	router.GET("/tests", ServiceHandler.GetAllTests)

	router.POST("/tests/newTest", ServiceHandler.AddTest)

	return router
}
