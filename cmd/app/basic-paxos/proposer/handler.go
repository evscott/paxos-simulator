package Proposer

import (
	"fmt"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/shared/constant"
	"github.com/paxos/cmd/pkg/shared/util"
)

// handleRequest processes request messages
func (c *Config) handleRequest(incomingMessage *message.Message) error {
	requestMessage := &message.Request{}
	if err := message.Unmarshal(incomingMessage.Payload, requestMessage); err != nil {
		return err
	}

	// Add the request to the proposers list of proposals
	c.Proposer.AddProposal(requestMessage.Value)

	// Construct the proposal message
	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    constant.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.Proposals[0].Nonce},
	}

	// Broadcast a prepare message to the proposals quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[0].Quorum {
		util.WriteToBasicFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, c.Proposer.Proposals[0].Nonce))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

// handlePromise processes promise messages
func (c *Config) handlePromise(incomingMessage *message.Message) error {
	promiseMessage := &message.Promise{}
	if err := message.Unmarshal(incomingMessage.Payload, promiseMessage); err != nil {
		return err
	}

	// Add the promise to the proposers list of promises
	c.Proposer.Proposals[0].RegisterPromise(*promiseMessage)

	// If the proposer has not received a sufficient number of promises for its current proposal, do nothing
	if c.Proposer.Proposals[0].HasInsufficientNumberOfPromises() {
		return nil
	}

	// Construct the accept message
	// If the proposer has learned that another proposal has already been accepted, share that with its proposals
	// quorum of acceptors
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
		Type:   constant.ACCEPT,
		Payload: payload,
	}

	// Broadcast an accept message to the proposals  quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[0].Quorum {
		util.WriteToBasicFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Accept: %s", c.Proposer.Port, acceptor, payload.Nonce, payload.Value))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}

// handleNack processes nack messages
func (c *Config) handleNack(incomingMessage *message.Message) error {
	nackMessage := &message.Nack{}
	if err := message.Unmarshal(incomingMessage.Payload, nackMessage); err != nil {
		return err
	}

	// If the nack is less than the proposers current nonce, it is outdated and can be ignored
	if nackMessage.Nonce < c.Proposer.CurrentNonce {
		return nil
	}

	// Construct a new prepare message with a nonce that is greater than the nonce supplied in the nack message
	c.Proposer.CurrentNonce = nackMessage.Nonce+1
	c.Proposer.Proposals[0].Nonce = c.Proposer.CurrentNonce
	c.Proposer.Proposals[0].Promises = []message.Promise{}
	outgoingMessage := &message.Message{
		Source:  c.Proposer.Port,
		Type:    constant.PREPARE,
		Payload: message.Prepare{Nonce: c.Proposer.Proposals[0].Nonce},
	}

	// Broadcast the updated prepare message to the proposals quorum of acceptors
	for _, acceptor := range c.Proposer.Proposals[0].Quorum {
		util.WriteToBasicFile(fmt.Sprintf("proposer %d->> acceptor %d:(%d) Prepare", c.Proposer.Port, acceptor, c.Proposer.Proposals[0].Nonce))
		if err := util.SendMessage(outgoingMessage, acceptor); err != nil {
			return err
		}
	}

	return nil
}
