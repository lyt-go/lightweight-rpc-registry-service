// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"lightweightrpc/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateService(x *model.Service) error
	GetService(id string) (*model.Service, error)
	GetServiceByName(v string) (*model.Service, error)
	ListServices() []*model.Service
	UpdateService(x *model.Service) error
	DeleteService(id string) error
	CreateMethod(x *model.Method) error
	GetMethod(id string) (*model.Method, error)
	GetMethodByName(v string) (*model.Method, error)
	ListMethods() []*model.Method
	UpdateMethod(x *model.Method) error
	DeleteMethod(id string) error
	CreateNode(x *model.Node) error
	GetNode(id string) (*model.Node, error)
	GetNodeByAddress(v string) (*model.Node, error)
	ListNodes() []*model.Node
	UpdateNode(x *model.Node) error
	DeleteNode(id string) error
	CreateCallLog(x *model.CallLog) error
	GetCallLog(id string) (*model.CallLog, error)
	ListCallLogs() []*model.CallLog
	DeleteCallLog(id string) error
	CreateCodec(x *model.Codec) error
	GetCodec(id string) (*model.Codec, error)
	GetCodecByName(v string) (*model.Codec, error)
	ListCodecs() []*model.Codec
	UpdateCodec(x *model.Codec) error
	DeleteCodec(id string) error
	CreateRetryPolicy(x *model.RetryPolicy) error
	GetRetryPolicy(id string) (*model.RetryPolicy, error)
	GetRetryPolicyByName(v string) (*model.RetryPolicy, error)
	ListRetryPolicys() []*model.RetryPolicy
	UpdateRetryPolicy(x *model.RetryPolicy) error
	DeleteRetryPolicy(id string) error
}
