package view

import (
	"github.com/gin-gonic/gin"
)

type Link struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

type PageLinkStruct struct {
	Links []Link `json:"links"`
}

func CreatePageLink() gin.H {

	PageLink := gin.H{
		"links": []gin.H{
			{"text": "DashBoard", "href": "/dashboard"},
			{"text": "Host", "href": "/host"},
			{"text": "Service", "href": "/service"},
			{"text": "Containers", "href": "/docker"},
			{"text": "Setting", "href": "/index"},
			{"text": "Docker-compose", "href": "/docker-compose"},
			{"text": "Feedback", "href": "/feedback"},
			{"text": "Check", "href": "/check"},
			{"text": "Logout", "href": "/logout"},
		},
	}
	return PageLink
}
