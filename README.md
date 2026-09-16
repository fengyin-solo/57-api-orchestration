# API 编排服务

纯 Go 标准库实现的 API 编排与执行服务，零第三方依赖。

## 运行说明

```bash
cd origin
go build ./...
go run ./cmd/server
```

环境变量：
- `PORT`：监听端口，默认 8080
- `ADDR`：监听地址（优先于 PORT）
- `MAX_PAGE_SIZE`：最大分页大小，默认 100
- `API_KEY`：API 鉴权密钥，默认 `orchestr-api-key`

前端访问：启动后打开 http://localhost:8080/

## 完整 API 表格

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/services | 创建下游服务 |
| GET | /api/services | 服务列表 |
| GET | /api/services/{id} | 服务详情 |
| PUT | /api/services/{id} | 更新服务 |
| DELETE | /api/services/{id} | 删除服务 |
| POST | /api/flows | 创建编排流程 |
| GET | /api/flows | 流程列表 |
| GET | /api/flows/{id} | 流程详情 |
| PUT | /api/flows/{id} | 更新流程 |
| DELETE | /api/flows/{id} | 删除流程 |
| POST | /api/steps | 创建编排步骤 |
| GET | /api/steps | 步骤列表 |
| GET | /api/steps/{id} | 步骤详情 |
| PUT | /api/steps/{id} | 更新步骤 |
| DELETE | /api/steps/{id} | 删除步骤 |
| POST | /api/steps/batch | 批量创建步骤 |
| POST | /api/executions | 创建执行实例 |
| GET | /api/executions | 执行实例列表 |
| GET | /api/executions/{id} | 执行详情 |
| PUT | /api/executions/{id} | 更新执行 |
| DELETE | /api/executions/{id} | 删除执行 |
| POST | /api/executions/{id}/run | 启动执行 |
| POST | /api/executions/{id}/complete | 完成执行 |
| POST | /api/executions/{id}/fail | 失败执行 |
| POST | /api/executions/{id}/timeout | 超时执行 |
| POST | /api/executions/execute-flow | 编排执行流程 |
| POST | /api/executions/batch-execute | 批量执行流程 |
| POST | /api/step-executions | 创建步骤执行 |
| GET | /api/step-executions | 步骤执行列表 |
| GET | /api/step-executions/{id} | 步骤执行详情 |
| PUT | /api/step-executions/{id} | 更新步骤执行 |
| DELETE | /api/step-executions/{id} | 删除步骤执行 |
| POST | /api/step-executions/{id}/transition | 状态迁移 |
| POST | /api/dependencies | 创建步骤依赖 |
| GET | /api/dependencies | 依赖列表 |
| GET | /api/dependencies/{id} | 依赖详情 |
| PUT | /api/dependencies/{id} | 更新依赖 |
| DELETE | /api/dependencies/{id} | 删除依赖 |
| POST | /api/retry-policies | 创建重试策略 |
| GET | /api/retry-policies | 重试策略列表 |
| GET | /api/retry-policies/{id} | 重试策略详情 |
| PUT | /api/retry-policies/{id} | 更新重试策略 |
| DELETE | /api/retry-policies/{id} | 删除重试策略 |
| POST | /api/circuit-breakers | 创建熔断配置 |
| GET | /api/circuit-breakers | 熔断配置列表 |
| GET | /api/circuit-breakers/{id} | 熔断配置详情 |
| PUT | /api/circuit-breakers/{id} | 更新熔断配置 |
| DELETE | /api/circuit-breakers/{id} | 删除熔断配置 |
| POST | /api/parallel-configs | 创建并行配置 |
| GET | /api/parallel-configs | 并行配置列表 |
| GET | /api/parallel-configs/{id} | 并行配置详情 |
| PUT | /api/parallel-configs/{id} | 更新并行配置 |
| DELETE | /api/parallel-configs/{id} | 删除并行配置 |
| POST | /api/conditions | 创建条件分支 |
| GET | /api/conditions | 条件分支列表 |
| GET | /api/conditions/{id} | 条件分支详情 |
| PUT | /api/conditions/{id} | 更新条件分支 |
| DELETE | /api/conditions/{id} | 删除条件分支 |
| POST | /api/audit-logs | 创建审计日志 |
| GET | /api/audit-logs | 审计日志列表 |
| GET | /api/audit-logs/{id} | 审计日志详情 |
| DELETE | /api/audit-logs/{id} | 删除审计日志 |
| GET | /api/stats/overview | 统计概览 |
| GET | /api/stats/success-rate | 成功率 |
| GET | /api/stats/by-flow | 按流程分布 |
| GET | /api/stats/by-service | 按服务分布 |
| GET | /api/stats/failure-reasons | 失败原因分布 |
| GET | /api/export | 全量快照导出 |
| GET | / | 前端首页 |
