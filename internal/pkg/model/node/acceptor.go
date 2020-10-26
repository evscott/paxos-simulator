package node

import (
	"github.com/paxos/internal/pkg/model/message"
)

type Acceptor struct {
	Port             uint16
	Learners         []uint16
	Promises         []message.Promise
	AcceptedProposal message.Proposal
}

func (a *Acceptor) HasPromisedGreaterNonceThan(nonce uint32) bool {
	for _, promise := range a.Promises {
		if nonce <= promise.Nonce {
			return true
		}
	}
	return false
}

func (a *Acceptor) HasAcceptedProposal() bool {
	return a.AcceptedProposal != (message.Proposal{})
}

func (a *Acceptor) RegisterPromise(promise message.Promise) {
	a.Promises = append(a.Promises, promise)
}