// Method 领域模型（RPC 方法定义）。
package model

import (
	"strings"
	"time"
)

const (
	MethodStatusActive     = "active"
	MethodStatusDeprecated = "deprecated"
)

type Method struct {
	ID         string    `json:"id"`
	ServiceID  string    `json:"service_id"`
	Name       string    `json:"name"`
	InputType  string    `json:"input_type"`
	OutputType string    `json:"output_type"`
	TimeoutMs  int       `json:"timeout_ms"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (x *Method) Validate() error {
	x.ServiceID = strings.TrimSpace(x.ServiceID)
	if x.ServiceID == "" {
		return NewValidationError("service_id", "服务ID不能为空")
	}
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "方法名不能为空")
	}
	x.InputType = strings.TrimSpace(x.InputType)
	if x.InputType == "" {
		return NewValidationError("input_type", "入参类型不能为空")
	}
	x.OutputType = strings.TrimSpace(x.OutputType)
	if x.OutputType == "" {
		return NewValidationError("output_type", "出参类型不能为空")
	}
	if x.TimeoutMs <= 0 {
		return NewValidationError("timeout_ms", "超时毫秒必须为正数")
	}
	if x.Status == "" {
		x.Status = MethodStatusActive
	}
	if !MethodValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func MethodValidStatus(s string) bool {
	switch s {
	case "active", "deprecated":
		return true
	default:
		return false
	}
}

var methodTransitions = map[string]map[string]bool{
	MethodStatusActive:     {"deprecated": true},
	MethodStatusDeprecated: {"active": true},
}

func MethodCanTransition(from, to string) bool {
	if m, ok := methodTransitions[from]; ok {
		return m[to]
	}
	return false
}

type MethodFilter struct {
	ServiceId string
	Status    string
}

func (f MethodFilter) Match(x *Method) bool {
	if f.ServiceId != "" && x.ServiceID != f.ServiceId {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
