package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NipunSood/Employee-Self-Service-Portal/src/model"
	"github.com/gin-gonic/gin"
)

func RegisterRouters() *gin.Engine {
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
		timesOff, ok := model.TimesOff[id]

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
		"nipun@trailerstop.me": "nipunrocks",
	}))
	admin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-overview.html",
			map[string]interface{}{
				"Employees": model.Employees,
			})
	})

	admin.GET("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "add" {
			c.HTML(http.StatusOK, "admin-employee-add.html", nil)
			return
		}

		employee, ok := model.Employees[id]

		if !ok {
			c.String(http.StatusNotFound, "404 - Not Found")
		}

		c.HTML(http.StatusOK, "admin-employee-edit.html",
			map[string]interface{}{
				"Employee": employee,
			})
	})

	admin.POST("/employees/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "add" {
			pto, err := strconv.ParseFloat(c.PostForm("pto"), 32)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}

			startDate, err := time.Parse("2006-01-02", c.PostForm("startDate"))
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}

			var emp model.Employee
			emp.ID = 42
			emp.FirstName = c.PostForm("firstName")
			emp.LastName = c.PostForm("lastName")
			emp.Position = c.PostForm("position")
			emp.Status = "Active"
			emp.TotalPTO = float32(pto)
			emp.StartDate = startDate
			model.Employees["42"] = emp

			c.Redirect(http.StatusMovedPermanently, "/admin/employees/42")

		}
	})

	r.Static("/public", "./public")

	return r
}
