package main

import (
	"github.com/paxos/cmd/app/basic-paxos"
	"github.com/paxos/cmd/app/multi-paxos"
)

func main() {
	BasicPaxos.Init()
	multiPaxos.Init()
}
