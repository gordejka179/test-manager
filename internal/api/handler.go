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

	testRepository, err := storage.NewSQLiteStorage("tmp.db")
	if err != nil {
		log.Fatalf("Failed to create SQLite storage: %v", err)
	}

	repService := service.NewTestService(testRepository)
	repServiceHandler := handler.NewTestServiceHandler(repService)

	router.GET("/home/tests", repServiceHandler.GetAllTests)
	router.POST("/home/tests/newTest", repServiceHandler.AddTest)
	router.POST("/home/tests/newConfig", repServiceHandler.AddConfig)
	router.POST("/home/tests/configsToTest", repServiceHandler.GetAllConfigsToTest)
	router.POST("/home/tests/configHistory", repServiceHandler.GetLogsToConfig)

	runService := service.NewRunService(testRepository)
	runServiceHandler := handler.NewRunServiceHandler(runService)
	router.POST("/home/tests/runTest", runServiceHandler.RunTest)

	return router
}

func (h *Handler) home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)

}
