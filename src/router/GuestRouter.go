package router

import (
	"cfa-suite/src/core"
	"cfa-suite/src/model"
	"fmt"
	"html"
	"os"

	"github.com/gin-gonic/gin"
)

type GuestRouter struct {
	Group *gin.RouterGroup
	Database *core.Database
}

func NewGuestRouter(group *gin.RouterGroup, database *core.Database) *GuestRouter {
	return &GuestRouter{
		Group: group,
		Database: database,
	}
}

func (router *GuestRouter) LoginPageRoute() {
	router.Group.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"LoginFormErr": html.EscapeString(c.Query("LoginFormErr")),
			"Banner": "CFA Suite",
			"Email": html.EscapeString(c.Query("Email")),
			"Password": html.EscapeString(c.Query("Password")),
		})
	})
}

func (router *GuestRouter) SignupPageRoute() {
	router.Group.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{
			"SignupFormErr": html.EscapeString(c.Query("SignupFormErr")),
			"Banner": "CFA Suite",
			"Email": html.EscapeString(c.Query("Email")),
			"Password": html.EscapeString(c.Query("Password")),
			"PasswordConfirmed": html.EscapeString(c.Query("PasswordConfirmed")),
		})
	})
}

func (router *GuestRouter) LoginUserRoute() {
	router.Group.POST("/api/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		user := model.NewUser()
		err := user.FindByEmail(router.Database, email)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/?LoginFormErr=%s&Email=%s&Password=%s", "invalid credentials", email, password))
			return
		}
		err = user.ComparePassword(password)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/?LoginFormErr=%s&Email=%s&Password=%s", err.Error(), email, password))
			return
		}
		err = user.DeleteSessionsByUser(router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		session := model.NewSession()
		err = session.Insert(router.Database, user.ID)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.SetCookie(os.Getenv("SESSION_TOKEN_KEY"), session.Token, 86400, "/", os.Getenv("SERVER_URL"), true, true)
		c.Redirect(303, "/app/home")
	})
}

func (router *GuestRouter) SignupUserRoute() {
	router.Group.POST("/api/user", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		passwordConfirmed := c.PostForm("password-confirmed")
		user := model.NewUser()
		user.SetEmail(email)
		err := user.ValidateEmail()
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/signup?SignupFormErr=%s&Email=%s&Password=%s&PasswordConfirmed=%s", err.Error(), email, password, passwordConfirmed))
			return
		}
		isUnique, err := user.IsUnique(router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		if !isUnique {
			c.Redirect(303, fmt.Sprintf("/signup?SignupFormErr=user already exists&Email=%s&Password=%s&PasswordConfirmed=%s", email, password, passwordConfirmed))
			return
		}
		if password != passwordConfirmed {
			c.Redirect(303, fmt.Sprintf("/signup?SignupFormErr=passwords must match&Email=%s&Password=%s&PasswordConfirmed=%s", email, password, passwordConfirmed))
			return
		}
		user.SetPassword(password)
		err = user.ValidatePassword()
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/signup?SignupFormErr=%s&Email=%s&Password=%s&PasswordConfirmed=%s", err.Error(), email, password, passwordConfirmed))
			return
		}
		err = user.HashPassword()
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/signup?SignupFormErr=%s&Email=%s&Password=%s&PasswordConfirmed=%s", err.Error(), email, password, passwordConfirmed))
			return
		}
		err = user.Insert(router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		emailKey := model.NewEmailKey()
		emailKey.SetUserID(user.ID)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		err = emailKey.Insert(router.Database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		err = emailKey.SendAccountVerificationEmail(user.Email)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.Redirect(303, fmt.Sprintf("/?Email=%s&Password=%s", user.Email, password))
	})
}