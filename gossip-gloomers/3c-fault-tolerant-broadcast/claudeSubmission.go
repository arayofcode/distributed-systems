package main

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"sync"
// 	"time"

// 	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
// )

// type Topology map[string][]string

// type NodeStruct struct {
// 	Node         *maelstrom.Node
// 	MessageStore sync.Map
// 	Topology     *Topology
// 	lock         sync.Mutex
// }

// type Request struct {
// 	Type     string    `json:"type"`
// 	Message  *int      `json:"message,omitempty"`
// 	Topology *Topology `json:"topology,omitempty"`
// }

// type Response struct {
// 	Type     string `json:"type"`
// 	Messages *[]int `json:"messages,omitempty"`
// }

// func main() {
// 	n := maelstrom.NewNode()

// 	nodeData := &NodeStruct{
// 		Node: n,
// 	}

// 	n.Handle("broadcast", nodeData.handleBroadcast)
// 	n.Handle("read", nodeData.handleRead)
// 	n.Handle("topology", nodeData.handleTopology)

// 	if err := n.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// func (data *NodeStruct) handleBroadcast(msg maelstrom.Message) error {
// 	var request Request
// 	if err := json.Unmarshal(msg.Body, &request); err != nil {
// 		return err
// 	}

// 	if _, exists := data.MessageStore.LoadOrStore(*request.Message, struct{}{}); exists || request.Message == nil {
// 		return nil
// 	}

// 	nodes := data.Node.NodeIDs()
// 	if data.Topology != nil {
// 		nodes = (*data.Topology)[data.Node.ID()]
// 	}

// 	// Create a wait group to track broadcast attempts
// 	var wg sync.WaitGroup
// 	unackedNodes := make(map[string]bool)
// 	data.lock.Lock()
// 	for _, node := range nodes {
// 		if node == msg.Src || node == data.Node.ID() {
// 			continue
// 		}
// 		unackedNodes[node] = true
// 	}
// 	data.lock.Unlock()

// 	// Retry mechanism
// 	go func() {
// 		retryInterval := time.Second
// 		maxRetries := 5

// 		for len(unackedNodes) > 0 && maxRetries > 0 {
// 			for node := range unackedNodes {
// 				wg.Add(1)
// 				go func(nodeID string) {
// 					defer wg.Done()
// 					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 					defer cancel()

// 					// Attempt to broadcast to this specific node
// 					rpcMsg, err := data.Node.SyncRPC(ctx, nodeID, msg.Body)
// 					if err != nil {
// 						log.Printf("Broadcast failed to %s: %v", nodeID, err)
// 						return
// 					}

// 					var response Response
// 					if err := json.Unmarshal(rpcMsg.Body, &response); err != nil {
// 						log.Printf("Failed to parse response from %s: %v", nodeID, err)
// 						return
// 					}

// 					// If successful, remove from unacked nodes
// 					if response.Type == "broadcast_ok" {
// 						data.lock.Lock()
// 						delete(unackedNodes, nodeID)
// 						data.lock.Unlock()
// 					}
// 				}(node)
// 			}

// 			// Wait for this round of broadcasts
// 			wg.Wait()

// 			// Exponential backoff
// 			if len(unackedNodes) > 0 {
// 				time.Sleep(retryInterval)
// 				retryInterval *= 2
// 				maxRetries--
// 			}
// 		}

// 		// Log any persistently unacked nodes
// 		if len(unackedNodes) > 0 {
// 			log.Printf("Some nodes could not be reached: %v", unackedNodes)
// 		}
// 	}()

// 	// Immediately respond to the original request
// 	return data.Node.Reply(msg, Response{
// 		Type: "broadcast_ok",
// 	})
// }

// func (data *NodeStruct) handleRead(msg maelstrom.Message) error {
// 	var response = Response{
// 		Type:     "read_ok",
// 		Messages: new([]int),
// 	}
// 	data.MessageStore.Range(func(key, value any) bool {
// 		*response.Messages = append(*response.Messages, key.(int))
// 		return true
// 	})

// 	return data.Node.Reply(msg, response)
// }

// func (data *NodeStruct) handleTopology(msg maelstrom.Message) error {
// 	var request Request
// 	if err := json.Unmarshal(msg.Body, &request); err != nil {
// 		return err
// 	}

// 	data.lock.Lock()
// 	data.Topology = request.Topology
// 	data.lock.Unlock()

// 	return data.Node.Reply(msg, Response{
// 		Type: "topology_ok",
// 	})
// }
