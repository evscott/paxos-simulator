package Proposer

import (
	"fmt"
	"github.com/paxos/internal/pkg/model/message"
	"github.com/paxos/internal/pkg/shared/util"
)

func (c *Config) handleRequest(incomingMessage *message.Message) error {
	requestMessage := &message.Request{}
	if err := message.Unmarshal(incomingMessage.Payload, requestMessage); err != nil {
		return err
	}

	c.Proposer.CurrentProposal.Value = requestMessage.Value
	c.Proposer.CurrentProposal.Nonce = c.Proposer.GetNonce()
	c.Proposer.CurrentProposal.Quorum = c.Proposer.GetQuorum()

	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    message.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.CurrentProposal.Nonce},
	}

	// Broadcast 'PREPARE' message to the selected quorum of acceptors
	for _, acceptor := range c.Proposer.CurrentProposal.Quorum {
		util.WriteToFile(fmt.Sprintf("%d->>%d: Prepare n:%d", c.Proposer.Port, acceptor, c.Proposer.CurrentProposal.Nonce))
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

	if c.Proposer.CurrentProposal.NonceDoesNotEqual(promiseMessage.Nonce) {
		// TODO add error
		return nil
	}

	c.Proposer.CurrentProposal.RegisterPromise(*promiseMessage)

	if c.Proposer.CurrentProposal.HasInsufficientNumberOfPromises() {
		// TODO exit gracefully
		return nil
	}

	payload := message.Accept{
		Nonce: c.Proposer.CurrentProposal.Nonce,
	}
	if c.Proposer.CurrentProposal.HasAcceptedValueToBroadcast() {
		payload.Value = c.Proposer.CurrentProposal.GetAcceptedValueToBroadcast()
	} else {
		payload.Value = c.Proposer.CurrentProposal.Value
	}

	outgoingMessage := &message.Message{
		Source: c.Proposer.Port,
		Type:   message.ACCEPT,
		Payload: payload,
	}

	// Broadcast 'ACCEPT' message to the proposal's associated quorum of acceptors
	for _, acceptor := range c.Proposer.CurrentProposal.Quorum {
		util.WriteToFile(fmt.Sprintf("%d->>%d: Accept n:%d v:%s", c.Proposer.Port, acceptor, payload.Nonce, payload.Value))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) handleAccepted(incomingMessage *message.Message) error {
	// TODO functionally to be determined
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
	c.Proposer.CurrentProposal.Nonce = c.Proposer.CurrentNonce
	c.Proposer.CurrentProposal.Promises = []message.Promise{}

	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    message.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.CurrentProposal.Nonce},
	}

	// Broadcast updated 'PREPARE' message to the selected quorum of acceptors
	for _, acceptor := range c.Proposer.CurrentProposal.Quorum {
		util.WriteToFile(fmt.Sprintf("%d->>%d: Prepare n:%d", c.Proposer.Port, acceptor, c.Proposer.CurrentProposal.Nonce))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}
