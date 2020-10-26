package message

type Type string

const (
	REQUEST  Type = "request"
	PREPARE  Type = "prepare"
	PROMISE  Type = "promise"
	ACCEPT   Type = "accept"
	ACCEPTED Type = "accepted"
	RESPONSE Type = "response"
	NACK     Type = "nack"
)
