# liarslie

Liars lie is a game where you can query a network of agents to determine the real value of a network given that some % of
agents are liars.

## Installation

`go install`

`liarslie start`


## Commands
You can use the flag `--help` to learn about possible parameters.

`play`:  Plays a round of liars lie. Determines the real value of the network.

`playexpert`:  Play a round only querying a subset of agents (default is 10)

`kill`:  Kill a node on the network

`extend`:  Extends the network.

## Protocol Design & Notes

- The protocol has N agent and a single registry (the config file)
- Each agent has on its local state the connection to the other nodes.
- We designed a protocol where each agent performs a consensus where it counts the most repeated value recieved from all its peers.
- On expert mode we can limit the consensus to be until a % of values is reached (given by the liar-ratio).
- Each node receives all values from all other peer nodes and determines the true value V locally. After that it responds to the requested resource with the final value V.
- We simulate the network using goroutines and go channels for communications.
- For expert mode, each node queries all of its peers and asks for its value. The query process occurs using a  req/response approach and the agent waits for all the nodes to respond, or timeouts after a certain period. This approach optimizes for reducing the state size of the agents but there might be other approaches to optimize for network latency.
- For standard play a single message is sent to each node asking for its value, and the clients performs the consensus and counting of values.
- The number of messages sent on the network is NXN where N is the number of online agents.
- We assumed the `extend` and `kill` commands are reconfigurations and not distributed proposals, so the changes to these values happen without voting or approval of agents.
