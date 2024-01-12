package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"pix-console/common"
	"sync"
	"syscall"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Utils struct {
	mu                  sync.Mutex
	stopListening       chan bool
	running             bool
	ConnCount           map[int]int
	UdpPackets          map[int]int
	Status              bool
	Running             bool
	interrupt           chan os.Signal
	serverWG            sync.WaitGroup
	serverClosed        bool
	packetCaptureClosed bool
	ctx                 context.Context
	cancel              context.CancelFunc

	closeSignal chan struct{}
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
	rule, err := casbin.NewEnforcer("/opt/pix-console/rbac/model.conf", "/opt/pix-console/rbac/policy.csv")
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

	ports := common.Config.MonitorPort

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
		// 假設你想返回一個簡單的HTML頁面
		htmlContent := `
		<html>
			<head>
				<title>Simple Web Page</title>
			</head>
		</html>
	`

		// Write the HTTP response
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html\r\n"+
			"Connection: close\r\n"+
			"\r\n"+
			"%s", fmt.Sprintf(htmlContent, conn.RemoteAddr()))

		_, err := io.WriteString(conn, response)
		if err != nil {
			fmt.Println("Error writing response:", err)
		}
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

// CaptureUDPPackets 在指定的端口范围内捕获 UDP 数据包
func (u *Utils) CaptureUDPPackets(device string, startPort, endPort int, timeout time.Duration) (map[int]int, error) {

	if u.packetCaptureClosed == true {

		u.mu.Lock()
		defer u.mu.Unlock()

		// 创建地图的副本进行返回
		test := make(map[int]int)
		for key, value := range u.UdpPackets {
			test[key] = value
		}

		return test, nil
	}

	u.packetCaptureClosed = true
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	//filter := fmt.Sprintf("udp portrange %d-%d or tcp port 5222 or tcp port 5269 or tcp port 443 or tcp port 7891 or tcp port 7891", startPort, endPort)

	filter := fmt.Sprintf("udp portrange %d-%d or tcp", startPort, endPort)

	err = handle.SetBPFFilter(filter)
	if err != nil {
		return nil, err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// 创建一个信号通道，以便在收到中断信号时优雅地关闭程序
	u.interrupt = make(chan os.Signal, 1)
	signal.Notify(u.interrupt, os.Interrupt, syscall.SIGTERM)

	// 创建map来存储UDP端口号和相应的数据包信息

	u.UdpPackets = make(map[int]int)

	// := time.NewTimer(timeout)

	fmt.Print("Start to record")
	for {
		select {
		case packet := <-packetSource.Packets():
			// 在这里处理数据包，可以输出或进一步处理
			transportLayer := packet.TransportLayer()
			var port int
			if transportLayer != nil {

				packet := transportLayer.LayerType()

				switch packet {

				case layers.LayerTypeUDP:
					udp, ok := transportLayer.(*layers.UDP)
					if !ok {
						fmt.Println("Failed to assert UDP layer")
						continue
					}

					port = int(udp.DstPort)

				case layers.LayerTypeTCP:
					tcp, ok := transportLayer.(*layers.TCP)
					if !ok {
						fmt.Println("Failed to assert TCP layer")
						continue
					}

					port = int(tcp.DstPort)

				}

			}

			if port > 0 && port < 60000 {
				u.mu.Lock()
				u.UdpPackets[port] += 1
				u.mu.Unlock()
			}

		case <-u.interrupt:
			fmt.Println("Received interrupt, stopping...")
			return u.UdpPackets, nil
			//case <-timeoutTimer.C:
			//	fmt.Println("Timeout reached, stopping...")
			//	return udpPackets, nil
		}
	}
}
func (u *Utils) CloseUDPPackets() {
	close(u.interrupt)
	u.packetCaptureClosed = false
	return
}
func (u *Utils) GetCaptureResult() map[int]int {

	if u.packetCaptureClosed == true {

		u.mu.Lock()
		defer u.mu.Unlock()
		test := make(map[int]int)
		for key, value := range u.UdpPackets {
			test[key] = value
		}
		return test
	}
	return nil
}
