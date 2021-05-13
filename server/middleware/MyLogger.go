package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"rickonono3/r-blog/helper/loggerhelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/logger"
	"rickonono3/r-blog/objects"
	"strings"
	"time"
)

func MyLogger() echo.MiddlewareFunc {
	logFile, err := os.Create(objects.CWD + objects.Config.MustGet("LogFile.WebLog").ValStr())
	if err != nil {
		logFile = os.Stdout
	}
	logger.L.Info("[Server]", "网络日志已启用")
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
			logStmt = append(logStmt, loggerhelper.EscapeComma(typehelper.MustItoa64(time.Now().Unix())))
			logStmt = append(logStmt, loggerhelper.EscapeComma(c.RealIP()))
			logStmt = append(logStmt, loggerhelper.EscapeComma(req.Method))
			logStmt = append(logStmt, loggerhelper.EscapeComma(req.RequestURI))
			logStmt = append(logStmt, loggerhelper.EscapeComma(req.UserAgent()))
			logStmt = append(logStmt, loggerhelper.EscapeComma(typehelper.MustAtoitoa(req.Header.Get(echo.HeaderContentLength))))
			logStmt = append(logStmt, loggerhelper.EscapeComma(typehelper.MustItoa64(res.Size)))
			logStmt = append(logStmt, loggerhelper.EscapeComma(typehelper.MustItoa(res.Status)))
			logStmt = append(logStmt, loggerhelper.EscapeComma(typehelper.MustItoa64(int64(stop.Sub(start)))))
			if err == nil {
				logStmt = append(logStmt, loggerhelper.EscapeComma(""))
			} else {
				logStmt = append(logStmt, loggerhelper.EscapeComma(err.Error()))
			}
			logStr := strings.Join(logStmt, ",")
			logFile.Write([]byte(logStr + "\n"))
			if res.Status != http.StatusOK {
				logger.L.Warn("[Server]", "URL访问异常(", res.Status, "): ", req.RequestURI)
			}
			return nil
		}
	}
}
