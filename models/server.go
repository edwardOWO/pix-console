package models

import "encoding/json"

type ServerMetaData struct {
	Region  string `json:"region"`
	Zone    string `json:"zone"`
	ShardId uint16 `json:"shard-id"`
	Weight  uint64 `json:"weight"`
}

func (m ServerMetaData) Bytes() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}
