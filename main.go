package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"pix-console/common"
	"pix-console/controllers"
	_ "pix-console/docs"
	"pix-console/view"

	v1 "pix-console/controllers"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/memberlist"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var logger *log.Logger

// Main manages main golang application
type Main struct {
	router     *gin.Engine
	memberlist *memberlist.Memberlist
}

func (m *Main) initServer() error {
	var err error
	// Load config file
	err = common.LoadConfig()
	if err != nil {
		return err
	}

	f, err := os.OpenFile("logs/memberlist.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("open file error=%v", err)
	}
	//defer f.Close()

	logger = log.New(f, "###################", log.Ldate|log.Ltime)

	// 啟動 memberlist 功能,並
	list, _, _ := controllers.StartMemberlist(logger, f)
	m.memberlist = list

	// Initialize mongo database
	//err = databases.Database.Init()
	//if err != nil {
	//	return err
	//}

	// Setting Gin Logger
	if common.Config.EnableGinFileLog {
		f, _ := os.Create("logs/gin.log")
		if common.Config.EnableGinConsoleLog {
			gin.DefaultWriter = io.MultiWriter(os.Stdout, f)
		} else {
			gin.DefaultWriter = io.MultiWriter(f)
		}
	} else {
		if !common.Config.EnableGinConsoleLog {
			gin.DefaultWriter = io.MultiWriter()
		}
	}

	m.router = gin.Default()

	return nil
}

func main() {

	m := Main{}
	if m.initServer() != nil {
		return
	}

	// 初始化 Server 物件
	c := v1.Server{}

	// 將 memberlist 賦予站台使用
	c.Memberlist = m.memberlist

	// 將使用者帳號讀進站台
	c.UserAcount = common.LoadAccountConfig()

	// 載入 HTML 目錄
	m.router.LoadHTMLGlob("templates/*")

	// 設定靜態目錄
	m.router.Static("/static", "static")

	// 設定 HTML 頁面
	m.router.GET("/", func(c *gin.Context) {
		// 檢查 URL 中是否包含 error 參數
		errorMessage := c.Query("error")

		// 將錯誤消息傳遞給 HTML 模板
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": errorMessage,
		})
	})

	// 使用者登入
	m.router.POST("/login", c.LoginHandler)

	// 使用者登出
	m.router.GET("/logout", c.LogoutHandler)

	// 驗證系統
	m.router.Use(c.JWTAuthMiddleware)

	// 在介面上產生站台連結
	PageLink := view.CreatePageLink()

	m.router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", PageLink)
	})

	m.router.GET("/host", func(c *gin.Context) {
		c.HTML(http.StatusOK, "host.html", PageLink)
	})

	m.router.GET("/docker", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docker.html", PageLink)
	})

	m.router.GET("/service", func(c *gin.Context) {

		c.HTML(http.StatusOK, "service.html", PageLink)
	})

	m.router.GET("/dashboard", func(c *gin.Context) {

		c.HTML(http.StatusOK, "dashboard.html", PageLink)
	})

	m.router.GET("/docker-compose", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docker-compose.html", PageLink)
	})

	m.router.GET("/feedback", func(c *gin.Context) {
		c.HTML(http.StatusOK, "feedback.html", PageLink)
	})

	// 設定 swagger
	url := ginSwagger.URL("http://60.199.173.12:8080/swagger/doc.json")
	m.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// apiV1 路由設定
	apiV1 := m.router.Group("/api/v1")
	{
		// 設定 API
		apiV1.GET("/createfile", v1.CreateFileHandler)
		apiV1.POST("/checkfile", v1.CheckFileHandler)
		apiV1.GET("/checkmemory", v1.CheckMemoryHandler)
		apiV1.POST("/startservice", v1.StartServiceHandler)
		apiV1.POST("/start_pix_compose", v1.StartPixComposeHandler)
		apiV1.POST("/stop_pix_compose", v1.StopPixComposeHandler)

		// Docker-compose 頁面
		apiV1.GET("/docker_compose", v1.DockerComposeHandler)
		apiV1.POST("/upload", v1.UploadDockerComposeHandler)
		apiV1.GET("/download", v1.DownloadConfigHandler)

		// 系統回報 回傳log 頁面
		apiV1.POST("/uploadToStune", v1.UploadToStune)
		apiV1.GET("/downloadFromStune", v1.DownloadFromStune)
		apiV1.GET("/clusterDownloadFromStune", v1.ClusterDownloadFromStune)

		// Containers 頁面
		apiV1.GET("/docker", v1.DockerHandler)
		apiV1.GET("/cluster_docker", c.ClusterDockerHandler)

		// Service 頁面
		apiV1.GET("/service", v1.ServiceHandler)
		apiV1.GET("/cluster_service", c.ClusterServiceHandler)

		// 取得主機群資訊
		apiV1.GET("/serverlist", c.ServerlistHandler)

		// 檢查主機 Port 號
		apiV1.GET("/monitor", c.MonitorHandler)

	}

	m.router.Run(common.Config.Port)
}
