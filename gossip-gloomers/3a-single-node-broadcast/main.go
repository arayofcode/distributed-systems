package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Topology map[string][]string

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

	var messageStore []int
	var topology *Topology
	// Useful for implementation of topology. Leaving it here for now
	_ = topology

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var request Request
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		messageStore = append(messageStore, *request.Message)

		response := Response{
			Type: "broadcast_ok",
		}

		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var request Request
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		response := Response{
			Type:     "read_ok",
			Messages: &messageStore,
		}

		return n.Reply(msg, response)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var request Request
		if err := json.Unmarshal(msg.Body, &request); err != nil {
			return err
		}

		topology = request.Topology

		response := Response{
			Type: "topology_ok",
		}

		return n.Reply(msg, response)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
