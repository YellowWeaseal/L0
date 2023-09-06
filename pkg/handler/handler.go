package handler

import (
	"TESTShop/pkg/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.OrderService
}

func NewHandler(services *service.OrderService) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Замените на ваш фронтенд URL
	router.Use(cors.New(config))
	router.GET("/orders/:orderUID", h.getOrderByUID)

	return router
}
