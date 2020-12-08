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

	proposal := c.Proposer.AddProposal(requestMessage.Value)

	outgoingMessage := &message.Message{
		Source: c.Proposer.Port,
		Type:   message.PREPARE,
		Payload: message.Prepare{
			Nonce: proposal.Nonce,
			Round: len(c.Proposer.Proposals),
		},
	}

	// Broadcast 'PREPARE' message to the selected quorum of acceptors
	for _, acceptor := range proposal.Quorum {
		util.WriteToMultiFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, proposal.Nonce))
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

	if promiseMessage.Round > len(c.Proposer.Proposals) {
		return nil
	}

	if c.Proposer.Proposals[promiseMessage.Round-1].NonceDoesNotEqual(promiseMessage.Nonce) {
		// TODO add error
		return nil
	}

	c.Proposer.Proposals[promiseMessage.Round-1].RegisterPromise(*promiseMessage)

	if c.Proposer.Proposals[promiseMessage.Round-1].HasInsufficientNumberOfPromises() {
		// TODO exit gracefully
		return nil
	}

	payload := message.Accept{
		Nonce: c.Proposer.Proposals[promiseMessage.Round-1].Nonce,
		Round: promiseMessage.Round,
	}
	if c.Proposer.Proposals[promiseMessage.Round-1].HasAcceptedValueToBroadcast() {
		payload.Value = c.Proposer.Proposals[promiseMessage.Round-1].GetAcceptedValueToBroadcast()
	} else {
		payload.Value = c.Proposer.Proposals[promiseMessage.Round-1].Value
	}

	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    message.ACCEPT,
		Payload: payload,
	}

	// Broadcast 'ACCEPT' message to the proposal's associated quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[promiseMessage.Round-1].Quorum {
		util.WriteToMultiFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Accept: %s", c.Proposer.Port, acceptor, payload.Nonce, payload.Value))
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

	c.Proposer.CurrentNonce = nackMessage.Nonce + 1
	c.Proposer.Proposals[nackMessage.Round-1].Nonce = c.Proposer.CurrentNonce
	c.Proposer.Proposals[nackMessage.Round-1].Promises = []message.Promise{}

	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    message.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.Proposals[nackMessage.Round-1].Nonce},
	}

	// Broadcast updated 'PREPARE' message to the selected quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[nackMessage.Round-1].Quorum {
		util.WriteToMultiFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, c.Proposer.Proposals[nackMessage.Round-1].Nonce))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}
