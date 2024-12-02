# Notes and Challenges

Okay I messed this one up. Didn't setup my bash script correctly, which showed me everything worked. SMH.

I have two approaches in my mind:
- Using broadcast workers to send messages asynchronously
- Setup goroutine for each message to be broadcasted. This blocks until the message has been successfully sent, or n tries have finished.

## Broadcast Workers

Architecture contains two types of processes: main process, and workers.

Main process -> Loop through all nodes where message needs to be sent. Send job (destination and message) to one of the workers through channel
Worker -> Fetch job through the channel, and send it infinitely until there's no error.

## Goroutine for each unsent message with explicit retry mechanism

Maintain a list of nodes where you didn't receive `broadcast_ok`. Until this list isn't empty, for each unsent node, create a new goroutine, send this message. If received `broadcast_ok`, remove this node from the list else continue trying.