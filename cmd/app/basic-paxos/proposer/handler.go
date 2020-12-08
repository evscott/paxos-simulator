package Proposer

import (
	"fmt"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/shared/util"
)

func (c *Config) handleRequest(incomingMessage *message.Message) error {
	requestMessage := &message.Request{}
	if err := message.Unmarshal(incomingMessage.Payload, requestMessage); err != nil {
		return err
	}

	c.Proposer.AddProposal(requestMessage.Value)

	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    message.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.Proposals[0].Nonce},
	}

	// Broadcast 'PREPARE' message to the selected quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[0].Quorum {
		util.WriteToBasicFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, c.Proposer.Proposals[0].Nonce))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) handlePromise(incomingMessage *message.Message) error {
	promiseMessage := &message.Promise{}
	if err := message.Unmarshal(incomingMessage.Payload, promiseMessage); err != nil {
		return err
	}

	if c.Proposer.Proposals[0].NonceDoesNotEqual(promiseMessage.Nonce) {
		// TODO add error
		return nil
	}

	c.Proposer.Proposals[0].RegisterPromise(*promiseMessage)

	if c.Proposer.Proposals[0].HasInsufficientNumberOfPromises() {
		return nil
	}

	payload := message.Accept{
		Nonce: c.Proposer.Proposals[0].Nonce,
	}
	if c.Proposer.Proposals[0].HasAcceptedValueToBroadcast() {
		payload.Value = c.Proposer.Proposals[0].GetAcceptedValueToBroadcast()
	} else {
		payload.Value = c.Proposer.Proposals[0].Value
	}

	outgoingMessage := &message.Message{
		Source: c.Proposer.Port,
		Type:   message.ACCEPT,
		Payload: payload,
	}

	// Broadcast 'ACCEPT' message to the proposal's associated quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[0].Quorum {
		util.WriteToBasicFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Accept: %s", c.Proposer.Port, acceptor, payload.Nonce, payload.Value))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) handleNack(incomingMessage *message.Message) error {
	nackMessage := &message.Nack{}
	if err := message.Unmarshal(incomingMessage.Payload, nackMessage); err != nil {
		return err
	}

	if nackMessage.Nonce < c.Proposer.CurrentNonce {
		return nil
	}

	c.Proposer.CurrentNonce = nackMessage.Nonce+1
	c.Proposer.Proposals[0].Nonce = c.Proposer.CurrentNonce
	c.Proposer.Proposals[0].Promises = []message.Promise{}

	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    message.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.Proposals[0].Nonce},
	}

	// Broadcast updated 'PREPARE' message to the selected quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[0].Quorum {
		util.WriteToBasicFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, c.Proposer.Proposals[0].Nonce))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}
