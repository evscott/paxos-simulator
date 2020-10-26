package node

type Learner struct {
	Port              uint16
	Acceptors         []uint16
	ProposalAcceptors []uint16
	Value             string
	Nonce             uint32
}

func (l *Learner) RegisterAcceptor(acceptor uint16) {
	l.ProposalAcceptors = append(l.ProposalAcceptors, acceptor)
}
