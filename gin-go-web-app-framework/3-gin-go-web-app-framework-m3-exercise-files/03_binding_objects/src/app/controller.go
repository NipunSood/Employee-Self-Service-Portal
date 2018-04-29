package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func registerRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
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
			map[string]interface{}{
				"TimesOff": timesOff,
			})
	})

	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "admin",
	}))
	admin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-overview.html",
			map[string]interface{}{
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
			map[string]interface{}{
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
			err = c.Bind(&emp)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
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
