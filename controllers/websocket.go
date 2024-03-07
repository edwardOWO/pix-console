package controllers

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有连接
		return true
	},
}

func (u *Server) LogsHandler(c *gin.Context) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	path := c.Query("path")

	container_name := c.Query("container_name")

	cmd := exec.Command("")
	if path != "" {
		cmd = exec.Command("journalctl", "-u", path, "-f")

	} else if container_name != "" {
		cmd = exec.Command("docker", "logs", "--tail", "10", "-f", container_name)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// 处理 cmd.Stderr
	go func() {
		defer stderr.Close()
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				fmt.Println("Error reading from stderr:", err)
				return
			}
			if n > 0 {
				// 可以在这里处理错误输出，比如打印到控制台或记录到日志文件
			}
		}
	}()

	// 将日志输出发送到 WebSocket 连接
	go func() {
		defer stdout.Close()
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				fmt.Println("Error reading from command:", err)
				return
			}
			if n > 0 {
				if err := conn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					fmt.Println("Error writing to WebSocket connection:", err)
					return
				}
			}
		}
	}()

	// 等待 WebSocket 连接关闭
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket connection closed:", err)
			return
		}
	}
}
