package store

import (
	"sync"

	"lightweightrpc/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	services      map[string]*model.Service
	methods       map[string]*model.Method
	nodes         map[string]*model.Node
	callLogs      map[string]*model.CallLog
	codecs        map[string]*model.Codec
	retryPolicies map[string]*model.RetryPolicy
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		services:      make(map[string]*model.Service),
		methods:       make(map[string]*model.Method),
		nodes:         make(map[string]*model.Node),
		callLogs:      make(map[string]*model.CallLog),
		codecs:        make(map[string]*model.Codec),
		retryPolicies: make(map[string]*model.RetryPolicy),
	}
}

var _ Store = (*MemoryStore)(nil)
