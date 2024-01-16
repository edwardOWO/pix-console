package models

type Image struct {
	ServiceName string `json:"servicename"`
	ImageName   string `json:"imagename"`
	ImageTag    string `json:"imagetag"`
}
type ServerInfo struct {
	ContainerInfo []Image `json:"containerinfo"`
}
type ServiceConfig struct {
	Services map[string]interface{} `yaml:"services"`
}
