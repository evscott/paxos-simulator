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

	go Proposer.Activate(8000, []uint16{9001, 9002, 9003, 9004, 9005})
	go Proposer.Activate(8999, []uint16{9001, 9002, 9003, 9004, 9005})
	go Proposer.Activate(9000, []uint16{9001, 9002, 9003, 9004, 9005})
	go Acceptor.Activate(9001, []uint16{9006})
	go Acceptor.Activate(9002, []uint16{9006})
	go Acceptor.Activate(9003, []uint16{9006})
	go Acceptor.Activate(9004, []uint16{9006})
	go Acceptor.Activate(9005, []uint16{9006})
	go Learner.Activate(9006, []uint16{9001, 9002, 9003, 9004, 9005})

	time.Sleep(time.Second/100)

	message1 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Poopy"},
	}
	util.SendMessage(message1, 9000)

	message2 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Banana"},
	}
	util.SendMessage(message2, 8999)

	time.Sleep(time.Second/10)

	message3 := &message.Message{
		Source:  0,
		Type:    message.REQUEST,
		Payload: message.Request{Value: "Blueberry"},
	}
	util.SendMessage(message3, 8000)

	time.Sleep(time.Second/10)
}
