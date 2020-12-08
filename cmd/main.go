package main

import (
	"github.com/paxos/cmd/app/basic-paxos"
	"github.com/paxos/cmd/app/multi-paxos"
)

// Main runner for the Paxos simulator
// Initializes an instance of Basic-Paxos and Multi-Paxos
// Output is written to the artifacts folder
func main() {
	BasicPaxos.Init()
	MultiPaxos.Init()
}
