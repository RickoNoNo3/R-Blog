package logger

import (
	"fmt"
	"io"
	"rickonono3/r-blog/helper/loggerhelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/objects"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	output io.Writer
	mutex  sync.Mutex
}

func (l *Logger) print(level string, i ...interface{}) {
	switch i[0].(type) {
	case nil:
		return
	case error:
		if i[0].(error) == nil {
			return
		}
	}
	logStmt := make([]string, 0)
	logStmt = append(logStmt, loggerhelper.EscapeComma(typehelper.MustItoa64(time.Now().Unix())))
	logStmt = append(logStmt, level)
	logStmt = append(logStmt, loggerhelper.EscapeComma(fmt.Sprint(i...)))
	logStr := strings.Join(logStmt, ",")
	l.mutex.Lock()
	l.output.Write([]byte(logStr + "\n"))
	l.mutex.Unlock()
}

func (l *Logger) Debug(i ...interface{}) {
	if objects.Config.MustGet("IsInDebug").ValBool() {
		l.print("DEBUG", i...)
	}
}

func (l *Logger) Info(i ...interface{}) {
	l.print("INFO", i...)
}

func (l *Logger) Warn(i ...interface{}) {
	l.print("WARN", i...)
}

func (l *Logger) Error(i ...interface{}) {
	l.print("ERROR", i...)
}

func (l *Logger) Panic(i ...interface{}) {
	l.print("PANIC", i...)
}

// ------------------------------------------------

var L *Logger

func InitLogger(logFile io.Writer) {
	L = &Logger{
		output: logFile,
	}
}
