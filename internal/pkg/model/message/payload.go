package message

import "encoding/json"

type Request struct {
	Value string `json:"value"`
}

type Prepare struct {
	Nonce uint32 `json:"nonce"`
}

type Proposal struct {
	Value string
	Nonce uint32
}

type Promise struct {
	Nonce    uint32   `json:"nonce"`
	Round    uint32   `json:"round"`
	Proposal Proposal `json:"proposal"`
}

type Accept struct {
	Nonce uint32 `json:"nonce"`
	Round uint32 `json:"round"`
	Value string `json:"value"`
}

type Accepted struct {
	Nonce uint32 `json:"nonce"`
	Round uint32 `json:"round"`
	Value string `json:"value"`
}

type Response struct {
	Value string `json:"value"`
}

type Nack struct {
	Nonce uint32 `json:"nonce"`
}

func Unmarshal(in interface{}, out interface{}) error {
	if raw, err := json.Marshal(in); err != nil {
		return err
	} else {
		return json.Unmarshal(raw, &out)
	}
}
