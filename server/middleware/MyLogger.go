package middleware

import (
	"github.com/labstack/echo/v4"
	"io"
	"regexp"
	"rickonono3/r-blog/helper/typehelper"
	"strings"
	"time"
)

func MyLogger(logFile io.Writer) echo.MiddlewareFunc {
	var escapeComma = func(str string) string {
		return regexp.MustCompile("\\s*,\\s*").ReplaceAllString(str, "|")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()
			logStmt := make([]string, 0)
			logStmt = append(logStmt, escapeComma(time.Now().Format(time.RFC3339)))
			logStmt = append(logStmt, escapeComma(c.RealIP()))
			logStmt = append(logStmt, escapeComma(req.Method))
			logStmt = append(logStmt, escapeComma(req.RequestURI))
			logStmt = append(logStmt, escapeComma(req.UserAgent()))
			logStmt = append(logStmt, escapeComma(typehelper.MustAtoitoa(req.Header.Get(echo.HeaderContentLength))))
			logStmt = append(logStmt, escapeComma(typehelper.MustItoa64(res.Size)))
			logStmt = append(logStmt, escapeComma(typehelper.MustItoa(res.Status)))
			logStmt = append(logStmt, escapeComma(typehelper.MustItoa64(int64(stop.Sub(start)))))
			if err == nil {
				logStmt = append(logStmt, escapeComma(""))
			} else {
				logStmt = append(logStmt, escapeComma(err.Error()))
			}
			logStr := strings.Join(logStmt, ",")
			logFile.Write([]byte(logStr + "\n"))
			return
		}
	}
}
