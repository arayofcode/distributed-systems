package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Topology map[string][]string

type NodeStruct struct {
	Node         *maelstrom.Node
	MessageStore sync.Map
	Topology     *Topology
}

type Request struct {
	Type     string    `json:"type"`
	Message  *int      `json:"message,omitempty"`
	Topology *Topology `json:"topology,omitempty"`
}

type Response struct {
	Type     string `json:"type"`
	Messages *[]int `json:"messages,omitempty"`
}

func main() {
	n := maelstrom.NewNode()

	nodeData := NodeStruct{
		Node: n,
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
		panic(err)
	}

	if _, exists := data.MessageStore.LoadOrStore(*request.Message, struct{}{}); exists || request.Message == nil {
		return nil
	}

	if data.Topology != nil {
		for _, node := range (*data.Topology)[data.Node.ID()] {
			go data.Node.Send(node, msg.Body)
		}
	} else {
		for _, node := range data.Node.NodeIDs() {
			go data.Node.Send(node, msg.Body)
		}
	}

	return data.Node.Reply(msg, Response{
		Type: "broadcast_ok",
	})
}

func (data *NodeStruct) handleRead(msg maelstrom.Message) error {
	var response = Response{
		Type:     "read_ok",
		Messages: new([]int),
	}
	data.MessageStore.Range(func(key, value any) bool {
		*response.Messages = append(*response.Messages, key.(int))
		return true
	})

	return data.Node.Reply(msg, response)
}

func (data *NodeStruct) handleTopology(msg maelstrom.Message) error {
	var request Request

	if err := json.Unmarshal(msg.Body, &request); err != nil {
		panic(err)
	}

	data.Topology = request.Topology

	return data.Node.Reply(msg, Response{
		Type: "topology_ok",
	})
}
