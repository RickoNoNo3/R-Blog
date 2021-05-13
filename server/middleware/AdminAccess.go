package middleware

import (
	"github.com/labstack/echo/v4"
	"rickonono3/r-blog/logger"

	"rickonono3/r-blog/helper/userhelper"
)

func AdminAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// do something to verify the admin access
		if userhelper.CheckAdmin(c) {
			return next(c)
		} else {
			logger.L.Warn("[Server]", "无权的管理员访问: ", c.RealIP())
			return echo.ErrBadRequest
		}
	}
}
