package main

import (
	"baseservice/common/remoteLogHook"
	"globalserver/module/platform"
	"io"
	"jarvis/base/database"
	"jarvis/base/database/redis"
	"jarvis/base/log"
	uTime "jarvis/util/time"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// 远程日志聚合地址
	LogRemoteAddress = "logserver:10000"
)

var (
	lh io.WriteCloser
)

func init() {
	// 新建远程日志钩子
	nlh, err := remoteLogHook.NewSocketRemoteHook(LogRemoteAddress)
	if err != nil {
		log.FatalF("New remote hook error : %s", err.Error())
		return
	}
	lh = nlh
	log.SetHook(lh)

	// 初始化 MySQL
	if err := database.InitializeMySQL(
		"frank",
		"frank123",
		"mysql-service",
		3306,
		"jarvis",
	); err != nil {
		log.FatalF("Initialize MySQL error : %s", err.Error())
		return
	}

	// 设置 MySQL
	database.SetUpMySQL(time.Minute*time.Duration(5), 10, 30)

	// 初始化 Redis
	redis.InitializeRedis(time.Minute*time.Duration(5), 10, 30, "redis-service", 6379, "frank123")

	// 初始化 Mongo
	if err := database.InitializeMongo(
		"frank",
		"frank123",
		"jarvis",
		"mongo-service",
		27017, time.Minute*time.Duration(5), 30); err != nil {
		log.FatalF("Initialize Mongo error : %s", err.Error())
		return
	}
}

func main() {
	// 1.启动全局平台信息更新
	ticker := uTime.NewTicker(time.Duration(30)*time.Minute, platform.LoadPlatformInfo)
	ticker.Run()

	monitorSystemSignal()
}

// 监听系统信号
// kill -SIGQUIT [进程号] : 杀死当前进程
func monitorSystemSignal() {
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGQUIT)
	select {
	case <-sc:
	case <-sc:
		_ = lh.Close()
		log.InfoF("Done")
	}
}
