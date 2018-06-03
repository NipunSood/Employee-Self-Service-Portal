package loginMiddleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var LoginCookies = map[string]*LoginCookie{}

var Identities = []identity{
	{EmployeeNumber: "1234", Password: "password"},
}

const LoginCookieName = "Identity"

type LoginCookie struct {
	Value      string
	Expiration time.Time
	Origin     string
}

type identity struct {
	EmployeeNumber string
	Password       string
}

func LoginMiddleware(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/login") ||
		strings.HasPrefix(c.Request.URL.Path, "/public") {
		return
	}

	cookieValue, err := c.Cookie(LoginCookieName)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	cookie, ok := LoginCookies[cookieValue]

	if !ok || cookie.Expiration.Unix() < time.Now().Unix() ||
		cookie.Origin != c.Request.RemoteAddr {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	c.Next()
}
