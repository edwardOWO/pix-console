package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Users struct {
	Account []User `json:"userAcount"`
}
