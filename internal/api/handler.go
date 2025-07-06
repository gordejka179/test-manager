package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/service"
	"github.com/gordejka179/test-manager/internal/storage"
	"github.com/gordejka179/test-manager/internal/transport/rest/handler"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.LoadHTMLGlob("internal/web/templates/*")
	router.Static("/images", "internal/web/static/images")
	router.Static("/styles", "internal/web/static/styles")
	router.Static("/scripts", "internal/web/static/scripts")

	router.GET("/home", h.home)

	// tests
	TestRepository, err := storage.NewSQLiteStorage("tmp.db")
	if err != nil {
		log.Fatalf("Failed to create SQLite storage: %v", err)
	}

	Service := service.NewTestService(TestRepository)
	ServiceHandler := handler.NewTestServiceHandler(Service)

	//Runner := storage.NewRunner(TestRepository.DB)

	router.GET("/home/tests", ServiceHandler.GetAllTests)

	router.POST("/home/tests/newTest", ServiceHandler.AddTest)
	router.POST("/home/tests/newConfig", ServiceHandler.AddConfig)

	return router
}

func (h *Handler) home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)

}
