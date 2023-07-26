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
			c.Redirect(303, "/401")
			return
		}
		session := model.NewSession()
		err = session.FindByToken(database, sessionToken)
		if err != nil {
			c.Redirect(303, "/401")
			return
		}
		user := model.NewUser()
		user.SetID(session.UserID)
		trueUser := model.NewUser()
		err = trueUser.FindById(database, session.UserID)
		if err != nil {
			c.Redirect(303, "/401")
		}
		c.Set("user", trueUser)
		c.Next()

	}
}