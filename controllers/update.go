package controllers

import (
	"fmt"
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

// 升級伺服器
func UpdateDocker() bool {
	if !backupService() {
		return false
	}
	if !updateDockerCompose() {
		return false
	}
	if !pullImage() {
		return false
	}
	if !restartService() {
		return false
	}
	if !checkService() {
		return false
	}
	return true
}
