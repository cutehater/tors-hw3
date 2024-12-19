package nodes

import (
	"context"
	"fmt"
	"hw3/internal/protos"
)

func (v *Node) ProcessUpdates(ctx context.Context, in *protos.ProcessUpdatesIn) (*protos.ProcessUpdatesOut, error) {
	v.processUpdates(in)
	return &protos.ProcessUpdatesOut{}, nil
}

func (v *Node) processUpdates(in *protos.ProcessUpdatesIn) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.msgBuffer = append(v.msgBuffer, *messageFromProto(in))

	hasUpdates := true
	for hasUpdates {
		hasUpdates = false

		for i := 0; i < len(v.msgBuffer); i++ {
			msg := v.msgBuffer[i]
			if msg.deps.lessOrEqual(v.delivered) {
				fmt.Printf("Deliver updates %v from sender %d with deps %v\n", msg.updates, msg.sender, msg.deps)
				hasUpdates = true
				v.deliverUpdates(msg.sender, msg.deps, msg.updates)
				v.delivered[msg.sender] = v.delivered[msg.sender] + 1
				v.msgBuffer = append(v.msgBuffer[:i], v.msgBuffer[i+1:]...)
				i--
			}
		}
	}
}
