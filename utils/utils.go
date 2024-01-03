package utils

import (
	"context"
	"fmt"
	"net"
	"pix-console/common"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Utils struct {
	stopListening chan bool
	running       bool
	ConnCount     map[int]int
	Status        bool
	Running       bool

	serverWG     sync.WaitGroup
	serverClosed bool
	ctx          context.Context
	cancel       context.CancelFunc
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
func (u *Utils) StartServer() {

	ports := []int{5222, 5269}

	if u.serverClosed == false {
		u.ConnCount = make(map[int]int)

		for _, port := range ports {
			u.ConnCount[port] = 0
		}
	} else {
		fmt.Print("Server Already start")
		return
	}

	u.serverClosed = true
	u.serverWG.Add(len(ports)) // 監聽兩個端口，您可以根據需求調整

	// 建立具有取消功能的上下文
	u.ctx, u.cancel = context.WithCancel(context.Background())

	for _, port := range ports {
		listenPort := fmt.Sprintf(":%d", port)
		go u.listenAndServe(listenPort, port)
	}

	// 等待 goroutine 完成
	u.serverWG.Wait()
}

func (u *Utils) listenAndServe(addr string, port int) {
	defer u.serverWG.Done()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("啟動伺服器錯誤 (%s): %v\n", addr, err)
		return
	}
	defer ln.Close()

	fmt.Printf("伺服器已啟動，正在監聽 %s\n", addr)

	for {
		select {
		case <-u.ctx.Done():
			fmt.Printf("伺服器已關閉 (%s) !!!!!!!!!!!!!!\n", addr)
			return
		default:
			ln.(*net.TCPListener).SetDeadline(time.Now().Add(100 * time.Millisecond))

			conn, err := ln.Accept()

			if err != nil {
				// 檢查錯誤是否是由於監聽器被關閉
				netErr, ok := err.(net.Error)
				if ok && netErr.Timeout() {
					continue // 這是非阻塞操作，繼續等待連線
				}

				fmt.Printf("接受連線錯誤 (%s): %v\n", addr, err)
				continue
			}

			go u.HandleConnection(conn, port)
		}
	}
}

func (u *Utils) HandleConnection(conn net.Conn, port int) {
	defer conn.Close()
	fmt.Printf("來自 %s 的連線已接受\n", conn.RemoteAddr())
	u.ConnCount[port]++
	// 檢查上下文的取消
	select {
	case <-u.ctx.Done():
		fmt.Println("伺服器關閉中，終止連線處理。")
		return
	default:
		// 繼續處理連線邏輯
		// 在這裡處理連接邏輯
	}
}

func (u *Utils) CloseServer() {

	if u.serverClosed == true {
		u.serverClosed = false
		u.cancel() // Cancel the context to signal the termination
		u.serverWG.Wait()
		fmt.Println("Server closed")
	}
}
