package controllers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/memberlist"

	tool "pix-console/StuneTool"
	"pix-console/common"
	"pix-console/models"
	"pix-console/utils"
)

type Server struct {
	utils      utils.Utils
	Memberlist *memberlist.Memberlist
	UserAcount *models.Users
}

type RequestData struct {
	Path string `json:"path"`
}

type ResponseData struct {
	Status int         `json:"message"`
	ErrMsg string      `json:"errmsg"`
	Data   *[]FileInfo `json:"data"`
}

type DirectoryListing struct {
	Files []string `json:"files"`
}
type FileInfo struct {
	Name    string `json:"name"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

type MemoryUsage struct {
	Free  string `json:"free"`
	Total string `json:"total"`
	Used  string `json:"used"`
}

func CreateFileHandler(c *gin.Context) {
	// 在GET請求時創建一個名為 "example.txt" 的文件
	err := CreateFile("example.txt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "無法創建檔案",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "文件已建立",
	})
}

func CreateFile(filename string) error {
	// 嘗試創建文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 寫入測試內容
	_, err = file.WriteString("這是一個測試範例\n")
	if err != nil {
		return err
	}

	return nil
}

func LsDirectory(directoryPath string) (string, error) {

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("Error:", err)
	}

	var fileInfoList []FileInfo

	test := ResponseData{}

	for _, file := range files {
		fileInfo := FileInfo{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Size:  file.Size(),
		}
		fileInfoList = append(fileInfoList, fileInfo)
	}

	test.Data = &fileInfoList
	test.Status = 0
	test.ErrMsg = "i not have k"

	jsonData, err := json.Marshal(test)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return string(jsonData), nil
}

// @Summary 取得檔案目錄結構
// @Produce json
// @Security BasicAuth
// @Param request body RequestData true "JSON請求數據" example={"path": "成功"}
// @Success 200 {object} ResponseData "成功"
// @Failure 400 {object} string "請求錯誤"
// @Failure 500 {object} string "內部錯誤"
// @Router /api/v1/checkfile [post]
func CheckFileHandler(c *gin.Context) {

	var requestData RequestData

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	directoryPath := requestData.Path // 替換為你要列出內容的目錄路徑
	result, err := LsDirectory(directoryPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	/*c.JSON(http.StatusOK, gin.H{
		"message": result,
	})*/

	c.JSON(http.StatusOK, result)
}

// @Summary 取得記憶體使用量
// @Produce json
// @Security BasicAuth
// @Success 200 {object} MemoryUsage "成功"
// @Failure 400 {object} string "請求錯誤"
// @Failure 500 {object} string "內部錯誤"
// @Router /api/v1/checkmemory [get]
func CheckMemoryHandler(c *gin.Context) {
	cmd := exec.Command("free", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 使用 awk 命令將輸出轉換為 JSON 格式
	awkCmd := exec.Command("awk", "/Mem:/{print \"{\\\"total\\\": \\\"\" $2 \"\\\", \\\"used\\\": \\\"\" $3 \"\\\", \\\"free\\\": \\\"\" $4 \"\\\"}\"}")
	awkCmd.Stdin = strings.NewReader(string(output))
	awkOutput, awkErr := awkCmd.CombinedOutput()
	if awkErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": awkErr.Error(),
		})
		return
	}

	// 解析 JSON 數據並返回
	//var result map[string]interface{}

	var result MemoryUsage
	if jsonErr := json.Unmarshal(awkOutput, &result); jsonErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": jsonErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary 啟動服務
// @Produce json
// @Security BasicAuth
// @Param request body RequestData true "JSON請求數據" example={"path": "成功"}
// @Success 200 {object} MemoryUsage "成功"
// @Failure 400 {object} string "請求錯誤"
// @Failure 500 {object} string "內部錯誤"
// @Router /api/v1/startservice [POST]
func StartServiceHandler(c *gin.Context) {

	var requestData RequestData
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 指定要執行的命令和參數
	command := exec.Command("systemctl start", string(requestData.Path))

	// 設置標準輸出和標準錯誤輸出
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// 執行命令
	err := command.Run()
	if err != nil {
		fmt.Printf("命令執行出錯: %s\n", err)
		return
	}

	fmt.Println("文件已成功創建: 1234")

	c.JSON(http.StatusOK, gin.H{
		"message": "文件已建利",
	})
}

// DownloadConfigHandler 下載 Config 檔案
// @Summary 下載 Config 檔案
// @Security BasicAuth
// @Success 200 {object} string "成功"
// @Failure 400 {object} string "請求錯誤"
// @Failure 404 {object} string "檔案未找到"
// @Router /api/v1/download [get]
func DownloadConfigHandler(c *gin.Context) {
	fileName := "/opt/pix/run/docker-compose-pro.yml"
	// 參數為空
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename parameter is required"})
		return
	}

	// 檢查檔案是否存在
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// 設置下載響應
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", ".env"))
	c.Header("Content-Type", "application/octet-stream")
	c.File(fileName)
}

// UploadConfigHandler 用於處理文件上傳。
// @Summary 上傳文件
// @Description 上傳文件到指定目錄
// @Accept multipart/form-data
// @Param file formData file true "上傳的文件"
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/upload [post]
func UploadConfigHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 上傳 docker-compose 檔案
	// err = c.SaveUploadedFile(file, "/tmp/pix/.env"+file.Filename)
	err = c.SaveUploadedFile(file, "/tmp/test.yml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "檔案上傳成功"})
}

// UploadDockerComposeHandler 用於處理文件上傳。
// @Summary 上傳文件
// @Description 上傳文件到指定目錄
// @Accept multipart/form-data
// @Param file formData file true "上傳的文件"
// @Produce json
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/upload [post]
func UploadDockerComposeHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 上傳 docker-compose 檔案
	// err = c.SaveUploadedFile(file, "/tmp/pix/.env"+file.Filename)
	err = c.SaveUploadedFile(file, "/opt/pix/run/docker-compose-pro.yml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "檔案上傳成功"})
}

func StartPixComposeHandler(c *gin.Context) {

	// 執行 "systemctl start pix-compose" 命令
	cmd := exec.Command("systemctl", "start", "pix-compose")

	// 執行並等待命令完成
	output, err := cmd.CombinedOutput()

	if err != nil {
		errorMessage := fmt.Errorf("執行命令失敗：%w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage.Error()})
		return
	}
	// 可選：將命令的輸出作為日誌或其他處理
	fmt.Println("命令輸出：", string(output))

	// 成功執行
	c.JSON(http.StatusOK, gin.H{"message": "成功啟動 pix-compose"})
}
func StopPixComposeHandler(c *gin.Context) {

	// 執行 "systemctl start pix-compose" 命令
	cmd := exec.Command("systemctl", "stop", "pix-compose")

	// 執行並等待命令完成
	output, err := cmd.CombinedOutput()

	if err != nil {
		errorMessage := fmt.Errorf("執行命令失敗：%w", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage.Error()})
		return
	}
	// 可選：將命令的輸出作為日誌或其他處理
	fmt.Println("命令輸出：", string(output))

	// 成功執行
	c.JSON(http.StatusOK, gin.H{"message": "成功關閉 pix-compose"})
}

// 使用者名稱和密碼（範例）

// JWT 驗證
func (u *Server) JWTAuthMiddleware(c *gin.Context) {
	// Try cookie-based authentication first
	cookie, err := c.Request.Cookie("jwt")

	bearerToken := ""

	// 先從 cookie 檢查如果沒有從 header 檢查
	if err == nil {
		bearerToken = cookie.Value

	} else {
		bearerToken = c.GetHeader("Authorization")
		bearerToken = strings.TrimPrefix(bearerToken, "Bearer ")

	}

	// 開始驗證
	jwtClaims, status := u.utils.AuthJWTToken(bearerToken)
	if status == true {
		if u.utils.CasbinAuthMiddleware(c, jwtClaims["username"].(string)) {
			c.Next()
		}
	}

	// 全通的 Tokken 只要是這組就直接放行
	if bearerToken == "sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)" {
		c.Next()
	}

	c.Redirect(http.StatusSeeOther, "/")
	c.Abort()
}

func (u *Server) LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var found bool

	for _, u := range u.UserAcount.Account {

		if u.Username == username && u.Password == password {

			found = true
			break
		}
	}

	if found {
		token := u.utils.GenerateJWTToken(username)
		c.SetCookie("jwt", token, 360000, "/", "localhost", false, true)
		c.SetCookie("jwt", token, 360000, "/", u.Memberlist.LocalNode().Addr.String(), false, true)
		c.SetCookie("jwt", token, 360000, "/", "60.199.173.12", false, true)
		c.Redirect(http.StatusSeeOther, "/index")
	} else {
		c.Redirect(http.StatusSeeOther, "/?error=InvalidCredentials")
	}
}

func (u *Server) LogoutHandler(c *gin.Context) {

	cookie, err := c.Request.Cookie("jwt")
	if err == nil {
		// 將 cookie 過期時間設為過去的時間，即立即過期
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		cookie.Path = "/"
		cookie.Domain = "localhost"
		cookie.Secure = false // 設為 true 如果你的應用啟用了 HTTPS
		cookie.HttpOnly = true
		c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	}
	cookie, err = c.Request.Cookie("jwt")
	if err == nil {
		// 將 cookie 過期時間設為過去的時間，即立即過期
		cookie.Expires = time.Now().AddDate(0, 0, -1)
		cookie.Path = "/"
		cookie.Domain = u.Memberlist.LocalNode().Addr.String()
		cookie.Secure = false // 設為 true 如果你的應用啟用了 HTTPS
		cookie.HttpOnly = true
		c.SetCookie("jwt", "", -1, "/", u.Memberlist.LocalNode().Addr.String(), false, true)
	}

	// 重定向到登入頁面或其他目標頁面
	c.Redirect(http.StatusSeeOther, "/")
}

type ContainerInfo struct {
	Host        string `json:"HOST"`
	ContainerID string `json:"CONTAINER ID"`
	Image       string `json:"IMAGE"`
	Command     string `json:"COMMAND"`
	Created     string `json:"CREATED"`
	Status      string `json:"STATUS"`
	Ports       string `json:"PORTS"`
	Names       string `json:"NAMES"`
}

func DockerHandler(c *gin.Context) {
	cmd := exec.Command("docker", "ps", "-a", "--format", `{{.ID}}#{{.Image}}#{{.Command}}#{{.RunningFor}}#{{.Status}}#{{json .Ports}}#{{.Names}}`)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing 'docker ps': %s\n", stderr.String())
		//return
	}

	// 解析 Docker ps 的結果並填入結構
	var dockerData []ContainerInfo
	parseDockerPSResult(stdout.String(), &dockerData)

	// 手動硬編碼的資料

	hardcodedData := []ContainerInfo{
		{
			Host:        common.Config.ServerName,
			ContainerID: "12345abcde67",
			Image:       "example:latest",
			Command:     "example-command",
			Created:     "2 weeks ago",
			Status:      "Up 2 weeks",
			Ports:       "8080/tcp",
			Names:       "example_container",
		},
		{
			Host:        common.Config.ServerName,
			ContainerID: "12345abcde6711",
			Image:       "example:latest",
			Command:     "example-command",
			Created:     "1 weeks ago",
			Status:      "Up 2 weeks",
			Ports:       "8080/tcp",
			Names:       "example_container",
		},
		// 其他手動編碼的資料...
	}

	combinedData := append(hardcodedData, dockerData...)

	// 設定回應標頭
	c.Header("Content-Type", "application/json")
	// 回傳合併後的結果
	c.JSON(http.StatusOK, combinedData)
}

