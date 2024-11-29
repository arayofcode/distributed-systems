package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Topology map[string][]string

type NodeStruct struct {
	node         *maelstrom.Node
	messageStore sync.Map
	topology     *Topology
}

func main() {
	n := maelstrom.NewNode()

	nodeData := NodeStruct{
		node: n,
	}

	n.Handle("broadcast", nodeData.handleBroadcast)
	n.Handle("read", nodeData.handleRead)
	n.Handle("topology", nodeData.handleTopology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func (data *NodeStruct) handleBroadcast(msg maelstrom.Message) error {
	request := make(map[string]any)
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		panic(err)
	}

	message := int(request["message"].(float64))

	if _, exists := data.messageStore.LoadOrStore(message, struct{}{}); exists {
		return nil
	}

	for _, node := range data.node.NodeIDs() {
		data.node.Send(node, request)
	}

	return data.node.Reply(msg, map[string]string{
		"type": "broadcast_ok",
	})
}

func (data *NodeStruct) handleRead(msg maelstrom.Message) error {
	messages := []int{}
	data.messageStore.Range(func(key, value any) bool {
		messages = append(messages, key.(int))
		return true
	})

	return data.node.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": messages,
	})
}

func (data *NodeStruct) handleTopology(msg maelstrom.Message) error {
	request := struct {
		Type     string   `json:"type"`
		Topology Topology `json:"topology"`
	}{}

	if err := json.Unmarshal(msg.Body, &request); err != nil {
		panic(err)
	}

	data.topology = &request.Topology

	return data.node.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}
