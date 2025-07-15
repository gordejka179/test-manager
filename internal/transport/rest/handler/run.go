package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/core"
)

type RunService interface {
	Run(ctx context.Context, configName string) (*core.Log, error)
}

type RunServiceHandler struct {
	service RunService
}

func NewRunServiceHandler(S RunService) *RunServiceHandler {
	return &RunServiceHandler{service: S}
}

func (h *RunServiceHandler) RunTest(c *gin.Context) error {
	configId := c.PostForm("configId")

	pkg.connectSSH()
}
