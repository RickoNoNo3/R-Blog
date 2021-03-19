package middleware

import (
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/helper/userhelper"
)

func AdminAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// do something to verify the admin access
		if userhelper.CheckAdmin(c) {
			return next(c)
		} else {
			return echo.ErrBadRequest
		}
	}
}
