package main

import (
	"cfa-suite/src/core"
	"cfa-suite/src/middleware"
	"cfa-suite/src/model"
	"fmt"
	"html"
	"log"
	"os"
	"strconv"

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
	r.LoadHTMLGlob("./templates/**/**/*")
	r.Static("/static", "./static")

	//==========================================================================
	// MIDDLEWARE
	//==========================================================================

	mw := middleware.NewMiddlware()
	protectedRoutes := r.Group("/", mw.Auth(database))
	guestRoutes := r.Group("/", mw.GuestRedirect())

	//==========================================================================
	// GUEST PAGES
	//==========================================================================

	guestRoutes.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"LoginFormErr": html.EscapeString(c.Query("LoginFormErr")),
			"Banner": "CFA Suite",
			"Email": html.EscapeString(c.Query("Email")),
			"Password": html.EscapeString(c.Query("Password")),
		})
	})
	
	guestRoutes.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", gin.H{
			"SignupFormErr": html.EscapeString(c.Query("SignupFormErr")),
			"Banner": "CFA Suite",
			"Email": html.EscapeString(c.Query("Email")),
			"Password": html.EscapeString(c.Query("Password")),
			"PasswordConfirmed": html.EscapeString(c.Query("PasswordConfirmed")),
		})
	})
	
	
	//==========================================================================
	// ERROR PAGES
	//==========================================================================
	
	r.GET("/500", func(c *gin.Context) {
		c.HTML(200, "500.html", gin.H{
			"ServerErr": html.EscapeString(c.Query("ServerErr")),
			"Banner": "CFA Suite",
		})
	})
	
	r.GET("/401", func(c *gin.Context) {
		c.HTML(200, "401.html", gin.H{
			"Banner": "CFA Suite",
		})
	})
	
	//==========================================================================
	// USER PAGES
	//==========================================================================
	
	protectedRoutes.GET("/app/home", func(c *gin.Context) {
		userData, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		user := userData.(*model.User)
		location := model.NewLocation()
		locations, err := location.GetLocationsByUserID(user.ID, database)
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

	protectedRoutes.GET("/app/create-location", func(c *gin.Context) {
		c.HTML(200, "create-location.html", gin.H{
			"Banner": "CFA Suite",
			"IsCreateLocationPage": true,
			"CreateLocationFormErr": html.EscapeString(c.Query("CreateLocationFormErr")),
		})
	})

	// could be improved by only getting the locations from within the index in the db call
	// however, then intention of the app is that users do not have an huge running list of locations
	// but wanted to leave the possibility open and implemented
	protectedRoutes.GET("/app/view-locations/:index", func(c *gin.Context) {

		// pulling user from middleware
		userData, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		user := userData.(*model.User)

		// getting the index from the params
		index, ok := c.Params.Get("index")
		if !ok {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		
		// getting our users locations
		location := model.NewLocation()
		locations, err := location.GetLocationsByUserID(user.ID, database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}

		// setting up variables to determine how location-gallery.html functions
		locationsPerPage := 5
		hasLocations := len(locations) > 0
		startingIndex, err := strconv.Atoi(index)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		endingIndex := startingIndex + locationsPerPage
		nextIndex := startingIndex + locationsPerPage
		previousIndex := startingIndex - locationsPerPage
		hasLessThanOnePage := len(locations) <= locationsPerPage
		
		// setting up rendering conditions for the back and forward buttons
		renderBackButton := true
		renderForwardButton := true
		if startingIndex == 0 {
			renderBackButton = false
		}
		if (startingIndex + locationsPerPage) >= len(locations) {
			renderForwardButton = false
		}

		// handling manual :index input from users
		if startingIndex >= len(locations) || startingIndex < 0 {
			c.Redirect(303, "/app/view-locations/0")
			return
		}
		
		// only showing a limited number of locations per page
		if endingIndex >= len(locations) {
			endingIndex = len(locations)
		}
		visibleLocations := locations[startingIndex:endingIndex]

		c.HTML(200, "view-locations.html", gin.H{
			"Banner": "View Locations",
			"IsViewAllLocationsPage": true,
			"Locations": visibleLocations,
			"SearchQuery": "",
			"RenderLocationGallery": true,
			"RenderSearchResults": false,
			"CurrentStartingIndex": startingIndex,
			"CurrentEndingIndex": endingIndex,
			"RenderBackButton": renderBackButton,
			"RenderForwardButton": renderForwardButton,
			"HasLocations": hasLocations,
			"HasNoLocations": !hasLocations,
			"NextIndex": nextIndex,
			"PreviousIndex": previousIndex,
			"HasLessThanOnePage": hasLessThanOnePage,
		})
	})
	
	
	//==========================================================================
	// API
	//==========================================================================

	r.GET("/api/logout", func(c *gin.Context) {
		c.SetCookie(os.Getenv("SESSION_TOKEN_KEY"), "", -1, "/", os.Getenv("SERVER_URL"), true, true)
		c.Redirect(303, "/")
	})
	
	r.POST("/api/login", func(c *gin.Context) {
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
		c.Redirect(303, "/app/home")
	})

	r.POST("/api/user", func(c *gin.Context) {
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

	protectedRoutes.POST("/api/location", func(c *gin.Context) {
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
		// hasThreeOrMoreLocations, err := location.LimitNumberOfLocations(userModel.ID, database)
		// if err != nil {
		// 	c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
		// 	return
		// }
		// if hasThreeOrMoreLocations {
		// 	c.Redirect(303, fmt.Sprintf("/app/create-location?CreateLocationFormErr=%s", "only 3 locations per user"))
		// 	return
		// }
		location.SetUserID(userModel.ID)
		err = location.Insert(database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		c.Redirect(303, "/app/home")
	})

	protectedRoutes.GET("/api/location/search", func(c *gin.Context) {
		userData, ok := c.Get("user")
		if !ok {
			c.Redirect(303, "/401")
			return
		}
		user := userData.(*model.User)
		query := c.PostForm("query")
		location := model.NewLocation()
		locations, err := location.GetLocationsBySearchAndUserID(user.ID, query, database)
		if err != nil {
			c.Redirect(303, fmt.Sprintf("/500?ServerErr=%s", err.Error()))
			return
		}
		fmt.Println(len(locations))
		c.HTML(200, "view-locations.html", gin.H{
			"Banner": "View Locations",
			"IsViewAllLocationsPage": true,
			"Locations": locations,
			"SearchQuery": query,
			"RenderLocationGallery": false,
			"RenderSearchResults": true,
			"CurrentStartingIndex": 0,
			"CurrentEndingIndex": 0,
			"RenderBackButton": false,
			"RenderForwardButton": false,
			"HasLocations": true,
			"HasNoLocations": false,
			"NextIndex": 0,
			"PreviousIndex": 0,
			"HasLessThanOnePage": true,
		})
	})

	//==========================================================================
	// RUNNING
	//==========================================================================

	r.Run()

}