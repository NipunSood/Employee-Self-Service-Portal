package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func registerRoutes() *gin.Engine {
	r := gin.Default()

	r.Use(loginMiddleware)

	r.LoadHTMLGlob("templates/**/*.html")
	r.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Any("/login", func(c *gin.Context) {
		employeeNumber := c.PostForm("employeeNumber")
		password := c.PostForm("password")

		for _, identity := range identities {
			if identity.employeeNumber == employeeNumber &&
				identity.password == password {
				lc := loginCookie{
					value:      employeeNumber,
					expiration: time.Now().Add(24 * time.Hour),
					origin:     c.Request.RemoteAddr,
				}

				loginCookies[lc.value] = &lc

				maxAge := lc.expiration.Unix() - time.Now().Unix()
				c.SetCookie(loginCookieName, lc.value, int(maxAge), "", "", false, true)

				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}
		}

		c.HTML(http.StatusOK, "login.html", nil)
	})

	r.GET("/employees/:id/vacation", func(c *gin.Context) {
		id := c.Param("id")
		timesOff, ok := TimesOff[id]

		if !ok {
			c.String(http.StatusNotFound, "404 - Page Not Found")
			return
		}

		c.HTML(http.StatusOK, "vacation-overview.html",
			gin.H{
				"TimesOff": timesOff,
			})
	})

	r.POST("/employees/:id/vacation/new", func(c *gin.Context) {
		var timeOff TimeOff
		err := c.BindJSON(&timeOff)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		id := c.Param("id")

		timesOff, ok := TimesOff[id]

		if !ok {
			TimesOff[id] = []TimeOff{}
		}

		TimesOff[id] = append(timesOff, timeOff)

		c.JSON(http.StatusCreated, &timeOff)
	})

	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))
	admin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-overview.html",
			gin.H{
				"Employees": employees,
			})
	})

	admin.GET("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "add" {
			c.HTML(http.StatusOK, "admin-employee-add.html", nil)
			return
		}
		employee, ok := employees[id]

		if !ok {
			c.String(http.StatusNotFound, "404 - Not Found")
			return
		}

		c.HTML(http.StatusOK, "admin-employee-edit.html",
			gin.H{
				"Employee": employee,
			})
	})

	admin.POST("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")

		if id == "add" {

			startDate, err := time.Parse("2006-01-02", c.PostForm("startDate"))
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}

			var emp Employee
			c.Bind(&emp)
			emp.ID = 42
			emp.Status = "Active"
			emp.StartDate = startDate
			employees["42"] = emp

			c.Redirect(http.StatusMovedPermanently, "/admin/employees/42")
		}
	})

	r.Static("/public", "./public")

	return r
}
