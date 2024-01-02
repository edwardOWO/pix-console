package utils

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"pix-console/common"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Utils struct {
	stopListening chan bool
	running       bool
	connCount     map[int]int
}
type SdtClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt_lib.StandardClaims
}

// RBAC 建構子
var (
	e *casbin.Enforcer
)

func init() {

	rule, err := casbin.NewEnforcer("./rbac/model.conf", "./rbac/policy.csv")
	e = rule
	if err != nil {
		panic(err)
	}

	e.AddGroupingPolicy("edward", "readonlyuser")
	if err != nil {
		fmt.Println("Failed to add user to group:", err)
		return
	}
}

// Create JWT Tokken
func (u *Utils) GenerateJWTToken(username string) string {

	expirationTime := time.Now().Add(20 * 365 * 24 * time.Hour).Unix()
	//expirationTime=time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime, // Token expiration time
	})

	tokenString, _ := token.SignedString([]byte(common.Config.JwtSecretPassword))
	return tokenString
}

// Create JWT Tokken
func (u *Utils) AuthJWTToken(tokenString string) (jwt.MapClaims, bool) {

	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(common.Config.JwtSecretPassword), nil
	})
	if err != nil {
		fmt.Print(err.Error())
		return nil, false
	}

	return claims, true
}

func (u *Utils) CasbinAuthMiddleware(c *gin.Context, username string) bool {

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

func (u *Utils) handleTCPConnection(conn net.Conn, p int) {
	defer conn.Close()

	// 处理TCP连接的代码
	conn.Write([]byte("TCP PASS!!!!"))
	u.connCount[p]++
}

func (u *Utils) ListenPortsAndExit(ports []int, start bool) (connMap map[int]int) {
	if u.running && start {
		fmt.Println("已經啟動過，不再重複啟動")
		return u.connCount
	}

	if !u.running && !start {
		fmt.Println("沒有啟動，無需停止")
		return u.connCount
	}

	// 如果是啟動，則初始化通道和標誌
	if start {
		u.connCount = make(map[int]int)
		u.stopListening = make(chan bool)
		u.running = true
	}

	// 如果是停止，則關閉通道並重置標誌
	if !start {
		close(u.stopListening)
		u.running = false
		return
	}

	// 啟動監聽多個端口的 goroutine
	for _, port := range ports {
		go func(p int) {
			// 監聽指定的端口
			var listener net.Listener
			var err error
			protocol := "tcp"

			if p < 40000 {
				protocol = "tcp"
			}

			address := fmt.Sprintf(":%d", p)
			listener, err = net.Listen(protocol, address)
			fmt.Printf("開始監聽端口 %d，協議：%s\n", p, protocol)

			if err != nil {
				fmt.Printf("無法監聽端口 %d: %s\n", p, err)
				u.stopListening <- true
				return
			}
			defer listener.Close()

			// 接收信號
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

			for {
				// 等待連線
				conn, err := listener.Accept()
				if err != nil {
					fmt.Printf("無法接受連線: %s\n", err)
					continue
				}

				// 啟動 goroutine 處理連線
				go u.handleTCPConnection(conn, p)
			}
		}(port)
	}

	// 等待所有 goroutine 完成
	for range ports {
		<-u.stopListening
	}

	if start {
		fmt.Println("所有監聽已退出")
	}

	return u.connCount
}
