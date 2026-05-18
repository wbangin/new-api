# New Features - Token 级请求/响应日志记录

## 功能概述

新增按 Token（API Key）维度控制的请求/响应内容记录功能。开启后，该 Token 的所有经过 relay 层的请求（包括聊天、图片、音频、embedding、rerank 等）都会被完整记录到本地日志文件中，并可在前端"使用日志"详情中查看。

## 核心特性

### 1. 按 Token 维度开关控制

- 每个 Token（API Key）可独立配置是否记录请求/响应内容
- 在 Token 编辑页面的"高级设置"中新增 **"Log Request Content"** 开关
- 默认关闭，仅对显式开启的 Token 生效

### 2. 请求/响应内容记录

- **记录内容**：完整的请求体（request body）+ 响应体（response body，非流式时）
- **存储格式**：JSONL（每行一个 JSON 对象，便于追加和检索）
- **存储位置**：`{log_dir}/request_logs/YYYY-MM-DD.jsonl`（按日期自动分割）
- **异步写入**：使用 goroutine 池异步写入，不阻塞请求处理

### 3. 前端查看入口

- 在"使用日志"详情对话框中，Request ID 行后新增 **"View Request/Response Detail"** 链接按钮
- 点击后弹出独立对话框，展示格式化的请求体和响应体 JSON

## 日志条目结构

```json
{
  "request_id": "req-xxxxxxxx",
  "timestamp": 1716000000,
  "token_id": 1,
  "token_name": "my-token",
  "user_id": 100,
  "model": "gpt-4",
  "relay_mode": 1,
  "request_body": "{...}",
  "response_body": "{...}",
  "status_code": 200,
  "is_stream": false
}
```

## 新增/修改的文件清单

### 后端

| 文件 | 操作 | 说明 |
|------|------|------|
| `model/token.go` | 修改 | Token 结构体新增 `LogRequestEnabled` 字段 |
| `middleware/auth.go` | 修改 | 将 `LogRequestEnabled` 写入 Gin Context |
| `relay/common/relay_info.go` | 修改 | RelayInfo 新增 `LogRequestEnabled` 和 `ResponseBody` 字段 |
| `controller/relay.go` | 修改 | 请求成功后异步调用日志记录 |
| `controller/token.go` | 修改 | Token CRUD 支持新字段 |
| `service/request_log.go` | **新增** | 日志写入/搜索核心服务 |
| `service/request_log_test.go` | **新增** | 单元测试 |
| `controller/request_log_controller.go` | **新增** | `GET /api/log/request_detail` 接口 |
| `router/api-router.go` | 修改 | 注册新路由 |

### 前端

| 文件 | 操作 | 说明 |
|------|------|------|
| `web/default/src/features/usage-logs/api.ts` | 修改 | 新增 `getRequestDetail` API |
| `web/default/src/features/usage-logs/components/dialogs/details-dialog.tsx` | 修改 | 添加"查看详情"链接按钮 |
| `web/default/src/features/usage-logs/components/dialogs/request-detail-dialog.tsx` | **新增** | 请求详情对话框组件 |
| `web/default/src/features/keys/types.ts` | 修改 | 类型定义新增字段 |
| `web/default/src/features/keys/lib/api-key-form.ts` | 修改 | 表单 schema 新增字段 |
| `web/default/src/features/keys/components/api-keys-mutate-drawer.tsx` | 修改 | Token 编辑表单添加开关 |

## API 接口

### GET /api/log/request_detail

获取指定 request_id 的完整请求/响应内容。

**参数**：
- `request_id` (query, required) — 请求 ID

**权限**：
- 需要登录认证
- 管理员可查看所有日志
- 普通用户仅可查看自己 Token 的日志

**响应示例**：
```json
{
  "success": true,
  "data": {
    "request_id": "req-xxxxxxxx",
    "timestamp": 1716000000,
    "token_id": 1,
    "token_name": "my-token",
    "user_id": 100,
    "model": "gpt-4",
    "relay_mode": 1,
    "request_body": "{\"messages\":[...]}",
    "response_body": "{\"choices\":[...]}",
    "status_code": 200,
    "is_stream": false
  }
}
```

## 日志文件存储详细说明

### 存储位置

日志文件存储在启动参数 `--log-dir` 指定目录下的 `request_logs/` 子目录中：

```
{log-dir}/request_logs/YYYY-MM-DD.jsonl
```

**默认路径示例**（未指定 `--log-dir` 时默认为 `./logs`）：

```
./logs/request_logs/2026-05-18.jsonl
./logs/request_logs/2026-05-19.jsonl
./logs/request_logs/2026-05-20.jsonl
```

### 配置方式

| 配置项 | 说明 |
|--------|------|
| `--log-dir` | 启动参数，指定日志基础目录（默认 `./logs`） |
| 子目录 | 固定为 `request_logs/`，自动创建 |
| 文件名格式 | `YYYY-MM-DD.jsonl`，按日期自动分割 |
| 文件格式 | JSONL（每行一个独立的 JSON 对象） |

### 启动示例

```bash
# 使用默认日志目录 ./logs
./new-api

# 自定义日志目录
./new-api --log-dir /var/log/new-api

# 日志文件将存储在 /var/log/new-api/request_logs/ 下
```

### 搜索机制

- API 接口 `GET /api/log/request_detail?request_id=xxx` 会从**今天开始往前搜索最近 7 天**的日志文件
- 搜索使用快速字符串匹配 + JSON 反序列化确认，性能较好
- 单条日志最大支持 **10MB**（覆盖超长请求/响应场景）

### 日志文件管理

- **自动分割**：每天自动生成新的日志文件，不会追加到旧文件
- **无自动清理**：系统不会自动删除旧日志文件
- **手动清理**：可直接删除旧的 `.jsonl` 文件，不影响系统运行
- **推荐方案**：配合 logrotate 或 cron 定时任务定期清理

**logrotate 配置示例**：

```
/path/to/logs/request_logs/*.jsonl {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
}
```

**cron 清理示例**（保留最近 30 天）：

```bash
0 2 * * * find /path/to/logs/request_logs/ -name "*.jsonl" -mtime +30 -delete
```

### 存储空间估算

日志大小取决于请求/响应体的大小：

| 场景 | 单条日志大小（约） | 1000 条/天 |
|------|------------------|------------|
| 普通聊天请求 | 2-5 KB | 2-5 MB |
| 长对话（含历史消息） | 10-50 KB | 10-50 MB |
| 图片生成（含 base64） | 100 KB - 1 MB | 100 MB - 1 GB |

## 注意事项

- 开启记录会增加磁盘 I/O 和存储占用，建议仅对需要调试的 Token 开启
- 流式响应（SSE）目前仅记录请求体，响应体字段为空（后续可扩展）
- 日志文件无自动清理机制，需要运维定期清理或配置 logrotate
- 日志文件可以在运行时安全删除，不会影响当前正在写入的文件（当天的文件除外）
- 如果磁盘空间不足，日志写入会静默失败（错误记录到系统日志），不会影响正常请求处理
