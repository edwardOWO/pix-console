package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"pix-console/models"
	"strings"
	"time"

	tool "pix-console/StuneTool"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// 更新 docker-compose 檔案
func backupService() bool {
	fmt.Print("update Docker Compose\n")
	return true
}

// 更新 docker-compose 檔案
func updateDockerCompose() bool {
	fmt.Print("update Docker Compose\n")
	return true
}

// 更新 Image
func pullImage() bool {
	fmt.Print("PullImage\n")
	return true
}

// 重啟服務
func restartService() bool {
	fmt.Print("RestartService\n")
	return true
}

// 檢查服務
func checkService() bool {
	fmt.Print("CheckService\n")
	return true
}
func getlatestVersion() bool {
	stunConfig := tool.StuneSetting{
		ClientID:     "pixCollector",
		ClientSecret: "(5vBX1Tu@DDPs0Om1Cfm",
		AuthURL:      "https://auth.tw.juiker.net",
		BrandID:      "juiker",
		Scope:        "tw:stune:basic",
	}

	err := tool.StuneDownload(stunConfig.GetAccessToken(), "container.txt", "edward")
	if err != nil {
		fmt.Print(err.Error())
		return false
	}
	return true
}

// 升級伺服器

func UpdateContainerHandler(c *gin.Context) {

	if !backupService() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}

	if !getlatestVersion() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}

	if !updateDockerCompose() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}
	if !pullImage() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}
	if !restartService() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}
	if !checkService() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})

}
func UpdateServerHandler(c *gin.Context) {

	// 開始更新 Server
	fmt.Printf("@@@@@@@@@Start Update Server@@@@@@@@@")

	patchServer()

	// 重啟服務
	fmt.Printf("@@@@@@@@@Restart Service@@@@@@@@@")
	command := exec.Command("systemctl", "restart", "pix-console")
	// 設置標準輸出和標準錯誤輸出
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// 執行命令
	err := command.Run()
	if err != nil {
		fmt.Printf("命令執行出錯: %s\n", err.Error())
		return
	}
}
func patchServer() {

	stunConfig := tool.StuneSetting{
		ClientID:     "pixCollector",
		ClientSecret: "(5vBX1Tu@DDPs0Om1Cfm",
		AuthURL:      "https://auth.tw.juiker.net",
		BrandID:      "juiker",
		Scope:        "tw:stune:basic",
	}

	err := tool.StuneDownload(stunConfig.GetAccessToken(), "version.txt", "edward")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	version := LoadVersion()

	err = tool.StuneDownload(stunConfig.GetAccessToken(), version, "edward")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	command := exec.Command("rpm", "-Uvh", version)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err = command.Run()
	if err != nil {
		fmt.Print(stderr.String())
		return
	}

}
func LoadVersion() string {
	// 打開文件，第二個參數是打開模式，這裡使用只讀模式
	file, err := os.Open("version.txt")
	if err != nil {
		fmt.Println("無法打開文件:", err)
		return ""
	}
	defer file.Close() // 確保在函數結束時關閉文件

	// 使用bufio.NewReader來讀取文件
	reader := bufio.NewReader(file)

	// 逐行讀取文件內容
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("讀取文件時發生錯誤:", err)
		return ""
	}

	line = strings.TrimRight(line, "\n")

	return line
}

// 提交各個服務的版本訊息
func CommitContainerHandler(c *gin.Context) {

	// 定義一個結構體來映射 JSON 中的屬性
	var requestData struct {
		CommitMessage string `json:"commitMessage"`
	}

	// 使用 BindJSON 方法將 JSON 參數綁定到 requestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	test := commitDockerfile(requestData.CommitMessage)

	stunConfig := tool.StuneSetting{
		ClientID:     "pixCollector",
		ClientSecret: "(5vBX1Tu@DDPs0Om1Cfm",
		AuthURL:      "https://auth.tw.juiker.net",
		BrandID:      "juiker",
		Scope:        "tw:stune:basic",
	}

	err := tool.StuneUpload(stunConfig.GetAccessToken(), "config/container.txt", "edward")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": test})
}

func commitDockerfile(commit string) models.ServerInfo {
	var containerInfo models.ServerInfo

	containerInfo.CommitMessage = commit

	file, err := os.Open("docker-compose.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 讀取文件內容
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	err = yaml.Unmarshal(content, &data)
	if err != nil {
		log.Fatal(err)
	}

	// 提取並打印 services 部分
	services, ok := data["services"].(map[interface{}]interface{})
	if !ok {
		log.Fatal("Services section not found")
	}

	now := time.Now()
	formattedTime := now.Format("2006-01-02 15:04:05")
	for serviceName, service := range services {
		// 提取並打印每個服務的 image 屬性
		image, ok := service.(map[interface{}]interface{})["image"]
		if !ok {
			log.Printf("Image not found for service %v\n", serviceName)
			continue
		}

		imageStr := fmt.Sprintf("%v", image)
		parts := strings.Split(imageStr, ":")
		if len(parts) != 2 {
			log.Printf("Invalid image format for service %v: %v\n", serviceName, imageStr)
			continue
		}
		imageName := parts[0]
		imageTag := parts[1]

		test := models.Image{
			UpdateTime:  formattedTime,
			ServiceName: fmt.Sprintf("%v", serviceName),
			ImageName:   imageName,
			ImageTag:    imageTag,
		}

		containerInfo.ContainerInfo = append(containerInfo.ContainerInfo, test)
	}

	// 將結構寫入 JSON 檔案
	err = writeJSONToFile(containerInfo, "config/container.txt")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return containerInfo
	}
	return containerInfo
}
