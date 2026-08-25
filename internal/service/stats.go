package service

import (
	"lightweightrpc/internal/model"
)

// Overview 返回全局统计总览：每个实体的总数与按状态分布。
func (s *Service) Overview() map[string]interface{} {
	return map[string]interface{}{
		"services":       s.countByStatus("Service", s.store.ListServices()),
		"methods":        s.countByStatus("Method", s.store.ListMethods()),
		"nodes":          s.countByStatus("Node", s.store.ListNodes()),
		"call-logs":      s.countByStatus("CallLog", s.store.ListCallLogs()),
		"codecs":         s.countByStatus("Codec", s.store.ListCodecs()),
		"retry-policies": s.countByStatus("RetryPolicy", s.store.ListRetryPolicys()),
	}
}

func (s *Service) countByStatus(kind string, list interface{}) map[string]interface{} {
	switch items := list.(type) {
	case []*model.Service:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.Method:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.Node:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.CallLog:
		return map[string]interface{}{"total": len(items)}
	case []*model.Codec:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.RetryPolicy:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	}
	return map[string]interface{}{"total": 0, "by_status": map[string]int{}}
}
