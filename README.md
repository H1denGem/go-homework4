# Go 个人博客系统

基于 Gin 框架和 GORM 的个人博客系统后端，实现用户认证、文章管理和评论功能。

## 运行环境

- Go 1.25.4+
- SQLite3 数据库（无需额外安装，Go 驱动内置）

## 项目结构

```
blog/
├── config/          # 配置文件
│   ├── config.go
│   └── config.yaml
├── handlers/        # HTTP 处理器
│   ├── user_handler.go
│   ├── post_handler.go
│   └── comment_handler.go
├── middlewares/     # 中间件
│   ├── auth.go      # JWT 认证中间件
│   ├── cors.go      # 跨域中间件
│   └── logger.go    # 日志中间件
├── models/          # 数据模型
│   ├── user.go
│   ├── post.go
│   ├── comment.go
│   └── request.go   # 请求体类型
├── services/        # 业务逻辑层
│   ├── user_service.go
│   ├── post_service.go
│   └── comment_service.go
├── utils/           # 工具函数
│   ├── errors.go    # 统一错误处理
│   ├── jwt.go       # JWT 工具
│   └── response.go  # 响应格式
├── logs/            # 日志文件（运行时生成）
│   ├── app.log      # 应用日志
│   └── access.log   # 访问日志
├── main.go          # 程序入口
├── go.mod           # 依赖管理
└── go.sum
```

## 依赖安装

### 1. 克隆项目

```bash
cd d:/GolandProjects/go-homework4/blog
```

### 2. 安装依赖

```bash
go mod tidy
```

或手动安装：

```bash
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v5
go get github.com/spf13/viper
go get golang.org/x/crypto
```

## 启动方式

### 开发模式启动

```bash
go run main.go
```

### 编译后启动

```bash
go build -o blog.exe
./blog.exe  # Windows: blog.exe
```

### 服务启动信息

- 默认地址：`http://0.0.0.0:8080`
- 配置文件：`config/config.yaml`
- 日志目录：`logs/`

## API 接口文档

### 基础信息

- 基础路径：`/api/v1`
- 认证方式：JWT Bearer Token

### 公开路由（无需认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/register` | 用户注册 |
| POST | `/api/v1/login` | 用户登录 |
| GET | `/api/v1/posts` | 获取文章列表 |
| GET | `/api/v1/posts/:id` | 获取单篇文章 |
| GET | `/api/v1/posts/:id/comments` | 获取文章评论列表 |
| GET | `/api/v1/posts/:id/comments/:comment_id` | 获取单条评论 |
| GET | `/health` | 健康检查 |

### 需认证的路由

需要在请求头中添加：`Authorization: Bearer <token>`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/posts` | 创建文章 |
| PUT | `/api/v1/posts/:id` | 更新文章（仅作者） |
| DELETE | `/api/v1/posts/:id` | 删除文章（仅作者） |
| POST | `/api/v1/posts/:id/comments` | 创建评论 |
| PUT | `/api/v1/posts/:id/comments/:comment_id` | 更新评论（仅作者） |
| DELETE | `/api/v1/posts/:id/comments/:comment_id` | 删除评论（仅作者） |

---

## 测试用例

### 使用 Postman 测试

#### 1. 用户注册

**请求：**
```
POST http://localhost:8080/api/v1/register
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456",
  "email": "test@example.com"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "message": "注册成功"
  }
}
```

**失败响应（用户已存在）：**
```json
{
  "code": 400,
  "message": "用户名已存在",
  "data": null
}
```

#### 2. 用户登录

**请求：**
```
POST http://localhost:8080/api/v1/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "123456"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

**失败响应（密码错误）：**
```json
{
  "code": 401,
  "message": "密码错误",
  "data": null
}
```

#### 3. 创建文章

**请求：**
```
POST http://localhost:8080/api/v1/posts
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "我的第一篇文章",
  "content": "这是文章内容..."
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "文章创建成功"
  }
}
```

#### 4. 获取文章列表

