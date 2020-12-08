package multiPaxos

import (
	"fmt"
	"github.com/paxos/internal/app/multi-paxos/acceptor"
	"github.com/paxos/internal/app/multi-paxos/learner"
	"github.com/paxos/internal/app/multi-paxos/proposer"
	"github.com/paxos/internal/pkg/model/message"
	"github.com/paxos/internal/pkg/shared/util"
	"time"
)

func Init() {

	fmt.Println("Multi Paxos!")

	go Proposer.Activate(9001, []int{9002, 9003, 9004})
	go Acceptor.Activate(9002, []int{9005, 9006})
	go Acceptor.Activate(9003, []int{9005, 9006})
	go Acceptor.Activate(9004, []int{9005, 9006})
	go Learner.Activate(9005)
	go Learner.Activate(9006)

	// Wait for nodes to activate
	time.Sleep(time.Second/100)

	// Request that proposer 9000 submit the value "Foo"
	message1 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Foo"},
	}
	util.SendMessage(message1, 9001)

	// Wait some time for Paxos to reach consensus, and then fire another message
	time.Sleep(time.Second/10)

	// Request that proposer 9000 submit the value "Foo"
	message2 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Bar"},
	}
	util.SendMessage(message2, 9001)

	// Wait some time for Paxos to reach consensus
	time.Sleep(time.Second)
}
