package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/core"
)

type RunService interface {
	RunTest(ctx context.Context, configId string) (*core.Log, error)
}

type RunServiceHandler struct {
	service RunService
}

func NewRunServiceHandler(S RunService) *RunServiceHandler {
	return &RunServiceHandler{service: S}
}

func (h *RunServiceHandler) RunTest(c *gin.Context) {
	configId := c.PostForm("configId")
	h.service.RunTest(c, configId)
}
