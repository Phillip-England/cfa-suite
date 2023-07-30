package router

import (
	"cfa-suite/src/core"
	"cfa-suite/src/model"
	"fmt"
	"html"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
	Group *gin.RouterGroup
	Database *core.Database
}

func NewUserRouter(group *gin.RouterGroup, database *core.Database) *UserRouter {
	return &UserRouter{
		Group: group,
		Database: database,
	}
}

func (router *UserRouter) HomeRoute() {
	router.Group.GET("/app/home", func(c *gin.Context) {
		userData, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		user := userData.(*model.User)
		location := model.NewLocation()
		locations, err := location.GetLocationsByUserID(user.ID, router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		maxLocationsToShow := 3
		if len(locations) > maxLocationsToShow {
			locations = locations[:maxLocationsToShow]
		}
		hasNoLocations := true
		if len(locations) != 0 {
			hasNoLocations = false
		}
		c.HTML(200, "home.html", gin.H{
			"Banner": "CFA Suite",
			"IsHomePage": "true",
			"Locations": locations,
			"HasNoLocations": hasNoLocations,
			"HasLocations": !hasNoLocations,
		})
	})
}

func (router *UserRouter) CreateLocationPageRoute() {
	router.Group.GET("/app/create-location", func(c *gin.Context) {
		c.HTML(200, "create-location.html", gin.H{
			"Banner": "CFA Suite",
			"IsCreateLocationPage": true,
			"CreateLocationFormErr": html.EscapeString(c.Query("CreateLocationFormErr")),
		})
	})
}

func (router *UserRouter) UserSettingsPageRoute() {
	router.Group.GET("/app/user-settings", func(c *gin.Context) {
		c.HTML(200, "user-settings.html", gin.H{
			"Banner": "CFA Suite",
			"IsUserSettingsPage": true,
		})
	})
}

func (router *UserRouter) LogoutRoute() {
	router.Group.GET("/api/logout", func(c *gin.Context) {
		c.SetCookie(os.Getenv("SESSION_TOKEN_KEY"), "", -1, "/", os.Getenv("SERVER_URL"), true, true)
		c.Redirect(303, "/")
	})
}

func (router *UserRouter) CreateLocationRoute() {
	router.Group.POST("/api/location", func(c *gin.Context) {
		user, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		userModel := user.(*model.User)
		name := c.PostForm("name")
		number := c.PostForm("number")
		location := model.NewLocation()
		err := location.SetName(name)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/app/create-location?CreateLocationFormErr=%s", err.Error()))
			return
		}
		err = location.SetNumber(number)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/app/create-location?CreateLocationFormErr=%s", err.Error()))
			return
		}
		hasThreeOrMoreLocations, err := location.LimitNumberOfLocations(userModel.ID, router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		if hasThreeOrMoreLocations {
			c.Redirect(303, fmt.Sprintf("/app/create-location?CreateLocationFormErr=%s", "only 3 locations per user"))
			return
		}
		location.SetUserID(userModel.ID)
		err = location.Insert(router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.Redirect(303, "/app/home")
	})
}

func (router *UserRouter) DeleteUserPageRoute() {
	router.Group.GET("/app/user-settings/delete", func(c *gin.Context) {
		user, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		userModel := user.(*model.User)
		c.HTML(303, "delete-user.html", gin.H{
			"Banner": "CFA Suite",
			"Email": userModel.Email,
			"DeleteUserFormErr": html.EscapeString(c.Query("DeleteUserFormErr")),
		})
	})
}

func (router *UserRouter) DeleteUserRoute() {
	router.Group.POST("/api/user/delete", func(c *gin.Context) {
		userData, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		user := userData.(*model.User)
		email := c.PostForm("email")
		if strings.ToLower(email) != user.Email {
			c.Redirect(303, "/app/user-settings/delete?DeleteUserFormErr=invalid email provided")
			return
		}
		err := user.Delete(router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.Redirect(303, "/api/logout")
	})
}


