package router

import (
	"html"

	"github.com/gin-gonic/gin"
)

type ErrorRouter struct {
	Group *gin.RouterGroup
}

func NewErrorRouter(group *gin.RouterGroup) *ErrorRouter {
	return &ErrorRouter{
		Group: group,
	}
}

func (router *ErrorRouter) InternalServerErrorRoute() {
	router.Group.GET("/500", func(c *gin.Context) {
		c.HTML(200, "500.html", gin.H{
			"ServerErr": html.EscapeString(c.Query("ServerErr")),
			"Banner": "CFA Suite",
		})
	})
}

func (router *ErrorRouter) UnauthorizedRoute() {
	router.Group.GET("/401", func(c *gin.Context) {
		c.HTML(200, "401.html", gin.H{
			"Banner": "CFA Suite",
		})
	})
}
