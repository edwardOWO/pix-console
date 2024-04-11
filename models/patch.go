package models

type Image struct {
	UpdateTime  string `json:"updatetime"`
	ServiceName string `json:"servicename"`
	ImageName   string `json:"imagename"`
	ImageTag    string `json:"imagetag"`
}
type ServerInfo struct {
	CommitMessage string  `json:"commitmessage"`
	ContainerInfo []Image `json:"containerinfo"`
}
type ServiceConfig struct {
	Services map[string]interface{} `yaml:"services"`
}

type PatchInfo struct {
	UpdateTime        string `json:"updatetime"`
	FileName          string `json:"patchfile"`
	PatchFilePath     string `json:"Patchfilepath"`
	RPMversion        string `json:"rpmVersion"`
	RPMpath           string `json:"rpmPath"`
	ContainerListPath string `json:"containerListPath"`
	PixConsoleUsed    bool   `json:"pixConsoleUsed"`
	PixComposeUsed    bool   `json:"pixComposeUsed"`
	PixConsoleSelect  bool   `json:"pixConsoleSelect"`
	PixComposeSelect  bool   `json:"pixComposeSelect"`
}
