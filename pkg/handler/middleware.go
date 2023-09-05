package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	userCtx = "userId"
)

func getUserId(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	idStr, ok := id.(string)
	if !ok {
		return "", errors.New("user id is of invalid type")
	}
	return idStr, nil
}
