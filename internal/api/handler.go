package api

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// tests
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
