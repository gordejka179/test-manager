package handler

import (
	"context"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RunService interface {
	RunTest(ctx context.Context, configId int, serverIp string, username string, commandTemplate string) error
}

type RunServiceHandler struct {
	service RunService
}

func NewRunServiceHandler(S RunService) *RunServiceHandler {
	return &RunServiceHandler{service: S}
}

func (h *RunServiceHandler) RunTest(c *gin.Context) {
	configIdStr := c.PostForm("config_id")
	configId, err := strconv.Atoi(configIdStr)
	if err != nil {
		log.Fatal("Ошибка преобразования типа в RunTest:", err)
	}

	serverIp := c.PostForm("server_ip")
	username := c.PostForm("username")
	commandTemplate := c.PostForm("commandTemplate")
	h.service.RunTest(c, configId, serverIp, username, commandTemplate)
}
