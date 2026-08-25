// Service 领域模型（RPC 服务注册）。
package model

import (
	"strings"
	"time"
)

const (
	ServiceStatusUp        = "up"
	ServiceStatusDown      = "down"
	ServiceProtocolJsonRpc = "json-rpc"
	ServiceProtocolGrpc    = "grpc"
	ServiceProtocolThrift  = "thrift"
)

type Service struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (x *Service) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "服务名不能为空")
	}
	x.Version = strings.TrimSpace(x.Version)
	if x.Version == "" {
		return NewValidationError("version", "版本不能为空")
	}
	x.Host = strings.TrimSpace(x.Host)
	if x.Host == "" {
		return NewValidationError("host", "主机不能为空")
	}
	if x.Port <= 0 {
		return NewValidationError("port", "端口必须为正数")
	}
	x.Protocol = strings.TrimSpace(x.Protocol)
	if x.Protocol == "" {
		return NewValidationError("protocol", "协议不能为空")
	}
	if x.Protocol != "" {
		switch x.Protocol {
		case "json-rpc", "grpc", "thrift":
		default:
			return NewValidationError("protocol", "协议不合法")
		}
	}
	if x.Status == "" {
		x.Status = ServiceStatusUp
	}
	if !ServiceValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func ServiceValidStatus(s string) bool {
	switch s {
	case "up", "down":
		return true
	default:
		return false
	}
}

var serviceTransitions = map[string]map[string]bool{
	ServiceStatusUp:   {"down": true},
	ServiceStatusDown: {"up": true},
}

func ServiceCanTransition(from, to string) bool {
	if m, ok := serviceTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ServiceFilter struct {
	Protocol string
	Status   string
}

func (f ServiceFilter) Match(x *Service) bool {
	if f.Protocol != "" && x.Protocol != f.Protocol {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
