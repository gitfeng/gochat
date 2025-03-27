# GoChat 技术设计文档
## 一、系统概述
```text
├── 实时聊天系统 (WebSocket)
├── RESTful API 服务 (Gin)
├── 智能对话引擎 (状态机 + 规则配置)
├── 情感分析集成 (腾讯云 NLP)
└── 日志追踪打点
 ```

## 二、技术选型 分类 技术栈 框架

Gin (HTTP)、Gorilla WebSocket、GORM (ORM) 配置管理

Viper 部署监控

Docker、数据库

MySQL (生产) 云服务

腾讯云 NLP (情感分析)

Testify、GoMock

## 三、核心模块设计
### 1. WebSocket 通信模块
处理流程：

```go
连接升级 → 会话验证 → 消息循环 → 持久化存储 → 业务处理 → 响应推送
```

关键类：

- `ServeWebSocket`
- `ChatBotEngine`

### 2. 规则引擎模块
```yaml
# chatbot_rules.yml 结构
状态机 → 意图识别 → 上下文管理 → 响应模板 → 异常处理
```

### 3. 认证中间件

```go
JWT解析 → 会话验证 → 权限上下文设置 → 超时控制
 ```

### 4. 配置管理系统

```bash
配置文件加载流程：
viper初始化 → 路径设置 → 类型绑定 → 热更新监听（预留）
```

## 四、数据模型设计
### 1. 消息实体
```go
type Message struct {
    ID          uint      `gorm:"primaryKey"`
    CustomerID  uint64    `gorm:"index;not null"`  
    Message     string    `gorm:"type:text;not null"`
    Sender      string    `gorm:"size:20;not null"`
    MessageType string    `gorm:"size:20;default:'normal'"`
    CreatedAt   time.Time
}
 ```

### 2. 反馈实体
```go
type Feedback struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	CustomerID uint64 `gorm:"index"`
	Score      uint   // requested, completed
	Comment    string `gorm:"type:text"`
	Sentiment  int    `gorm:"type:tinyint;default:0" json:"sentiment"`

	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

 ```

## 五、API 设计

### 1. RESTful API
| 端点               | 方法   | 参数                  | 请求示例                          | 描述                     |
|--------------------|--------|-----------------------|-----------------------------------|------------------------|
| `/healthcheck`     | GET    | -                     | `curl http://localhost:8080/healthcheck` | 服务健康检查            |
| `/message/list`    | GET    | `customer_id`         | `?customer_id=1&page=2`           | 分页获取消息记录        |

### 2. WebSocket 接口
```text
ws://host:port/ws?token=<JWT>

协议规范：
1. 连接需携带有效 JWT Token
2. 心跳机制：每30秒发送空消息维持连接

### 3. 认证机制
```http
请求头示例：

### 4. 错误代码表

完整错误代码参见：<mcsymbol name="ErrorCode" filename="errcode.go" path="d:\workspace\gochat\internal\service\errcode.go" startline="15" type="class"></mcsymbol>
```go
const (
	ErrCodeSuccess        = 0
	ErrCodeInternalServer = 500
	ErrCodeInvalidRequest = 400
	ErrCodeUnauthorized   = 401
	ErrCodeForbidden      = 403
	ErrCodeNotFound       = 404
	ErrCodeConflict       = 409
	// 可以根据业务需求添加更多自定义错误码
	ErrCodeUserExists      = 1001
	ErrCodeUserNotFound    = 1002
	ErrCodeInvalidPassword = 1003
)
```

## 五、可改进方向
1. 消息队列优化：
   - 异步处理：使用消息队列如 RabbitMQ、Kafka 处理消息发送、存储等操作，提高系统性能。
   - 场景：用户消息的情绪分析，可以在用户留言出发feedback消息后，发送mq，下游消费分析情绪后把结果再落库，异步处理，减少用户等待时间。
2. 引入缓存：
   - 会话缓存：存储用户会话状态，减少数据库查询
   - 场景：全局化请求限流信息存储，防止恶意请求。
   - 场景：存储websocket连接会话id，用于后续对话。
3. 分布式部署：
   - 负载均衡：使用 Nginx 等负载均衡器
   - 集群化：使用 Kubernetes 等容器编排工具
4. 引入AI增强：
   - 智能客服：基于规则引擎和 AI 模型，提供智能客服服务
   - 场景：基于 AI 模型，提供智能回复功能
5. 情感分析集成：
    - 初创项目建议直接使用腾讯云API
    - 数据敏感场景推荐HuggingFace+Flask API+Go调用
    - 高并发简单场景可用SnowNLP+Python微服务
6. 安全优化：
   - 数据加密：对敏感数据进行加密存储
   - 访问控制：使用 RBAC 等权限管理机制
   - 安全审计：记录用户操作日志
7. 数据仓库建设：
   - 数据仓库：使用数据仓库工具如 ClickHouse、Druid 等
   - 场景：存储用户消息、会话等数据，用于数据分析和报表生成。
8. 性能优化：
    - 控制websocket保持连接的时间：减少无用连接的资源占用
    - 控制websocket连接的数量：防止连接过多导致服务器负载过高
    - 控制websocket连接的频率：防止连接过快导致服务器负载过高

## 六、部署与监控

### 本地安装

#### 安装依赖
go mod tidy
go build cmd/server.go

#### 配置

#### 数据库
- 安装MySQL
- 导入数据库文件，config/sql/init_db.mysql.sql
- 修改配置文件 config/config.yml

#### 启动服务
- go run cmd/server.go

### 基于 docker 安装【由于环境问题，docker安装并没有测试】
#### 1. 构建镜像（在项目根目录执行）
```shell
docker build -t gochat:latest -f Dockerfile .
```

#### 2. 运行容器（带环境变量）
```shell
docker run -d \
  -p 8080:8080 \
  -e DB_DSN="root:gochat@tcp(host.docker.internal:3306)/gochat" \
  -e CONFIG_DIR=/app/config \
  -v d:\workspace\gochat\config:/app/config \
  --name gochat \
  gochat:latest
```

#### 3.全栈启动（包含MySQL）
```shell
docker-compose -f d:\workspace\gochat\docker-compose.yml up -d
# 启动服务后验证表结构
docker exec -it gochat-mysql mysql -ugochat -pgochat -e "USE gochat; SHOW TABLES;"

#### 查看实时日志
docker logs -f gochat

```

### 单元自动化测试

####  运行全部测试
sh scripts/test.sh

####  运行特定测试
go test -v ./internal/handler -run TestCreateCustomer

### 接口测试
#### 验证服务状态
curl http://localhost:8080/healthcheck?full=1

#### websocket发送消息
说明：这里的token是经过加密的用户的id，需要在数据库customers表中创建一个id为1的用户，然后在websocket连接时传递token参数。以这种方式简单实现登录用户信息，以及鉴权操作。
```shell
wscat -c "ws://localhost:8080/ws?token=5XRyxivsjCvCp75cDVbVgUf8jdzhYH1wxHexRWo="
> nihao
< 抱歉，我还在学习中，暂时无法回答这个问题
> hello
< 您好，我是${bot_name}，请问需要什么帮助？
> feedback: sing nice
< Tanks for your feedback, we will deal in 3 days.
```
如果本机没有wscat，需要安装
```shell
npm install -g wscat
```

#### 获取message/list

- 按客户过滤
curl "http://localhost:8080/message/list?customer_id=1&page=1&limit=5"
