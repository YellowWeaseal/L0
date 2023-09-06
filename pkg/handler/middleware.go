package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	orderCtx = "orderUID"
)

func getOrderUID(c *gin.Context) (string, error) {
	id := c.Param(orderCtx)
	if id == "" {
		return "", errors.New("orderUID not found")
	}

	return id, nil
}
