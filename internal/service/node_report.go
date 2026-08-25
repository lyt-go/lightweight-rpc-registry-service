package service

import (
	"fmt"

	"lightweightrpc/internal/store"
)

type NodeReporter struct{ nodes *store.NodeLeaseStore }

func NewNodeReporter(nodes *store.NodeLeaseStore) *NodeReporter { return &NodeReporter{nodes: nodes} }

func (r *NodeReporter) DelayedReport(captured chan<- struct{}, release <-chan struct{}) <-chan []string {
	result := make(chan []string, 1)
	snapshot := r.nodes.Snapshot()
	go func() {
		captured <- struct{}{}
		<-release
		lines := make([]string, 0, len(snapshot))
		for _, node := range snapshot {
			lines = append(lines, fmt.Sprintf("%s@%s#%d", node.ID, node.Endpoint, node.Revision))
		}
		result <- lines
	}()
	return result
}
