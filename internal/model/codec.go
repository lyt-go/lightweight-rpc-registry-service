// Codec 领域模型（编解码器）。
package model

import (
	"strings"
	"time"
)

const (
	CodecStatusActive   = "active"
	CodecStatusInactive = "inactive"
	CodecTypeJson       = "json"
	CodecTypeProtobuf   = "protobuf"
	CodecTypeMsgpack    = "msgpack"
)

type Codec struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Schema    string    `json:"schema"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (x *Codec) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "名称不能为空")
	}
	x.Type = strings.TrimSpace(x.Type)
	if x.Type == "" {
		return NewValidationError("type", "编码类型不能为空")
	}
	if x.Type != "" {
		switch x.Type {
		case "json", "protobuf", "msgpack":
		default:
			return NewValidationError("type", "编码类型不合法")
		}
	}
	x.Version = strings.TrimSpace(x.Version)
	x.Schema = strings.TrimSpace(x.Schema)
	if x.Status == "" {
		x.Status = CodecStatusActive
	}
	if !CodecValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func CodecValidStatus(s string) bool {
	switch s {
	case "active", "inactive":
		return true
	default:
		return false
	}
}

var codecTransitions = map[string]map[string]bool{
	CodecStatusActive:   {"inactive": true},
	CodecStatusInactive: {"active": true},
}

func CodecCanTransition(from, to string) bool {
	if m, ok := codecTransitions[from]; ok {
		return m[to]
	}
	return false
}

type CodecFilter struct {
	Type   string
	Status string
}

func (f CodecFilter) Match(x *Codec) bool {
	if f.Type != "" && x.Type != f.Type {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
