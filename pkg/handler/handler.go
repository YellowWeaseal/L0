package handler

import (
	"TESTShop/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	shop := router.Group("shop")
	{
		shop.GET("/shop/order", h.getOrderByUID)
	}

	return router
}
