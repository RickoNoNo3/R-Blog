package userhelper

import (
	"github.com/labstack/echo/v4"

	"rickonono3/r-blog/objects"
)

func CheckAdmin(c echo.Context) (isAdmin bool) {
	if cookie, err := c.Cookie("blog-login"); err != nil {
		return false
	} else {
		if cookie.Value != objects.RuntimeEnv.MustGet("AdminHash").ValStr() {
			return false
		} else {
			return true
		}
	}
}
