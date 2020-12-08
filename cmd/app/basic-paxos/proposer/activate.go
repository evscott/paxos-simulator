package Proposer

import (
	"encoding/json"
	"fmt"
	"github.com/paxos/cmd/pkg/model/message"
	"github.com/paxos/cmd/pkg/model/node"
	"log"
	"net"
)

type Config struct {
	Proposer node.Proposer
}

func Activate(port int, acceptors []int) {
	c := &Config{
		Proposer: node.Proposer{
			Port:      port,
			Acceptors: acceptors,
		},
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", c.Proposer.Port))
	if err != nil {
		log.Printf("Failed to connect to port: %d, error: %v\n ", c.Proposer.Port, err)
		return
	}

	log.Printf("Accepting messages on: 127.0.0.1:%d\n", c.Proposer.Port)
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				log.Printf("Error received while listening 127.0.0.1:%d\n", c.Proposer.Port)
			}
		}

		msg := &message.Message{}
		if err := json.NewDecoder(connIn).Decode(msg); err != nil {
			log.Printf("Error decoding %v\n", err)
		}

		switch msg.Type {
		case message.REQUEST:
			c.handleRequest(msg)
			break
		case message.PROMISE:
			c.handlePromise(msg)
			break
		case message.NACK:
			c.handleNack(msg)
			break
		}
	}
}
