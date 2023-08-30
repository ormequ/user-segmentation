package http

import (
	"github.com/gin-gonic/gin"
	"user-segmentation/internal/service"
)

func SetRoutes(r gin.IRouter, svc service.Service) {
	r.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/segments", createSegment(svc))
	r.DELETE("/segments", deleteSegment(svc))

	r.GET("/history/:year/:month", getHistory(svc))
	r.GET("/users/:user_id", getUserSegments(svc))
	r.POST("/users/:user_id", changeUserSegments(svc))
}
