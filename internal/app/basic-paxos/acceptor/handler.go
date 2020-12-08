package Acceptor

import (
	"fmt"
	"github.com/paxos/internal/pkg/model/message"
	"github.com/paxos/internal/pkg/shared/util"
)

func (c *Config) handlePrepare(incomingMessage *message.Message) error {
	prepareMessage := &message.Prepare{}
	if err := message.Unmarshal(incomingMessage.Payload, prepareMessage); err != nil {
		return err
	}

	for _, promise := range c.Acceptor.Promises {
		if prepareMessage.Nonce <= promise.Nonce {
			outgoingMessage := &message.Message{
				Source: c.Acceptor.Port,
				Type:   message.NACK,
				Payload: message.Nack{
					Nonce: promise.Nonce,
				},
			}
			util.WriteToFile(fmt.Sprintf("%d-->>%d:(%d) Nack", c.Acceptor.Port, incomingMessage.Source, prepareMessage.Nonce))
			if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
				return err
			}
			return nil
		}
	}

	promise := message.Promise{
		Nonce: prepareMessage.Nonce,
	}
	if c.Acceptor.HasAcceptedProposal(1) {
		promise.Proposal = c.Acceptor.AcceptedProposals[0]
	}
	c.Acceptor.AddPromise(promise)

	outgoingMessage := &message.Message{
		Source:  c.Acceptor.Port,
		Type:    message.PROMISE,
		Payload: promise,
	}

	// Send 'PROMISE' message to proposer
	if c.Acceptor.HasAcceptedProposal(1) {
		util.WriteToFile(fmt.Sprintf("%d-->>%d:(%d) Promise: %+v", c.Acceptor.Port, incomingMessage.Source, prepareMessage.Nonce, promise.Proposal))
		if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
			return err
		}
	} else {
		util.WriteToFile(fmt.Sprintf("%d-->>%d:(%d) Promise", c.Acceptor.Port, incomingMessage.Source, prepareMessage.Nonce))
		if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) handleAccept(incomingMessage *message.Message) error {
	acceptMessage := &message.Accept{}
	if err := message.Unmarshal(incomingMessage.Payload, acceptMessage); err != nil {
		return err
	}

	for _, promise := range c.Acceptor.Promises {
		if acceptMessage.Nonce < promise.Nonce {
			outgoingMessage := &message.Message{
				Source: c.Acceptor.Port,
				Type:   message.NACK,
				Payload: message.Nack{
					Nonce: promise.Nonce,
				},
			}
			util.WriteToFile(fmt.Sprintf("%d-->>%d:(%d) Nack", c.Acceptor.Port, incomingMessage.Source, acceptMessage.Nonce))
			if err := util.SendMessage(outgoingMessage, incomingMessage.Source); err != nil {
				return err
			}
		}
	}

	c.Acceptor.AddAcceptedProposal(acceptMessage.Value, acceptMessage.Nonce)

	outgoingMessage := &message.Message{
		Source: c.Acceptor.Port,
		Type:   message.ACCEPTED,
		Payload: message.Accepted{
			Nonce: acceptMessage.Nonce,
			Value: acceptMessage.Value,
		},
	}

	// Broadcast 'ACCEPTED' message to all learners
	for _, learner := range c.Acceptor.Learners {
		util.WriteToFile(fmt.Sprintf("%d-->>%d:(%d) Accepted: %s", c.Acceptor.Port, learner, acceptMessage.Nonce, acceptMessage.Value))
		if err := util.SendMessage(outgoingMessage, learner); err != nil {
			return err
		}
	}

	return nil
}
