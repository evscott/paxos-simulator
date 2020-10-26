package Learner

import (
	"fmt"
	"github.com/paxos/internal/pkg/model/message"
	"github.com/paxos/internal/pkg/shared/util"
)

func (c *Config) handleAccepted(incomingMessage *message.Message) error {

	acceptedMessage := &message.Accepted{}
	if err := message.Unmarshal(incomingMessage.Payload, acceptedMessage); err != nil {
		return err
	}

	c.Learner.RegisterAcceptor(incomingMessage.Source)
	c.Learner.Value = acceptedMessage.Value
	c.Learner.Nonce = acceptedMessage.Nonce

	if len(c.Learner.ProposalAcceptors) > 2 {
		util.WriteToFile(fmt.Sprintf("%d->>+%d: %s was accepted as the value!", c.Learner.Port, c.Learner.Port, c.Learner.Value))
	}

	return nil
}
