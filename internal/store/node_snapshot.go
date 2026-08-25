package store

import "sync"

type NodeLease struct {
	ID       string
	Endpoint string
	Revision int
}

type NodeLeaseStore struct {
	mu    sync.RWMutex
	nodes map[string]*NodeLease
}

func NewNodeLeaseStore(nodes ...NodeLease) *NodeLeaseStore {
	s := &NodeLeaseStore{nodes: make(map[string]*NodeLease)}
	for i := range nodes {
		n := nodes[i]
		s.nodes[n.ID] = &n
	}
	return s
}

func (s *NodeLeaseStore) Snapshot() []*NodeLease {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*NodeLease, 0, len(s.nodes))
	for _, node := range s.nodes {
		out = append(out, node)
	}
	return out
}

func (s *NodeLeaseStore) Update(id, endpoint string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	node := s.nodes[id]
	node.Endpoint = endpoint
	node.Revision++
}
