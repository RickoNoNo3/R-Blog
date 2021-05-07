package cleaner

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	// queChan
	//
	// 请求队列, 传递的是FilePath
	queChan = make(chan string, 1000)

	// exitChan
	//
	// 想退出程序时，应该给exitChan发送一个byte值, 触发停摆机制,
	// 然后准备再接收一个byte值, 接收到时代表清理器确实已退出.
	exitChan = make(chan byte)

	exiting      = false
	exitingMutex sync.Mutex
	exited       = false
	exitedMutex  sync.Mutex
)

// SendRequest
//
// 给清理器发送一条清理请求，成功返回true，不成功返回false（清理器已经准备关闭）
func SendRequest(req string) bool {
	// 一旦正在关闭，请求就不能成功发送了
	exitingMutex.Lock()
	defer exitingMutex.Unlock()
	if !exiting {
		queChan <- req
		return true
	} else {
		return false
	}
}

func Exit() bool {
	exitChan <- 0
	ch := make(chan byte)
	go func() {
		t := <-exitChan
		ch <- t
	}()
	go func() {
		time.Sleep(10 * time.Second)
		ch <- 0
	}()
	_ = <-ch
	exitedMutex.Lock()
	defer exitedMutex.Unlock()
	return exited
}

func Run() {
	for {
		var req string
		// 在未触发退出事件时，阻塞接收清理队列数据
		// 在已触发退出事件，未处理结束时，非阻塞接收清理队列数据
		// 非阻塞接收一旦为空，即视为处理结束，退出清理队列线程
		exitingMutex.Lock()
		if !exiting {
			exitingMutex.Unlock()
			req = <-queChan
		} else {
			exitingMutex.Unlock()
			var ok bool
			req, ok = <-queChan
			if !ok {
				exitedMutex.Lock()
				exited = true
				exitedMutex.Unlock()
				break
			}
		}
		fmt.Println("cleaning for: " + req)
		if _, ok := <-exitChan; !ok {
			err := os.Remove(req)
			if err != nil && err != os.ErrNotExist {
				SendRequest(req)
			}
		} else {
			exitingMutex.Lock()
			exiting = true
			exitingMutex.Unlock()
		}
	}
}
