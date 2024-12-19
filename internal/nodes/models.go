package nodes

import (
	"hw3/internal/protos"
	"sync"
)

type Key string
type Value string
type Update struct {
	Key   Key
	Value Value
}
type Node struct {
	protos.UnimplementedNodeServer
	mu               sync.RWMutex
	id               int64
	sendCnt          int64
	delivered        vectorClock
	msgBuffer        []message
	data             map[Key]Value
	dataMeta         map[Key]meta
	grpcNodesClients []protos.NodeClient
}
type vectorClock []int64
type message struct {
	sender  int64
	deps    vectorClock
	updates []Update
}
type meta struct {
	sender int64
	deps   vectorClock
}

func (u *Update) toProto() *protos.Update {
	return &protos.Update{
		Key:   string(u.Key),
		Value: string(u.Value),
	}
}

func messageFromProto(in *protos.ProcessUpdatesIn) *message {
	msg := &message{
		sender:  in.Sender,
		deps:    in.Deps,
		updates: make([]Update, len(in.Updates)),
	}
	for i, u := range in.Updates {
		msg.updates[i] = Update{
			Key:   Key(u.Key),
			Value: Value(u.Value),
		}
	}
	return msg
}

func (lc vectorClock) lessOrEqual(rc vectorClock) bool {
	for i := range lc {
		if lc[i] > rc[i] {
			return false
		}
	}
	return true
}

func (lc vectorClock) less(rc vectorClock) bool {
	strictly := false
	for i := range lc {
		if lc[i] < rc[i] {
			strictly = true
		}
		if lc[i] > rc[i] {
			return false
		}
	}
	return strictly
}
