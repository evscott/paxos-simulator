package BasicPaxos

import (
	"fmt"
	Acceptor "github.com/paxos/internal/app/basic-paxos/acceptor"
	Learner "github.com/paxos/internal/app/basic-paxos/learner"
	Proposer "github.com/paxos/internal/app/basic-paxos/proposer"
	"github.com/paxos/internal/pkg/model/message"
	"github.com/paxos/internal/pkg/shared/util"
	"time"
)

func Init() {

	fmt.Println("Basic Paxos!")

	go Proposer.Activate(9001, []uint16{9003, 9004, 9005})
	go Proposer.Activate(9002, []uint16{9003, 9004, 9005})
	go Acceptor.Activate(9003, []uint16{9006, 9007})
	go Acceptor.Activate(9004, []uint16{9006, 9007})
	go Acceptor.Activate(9005, []uint16{9006, 9007})
	go Learner.Activate(9006)
	go Learner.Activate(9007)

	// Wait for nodes to activate
	time.Sleep(time.Second/100)

	// Request that proposer 9000 submit the value "Foo"
	message1 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Foo"},
	}
	util.SendMessage(message1, 9001)

	// Request that proposer 9001 submit the value "Foo"
	message2 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Bar"},
	}
	util.SendMessage(message2, 9002)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second/10)

	fmt.Println("What the hell...")
}
