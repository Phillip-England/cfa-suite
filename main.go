package main

import (
	"cfa-suite/src/core"
	"cfa-suite/src/middleware"
	"cfa-suite/src/model"
	"fmt"
	"html"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//==========================================================================
	// CONFIG
	//==========================================================================

	_ = godotenv.Load()

	
	//==========================================================================
	// DATABASE
	//==========================================================================

	database := core.NewDatabase()
	err := database.InitTables()
	if err != nil {
		log.Fatal(err.Error())
	}

	//==========================================================================
	// ROUTER
	//==========================================================================

	r := gin.Default()
	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./static")

	//==========================================================================
	// MIDDLEWARE
	//==========================================================================

	mw := middleware.NewMiddlware()
	authGroup := r.Group("/", mw.Auth(database))
	fmt.Println(authGroup)

	//==========================================================================
	// PAGES
	//==========================================================================

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"LoginFormErr": html.EscapeString(c.Query("LoginFormErr")),
			"Email": html.EscapeString(c.Query("Email")),
			"Password": html.EscapeString(c.Query("Password")),
		})
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{
			"SignupFormErr": html.EscapeString(c.Query("SignupFormErr")),
			"Email": html.EscapeString(c.Query("Email")),
			"Password": html.EscapeString(c.Query("Password")),
			"PasswordConfirmed": html.EscapeString(c.Query("PasswordConfirmed")),
		})
	})

	r.GET("/500", func(c *gin.Context) {
		c.HTML(200, "500.html", gin.H{
			"ServerErr": html.EscapeString(c.Query("ServerErr")),
		})
	})

	r.GET("/401", func(c *gin.Context) {
		c.HTML(200, "401.html", nil)
	})

	authGroup.GET("/home", func(c *gin.Context) {
		user, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		fmt.Println(user)
		c.HTML(200, "home.html", nil)
	})

	//==========================================================================
	// ACTIONS
	//==========================================================================

	r.POST("/action/login", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		user := model.NewUser()
		err := user.FindByEmail(database, email)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/?LoginFormErr=%s&Email=%s&Password=%s", "invalid credentials", email, password))
			return
		}
		err = user.ComparePassword(password)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/?LoginFormErr=%s&Email=%s&Password=%s", err.Error(), email, password))
			return
		}
		err = user.DeleteSessionsByUser(database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		session := model.NewSession()
		err = session.Insert(database, user.ID)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.SetCookie(os.Getenv("SESSION_TOKEN_KEY"), session.Token, 86400, "/", os.Getenv("SERVER_URL"), true, true)
		c.Redirect(303, "/home")
	})

	r.POST("/action/user", func(c *gin.Context) {
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
		isUnique, err := user.IsUnique(database)
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
		err = user.Insert(database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.Redirect(303, "/")
	})

	//==========================================================================
	// RUNNING
	//==========================================================================

	r.Run()

}