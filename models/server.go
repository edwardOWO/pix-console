package models

import "encoding/json"

type ServerMetaData struct {
	Region         string `json:"region"`
	Zone           string `json:"zone"`
	ShardId        uint16 `json:"shard-id"`
	Weight         uint64 `json:"weight"`
	ServerStatus   uint64 `json:"server-status"`
	ServiceStatus  uint64 `json:"service-status"`
	ServerName     string `json:"servername"`
	ServiceVersion string `json:"serviceversion"`
}

func (m ServerMetaData) Bytes() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}
