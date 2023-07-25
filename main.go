package main

import (
	"cfa-suite/src/core"
	"fmt"
	"log"

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
	// PAGES
	//==========================================================================

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})

	//==========================================================================
	// ACTIONS
	//==========================================================================

	r.POST("/action/login", func(c *gin.Context) {
		c.Redirect(303, "/")
	})

	r.POST("/action/user", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")
		passwordConfirmed := c.PostForm("password-confirmed")
		fmt.Println(email, password, passwordConfirmed)
		c.Redirect(303, "/")
	})

	//==========================================================================
	// RUNNING
	//==========================================================================

	r.Run()

}