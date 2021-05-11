package main

import (
	"context"
	"os"
	"rickonono3/r-blog/cleaner"
	"rickonono3/r-blog/cleaner/monitor"
	"rickonono3/r-blog/data"
	"rickonono3/r-blog/helper/cmdhelper"
	"rickonono3/r-blog/helper/typehelper"
	"rickonono3/r-blog/logger"
	"rickonono3/r-blog/objects"
	"rickonono3/r-blog/server"
	"time"
)

var (
	commandCycle = true
	// exitCode
	//
	//  -1 stop unexpectedly
	//   0 stop normally
	//   1 stop and need to restart
	exitCode = 0
)

// TODO: 当前版本：部署工具
// TODO: 后续版本：日志、导航、标签、搜索、打印、统计、分享、捐赠、强数据收集、防爆破
func main() {
	defer func() {
		if r := recover(); r != nil {
			os.Exit(-1)
		} else {
			os.Exit(exitCode)
		}
	}()

	// objects
	objects.Init()

	// pid
	pid := typehelper.MustItoa(os.Getpid())
	defer func() {
		logger.L.Info("[Main]", "进程已退出(PID:"+pid+")")
	}()

	// console log
	{
		var (
			logFile *os.File
			err     error
		)
		if objects.Config.MustGet("IsInDebug").ValBool() {
			logFile = os.Stdout
		} else {
			logFile, err = os.Create(objects.CWD + objects.Config.MustGet("LogFile.ConsoleLog").ValStr())
			if err != nil {
				logFile = os.Stdout
			}
		}
		logger.InitLogger(logFile)
		logger.L.Info("[Main]", "进程启动(PID:"+pid+")")
		logger.L.Info("[Main]", "控制台日志已启用")
	}

	// db
	logger.L.Info("[Main]", "挂载数据库")
	data.OpenDB("blog.db")
	defer func() {
		data.CloseDB()
		logger.L.Info("[Main]", "卸载数据库")
	}()

	// cleaner
	logger.L.Info("[Main]", "启动清理队列")
	go cleaner.Run()
	defer func() {
		cleaner.Exit()
		logger.L.Info("[Main]", "清理队列已关闭")
	}()

	// cleaner monitor
	logger.L.Info("[Main]", "启动清理监控器")
	go monitor.Run()
	defer func() {
		monitor.Exit()
		logger.L.Info("[Main]", "清理监控器已关闭")
	}()

	// server
	server.Init()
	port := objects.Config.MustGet("ServerPort").ValInt()
	logger.L.Info("[Main]", "启动服务器于", port, "端口")
	go func() {
		logger.L.Panic("[Server]", server.E.Start(":"+typehelper.MustItoa(port)))
		cmdhelper.CloseInput()
	}()
	defer func() {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
		defer cancel()
		if err := server.E.Shutdown(ctx); err != nil {
			logger.L.Error(err)
			server.E.Close()
		}
	}()

	// command
	logger.L.Info("[Main]", "监听控制台...")
	cmdhelper.InitCmd()
	for commandCycle {
		str := cmdhelper.GetInput()
		logger.L.Debug("[Main]", "控制台输入: ", str)
		switch str {
		case "clean": // 强制执行文件清理
			monitor.Manually()
		case "exit": // 退出程序
			commandCycle = false
		case "restart": // 重启程序
			exitCode = 1
			commandCycle = false
		}
	}

	logger.L.Info("[Main]", "准备退出进程")
}
