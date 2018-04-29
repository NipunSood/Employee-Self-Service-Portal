package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var loginCookies = map[string]*loginCookie{}
var identities = []identity{
	{employeeNumber: "1234", password: "password"},
}

const loginCookieName = "Identity"

func loginMiddleware(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/login") ||
		strings.HasPrefix(c.Request.URL.Path, "/public") {
		return
	}

	cookieValue, err := c.Cookie(loginCookieName)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	cookie, ok := loginCookies[cookieValue]

	if !ok ||
		cookie.expiration.Unix() < time.Now().Unix() ||
		cookie.origin != c.Request.RemoteAddr {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	c.Next()
}

type loginCookie struct {
	value      string
	expiration time.Time
	origin     string
}

type identity struct {
	employeeNumber string
	password       string
}
