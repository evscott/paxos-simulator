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

	if acceptedMessage.Round > len(c.Learner.Logs) {
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
		util.WriteToMultiFile(fmt.Sprintf("learner %d->>+ client: %s was accepted as the value!", c.Learner.Port, c.Learner.Logs[acceptedMessage.Round-1]))
	}

	return nil
}
