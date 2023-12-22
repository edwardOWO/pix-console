package StuneTool

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type StuneSetting struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	BrandID      string
	Scope        string
}
type Payload struct {
	RemotePath   string  `json:"remote_path"`
	FileName     string  `json:"file_name"`
	ContentType  string  `json:"content_type"`
	FileSize     int64   `json:"file_size"`
	ExpireMinute *int    `json:"expire_minute,omitempty"`
	EndpointName *string `json:"endpoint_name,omitempty"`
}

func (y *StuneSetting) GetAccessToken() string {
	// Default
	credential := y.ClientID + ":" + y.ClientSecret

	// Specify
	if y.ClientID != "" && y.ClientSecret != "" {
		credential = y.ClientID + ":" + y.ClientSecret
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(credential))
	headers := map[string]string{
		"Authorization": "Basic " + encoded,
		"Brand-Id":      y.BrandID,
	}

	payload := map[string]string{"scope": y.Scope}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	req, err := http.NewRequest("POST", y.AuthURL+"/oauth2/getDeveloperToken", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return ""
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ""
	}

	if result["returnCode"].(float64) != 0 {
		return ""
	}

	token, ok := result["accessToken"].(string)
	if !ok {
		return ""
	}

	return token
}
func StuneUpload(token string, localFilePath string, remotePath string) error {
	url := "https://jstune.tw.juiker.net/api/Put.php"
	method := "POST"

	contentType := getContentType(localFilePath)
	fmt.Printf("Content Type: %s\n", contentType)

	fileSize := getFileSize(localFilePath)
	fmt.Printf("File Size: %d bytes\n", fileSize)

	fileName := getFileName(localFilePath)
	fmt.Printf("File Name: %s\n", fileName)

	payloadString := fmt.Sprintf(`{
		"file_name": "%s",
		"remote_path": "%s",
		"file_size": %d,
		"content_type": "%s"
	}`, fileName, remotePath, fileSize, contentType)

	payloadReader := strings.NewReader(payloadString)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)

	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// API 網址
	apiURL := string(body)

	// 讀取檔案
	file, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("無法開啟檔案:", err)
		return err
	}
	defer file.Close()

	// 讀取檔案內容
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("無法讀取檔案內容:", err)
		return err
	}

	// 發送 HTTP 請求
	client = &http.Client{}
	req, err = http.NewRequest("PUT", apiURL, bytes.NewBuffer(fileContents))
	if err != nil {
		fmt.Println("建立請求時發生錯誤:", err)
		return err
	}

	// 設定 Content-Type 標頭
	req.Header.Set("Content-Type", contentType)

	// 執行請求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("發送請求時發生錯誤:", err)
		return err
	}
	defer resp.Body.Close()

	// 讀取回應內容
	respBody, err := ioutil.ReadAll(resp.Body)

	fmt.Print(string(respBody))
	if err != nil {
		fmt.Println("讀取回應時發生錯誤:", err)
		return err
	}

	return err
}
func StuneDownload(token string, localPath string, remotePath string) error {

	url := "https://jstune.tw.juiker.net/api/Get.php"
	method := "POST"

	fileName := getFileName(localPath)
	fmt.Printf("File Name: %s\n", fileName)

	payloadReader := strings.NewReader(fmt.Sprintf(`{"file_name": "%s","remote_path": "%s"}`, fileName, remotePath))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadReader)
	if err != nil {
		fmt.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	resp, err := http.Get(string(body))
	if err != nil {
		fmt.Println("發送請求時發生錯誤:", err)
		return err
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼
	if resp.StatusCode != http.StatusOK {
		fmt.Println("下載失敗，狀態碼:", resp.Status)
		return err
	}

	// 創建目標檔案
	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println("無法建立檔案:", err)
		return err
	}
	defer file.Close()

	// 將回應內容寫入檔案
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("寫入檔案時發生錯誤:", err)
		return err
	}

	fmt.Print("檔案下載位置:", string(body))
	fmt.Println("檔案下載完成:", localPath)
	return nil
}

func getContentType(localFilePath string) string {
	// 使用 mime 包获取 MIME 类型
	contentType := mime.TypeByExtension(filepath.Ext(localFilePath))
	return contentType
}

func getFileSize(localFilePath string) int64 {
	// 使用 os 包获取文件大小
	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return -1
	}
	fileSize := fileInfo.Size()
	return fileSize
}

func getFileName(localFilePath string) string {
	// 使用 filepath 包获取不带路径的文件名
	fileName := filepath.Base(localFilePath)
	return fileName
}
