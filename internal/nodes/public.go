package nodes

import (
	"context"
	"fmt"
	"hw3/internal/protos"
	"sync"
)

func New(id int64, GRPCNodesAddress []string) *Node {
	node := &Node{
		mu:               sync.RWMutex{},
		id:               id,
		sendCnt:          0,
		delivered:        make([]int64, len(GRPCNodesAddress)),
		data:             make(map[Key]Value),
		dataMeta:         make(map[Key]meta),
		grpcNodesClients: make([]protos.NodeClient, len(GRPCNodesAddress)),
	}

	wg := sync.WaitGroup{}
	for i, address := range GRPCNodesAddress {
		go func() {
			defer wg.Done()
			wg.Add(1)
			node.createGRPCClient(i, address)
		}()
	}
	wg.Wait()

	return node
}

func (v *Node) Broadcast(updates []Update) {
	v.mu.Lock()

	msg := protos.ProcessUpdatesIn{
		Sender:  v.id,
		Deps:    make([]int64, len(v.delivered)),
		Updates: make([]*protos.Update, 0, len(updates)),
	}
	for _, update := range updates {
		msg.Updates = append(msg.Updates, update.toProto())
	}
	copy(msg.Deps, v.delivered)

	msg.Deps[v.id] = v.sendCnt
	v.sendCnt++

	v.mu.Unlock()

	fmt.Printf("Broadcast updates %v with deps %v\n", msg.Updates, msg.Deps)
	for i, node := range v.grpcNodesClients {
		if i == int(v.id) {
			go func() {
				v.processUpdates(&msg)
			}()
		} else {
			go func() {
				for {
					_, err := node.ProcessUpdates(context.Background(), &msg)
					if err == nil {
						break
					}
				}
			}()
		}
	}
}

func (v *Node) Get() map[Key]Value {
	v.mu.RLock()
	defer v.mu.RUnlock()
	res := make(map[Key]Value, len(v.data))
	for key, value := range v.data {
		res[key] = value
	}
	return res
}
