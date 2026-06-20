# goworks-product

一个基于 Go + Iris 的商品订单管理后端学习项目，包含用户、商品、订单三个基础模块，并补充了轻量级接口自动化测试和 CI 检查。

## 技术栈

- 后端：Go、Iris
- 数据库：MySQL
- 缓存：Redis
- 测试：Python 标准库接口冒烟测试
- CI：GitHub Actions

## 项目结构

```text
backend/        Web 启动入口、控制器、页面模板和静态资源
common/         MySQL、Redis、请求参数解析和统一返回结构
datamodels/     用户、商品、订单数据模型
repositories/   数据访问层
services/       业务逻辑层
tests/          接口测试用例、冒烟测试脚本和测试报告
.github/        GitHub Actions CI 配置
```

## 已实现功能

- 用户模块：登录校验、新增用户、查询用户、更新用户、删除用户。
- 商品模块：商品列表查询、商品创建、商品更新、Redis 缓存和缓存失效。
- 订单模块：订单创建、订单查询、用户和商品依赖校验。
- 测试补充：围绕用户、商品、订单接口补充参数校验和异常响应类冒烟测试。

## 运行方式

项目默认从环境变量读取 MySQL 和 Redis 配置；未设置时使用本地开发默认值。

```bash
go run ./backend
```

常用环境变量：

```text
MYSQL_HOST
MYSQL_PORT
MYSQL_USER
MYSQL_PASSWORD
MYSQL_DATABASE
REDIS_HOST
REDIS_PORT
```

## 接口自动化测试

启动后端服务后运行：

```bash
python tests/api_smoke_test.py
```

当前测试重点覆盖不依赖固定数据库种子数据的场景：

- 用户登录缺少用户名或密码
- 用户创建缺少必要字段
- 商品更新缺少 ID 或库存格式错误
- 订单创建缺少用户 ID 或商品 ID

详细用例和测试报告见：

- `tests/api_test_cases.md`
- `tests/test_report.md`

## CI 说明

`.github/workflows/api-smoke.yml` 会在 push / pull request 时执行：

- `go mod download`
- `go build ./...`
- `python -m py_compile tests/api_smoke_test.py`

由于当前项目还没有完整的数据库 schema 和 seed 数据，CI 暂不启动 MySQL / Redis 运行真实接口请求。后续补齐初始化脚本后，可以将 `python tests/api_smoke_test.py` 接入 CI 运行时检查。
