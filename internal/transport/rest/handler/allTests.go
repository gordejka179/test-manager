package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/core"
)

type TestService interface {
	AddTest(ctx context.Context, test *core.Test) error
	GetTestByID(ctx context.Context, testID string) (*core.Test, error)
	GetAllTests(ctx context.Context) ([]core.Test, error)
	DeleteTest(ctx context.Context, id string) error
	AddConfig(ctx context.Context, testID string, config *core.Config) error
	GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error)
	GetAllConfigs(ctx context.Context) ([]core.Config, error)
	GetAllConfigsToTest(ctx context.Context, testID string) ([]core.Config, error)
	DeleteConfig(ctx context.Context, testID string) error
	GetLogs(ctx context.Context, testID string, configID string) ([]core.Log, error)
}

type TestServiceHandler struct {
	service TestService
}

func NewTestServiceHandler(S TestService) *TestServiceHandler {
	return &TestServiceHandler{service: S}
}

func (h *TestServiceHandler) GetAllTests(c *gin.Context) {
	tests, err := h.service.GetAllTests(c.Request.Context())
	c.JSON(http.StatusOK, tests)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, tests)
}

func (h *TestServiceHandler) AddTest(c *gin.Context) {

}
