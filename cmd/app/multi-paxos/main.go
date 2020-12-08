package MultiPaxos

import (
	"fmt"
	Acceptor "github.com/paxos/cmd/app/multi-paxos/acceptor"
	Learner "github.com/paxos/cmd/app/multi-paxos/learner"
	Proposer "github.com/paxos/cmd/app/multi-paxos/proposer"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/shared/util"
	"time"
)

// Initializes an instance of Multi-Paxos with several nodes: one proposer, three acceptors, and one learner
// The instance simulates a scenario where a client submits two requests to the same proposer with different values
// The requests are processed in rounds by the network, and the network arrives to a consensus on both values in their
// respective rounds
func Init() {

	fmt.Println("Initializing Multi-Paxos...")
	util.CreateNewFile("multi")

	go Proposer.Activate(9001, []int{9002, 9003, 9004})
	go Acceptor.Activate(9002, []int{9005})
	go Acceptor.Activate(9003, []int{9005})
	go Acceptor.Activate(9004, []int{9005})
	go Learner.Activate(9005)

	// Wait for nodes to activate
	time.Sleep(time.Second / 100)

	// Request that proposer 9000 propose the value "Foo"
	message1 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Foo"},
	}

	util.WriteToMultiFile(fmt.Sprintf("client ->> proposer 9001: Request: %v", "Foo"))
	util.WriteToMultiFile(fmt.Sprintf("Note over client,proposer 9001: Initialize round 1\n"))
	util.SendMessage(message1, 9001)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)

	// Request that proposer 9000 propose the value "Bar"
	message2 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Bar"},
	}

	util.WriteToMultiFile(fmt.Sprintf("client ->> proposer 9001: Request: %v", "Bar"))
	util.WriteToMultiFile(fmt.Sprintf("Note over client,proposer 9001: Initialize round 2\n"))
	util.SendMessage(message2, 9001)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second / 10)
}
