package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"pix-console/common"
	"pix-console/models"
	"syscall"

	"github.com/hashicorp/memberlist"
)

type MyDelegate struct {
	MsgCh      chan []byte
	Broadcasts *memberlist.TransmitLimitedQueue
	Meta       models.ServerMetaData
}

func (d *MyDelegate) NotifyMsg(msg []byte) {
	d.MsgCh <- msg
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	return d.Broadcasts.GetBroadcasts(overhead, limit)
}

func (d *MyDelegate) GetNodeMeta(overhead, limit int) uint64 {
	return d.Meta.Weight
}

func (d *MyDelegate) NodeMeta(limit int) []byte {
	return d.Meta.Bytes()
}
func (d *MyDelegate) LocalState(join bool) []byte {
	// not use, noop
	return []byte("")
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
	// not use
}

type MyEventDelegate struct {
	UpdateCh chan struct{}
	Num      int
}

func (d *MyEventDelegate) NotifyJoin(node *memberlist.Node) {
	d.notifyUpdate()
}
func (d *MyEventDelegate) NotifyLeave(node *memberlist.Node) {
	d.notifyUpdate()
}
func (d *MyEventDelegate) NotifyUpdate(node *memberlist.Node) {
	d.notifyUpdate()
}
func (d *MyEventDelegate) notifyUpdate() {
	select {
	case d.UpdateCh <- struct{}{}:
	default:
	}
}

type MyBroadcastMessage struct {
	Key   string `json:"key"`
	Value uint64 `json:"value"`
}

func (m MyBroadcastMessage) Invalidates(other memberlist.Broadcast) bool {
	return false
}
func (m MyBroadcastMessage) Finished() {
	// nop
}
func (m MyBroadcastMessage) Message() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}

func ParseMyBroadcastMessage(data []byte) (*MyBroadcastMessage, bool) {
	msg := new(MyBroadcastMessage)
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, false
	}
	return msg, true
}
func ParseMyMetaData(data []byte) (models.ServerMetaData, bool) {
	meta := models.ServerMetaData{}
	if err := json.Unmarshal(data, &meta); err != nil {
		return meta, false
	}
	return meta, true
}

func wait_signal(cancel context.CancelFunc) {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGINT)
	for {
		select {
		case s := <-signal_chan:
			log.Printf("signal %s happen", s.String())
			cancel()
		}
	}
}

func printMemberlistStatus(list *memberlist.Memberlist) {
	for _, node := range list.Members() {
		meta, ok := ParseMyMetaData(node.Meta)
		if ok {
			log.Printf(
				"%s region: %s, zone: %s, shard: %d, weight: %d",
				node.Name,
				meta.Region,
				meta.Zone,
				meta.ShardId,
				meta.Weight,
			)
		}
	}
}
func StartMemberlist() (*memberlist.Memberlist, *MyDelegate, *MyEventDelegate) {

	msgCh := make(chan []byte)

	d := &MyDelegate{
		Meta:       models.ServerMetaData{Region: "ap-northeast-1", Zone: "1a", ShardId: 100, Weight: 0},
		MsgCh:      msgCh,
		Broadcasts: new(memberlist.TransmitLimitedQueue),
	}
	d2 := MyEventDelegate{UpdateCh: make(chan struct{}, 1)}

	conf := memberlist.DefaultLocalConfig()
	conf.Name = common.Config.ServerName
	conf.BindPort = common.Config.MemberlistPort
	conf.AdvertisePort = conf.BindPort
	conf.Delegate = d
	conf.Events = &d2

	list, err := memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	//local := list.LocalNode()
	list.Join(common.Config.ServerHost)

	go func() {
		run := true
		for run {
			select {
			case data := <-d.MsgCh:
				printMemberlistStatus(list)

				fmt.Print(string(data))
				log.Printf("------------------")

			case _ = <-d2.UpdateCh:
				fmt.Print("cluster status update")
				continue

			}
		}
		log.Printf("bye.")
	}()

	return list, d, &d2
}
func getMemberlistStatus(list *memberlist.Memberlist) []map[string]interface{} {

	var hostData []map[string]interface{}

	for _, node := range list.Members() {
		meta, ok := ParseMyMetaData(node.Meta)
		if ok {
			memberInfo := map[string]interface{}{
				"HOST":      node.Addr,
				"CONTAINER": meta.Region,
				"IMAGE":     "example:latest",
				"COMMAND":   "example-command",
				"CREATED":   "1 weeks ago",
				"STATUS":    "UP",
				"PORTS":     "8080/tcp",
				"NAMES":     meta.Region,
			}
			hostData = append(hostData, memberInfo)
		}
	}

	return hostData
}
