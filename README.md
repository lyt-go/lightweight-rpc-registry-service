# 轻量 RPC 框架

纯 Go 标准库（`net/http`）实现的后端服务，零第三方依赖，标准分层（cmd/internal/pkg），开箱即跑。

## 运行

```bash
cd origin
go run ./cmd/server
# 默认监听 :8080，可用 PORT/ADDR/MAX_PAGE_SIZE 环境变量覆盖
```

## 业务实体

Service（RPC 服务注册）、Method（RPC 方法定义）、Node（机器节点注册）、CallLog（调用记录）、Codec（编解码器）、RetryPolicy（超时重试策略）

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`。

| 模块 | 接口 | 说明 |
|------|------|------|
| services | POST /api/services | 创建Service |
| services | GET /api/services | 列表查询（分页+筛选） |
| services | GET /api/services/{id} | 详情 |
| services | PUT /api/services/{id} | 更新 |
| services | DELETE /api/services/{id} | 删除 |
| services | PATCH /api/services/{id}/status | 状态流转 |
| methods | POST /api/methods | 创建Method |
| methods | GET /api/methods | 列表查询（分页+筛选） |
| methods | GET /api/methods/{id} | 详情 |
| methods | PUT /api/methods/{id} | 更新 |
| methods | DELETE /api/methods/{id} | 删除 |
| methods | PATCH /api/methods/{id}/status | 状态流转 |
| nodes | POST /api/nodes | 创建Node |
| nodes | GET /api/nodes | 列表查询（分页+筛选） |
| nodes | GET /api/nodes/{id} | 详情 |
| nodes | PUT /api/nodes/{id} | 更新 |
| nodes | DELETE /api/nodes/{id} | 删除 |
| nodes | PATCH /api/nodes/{id}/status | 状态流转 |
| call-logs | POST /api/call-logs | 创建CallLog |
| call-logs | GET /api/call-logs | 列表查询（分页+筛选） |
| call-logs | GET /api/call-logs/{id} | 详情 |
| call-logs | DELETE /api/call-logs/{id} | 删除 |
| codecs | POST /api/codecs | 创建Codec |
| codecs | GET /api/codecs | 列表查询（分页+筛选） |
| codecs | GET /api/codecs/{id} | 详情 |
| codecs | PUT /api/codecs/{id} | 更新 |
| codecs | DELETE /api/codecs/{id} | 删除 |
| codecs | PATCH /api/codecs/{id}/status | 状态流转 |
| retry-policies | POST /api/retry-policies | 创建RetryPolicy |
| retry-policies | GET /api/retry-policies | 列表查询（分页+筛选） |
| retry-policies | GET /api/retry-policies/{id} | 详情 |
| retry-policies | PUT /api/retry-policies/{id} | 更新 |
| retry-policies | DELETE /api/retry-policies/{id} | 删除 |
| retry-policies | PATCH /api/retry-policies/{id}/status | 状态流转 |
| stats | GET /api/stats/overview | 全局统计总览 |

## 说明

- 数据存储为内存实现，重启即清空。
- 金额/时间等数值字段统一采用整数（分 / 秒 / 毫秒 / Unix 时间戳），避免浮点精度问题。
