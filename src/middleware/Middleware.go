package middleware

import (
	"cfa-suite/src/core"
	"cfa-suite/src/model"
	"os"

	"github.com/gin-gonic/gin"
)

type Middleware struct {}

func NewMiddlware() *Middleware {
	return &Middleware{}
}

func (mw *Middleware) Auth(database *core.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie(os.Getenv("SESSION_TOKEN_KEY"))
		if err != nil {
			c.Redirect(303, "/")
			return
		}
		session := model.NewSession()
		err = session.FindByToken(database, sessionToken)
		if err != nil {
			c.SetCookie(os.Getenv("SESSION_TOKEN_KEY"), "", -1, "/", os.Getenv("SERVER_URL"), true, true)
			c.Redirect(303, "/")
			return
		}
		user := model.NewUser()
		user.SetID(session.UserID)
		trueUser := model.NewUser()
		err = trueUser.FindById(database, session.UserID)
		if err != nil {
			c.Redirect(303, "/401")
		}
		// gotcha check if we have the correct user here
		c.Set("user", trueUser)
		c.Next()
	}
}

func (mw *Middleware) GuestRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, err := c.Cookie(os.Getenv("SESSION_TOKEN_KEY"))
		// gotta make sure logged in users are redirected home
		c.Next()
	}
}