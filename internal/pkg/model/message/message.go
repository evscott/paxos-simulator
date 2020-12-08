package message

type Message struct {
	Source  int      `json:"source"`
	Type    Type        `json:"type"`
	Payload interface{} `json:"payload"`
}
