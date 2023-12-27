package utils

import (
	"fmt"
	"pix-console/common"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Utils struct {
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
