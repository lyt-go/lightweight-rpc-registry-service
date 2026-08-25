// Node 领域模型（机器节点注册）。
package model

import (
	"strings"
	"time"
)

const (
	NodeStatusOnline  = "online"
	NodeStatusOffline = "offline"
)

type Node struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Region    string    `json:"region"`
	Weight    int       `json:"weight"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (x *Node) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "节点名不能为空")
	}
	x.Address = strings.TrimSpace(x.Address)
	if x.Address == "" {
		return NewValidationError("address", "地址不能为空")
	}
	x.Region = strings.TrimSpace(x.Region)
	if x.Weight <= 0 {
		return NewValidationError("weight", "权重必须为正数")
	}
	if x.Status == "" {
		x.Status = NodeStatusOnline
	}
	if !NodeValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func NodeValidStatus(s string) bool {
	switch s {
	case "online", "offline":
		return true
	default:
		return false
	}
}

var nodeTransitions = map[string]map[string]bool{
	NodeStatusOnline:  {"offline": true},
	NodeStatusOffline: {"online": true},
}

func NodeCanTransition(from, to string) bool {
	if m, ok := nodeTransitions[from]; ok {
		return m[to]
	}
	return false
}

type NodeFilter struct {
	Region string
	Status string
}

func (f NodeFilter) Match(x *Node) bool {
	if f.Region != "" && x.Region != f.Region {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
