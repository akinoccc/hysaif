package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akinoccc/hysaif/api/config"
	"github.com/akinoccc/hysaif/api/models"
	"github.com/akinoccc/hysaif/api/packages/notification"
	"github.com/akinoccc/hysaif/api/packages/permission"
	"github.com/akinoccc/hysaif/api/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("config.json"); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库
	models.InitDB()

	// 初始化Casbin权限管理器
	permission.GetCasbinManager(models.DB)

	// 启动定时任务服务
	notification.Start()
	defer notification.Stop()

	// 创建Gin路由
	r := gin.Default()

	router.InitRouter(r)

	// 启动服务器
	host := config.AppConfig.Server.Host
	port := config.AppConfig.Server.Port
	if host == "" {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 50010
	}

	log.Printf("服务器启动在 %s:%d", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r))
}
