# Notes and Challenges

- When topology is missing, we SEND the message to all nodes directly.
- If in a topology, if A -> B -> C, how do we send a message such that each node receives only once (avoid cycle)?
  - For each node, check all neighbours using read. If neighbour has already received the message, skip otherwise send.
  - Create a new RPC for handling such responses?

## Implementation and corner cases
- Broadcasting with topology not implemented yet.
- Store the messages in an array or map, and use mutexs to avoid race conditions for both read and writes.
  - I used `sync.Map` which seemed to be good for given usecase, given a map could also prevent addition of duplicate entries.
- For broadcasting, add append to list of messages, then send the body received in maelstrom message to all nodes. 
- Corner cases: Sending message to another node leads to an infinite loop where that node will keep sending the message to all other nodes. Fix it by returning nil if the message is already present.