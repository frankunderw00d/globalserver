package platform

import (
	"baseservice/model/platform"
	"context"
	"encoding/json"
	"fmt"
	"jarvis/base/database"
	"log"
	"time"
)

// 加载平台信息
func LoadPlatformInfo(now time.Time) bool {
	log.Println("======================================================================")
	log.Printf("Start loading platform information from `jarvis`.`static_platform`")
	// 1.获取 MySQL 数据库连接
	conn, err := database.GetMySQLConn()
	if err != nil {
		log.Fatalf("database.GetMySQLConn error : %s", err.Error())
		return true
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("database.GetMySQLConn close error : %s", err.Error())
			return
		}
	}()

	platformList := make(platform.PlatformList, 0)

	// 查询 MySQL
	rows, err := conn.QueryContext(context.Background(), platformList.QueryOrder())
	if err != nil {
		log.Fatalf("database.GetMySQLConn query error : %s", err.Error())
		return true
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("database.GetMySQLConn rows close error : %s", err.Error())
			return
		}
	}()

	// 扫描
	for rows.Next() {
		platform := platform.Platform{}
		err := rows.Scan(&platform.ID, &platform.Name, &platform.Link, &platform.Owner, &platform.CreateAt, &platform.UpdateAt)
		if err != nil {
			log.Printf("database.GetMySQLConn query rows scan error : %s", err.Error())
			return true
		}
		platformList = append(platformList, platform)
	}

	// 遍历写入 Redis
	log.Printf("Now we having %d platform:", len(platformList))
	for _, p := range platformList {
		data, err := json.Marshal(&p)
		if err != nil {
			log.Printf("Marshal [%+v] to []byte error : %s", p, err.Error())
			break
		}

		log.Printf("Platform : %s", string(data))

		if err := platform.HSetPlatformInfoByID(fmt.Sprintf("%d", p.ID), string(data)); err != nil {
			log.Printf("HSetPlatformInfoByID [%s] to []byte error : %s", string(data), err.Error())
			break
		}
	}
	log.Printf("Finish loading platform information from `jarvis`.`static_platform`")
	log.Println("======================================================================")
	return true
}
