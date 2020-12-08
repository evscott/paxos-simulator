package BasicPaxos

import (
	"fmt"
	Acceptor "github.com/paxos/cmd/app/basic-paxos/acceptor"
	Learner "github.com/paxos/cmd/app/basic-paxos/learner"
	Proposer "github.com/paxos/cmd/app/basic-paxos/proposer"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/shared/util"
	"time"
)

// Initializes an instance of Basic-Paxos with several nodes: two proposers, three acceptors, and one learner
// The instance simulates a scenario where a client submits two requests to different proposers with different values
// The first request submitted *should* be accepted as the value on the network
func Init() {

	fmt.Println("Initializing Basic-Paxos...")
	util.CreateNewFile("basic")

	go Proposer.Activate(8001, []int{8003, 8004, 8005})
	go Proposer.Activate(8002, []int{8003, 8004, 8005})
	go Acceptor.Activate(8003, []int{8006})
	go Acceptor.Activate(8004, []int{8006})
	go Acceptor.Activate(8005, []int{8006})
	go Learner.Activate(8006)

	// Wait for nodes to activate
	time.Sleep(time.Second / 100)

	// Request that proposer 8001 propose the value "Foo"
	message1 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Foo"},
	}
	util.WriteToBasicFile(fmt.Sprintf("client ->> proposer 8001: Request: %v", "Foo"))
	util.WriteToBasicFile(fmt.Sprintf("Note over client,proposer 8001: Initialize round 1\n"))
	util.SendMessage(message1, 8001)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)

	// Request that proposer 8002 propose the value "Bar"
	message2 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Bar"},
	}
	util.WriteToBasicFile(fmt.Sprintf("client ->> proposer 8002: Request: %v", "Bar"))
	util.WriteToBasicFile(fmt.Sprintf("Note over client,proposer 8002: Initialize round 2\n"))
	util.SendMessage(message2, 8002)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)
}
