package platform

import (
	"baseservice/model/platform"
	"context"
	"encoding/json"
	"fmt"
	"jarvis/base/database"
	"jarvis/base/log"
	"time"
)

// 加载平台信息
func LoadPlatformInfo(now time.Time) bool {
	log.InfoLn("======================================================================")
	log.InfoLn("Start loading platform information from `jarvis`.`static_platform`")
	// 1.获取 MySQL 数据库连接
	conn, err := database.GetMySQLConn()
	if err != nil {
		log.FatalF("database.GetMySQLConn error : %s", err.Error())
		return true
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.FatalF("database.GetMySQLConn close error : %s", err.Error())
			return
		}
	}()

	platformList := make(platform.PlatformList, 0)

	// 查询 MySQL
	rows, err := conn.QueryContext(context.Background(), platformList.QueryOrder())
	if err != nil {
		log.FatalF("database.GetMySQLConn query error : %s", err.Error())
		return true
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.FatalF("database.GetMySQLConn rows close error : %s", err.Error())
			return
		}
	}()

	// 扫描
	for rows.Next() {
		platform := platform.Platform{}
		err := rows.Scan(&platform.ID, &platform.Name, &platform.Link, &platform.Owner, &platform.CreateAt, &platform.UpdateAt)
		if err != nil {
			log.ErrorF("database.GetMySQLConn query rows scan error : %s", err.Error())
			return true
		}
		platformList = append(platformList, platform)
	}

	// 遍历写入 Redis
	log.InfoF("Now we having %d platform:", len(platformList))
	for _, p := range platformList {
		data, err := json.Marshal(&p)
		if err != nil {
			log.ErrorF("Marshal [%+v] to []byte error : %s", p, err.Error())
			break
		}

		log.InfoF("Platform : %s", string(data))

		if err := platform.HSetPlatformInfoByID(fmt.Sprintf("%d", p.ID), string(data)); err != nil {
			log.ErrorF("HSetPlatformInfoByID [%s] to []byte error : %s", string(data), err.Error())
			break
		}
	}
	log.InfoLn("Finish loading platform information from `jarvis`.`static_platform`")
	log.InfoLn("======================================================================")
	return true
}