func (u *Server) ClusterDockerHandler(c *gin.Context) {

	node := u.Memberlist.Members()
	var mergedData []map[string]interface{}
	for _, address := range node {

		apiUrl := fmt.Sprintf("http://%s%s/api/v1/docker", address.Addr, common.Config.Port)
		data, _ := getDockerData(apiUrl)
		mergedData = append(mergedData, data...)
	}

	c.JSON(http.StatusOK, mergedData)
}
func getDockerData(url string) ([]map[string]interface{}, error) {

	client := &http.Client{
		Timeout: time.Second * 1, // 設定超時時間為 5 秒
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加 Bearer Token 到標頭
	req.Header.Set("Authorization", "Bearer "+"sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func ServiceHandler(c *gin.Context) {

	services := []string{"mysqld", "mongod", "cassandra", "pix-compose", "pixd", "pix-onlyoffice", "crond", "rsyslog", "sshd"}
	var dockerData []ContainerInfo
	for _, serviceName := range services {
		///cmd := exec.Command("systemctl", "show", "--property=Names,ActiveState,SubState", "--value", "--no-pager", serviceName)
		cmd := exec.Command("systemctl", "show", "--property=Names,ActiveState,SubState", "--no-pager", serviceName)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error executing 'systemctl': %s\n", stderr.String())
			continue
		}

		// 解析 Docker ps 的結果並填入結構

		parseServicePSResult(stdout.String(), &dockerData)
	}

	// 合併兩個資料集
	var combinedData []ContainerInfo

	hardcodedData := []ContainerInfo{
		{
			Host:        common.Config.ServerName,
			ContainerID: "12345abcde67",
			Image:       "example:latest",
			Command:     "example-command",
			Created:     "2 weeks ago",
			Status:      "Up 2 weeks",
			Ports:       "8080/tcp",
			Names:       "example_container",
		},
		{
			Host:        common.Config.ServerName,
			ContainerID: "12345abcde6711",
			Image:       "example:latest",
			Command:     "example-command",
			Created:     "1 weeks ago",
			Status:      "Up 2 weeks",
			Ports:       "8080/tcp",
			Names:       "example1234_container",
		},
		// 其他手動編碼的資料...
	}

	combinedData = append(combinedData, dockerData...)
	combinedData = append(combinedData, hardcodedData...)
	// 設定回應標頭
	c.Header("Content-Type", "application/json")
	// 回傳合併後的結果
	c.JSON(http.StatusOK, combinedData)
}

func (u *Server) ClusterServiceHandler(c *gin.Context) {
	node := u.Memberlist.Members()
	var mergedData []map[string]interface{}
	for _, address := range node {

		apiUrl := fmt.Sprintf("http://%s%s/api/v1/service", address.Addr, common.Config.Port)
		data, _ := getServiceData(apiUrl)
		mergedData = append(mergedData, data...)
	}
	c.JSON(http.StatusOK, mergedData)

}
func (u *Server) ServerlistHandler(c *gin.Context) {
	memberlistStatus := getMemberlistStatus(u.Memberlist)
	c.JSON(http.StatusOK, memberlistStatus)
}

func (u *Server) MoniotrListenPort(c *gin.Context) {

	status := c.Query("status")
	setting, err := strconv.ParseBool(status)

	if err == nil {
		if setting == true {
			u.utils.StartServer()
		} else {
			u.utils.CloseServer()
		}
	}
	c.JSON(http.StatusOK, u.utils.ConnCount)

}
func (u *Server) MonitorHandler(c *gin.Context) {

	status := c.Query("status")
	setting, err := strconv.ParseBool(status)

	if err == nil {
		if setting == true {
			portRangeStart := 0
			portRangeEnd := 60000
			device := "enp0s3"
			captureResult, _ := u.utils.CaptureUDPPackets(device, portRangeStart, portRangeEnd, 50000000000000)
			c.JSON(http.StatusOK, captureResult)
		} else {
			u.utils.CloseUDPPackets()

		}
	}
	c.JSON(http.StatusOK, "ok")
}

func (u *Server) GetMonitorHandler(c *gin.Context) {

	c.JSON(http.StatusOK, u.utils.GetCaptureResult())
}

// KeyValuePair 用于存储键值对
type KeyValuePair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func getServiceData(url string) ([]map[string]interface{}, error) {

	client := &http.Client{
		Timeout: time.Second * 1, // 設定超時時間為 5 秒
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 添加 Bearer Token 到標頭
	req.Header.Set("Authorization", "Bearer "+"sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// 解析 Docker ps 的結果並填入結構
func parseDockerPSResult(output string, data *[]ContainerInfo) {
	lines := strings.Split(output, "\n")

	// 叠代每一行
	for _, line := range lines {
		// 忽略空行
		if line == "" {
			continue
		}

		// 去除頭尾的單引號
		line = strings.Trim(line, "'")

		// 使用逗號分割字段
		fields := strings.Split(line, "#")

		// 檢查字段數量是否足夠
		if len(fields) > 6 {
			containerInfo := ContainerInfo{
				Host:        common.Config.ServerName,
				ContainerID: fields[0],
				Image:       fields[1],
				Command:     fields[2],
				Created:     fields[3],
				Status:      fields[4],
				Ports:       fields[5],
				Names:       fields[6],
			}

			// 添加到數據切片
			*data = append(*data, containerInfo)
		}
	}
}

func parseServicePSResult(output string, data *[]ContainerInfo) {

	properties := strings.Split(output, "\n")

	containerInfo := ContainerInfo{
		Host:        common.Config.ServerName,
		ContainerID: "test",
		Image:       "mongod",
		Command:     "test",
		Created:     properties[2],
		Status:      properties[1],
		Ports:       "27017",
		Names:       properties[0],
	}

	// 添加到數據切片
	*data = append(*data, containerInfo)

}

func DockerComposeHandler(c *gin.Context) {
	// 讀取文件內容
	filePath := "/opt/pix/run/docker-compose-pro.yml"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 如果讀取文件出現錯誤，返回錯誤信息
		c.String(http.StatusInternalServerError, "Error reading file: %s", err.Error())
		return
	}

	// 直接回傳文件內容
	c.Data(http.StatusOK, "application/yml", content)
}

func ClusterDownloadFromStune(c *gin.Context) {

	Service := c.Query("service")
	StartTime := c.Query("startTime")
	EndTime := c.Query("endTime")

	addresses := []string{
		"http://192.168.70.111:8080/api/v1/downloadFromStune?service=" + Service + "&startTime=" + StartTime + "&endTime=" + EndTime + "&time=1",
		"http://192.168.70.112:8080/api/v1/downloadFromStune?service=" + Service + "&startTime=" + StartTime + "&endTime=" + EndTime + "&time=1",
		"http://192.168.70.113:8080/api/v1/downloadFromStune?service=" + Service + "&startTime=" + StartTime + "&endTime=" + EndTime + "&time=1",
	}

	zipFilename := "/tmp/test.zip"
	zipFile, err := os.Create(zipFilename)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error creating zip file: %s", err))
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)

	var test = 1
	for _, address := range addresses {
		getLog(address, "/tmp/"+strconv.Itoa(test)+".zip")

		// Add file to zip
		err := addFileToZip(zipWriter, "/tmp/"+strconv.Itoa(test)+".zip")
		if err != nil {
			continue
		}
		test += 1
	}
	c.Header("Content-Disposition", "attachment; filename=test.zip")
	c.File("/tmp/test.zip")
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个新文件头
	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// 修改头部以使用指定的名称
	header.Name = filepath.Base(filename)

	// 创建一个新的 zip 文件条目
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// 将文件数据写入 zip 条目
	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}
func getLog(url string, filename string) error {
	client := &http.Client{
		Timeout: time.Second * 5, // 設定超時時間為 5 秒
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 添加 Bearer Token 到標頭
	req.Header.Set("Authorization", "Bearer "+"sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)")

	// 發送請求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 創建檔案以保存回應內容
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 將回應內容寫入檔案
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFromStune(c *gin.Context) {

	Service := c.Query("service")
	StartTime := c.Query("startTime")
	EndTime := c.Query("endTime")

	startTime, err := parseTime(StartTime)
	if err != nil {
		fmt.Println("無法解析起始日期:", err)
		return
	}

	endTime, err := parseTime(EndTime)
	if err != nil {
		fmt.Println("無法解析結束日期:", err)
		return
	}
	fileName := Service + "_" + StartTime + "_" + EndTime
	fileName = strings.ReplaceAll(fileName, "/", "_")
	fileName += ".zip"

	switch Service {
	case "IM":
		directories := []string{"/data/docker-data/volumes/run_im_log", "/data/docker-data/volumes/run_im2_log", "/data/docker-data/volumes/run_im3_log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}
	case "SIP":
		directories := []string{"/data/docker-data/volumes/run_sorrel_api_log", "/data/docker-data/volumes/run_sorrel_rose_log", "/data/docker-data/volumes/run_sorrel_sbcallinone_log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}
	case "DB":
		directories := []string{"/var/log/mongodb", "/var/log/cassandra", "/var/log/mysqld.log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}
	case "STUNE":
		directories := []string{"/data/docker-data/volumes/run_stune_log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}

	default:
		fmt.Println("Unknown service")
	}

	c.File(fileName)
}

func UploadToStune(c *gin.Context) {

	var requestPayload struct {
		// 定義結構體的字段，以匹配 JSON 中的屬性
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		Service   string `json:"service"`
	}

	// 使用 ShouldBindJSON 來將請求主體映射到結構體
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileName := requestPayload.Service + "_" + requestPayload.StartTime + "_" + requestPayload.EndTime

	startTime, err := parseTime(requestPayload.StartTime)
	if err != nil {
		fmt.Println("無法解析起始日期:", err)
		return
	}

	endTime, err := parseTime(requestPayload.EndTime)
	if err != nil {
		fmt.Println("無法解析結束日期:", err)
		return
	}

	fileName = strings.ReplaceAll(fileName, "/", "_")

	fileName += ".zip"

	switch requestPayload.Service {
	case "IM":
		directories := []string{"/data/docker-data/volumes/run_im_log", "/data/docker-data/volumes/run_im2_log", "/data/docker-data/volumes/run_im3_log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}
	case "SIP":
		directories := []string{"/data/docker-data/volumes/run_sorrel_api_log", "/data/docker-data/volumes/run_sorrel_rose_log", "/data/docker-data/volumes/run_sorrel_sbcallinone_log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}
	case "DB":
		directories := []string{"/var/log/mongodb", "/var/log/cassandra", "/var/log/mysqld.log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}
	case "STUNE":
		directories := []string{"/data/docker-data/volumes/run_stune_log"}
		err := CompressFiles(directories, fileName, startTime, endTime)
		if err != nil {
			fmt.Printf("Error compressing files: %v\n", err)
		}

	default:
		fmt.Println("Unknown service")
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "檔案產生失敗"})
		return
	}

	stunConfig := tool.StuneSetting{
		ClientID:     "pixCollector",
		ClientSecret: "(5vBX1Tu@DDPs0Om1Cfm",
		AuthURL:      "https://auth.tw.juiker.net",
		BrandID:      "juiker",
		Scope:        "tw:stune:basic",
	}

	err = tool.StuneUpload(stunConfig.GetAccessToken(), fileName, "edward")
	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusOK, gin.H{"message": "檔案同步失敗"})
	}

	/*
		err = tool.StuneDownload(stunConfig.GetAccessToken(), fileName, "edward")
		if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusOK, gin.H{"message": "檔案同步失敗"})
		}
	*/

	c.JSON(http.StatusOK, gin.H{"message": fileName + " " + "檔案同步成功"})
}

// parseTime 將格式為 "2006/01/02" 的字串轉換為 time.Time
func parseTime(dateString string) (time.Time, error) {
	return time.Parse("2006/01/02", dateString)
}

func CompressFiles(directories []string, zipFilePath string, startTime, endTime time.Time) error {
	// Create ZIP file
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("unable to create ZIP file: %v", err)
	}
	defer zipFile.Close()

	// Create ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through files in the directory
	for _, directoryPath := range directories {
		err = filepath.Walk(directoryPath, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Check if the file is a regular file and within the specified time range
			if !fileInfo.IsDir() && fileInfo.ModTime().After(startTime) && fileInfo.ModTime().Before(endTime) {
				// Get relative path
				relativePath, err := filepath.Rel(directoryPath, filePath)
				if err != nil {
					return err
				}

				// Open the file
				file, err := os.Open(filePath)
				if err != nil {
					return err
				}
				defer file.Close()

				// Create file header with original times
				fileHeader, err := zip.FileInfoHeader(fileInfo)
				if err != nil {
					return err
				}
				fileHeader.Name = relativePath

				// Set the modified and accessed times to match the original file
				fileHeader.Modified = fileInfo.ModTime()

				// Create file in ZIP with the original file header
				zipFile, err := zipWriter.CreateHeader(fileHeader)
				if err != nil {
					return err
				}

				// Copy file content to ZIP
				_, err = io.Copy(zipFile, file)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	return err
}
