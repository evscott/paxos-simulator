package Learner

import (
	"fmt"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/shared/util"
)

func (c *Config) handleAccepted(incomingMessage *message.Message) error {

	acceptedMessage := &message.Accepted{}
	if err := message.Unmarshal(incomingMessage.Payload, acceptedMessage); err != nil {
		return err
	}

	if len(c.Learner.Logs) == 0 {
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
		util.WriteToBasicFile(fmt.Sprintf("learner %d->>+ client: %s was accepted as the value!", c.Learner.Port, c.Learner.Logs[0]))
	}

	return nil
}
