package main

import (
	logModule "globalserver/module/log"
	"globalserver/module/platform"
	"jarvis/base/database"
	"jarvis/base/database/redis"
	"jarvis/base/log"
	"jarvis/base/network"
	uTime "jarvis/util/time"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	// 自定义服务最大连接数量
	CustomMaxConnection = 5000
	// 自定义服务消息管道大小
	CustomIntoStreamSize = 10000
	// 日志 Socket 接收地址
	LogSocketAddress = ":10000"
)

var (
	service network.Service
)

func init() {
	// 实例化服务
	service = network.NewService(
		CustomMaxConnection,
		CustomIntoStreamSize,
	)

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

	// 2.注册模块
	if err := service.RegisterModule(logModule.NewModule()); err != nil {
		log.ErrorF("Register module error : %s", err)
		return
	}

	// 3.启动
	if err := service.Run(
		network.NewSocketGate(LogSocketAddress), // Socket 入口
	); err != nil {
		log.ErrorF("Register observer error : %s", err)
		return
	}

	// 4.监听系统信号
	monitorSystemSignal()
}

// 监听系统信号
// kill -SIGQUIT [进程号] : 杀死当前进程
func monitorSystemSignal() {
	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGQUIT)
	select {
	case <-sc:
		log.InfoF("Done")
	}
}
