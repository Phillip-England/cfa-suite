package main

import (
	"cfa-suite/src/core"
	"cfa-suite/src/middleware"
	"cfa-suite/src/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	//==========================================================================
	// CONFIG
	//==========================================================================

	// ignored the error here on purpose
	_ = godotenv.Load()

	
	//==========================================================================
	// DATABASE
	//==========================================================================

	// the database runs table creation on startup
	// all migrations are manually ran at this point
	// no intentions on using an ORM at this point
	database := core.NewDatabase()
	err := database.InitTables()
	if err != nil {
		log.Fatal(err.Error())
	}

	//==========================================================================
	// ROUTER
	//==========================================================================

	r := gin.Default()
	r.LoadHTMLGlob("./templates/**/**/*")
	r.Static("/static", "./static")

	//==========================================================================
	// MIDDLEWARE
	//==========================================================================
	
	mw := middleware.NewMiddlware()
	authGroup := r.Group("/", mw.Auth(database))
	guestGroup := r.Group("/", mw.GuestRedirect())

	//==========================================================================
	// GUEST ROUTES
	//==========================================================================
	
	guestRouter := router.NewGuestRouter(guestGroup, database)
	guestRouter.LoginPageRoute()
	guestRouter.SignupPageRoute()
	guestRouter.LoginUserRoute()
	guestRouter.SignupUserRoute()

	//==========================================================================
	// ERROR ROUTES
	//==========================================================================
	
	errorRouter := router.NewErrorRouter(r)
	errorRouter.InternalServerErrorRoute()
	errorRouter.UnauthorizedRoute()
	
	//==========================================================================
	// USER ROUTES
	//==================r=======================================================
	
	userRouter := router.NewUserRouter(r, authGroup, database)
	userRouter.HomeRoute()
	userRouter.UserSettingsPageRoute()
	userRouter.CreateLocationPageRoute()
	userRouter.LogoutRoute()
	userRouter.CreateLocationRoute()
	userRouter.DeleteUserPageRoute()
	userRouter.DeleteUserRoute()
	userRouter.VerifyAccountRoute()

	//==========================================================================
	// RUNNING
	//==========================================================================

	r.Run()
	

}