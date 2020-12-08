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

	if acceptedMessage.Round > len(c.Learner.Logs) {
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
		util.WriteToFile(fmt.Sprintf("%d->>+%d: %s was accepted as the value!", c.Learner.Port, c.Learner.Port, c.Learner.Logs[acceptedMessage.Round-1]))
	}

	return nil
}
