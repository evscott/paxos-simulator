package message

import "github.com/paxos/cmd/pkg/shared/constant"

type Message struct {
	Source  int           `json:"source"`
	Type    constant.Type `json:"type"`
	Payload interface{}   `json:"payload"`
}
