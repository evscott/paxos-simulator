package Acceptor

import (
	"encoding/json"
	"fmt"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/model/node"
	"log"
	"net"
)

type Config struct {
	Acceptor node.Acceptor
}

func Activate(port int, learners []int) {
	c := &Config{
		Acceptor: node.Acceptor{
			Port:     port,
			Learners: learners,
		},
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", c.Acceptor.Port))
	if err != nil {
		log.Printf("Failed to connect to port: %d, error: %v\n ", c.Acceptor.Port, err)
		return
	}

	log.Printf("Accepting messages on: 127.0.0.1:%d\n", c.Acceptor.Port)
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening 127.0.0.1:%d\n", c.Acceptor.Port)
			}
		}

		msg := &message.Message{}
		if err := json.NewDecoder(connIn).Decode(msg); err != nil {
			log.Printf("Error decoding %v\n", err)
		}

		switch msg.Type {
		case message.PREPARE:
			c.handlePrepare(msg)
			break
		case message.ACCEPT:
			c.handleAccept(msg)
			break
		}
	}
}
