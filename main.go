package main

import (
	"net/http"

	_ "pix-console/docs"
	v1 "pix-console/v1/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {

	memberCluster := v1.EventDelegate{}
	go memberCluster.Start([]string{"192.168.1.104", "192.168.1.107"}, &memberCluster, "Node-1", 7946)
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	// 載入 HTML 目錄
	router.LoadHTMLGlob("templates/*")

	// 設定靜態目錄
	router.Static("/static", "static")

	// 設定 HTML 頁面
	router.GET("/", func(c *gin.Context) {
		// 檢查 URL 中是否包含 error 參數
		errorMessage := c.Query("error")

		// 將錯誤消息傳遞給 HTML 模板
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": errorMessage,
		})
	})

	router.POST("/login", v1.LoginHandler)

	router.GET("/logout", v1.LogoutHandler)

	// 驗證
	router.Use(v1.JWTAuthMiddleware)

	PageLink := gin.H{
		"links": []gin.H{
			{"text": "DashBoard", "href": "/dashboard", "class": "dashboard"},
			{"text": "Containers", "href": "/docker"},
			{"text": "Setting", "href": "/index"},
			{"text": "Docker-compose", "href": "/docker-compose"},
			{"text": "Service", "href": "/service"},
			{"text": "Feedback", "href": "/feedback"},
			{"text": "Logout", "href": "/logout"},
		},
	}

	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", PageLink)
	})

	router.GET("/docker", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docker.html", PageLink)
	})

	router.GET("/service", func(c *gin.Context) {

		c.HTML(http.StatusOK, "service.html", PageLink)
	})

	router.GET("/dashboard", func(c *gin.Context) {

		c.HTML(http.StatusOK, "dashboard.html", PageLink)
	})

	router.GET("/docker-compose", func(c *gin.Context) {
		c.HTML(http.StatusOK, "docker-compose.html", PageLink)
	})

	router.GET("/feedback", func(c *gin.Context) {
		c.HTML(http.StatusOK, "feedback.html", PageLink)
	})

	// 設定 swagger
	url := ginSwagger.URL("http://60.199.173.12:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// apiV1 路由設定
	apiV1 := router.Group("/api/v1")

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
	apiV1.GET("/cluster_docker", v1.ClusterDockerHandler)

	// Service 頁面
	apiV1.GET("/service", v1.ServiceHandler)
	apiV1.GET("/cluster_service", v1.ClusterServiceHandler)
	apiV1.GET("/server", v1.ServerHandler)

	router.Run(":8080")
}
