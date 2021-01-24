package main

import (
	"globalserver/module/platform"
	"jarvis/base/database"
	uTime "jarvis/util/time"
	"log"
	"time"
)

func init() {
	// 初始化 MySQL
	if err := database.InitializeMySQL(
		"frank",
		"frank123",
		"localhost",
		7000,
		"jarvis",
	); err != nil {
		log.Panicf("Initialize MySQL error : %s", err.Error())
		return
	}

	// 设置 MySQL
	database.SetUpMySQL(time.Minute*time.Duration(5), 10, 30)

	// 初始化 Redis
	database.InitializeRedis(time.Minute*time.Duration(5), 10, 30, "localhost", 8000, "frank123")

	// 初始化 Mongo
	if err := database.InitializeMongo(
		"frank",
		"frank123",
		"jarvis",
		"localhost",
		9000, time.Minute*time.Duration(5), 30); err != nil {
		log.Panicf("Initialize Mongo error : %s", err.Error())
		return
	}
}

func main() {
	closeChannel := make(chan struct{})

	ticker := uTime.NewTicker(time.Duration(30)*time.Minute, platform.LoadPlatformInfo)
	ticker.Run()

	<-closeChannel
}
