package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func apiRegister() {
	api := E.Group("/api")
	apiAdmin := api.Group("/admin", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// do something to verify the admin access
			fmt.Println("admin")
			return next(c)
		}
	})
	apiAdmin.POST("/test", func(c echo.Context) error {
		return c.String(200, "success!")
	})
	// TODO: login
}
