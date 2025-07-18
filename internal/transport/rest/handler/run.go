package handler

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RunService interface {
	RunTest(ctx context.Context, configId int) error
}

type RunServiceHandler struct {
	service RunService
}

func NewRunServiceHandler(S RunService) *RunServiceHandler {
	return &RunServiceHandler{service: S}
}

func (h *RunServiceHandler) RunTest(c *gin.Context) {
	configIdStr := c.PostForm("configId")
	configId, err := strconv.Atoi(configIdStr)
	if err != nil {
		log.Fatal("Ошибка преобразования типа в RunTest:", err)
	}
	h.service.RunTest(c, configId)
}
