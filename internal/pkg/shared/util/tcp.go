package util

import (
	"encoding/json"
	"fmt"
	"github.com/paxos/internal/pkg/model/message"
	"net"
	"time"
)

func SendMessage(msg *message.Message, dest uint16) error {
	connOut, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d",dest), time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Could not connect with 127.0.0.1:%d\n", dest)
			return err
		}
	}

	if err := json.NewEncoder(connOut).Encode(msg); err != nil {
		//fmt.Printf("Could not enncode message: %v\n", msg)
		return err
	}
	return nil
}
