package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getOrderByUID(c *gin.Context) {
	orderUID, err := getOrderUID(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())

	}

	order, err := h.services.GetOrderByUID(orderUID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}
