package cleaner

import (
	"os"
	"rickonono3/r-blog/logger"
	"sync"
	"syscall"
	"time"
)

var (
	// queChan 是与 Cleaner 通信的主通道, 可传递"__EXIT__"消息或要删除的文件路径,
	// 在未接到 exiting 信号时循环使用阻塞接收，否则非阻塞接收。每次接收后新开线程进行删除，不阻塞下一次循环
	queChan = make(chan string, 1000)

	// exitChan 是专为 Exit 方法接收确实退出信号而设，exitChan 有信号则代表 Cleaner 已确实退出
	exitChan = make(chan byte)

	// queMap 是为了缓解同一文件多次加入清理队列造成的冗余而使用的去重机制。在开始处理一个文件时加入map，结束处理时从map中移除。此机制不能完全避免重复，但可有效减少。同时，在已接到退出信号后，`len(queMap)==0`是 Cleaner 从【准备退出】到【确实退出】状态转移的必要条件。
	queMap      = make(map[string]byte)
	queMapMutex sync.Mutex

	// exiting 信号标记当前是否处于【准备退出】模式，在准备退出模式下，发送信号(SendRequest)不再有效，Cleaner 主线程立即清理所有 queChan 内容。
	exiting      = false
	exitingMutex sync.Mutex

	// exited 信号标记当前是否已经【确实退出】，也即 queChan 内再无任何新信号，而所有的删除任务也已全部结束。
	exited      = false
	exitedMutex sync.Mutex
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
		logger.L.Debug("[Cleaner]", "添加["+req+"]到清理队列")
		return true
	} else {
		return false
	}
}

func Exit() bool {
	queChan <- "__EXIT__"
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
	logger.L.Debug("[Cleaner]", "已退出Cleaner")
	return exited
}

func Run() {
	var cycle = true
	for cycle {
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
			select {
			case req = <-queChan:
				// do nothing
			default:
				logger.L.Debug("[Cleaner]", "队列处理完毕，即将退出Cleaner")
				cycle = false
				continue
			}
		}
		queMapMutex.Lock()
		if _, ok := queMap[req]; !ok {
			queMap[req] = 0
			queMapMutex.Unlock()
			if req == "__EXIT__" {
				exitingMutex.Lock()
				exiting = true
				exitingMutex.Unlock()
				queMapMutex.Lock()
				delete(queMap, req)
				queMapMutex.Unlock()
			} else {
				go func() {
					logger.L.Info("[Cleaner]", "正在清理: "+req)
					err := os.Remove(req)
					if err == nil ||
						err == os.ErrNotExist ||
						err.(*os.PathError).Err == syscall.ERROR_FILE_NOT_FOUND ||
						err.(*os.PathError).Err == syscall.ERROR_PATH_NOT_FOUND {
						logger.L.Info("[Cleaner]", "清理成功")
						queMapMutex.Lock()
						delete(queMap, req)
						queMapMutex.Unlock()
					} else {
						logger.L.Info("[Cleaner]", "清理失败")
						queMapMutex.Lock()
						delete(queMap, req)
						queMapMutex.Unlock()
						time.Sleep(2 * time.Second)
						SendRequest(req)
					}
				}()
			}
		} else {
			queMapMutex.Unlock()
		}
	}
	for {
		queMapMutex.Lock()
		if len(queMap) == 0 {
			break
		}
		queMapMutex.Unlock()
	}
	exitedMutex.Lock()
	exited = true
	exitedMutex.Unlock()
	exitChan <- 1
}
