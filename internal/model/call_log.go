// CallLog 领域模型（调用记录）。
package model

import (
	"strings"
	"time"
)

type CallLog struct {
	ID        string    `json:"id"`
	ServiceID string    `json:"service_id"`
	MethodID  string    `json:"method_id"`
	Caller    string    `json:"caller"`
	LatencyMs int64     `json:"latency_ms"`
	Code      int       `json:"code"`
	ErrMsg    string    `json:"err_msg"`
	Timestamp int64     `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (x *CallLog) Validate() error {
	x.ServiceID = strings.TrimSpace(x.ServiceID)
	if x.ServiceID == "" {
		return NewValidationError("service_id", "服务ID不能为空")
	}
	x.MethodID = strings.TrimSpace(x.MethodID)
	if x.MethodID == "" {
		return NewValidationError("method_id", "方法ID不能为空")
	}
	x.Caller = strings.TrimSpace(x.Caller)
	if x.LatencyMs < 0 {
		return NewValidationError("latency_ms", "耗时毫秒不能为负数")
	}
	x.ErrMsg = strings.TrimSpace(x.ErrMsg)
	if x.Timestamp <= 0 {
		return NewValidationError("timestamp", "时间戳必须为正数")
	}
	return nil
}

type CallLogFilter struct {
	ServiceId string
	Caller    string
}

func (f CallLogFilter) Match(x *CallLog) bool {
	if f.ServiceId != "" && x.ServiceID != f.ServiceId {
		return false
	}
	if f.Caller != "" && x.Caller != f.Caller {
		return false
	}
	return true
}
