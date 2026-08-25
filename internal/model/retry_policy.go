// RetryPolicy 领域模型（超时重试策略）。
package model

import (
	"strings"
	"time"
)

const (
	RetryPolicyStatusEnabled      = "enabled"
	RetryPolicyStatusDisabled     = "disabled"
	RetryPolicyBackoffFixed       = "fixed"
	RetryPolicyBackoffExponential = "exponential"
)

type RetryPolicy struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	MaxRetries    int       `json:"max_retries"`
	BaseTimeoutMs int       `json:"base_timeout_ms"`
	Backoff       string    `json:"backoff"`
	Jitter        float64   `json:"jitter"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (x *RetryPolicy) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "策略名不能为空")
	}
	if x.MaxRetries <= 0 {
		return NewValidationError("max_retries", "最大重试必须为正数")
	}
	if x.BaseTimeoutMs <= 0 {
		return NewValidationError("base_timeout_ms", "基础超时必须为正数")
	}
	x.Backoff = strings.TrimSpace(x.Backoff)
	if x.Backoff == "" {
		return NewValidationError("backoff", "退避策略不能为空")
	}
	if x.Backoff != "" {
		switch x.Backoff {
		case "fixed", "exponential":
		default:
			return NewValidationError("backoff", "退避策略不合法")
		}
	}
	if x.Jitter < 0 {
		return NewValidationError("jitter", "抖动系数不能为负数")
	}
	if x.Status == "" {
		x.Status = RetryPolicyStatusEnabled
	}
	if !RetryPolicyValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func RetryPolicyValidStatus(s string) bool {
	switch s {
	case "enabled", "disabled":
		return true
	default:
		return false
	}
}

var retry_policyTransitions = map[string]map[string]bool{
	RetryPolicyStatusEnabled:  {"disabled": true},
	RetryPolicyStatusDisabled: {"enabled": true},
}

func RetryPolicyCanTransition(from, to string) bool {
	if m, ok := retry_policyTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RetryPolicyFilter struct {
	Backoff string
	Status  string
}

func (f RetryPolicyFilter) Match(x *RetryPolicy) bool {
	if f.Backoff != "" && x.Backoff != f.Backoff {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
