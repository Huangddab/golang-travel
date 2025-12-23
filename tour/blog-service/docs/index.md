<!-- 项目总览：blog-service（详版） -->

## 1. 核心技术栈与外部包
- Web：Gin (`github.com/gin-gonic/gin`)
- ORM：GORM + MySQL (`gorm.io/gorm`, `gorm.io/driver/mysql`)
- 校验与翻译：Gin binding + 自定义翻译中间件（validator zh）
- 文档：Swagger (`github.com/swaggo/gin-swagger`, `github.com/swaggo/files`)
- 日志：自定义封装 + `gopkg.in/natefinch/lumberjack.v2`
- 限流：令牌桶 (`github.com/juju/ratelimit`)
- 错误码与响应：`pkg/errcode` + `pkg/app`
- 其他工具：`pkg/convert`（类型转换）、`pkg/setting`（配置加载）

## 2. 分层与调用链（目录/职责）
- 路由与 Handler：`internal/routers`、`internal/routers/api/v1`
- 中间件：`internal/middleware`（日志、Recovery、翻译、限流、超时等）
- 业务层：`internal/service`（DTO + 业务规则）
- 数据访问：`internal/dao`（封装 DB 调用）
- 模型层：`internal/model`（表结构、GORM 回调：时间戳/软删除）
- 公共：`pkg/*`（配置、日志、错误码、响应、限流器等）
- 全局单例：`global/*`（配置、DB、Logger）

调用链：  
HTTP → Gin 路由/中间件 → Handler → Service → DAO → Model/GORM → MySQL

## 3. 启动流程（`main.go`）
1) `init()`  
   - `setupSetting`：读 `configs/config.yaml` 填充 `global.Server/App/Database`，超时转换为 `time.Duration`。  
   - `setupDBEngine`：`model.NewDBEngine` 建立 GORM，配置连接池、SQL 日志（debug）、注册回调（时间戳、软删除），写入 `global.DBEngin`。  
   - `setupLogger`：lumberjack 滚动日志，写入 `global.Logger`。  
2) `main()`  
   - 设置 Gin 模式 → `routers.NewRouters()` 注册路由/中间件 → 构造 `http.Server`（端口/超时来自配置）→ `ListenAndServe`。

## 4. 路由与模块
- Swagger：`GET /swagger/*any`
- 标签：`/api/v1/tags`（List/Create/Update/Delete，已通路 Service/DAO/Model）
- 文章：`/api/v1/articles`（Handler 占位，待接入 Service/DAO/Model）
- 中间件入口：`internal/middleware`
  - `translations`：validator 错误翻译（中文）
  - `recovery`：统一 panic 处理，防止崩溃
  - `limiter`：基于令牌桶的接口级限流（可按 URI）
  - `context_timeout`：为请求注入超时上下文（示例代码已在中间件目录）

## 5. 业务与数据层
- Service（`internal/service`）  
  - 定义请求 DTO（含 binding 校验）  
  - 封装业务流程，调用 DAO  
  - 示例：标签模块 `CountTag/GetTagList/Create/Update/Delete`
- DAO（`internal/dao`）  
  - 仅做数据访问封装，使用 `global.DBEngin`  
  - 示例：`CountTag/GetTagList/CreateTag/UpdateTag/DeleteTag`
- Model（`internal/model`）  
  - 通用 `Model`：`ID/CreatedOn/ModifiedOn/DeletedOn/IsDel`  
  - 回调：创建/更新自动时间戳，删除走软删除（`deleted_on/is_del`）  
  - 表：`Tag`、`Article`、`ArticleTag`

## 6. 配置与日志
- 配置文件：`configs/config.yaml`  
  - `Server`：`RunMode/HttpPort/ReadTimeout/WriteTimeout`  
  - `App`：分页默认/最大、日志路径/文件名/后缀  
  - `Database`：`DBType/Username/Password/Host/DBName/TablePrefix/Charset/ParseTime/MaxIdleConns/MaxOpenConns`
- 日志：lumberjack 按大小/时间切割，路径由 `App.LogSavePath` 等决定；封装在 `pkg/logger`，全局引用 `global.Logger`。

## 7. 错误码与响应
- 错误码：`pkg/errcode` 定义通用（成功、参数错误、服务错误、鉴权、限流）及标签业务错误（增删改查/统计）。  
- 响应：`pkg/app.Response` 统一 `ToResponse`、`ToResponseList`（含 `pager`）、`ToErrorResponse`（`code/msg/details`）。

## 8. 限流与中间件要点
- 限流器接口：`pkg/limiter.LimiterIface`，实现需提供 `Key/GetBucket/AddBuckets`。  
- `MethodLimiter`：按 URI（去掉 query）分桶，可通过中间件在路由组绑定；配置 `LimiterBucketsRule`（填充间隔/容量/量子）。  
- 超时：`context_timeout` 中间件示例可为请求注入超时上下文。  
- Recovery：`recovery` 中间件兜底 panic，避免服务崩溃。  
- 翻译：`translations` 注入 validator 翻译器，使 `binding` 校验错误返回中文提示。

## 9. 文档与开发辅助
- Swagger：`docs/` 目录包含生成的 `swagger.yaml/json`，`main.go` 通过 `_ "blog-service/docs"` 注册文档。  
- 调试：`/swagger/index.html` 直连当前服务端口，无需硬编码 host。  
- 示例调用（标签）：`GET /api/v1/tags?page=1&page_size=10&state=1`

## 10. 待完善与可扩展
- 文章模块：补齐 Service/DAO/Model 增删改查与列表/统计（可参考标签实现）。  
- 鉴权/认证：目前未接入，可按需增加 JWT 或 AppKey/Auth 接口。  
- 限流策略：可按 IP/用户/URI 组合 key；可在中间件扩展白名单、熔断/降级策略。  
- 监控与追踪：可后续接入 Prometheus/Tracing。  
- 配置分环境：可增加 dev/prod 配置文件与装载策略。