package memberlist

import (
	"fmt"
	"log"

	"github.com/hashicorp/memberlist"
	"github.com/serialx/hashring"
)

type EventDelegate struct {
	consistent *hashring.HashRing
	node       []*memberlist.Node
	updateCh   chan struct{} // 新增一個通知事件的通道
}

func (d *EventDelegate) NotifyJoin(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	log.Printf("join %s", hostPort)
	if d.consistent == nil {
		d.consistent = hashring.New([]string{hostPort})
	} else {
		d.consistent = d.consistent.AddNode(hostPort)
	}
	d.notifyUpdate()
}

func (d *EventDelegate) NotifyLeave(node *memberlist.Node) {
	hostPort := fmt.Sprintf("%s:%d", node.Addr.To4().String(), node.Port)
	log.Printf("leave %s", hostPort)
	if d.consistent != nil {
		d.consistent = d.consistent.RemoveNode(hostPort)
	}
	d.notifyUpdate()
}

func (d *EventDelegate) NotifyUpdate(node *memberlist.Node) {
	d.notifyUpdate()
}

func (d *EventDelegate) notifyUpdate() {
	select {
	case d.updateCh <- struct{}{}:
	default:
	}
}
func (d *EventDelegate) GetMemblist() []*memberlist.Node {

	for _, member := range d.node {
		fmt.Print(member.String())
	}

	return d.node
}

func (d *EventDelegate) Start(conf *memberlist.Config, node []string) {

	conf.Events = &EventDelegate{updateCh: make(chan struct{}, 1)}

	list, err := memberlist.Create(conf)
	if err != nil {
		log.Fatal(err)
	}

	list.Join(node)

	run := true
	for run {
		select {
		case <-conf.Events.(*EventDelegate).updateCh:
			/*
				devt := conf.Events.(*EventDelegate)
				if devt == nil {
					log.Printf("consistent isnt initialized")
					continue
				}
				count := devt.consistent.Size()
				data, ok := devt.consistent.GetNodes("", count)
				if ok {
					d.node = data
				}*/

			d.node = list.Members()

		}

	}
	log.Printf("bye.")
	close(conf.Events.(*EventDelegate).updateCh)
}
