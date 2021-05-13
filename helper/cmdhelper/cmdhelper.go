// Package cmdhelper
//
// ------------------
//
// 控制台操作
//
// 1. 用户可从 os.Stdin 进行命令输入
//
// 2. 所有的kill信号都自动视为输入 "exit"
//
// 3. 程序内部本身可以通过 Input 函数发送命令输入
package cmdhelper

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	sysChan   = make(chan os.Signal, 1)
	inputChan = make(chan string, 100)
)

func InitCmd() {
	signal.Notify(sysChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	// syscall
	go func() {
		for {
			<-sysChan
			inputChan <- "exit"
		}
	}()
	// stdin
	go func() {
		var str string
		for {
			n, err := fmt.Scanln(&str)
			if err != nil {
				break
			}
			if n == 0 {
				continue
			}
			str = strings.TrimSpace(str)
			if str == "" {
				continue
			}
			inputChan <- str
		}
	}()
}

func Input(str string) {
	inputChan <- str
}

func GetInput() (str string) {
	return <-inputChan
}

func CloseInput() {
	inputChan <- "exit"
}
