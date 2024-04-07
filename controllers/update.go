package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"pix-console/common"
	"pix-console/models"
	"pix-console/utils"
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
func updateDockerCompose() (bool, error) {

	filePath := "config/container.json"
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return false, err
	}
	var containerInfo models.ServerInfo
	err = json.Unmarshal(jsonData, &containerInfo)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return false, err
	}

	// Read the contents of the file
	filePath = "/opt/pix/run/docker-compose-pro.yml"
	yamlData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)

	}

	var data map[string]interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 提取並打印 services 部分
	services, ok := data["services"].(map[interface{}]interface{})
	if !ok {
		fmt.Print("not get services")
	}

	for _, containerInfo := range containerInfo.ContainerInfo {

		service, exists := services[containerInfo.ServiceName].(map[interface{}]interface{})
		if !exists {
			continue
		}

		service["image"] = containerInfo.ImageName + ":" + containerInfo.ImageTag // Replace with the desired image and tag

	}

	// Marshal the modified data back to YAML
	updatedYAMLData, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Write the updated content back to the file
	err = ioutil.WriteFile(filePath, updatedYAMLData, 0644)
	if err != nil {
		fmt.Print(err.Error())
	}
	return true, nil
}

// 更新 Image
func pullImage() bool {

	fmt.Print("RestartService\n")

	// 設置標準輸出和標準錯誤輸出
	command := exec.Command("docker-compose", "-f", "/opt/pix/run/docker-compose-pro.yml", "pull")
	// 設置標準輸出和標準錯誤輸出
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// 執行命令
	err := command.Run()
	if err != nil {
		fmt.Printf("命令執行出錯: %s\n", err.Error())
		return false
	}
	return true
}

// 重啟服務
func restartService() bool {
	fmt.Print("RestartService\n")

	// 設置標準輸出和標準錯誤輸出
	command := exec.Command("systemctl", "stop", "pix-compose")
	// 設置標準輸出和標準錯誤輸出
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// 執行命令
	err := command.Run()
	if err != nil {
		fmt.Printf("命令執行出錯: %s\n", err.Error())
		return false
	}

	command = exec.Command("systemctl", "start", "pix-compose")
	// 設置標準輸出和標準錯誤輸出
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// 執行命令
	err = command.Run()
	if err != nil {
		fmt.Printf("命令執行出錯: %s\n", err.Error())
		return false
	}
	return true
}

// 檢查服務
func checkService() bool {
	fmt.Print("CheckService\n")
	return true
}
func getlatestVersion() (bool, error) {
	stunConfig := tool.StuneSetting{
		ClientID:     "pixCollector",
		ClientSecret: "(5vBX1Tu@DDPs0Om1Cfm",
		AuthURL:      "https://auth.tw.juiker.net",
		BrandID:      "juiker",
		Scope:        "tw:stune:basic",
	}

	err := tool.StuneDownload(stunConfig.GetAccessToken(), "config/container.json", "edward")
	if err != nil {
		fmt.Print(err.Error())
		return false, err
	}
	return true, nil
}

// 升級伺服器

