package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Topology map[string][]string

type Request struct {
	Type      string   `json:"type"`
	MessageId int      `json:"msg_id"`
	Message   int      `json:"message"`
	Topology  Topology `json:"topology"`
}

type Response struct {
	Type     string `json:"type"`
	Messages []int `json:"messages,omitempty"`
}

type NodeStruct struct {
	node         *maelstrom.Node
	messageStore map[int]struct{}
	storeMutex   sync.RWMutex
	topology     *Topology
}

func main() {
	n := maelstrom.NewNode()

	nodeData := NodeStruct{
		node:         n,
		messageStore: make(map[int]struct{}),
	}

	n.Handle("broadcast", nodeData.handleBroadcast)
	n.Handle("read", nodeData.handleRead)
	n.Handle("topology", nodeData.handleTopology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (data *NodeStruct) handleBroadcast(msg maelstrom.Message) error {
	var request Request
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	data.storeMutex.Lock()
	defer data.storeMutex.Unlock()

	log.Print(request.Message)
	log.Println(data.messageStore)
	data.messageStore[request.Message] = struct{}{}

	log.Printf("Funny broadcast: %#v\n", data.messageStore)
	for _, node := range data.node.NodeIDs() {
		if node != data.node.ID() && node != msg.Src {
			go func() {
				if err := data.node.Send(node, request); err != nil {
					panic(err)
				}
			}()
		}
	}

	response := Response{
		Type: "broadcast_ok",
	}

	var err error
	if request.MessageId != 0 {
		err = data.node.Reply(msg, response)
	}

	return err
}

func (data *NodeStruct) handleRead(msg maelstrom.Message) error {
	var request Request
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	response := Response{
		Type:     "read_ok",
		Messages: data.getAllIDs(),
	}

	return data.node.Reply(msg, response)
}

func (data *NodeStruct) handleTopology(msg maelstrom.Message) error {
	var request Request
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		return err
	}

	data.topology = &request.Topology

	response := Response{
		Type: "topology_ok",
	}

	return data.node.Reply(msg, response)
}

func (data *NodeStruct) getAllIDs() []int {
	data.storeMutex.Lock()
	defer data.storeMutex.Unlock()

	copy := make([]int, len(data.messageStore))
	for message, _ := range data.messageStore {
		copy = append(copy, message)
	}
	return copy
}
