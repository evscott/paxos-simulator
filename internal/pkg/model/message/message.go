package message

type Message struct {
	Source  uint16      `json:"source"`
	Type    Type        `json:"type"`
	Payload interface{} `json:"payload"`
}
