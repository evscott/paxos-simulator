package Learner

import (
	"fmt"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/shared/util"
)

// handleAccepted processes accepted messages
func (c *Config) handleAccepted(incomingMessage *message.Message) error {

	acceptedMessage := &message.Accepted{}
	if err := message.Unmarshal(incomingMessage.Payload, acceptedMessage); err != nil {
		return err
	}

	// If the learner has no log of any other value being accepted by the network, log it and inform the client
	// of the accepted value
	if len(c.Learner.Logs) == 0 {
		c.Learner.Logs = append(c.Learner.Logs, acceptedMessage.Value)
		util.WriteToBasicFile(fmt.Sprintf("learner %d->>+ client: %s was accepted as the value!", c.Learner.Port, c.Learner.Logs[0]))
	}

	return nil
}
