package controller

import (
	"net/http"
	"time"

	"github.com/NipunSood/Employee-Self-Service-Portal/src/loginMiddleware"
	"github.com/NipunSood/Employee-Self-Service-Portal/src/model"
	"github.com/gin-gonic/gin"
)

func RegisterRouters() *gin.Engine {
	r := gin.Default()
	r.Use(loginMiddleware.LoginMiddleware)
	r.LoadHTMLGlob("templates/**/*.html")
	r.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.Any("/login", func(c *gin.Context) {
		employeeNumber := c.PostForm("employeeNumber")
		password := c.PostForm("password")

		for _, identity := range loginMiddleware.Identities {
			if identity.EmployeeNumber == employeeNumber &&
				identity.Password == password {
				lc := loginMiddleware.LoginCookie{
					Value:      employeeNumber,
					Expiration: time.Now().Add(24 * time.Hour),
					Origin:     c.Request.RemoteAddr,
				}
				loginMiddleware.LoginCookies[lc.Value] = &lc
				maxAge := lc.Expiration.Unix() - time.Now().Unix()
				c.SetCookie(loginMiddleware.LoginCookieName, lc.Value, int(maxAge), "", "", false, true)

				c.Redirect(http.StatusTemporaryRedirect, "/")
				return
			}
		}
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
			gin.H{
				"TimesOff": timesOff,
			})

	})

	r.POST("/employees/:id/vacation/new", func(c *gin.Context) {
		var timeOff model.TimeOff
		err := c.BindJSON(&timeOff)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		id := c.Param("id")

		timesOff, ok := model.TimesOff[id]
		if !ok {
			model.TimesOff[id] = []model.TimeOff{}
		}

		model.TimesOff[id] = append(timesOff, timeOff)

		c.JSON(http.StatusCreated, &timeOff)
	})

	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"nipun@trailerstop.me": "nipunrocks",
	}))
	admin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin-overview.html",
			gin.H{
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

			var emp model.Employee
			err = c.Bind(&emp)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				return
			}
			emp.ID = 42
			emp.Status = "Active"
			emp.StartDate = startDate
			model.Employees["42"] = emp

			c.Redirect(http.StatusMovedPermanently, "/admin/employees/42")

		}
	})

	r.Static("/public", "./public")

	return r
}