**请求：**
```
GET http://localhost:8080/api/v1/posts
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "title": "我的第一篇文章",
      "content": "这是文章内容...",
      "user_id": 1,
      "created_at": "2026-03-09T10:30:00+08:00",
      "updated_at": "2026-03-09T10:30:00+08:00",
      "user": {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com"
      }
    }
  ]
}
```

#### 5. 获取单篇文章

**请求：**
```
GET http://localhost:8080/api/v1/posts/1
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "title": "我的第一篇文章",
    "content": "这是文章内容...",
    "user_id": 1,
    "created_at": "2026-03-09T10:30:00+08:00",
    "updated_at": "2026-03-09T10:30:00+08:00",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    }
  }
}
```

#### 6. 更新文章

**请求：**
```
PUT http://localhost:8080/api/v1/posts/1
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "更新后的文章标题",
  "content": "更新后的文章内容..."
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "文章更新成功"
  }
}
```

**失败响应（无权操作）：**
```json
{
  "code": 403,
  "message": "权限不足",
  "data": null
}
```

#### 7. 删除文章

**请求：**
```
DELETE http://localhost:8080/api/v1/posts/1
Authorization: Bearer <token>
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "文章删除成功"
  }
}
```

#### 8. 创建评论

**请求：**
```
POST http://localhost:8080/api/v1/posts/1/comments
Content-Type: application/json
Authorization: Bearer <token>

{
  "content": "好文章！"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "评论创建成功"
  }
}
```

#### 9. 获取文章评论列表

**请求：**
```
GET http://localhost:8080/api/v1/posts/1/comments
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": 1,
      "content": "好文章！",
      "post_id": 1,
      "user_id": 1,
      "created_at": "2026-03-09T10:35:00+08:00",
      "user": {
        "id": 1,
        "username": "testuser",
        "email": "test@example.com"
      }
    }
  ]
}
```

#### 10. 更新评论

**请求：**
```
PUT http://localhost:8080/api/v1/posts/1/comments/1
Content-Type: application/json
Authorization: Bearer <token>

{
  "content": "更新后的评论内容"
}
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "评论更新成功"
  }
}
```

#### 11. 删除评论

**请求：**
```
DELETE http://localhost:8080/api/v1/posts/1/comments/1
Authorization: Bearer <token>
```

**成功响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "message": "评论删除成功"
  }
}
```

---

## 使用 curl 测试

### 1. 注册用户
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser\",\"password\":\"123456\",\"email\":\"test@example.com\"}"
```

### 2. 登录获取 Token
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser\",\"password\":\"123456\"}"
```

### 3. 创建文章（替换 <token>）
```bash
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d "{\"title\":\"我的文章\",\"content\":\"文章内容\"}"
```

### 4. 获取文章列表
```bash
curl http://localhost:8080/api/v1/posts
```

### 5. 创建评论
```bash
curl -X POST http://localhost:8080/api/v1/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d "{\"content\":\"好文章！\"}"
```

---

## 错误处理

系统使用统一的错误响应格式：

```json
{
  "code": <HTTP状态码>,
  "message": "<错误信息>",
  "data": null
}
```

### 常见错误码

| 错误码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权/令牌无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 日志系统

### 日志文件

| 文件 | 说明 |
|------|------|
| `logs/app.log` | 应用日志（数据库、服务、认证等）|
| `logs/access.log` | HTTP 请求日志（访问记录）|

### 日志输出

日志同时输出到：
- **终端** - 开发调试时实时查看
- **文件** - 持久化存储

---

## 技术栈

- **Web 框架**：Gin
- **ORM**：GORM
- **数据库**：SQLite3
- **认证**：JWT
- **配置管理**：Viper
- **密码加密**：bcrypt

---

## 安全特性

- 密码使用 bcrypt 加密（自动加盐）
- JWT 令牌认证
- 权限控制（文章/评论只能由作者修改删除）
- SQL 注入防护（GORM 参数化查询）
- CORS 跨域支持
