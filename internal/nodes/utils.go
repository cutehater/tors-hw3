package nodes

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"hw3/internal/protos"
	"log"
)

func (v *Node) createGRPCClient(id int, address string) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to node %d at %s: %v", id, address, err)
	}

	v.mu.Lock()
	defer v.mu.Unlock()
	v.grpcNodesClients[id] = protos.NewNodeClient(conn)
}

// call under v.mu.Lock() only
func (v *Node) deliverUpdates(sender int64, deps vectorClock, updates []Update) {
	for _, update := range updates {
		lastUpdateMeta, ok := v.dataMeta[update.Key]
		if !ok || lastUpdateMeta.deps.less(deps) || lastUpdateMeta.sender < sender {
			v.data[update.Key] = update.Value
			v.dataMeta[update.Key] = meta{
				sender: sender,
				deps:   deps,
			}
		}
	}
}