func UpdateContainerHandler(c *gin.Context) {

	if !backupService() {
		c.JSON(http.StatusOK, gin.H{"message": "error"})
		return
	}

	status, err := getlatestVersion()

	if !status {
		c.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	status, err = updateDockerCompose()

	if !status {
		c.JSON(http.StatusOK, gin.H{"message": "update Docker Compose Error" + err.Error()})
		return
	}
	if !pullImage() {
		c.JSON(http.StatusOK, gin.H{"message": "pull image error"})
		return
	}
	if !restartService() {
		c.JSON(http.StatusOK, gin.H{"message": "restart service error"})
		return
	}
	if !checkService() {
		c.JSON(http.StatusOK, gin.H{"message": "check service error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})

}

func ClusterUpdateContainerHandler(c *gin.Context) {

	// 定義一個結構體來映射 JSON 中的屬性
	var requestData struct {
		UpdateHost string `json:"updatehost"`
	}

	// 使用 BindJSON 方法將 JSON 參數綁定到 requestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiUrl := fmt.Sprintf("http://%s%s/api/v1/updateContainer", requestData.UpdateHost, common.Config.Port)
	client := &http.Client{
		Timeout: time.Second * 120, // 設定超時時間為 5 秒
	}

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	// 添加 Bearer Token 到標頭
	req.Header.Set("Authorization", "Bearer "+"sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)")

	response, err := client.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	} else {
		c.JSON(http.StatusOK, "更新主機: "+requestData.UpdateHost+"更新結果: "+string(body))
	}

	response.Body.Close()

}
func (s *Server) UpdateServerHandler(c *gin.Context) {

	err := patchServer()

	if err != nil {
		utils.Log(s.Logger, "Error", utils.Trace()+" URL: "+c.Request.URL.Path, "Update Pix-console Error: "+err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	command := exec.Command("systemctl", "restart", "pix-console")
	output, err := command.CombinedOutput()

	if err != nil {
		utils.Log(s.Logger, "Error", utils.Trace()+" URL: "+c.Request.URL.Path, "Update Pix-console Error ErrorCode:"+err.Error()+" "+string(output))
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "OK")
}

func (s *Server) ClusterUpdateServerHandler(c *gin.Context) {

	// 定義一個結構體來映射 JSON 中的屬性
	var requestData struct {
		UpdateHost string `json:"updatehost"`
	}

	// 使用 BindJSON 方法將 JSON 參數綁定到 requestData
	if err := c.BindJSON(&requestData); err != nil {
		utils.Log(s.Logger, "Error", utils.Trace()+" URL: "+c.Request.URL.Path, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiUrl := fmt.Sprintf("http://%s%s/api/v1/updateServer", requestData.UpdateHost, common.Config.Port)
	client := &http.Client{
		Timeout: time.Second * 10, // 設定超時時間為 5 秒
	}

	req, err := http.NewRequest("POST", apiUrl, nil)
	if err != nil {
		utils.Log(s.Logger, "Error", utils.Trace()+" URL: "+c.Request.URL.Path, err.Error())
	}

	// 添加 Bearer Token 到標頭
	req.Header.Set("Authorization", "Bearer "+"sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)")

	response, err := client.Do(req)
	if err != nil {
		utils.Log(s.Logger, "Error", utils.Trace()+" URL: "+c.Request.URL.Path, err.Error())
	}
	defer response.Body.Close()

	//c.JSON(http.StatusOK, mergedData)
}
func patchServer() error {

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
		return err
	}

	version := LoadVersion()

	err = tool.StuneDownload(stunConfig.GetAccessToken(), version, "edward")
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	command := exec.Command("rpm", "-Uvh", version)

	_, err = command.CombinedOutput()

	if err != nil {
		return err
	}

	return nil

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

	err := tool.StuneUpload(stunConfig.GetAccessToken(), "config/container.json", "edward")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": test})
}

func commitDockerfile(commit string) models.ServerInfo {
	var containerInfo models.ServerInfo

	containerInfo.CommitMessage = commit

	file, err := os.Open("/opt/pix/run/docker-compose-pro.yml")
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

		container := models.Image{
			UpdateTime:  formattedTime,
			ServiceName: fmt.Sprintf("%v", serviceName),
			ImageName:   imageName,
			ImageTag:    imageTag,
		}

		if strings.HasPrefix(fmt.Sprintf("%v", serviceName), "im") {
			container.ServiceName = "im"
			containerInfo.ContainerInfo = append(containerInfo.ContainerInfo, container)
			container.ServiceName = "im2"
			containerInfo.ContainerInfo = append(containerInfo.ContainerInfo, container)
			container.ServiceName = "im3"
			containerInfo.ContainerInfo = append(containerInfo.ContainerInfo, container)
		} else {
			containerInfo.ContainerInfo = append(containerInfo.ContainerInfo, container)
		}

		if exportDockerImage("/tmp/"+container.ServiceName+".tar", imageName+":"+imageTag) == nil {
			fmt.Println("Generate docker image :", imageName+":"+imageTag)
		}

	}

	// 將結構寫入 JSON 檔案
	err = writeJSONToFile(containerInfo, "config/container.json")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return containerInfo
	}
	return containerInfo
}

func exportDockerImage(path string, DockerImageName string) error {

	cmd := exec.Command("docker", "save", "-o", path, DockerImageName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("Docker image exported successfully to", path+"/"+DockerImageName+".tar")
	return nil

}
