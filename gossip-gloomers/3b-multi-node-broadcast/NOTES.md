# Notes and Challenges

- When topology is missing, we SEND the message to all nodes directly.
- If in a topology, if A -> B -> C, how do we send a message such that each node receives only once (avoid cycle)?
  - For each node, check all neighbours using read. If neighbour has already received the message, skip otherwise send.
  - Create a new RPC for handling such responses?

## Implementation and corner cases
- Store the messages in an array or map, and use mutexs to avoid race conditions for both read and writes.
  - I used `sync.Map` which seemed to be good for given usecase of avoiding duplicates and handling mutexes.
- For broadcasting, append to list of messages, then send the body received in maelstrom message to all nodes. 
- While broadcasting, read kept showing `broadcast 0` too many times. This happened because node `a` broadcasted to `b` and `b` sent the same message back to `a`, and so on. 
  - Fix is to not send message to message source, or the current node.
- Sending message kept showing no handler error. This is because when broadcasting with Send, the other node replied but the current node didn’t handle it. 
  - A way around was to use node.RPC, and create empty handler function.
  - Another is to check if a message has msg_id. If it doesn't, this one doesn't need a reply. return nil instead of `Reply()` in such cases.