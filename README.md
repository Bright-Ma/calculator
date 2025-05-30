# 小学生口算题系统

## 快速开始

### 启动服务

```bash
go run main.go
```

服务将启动在 `http://localhost:8080`

### API使用示例

**获取题目**:
```bash
curl "http://localhost:8080/api/questions?difficulty=2"
```

响应示例:
```json
{
  "question": "12 × 5",
  "difficulty": 2
}
```

**提交答案**:
```bash
curl -X POST "http://localhost:8080/api/answers" \
  -H "Content-Type: application/json" \
  -d '{"question_id":"123", "answer":60}'
```

响应示例:
```json
{
  "correct": true,
  "score": 10,
  "message": "回答正确!"
}
```

### 运行测试

```bash
go test -v ./internal/drill/...
```

## 项目概述

小学生口算题系统是一个专注于提升小学生口算能力的在线练习平台。系统根据学生能力提供不同难度的题目，记录学习历史并生成成绩统计，同时通过热度排行榜激发学习兴趣。

## 主要功能

### 1. 题目难度分级
- **低难度**：10以内加减法
- **中难度**：20以内加减法、简单乘法
- **高难度**：100以内加减法、乘除法混合运算

### 2. 学习历史记录
- 使用MySQL存储每位学生的做题记录
- 记录包括：题目内容、答题结果、用时、时间戳
- 可按日期、难度筛选查看历史记录

### 3. 成绩统计分析
- 正确率统计
- 各难度题目掌握程度
- 生成可视化学习报告

### 4. 用户认证系统
- 学生/教师账号注册与登录
- 基于JWT的身份验证
- 密码加密存储
- 会话管理

### 5. 热度排行榜
- 基于时间戳和做题数量计算热度值
- 使用Redis存储和实时更新排行榜
- 每日/每时热度榜单
- 热度算法：`热度 = 做题数量 * 时间衰减因子`

## 技术架构

### 后端
- 编程语言：Go
- Web框架：Gin
- 数据库：MySQL（主存储）
- 缓存：Redis（排行榜）
- 认证：JWT
- ORM：GORM

### 前端
- HTML5/CSS3/JavaScript
- 响应式设计
- 图表库：Chart.js（数据可视化）

## 安装指南

### 系统要求
- Go 1.16+
- MySQL 5.7+
- Redis 6.0+
- Node.js 14+ (前端开发)

### 安装步骤

1. 克隆代码库
```bash
git clone https://github.com/Bright-Ma/calculator.git
cd calculator
```

2. 安装依赖
```bash
go mod download
```


4. 运行应用
```bash
go run .
```

5. 访问应用
```
http://localhost:8080
```

## 使用指南

### 学生使用流程
1. 注册/登录账号
2. 选择难度级别开始练习
3. 完成练习后查看成绩和错题
4. 查看个人学习历史和准确率
5. 查看热度排行榜


## 项目结构
```
calculator/
├── internal/             # 后端代码
│   ├── handlers/         # HTTP请求处理器
│   ├── redis/            # Redis缓存
│   ├── model/            # 数据模型
│   ├── router/           # 路由配置
│   ├── middleware/       # 中间件
│   ├── drill/            # 口算题生成
│   └── database/         # 数据库操作
├── main.go               # 程序入口
├── go.mod                # Go模块文件
├── go.sum                # Go依赖版本锁定
├── redis.conf            # Redis配置文件
├── frontend/             # 前端资源
```