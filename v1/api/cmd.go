package v1

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
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	tool "pix-console/StuneTool"
)

// 設定 jwtkey 密鑰
var jwtKey = []byte("sdflkk;lkfds@@31JKKKhdfhkk123*!*&91`2387")

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
type User struct {
	Username string
	Password string
}

var users = []User{
	{Username: "admin", Password: "password"},
	{Username: "edward", Password: "password"},
	{Username: "user1", Password: "password"},
}
var e *casbin.Enforcer

var server string

func init() {

	// 指定檔案路徑
	filePath := "./host.ini"
	var err error
	// 讀取檔案內容
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("無法讀取檔案:", err)
		return
	}

	// 將字串轉換為Go語言中的字串型態
	server = string(content)

	// Initialize Casbin enforcer with the RBAC model and policy file

	e, err = casbin.NewEnforcer("./rbac/model.conf", "./rbac/policy.csv")
	if err != nil {
		panic(err)
	}

	e.AddGroupingPolicy("edward", "readonlyuser")
	if err != nil {
		fmt.Println("Failed to add user to group:", err)
		return
	}
}

func CasbinAuthMiddleware(c *gin.Context, username string, e *casbin.Enforcer) bool {

	// Get the username from the JWT token or your authentication mechanism
	// Get the request path
	obj := c.Request.URL.Path
	// Get the request method
	act := c.Request.Method

	// Check the permission
	e.Enforce()

	status, _ := e.Enforce(username, obj, act)

	fmt.Printf("%s,%s,%s,%t \n", username, obj, act, status)
	return status

}

// JWTAuthMiddleware 是基本身份驗證中介軟體
func JWTAuthMiddleware(c *gin.Context) {
	// Try cookie-based authentication first
	cookie, err := c.Request.Cookie("jwt")
	if err == nil {
		tokenString := cookie.Value
		claims := jwt.MapClaims{}

		_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err == nil {
			_, ok := claims["username"].(string)
			if ok {
				// Your existing session check logic can go here if needed

				if CasbinAuthMiddleware(c, claims["username"].(string), e) {
					c.Next()
					return
				}

			}
		} else {
			fmt.Print(err.Error())
		}
	}

	// If no cookie or cookie authentication fails, try Bearer token
	bearerToken := c.GetHeader("Authorization")

	if bearerToken == "Bearer "+"sdklkkfkj!2323dfj92083DKKD2**!*@#&&#!(#&1-9dfg,mzx//v)" {
		c.Next()
		return
	}

	if bearerToken != "" {
		tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err == nil {
			_, ok := claims["username"].(string)
			if ok {
				// Your existing session check logic can go here if needed
				if CasbinAuthMiddleware(c, claims["username"].(string), e) {
					c.Next()
					return
				}
			}
		} else {
			fmt.Print(err.Error())
		}
	}

	// If both cookie and Bearer token authentication fail, redirect or handle as needed
	c.Redirect(http.StatusSeeOther, "/")
	c.Abort()
}

func LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var user User
	var found bool

	for _, u := range users {
		if u.Username == username && u.Password == password {
			user = u
			found = true
			break
		}
	}

	if found {
		token := generateJWTToken(user.Username)
		c.SetCookie("jwt", token, 360000, "/", "localhost", false, true)
		c.SetCookie("jwt", token, 360000, "/", "60.199.173.12", false, true)
		c.Redirect(http.StatusSeeOther, "/index")
	} else {
		c.Redirect(http.StatusSeeOther, "/?error=InvalidCredentials")
	}
}

func LogoutHandler(c *gin.Context) {
	// 假設你的登入資訊保存在 cookie 中
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
		cookie.Domain = "localhost"
		cookie.Secure = false // 設為 true 如果你的應用啟用了 HTTPS
		cookie.HttpOnly = true
		c.SetCookie("jwt", "", -1, "/", "60.199.173.12", false, true)
	}

	// 重定向到登入頁面或其他目標頁面
	c.Redirect(http.StatusSeeOther, "/")
}

func generateJWTToken(username string) string {

	expirationTime := time.Now().Add(20 * 365 * 24 * time.Hour).Unix()
	//expirationTime=time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime, // Token expiration time
	})

	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
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
			Host:        server,
			ContainerID: "12345abcde67",
			Image:       "example:latest",
			Command:     "example-command",
			Created:     "2 weeks ago",
			Status:      "Up 2 weeks",
			Ports:       "8080/tcp",
			Names:       "example_container",
		},
		{
			Host:        server,
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

func ClusterDockerHandler(c *gin.Context) {
	addresses := []string{
		"http://192.168.70.111:8080/api/v1/docker",
		"http://192.168.70.112:8080/api/v1/docker",
		"http://192.168.70.113:8080/api/v1/docker",
		//"http://localhost:8080/api/v1/service",
	}

	var mergedData []map[string]interface{}
	for _, address := range addresses {
		data, _ := getDockerData(address)
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
		cmd := exec.Command("systemctl", "show", "--property=Names,ActiveState,SubState", "--value", "--no-pager", serviceName)

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
			Host:        server,
			ContainerID: "12345abcde67",
			Image:       "example:latest",
			Command:     "example-command",
			Created:     "2 weeks ago",
			Status:      "Up 2 weeks",
			Ports:       "8080/tcp",
			Names:       "example_container",
		},
		{
			Host:        server,
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

func ClusterServiceHandler(c *gin.Context) {

	addresses := []string{
		"http://192.168.70.111:8080/api/v1/service",
		"http://192.168.70.112:8080/api/v1/service",
		"http://192.168.70.113:8080/api/v1/service",
		//"http://localhost:8080/api/v1/service",
	}

	var mergedData []map[string]interface{}
	for _, address := range addresses {
		data, _ := getServiceData(address)
		mergedData = append(mergedData, data...)
	}

	c.JSON(http.StatusOK, mergedData)
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
				Host:        server,
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
		Host:        server,
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
