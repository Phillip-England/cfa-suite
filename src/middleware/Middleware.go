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
		err = user.FindById(database, session.UserID)
		if err != nil {
			c.SetCookie(os.Getenv("SESSION_TOKEN_KEY"), "", -1, "/", os.Getenv("SERVER_URL"), true, true)
			c.Redirect(303, "/401")
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func (mw *Middleware) GuestRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionToken, _ := c.Cookie(os.Getenv("SESSION_TOKEN_KEY"))
		if sessionToken != "" {
			c.Redirect(303, "/app/home")
			return
		}
		c.Next()
	}
}