# Notes and Challenges

- When topology is missing, we SEND the message to all nodes directly.
- If in a topology, if A -> B -> C, how do we send a message such that each node receives only once (avoid cycle)?
  - For each node, check all neighbours using read. If neighbour has already received the message, skip otherwise send.
  - Create a new RPC for handling such responses?