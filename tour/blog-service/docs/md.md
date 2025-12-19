<!-- 项目整体流程梳理 -->

## 一、基础组件

1. **错误码标准化**
   - `pkg/errcode` 定义统一的业务错误码（如：成功、入参错误、服务内部错误、标签相关错误等）。
   - 对外暴露 `Error` 类型，提供 `Code()`、`Message()`、`StatusCode()`、`WithDetails()` 等方法。
   - `pkg/app.Response.ToErrorResponse` 使用这些错误码，输出统一格式：`{code, msg, details?}`。

2. **配置管理**
   - 配置文件：`configs/config.yaml`，包含 `Server`、`App`、`Database` 三部分。
   - 配置模块：`pkg/setting`，通过 `NewSetting()` 加载 YAML，再用 `ReadSection("Server", &global.ServerSetting)` 等方式填充结构体。
   - 全局配置对象：`global.ServerSetting`、`global.AppSetting`、`global.DatabaseSetting`，在 `init()` 里完成赋值。

3. **数据库连接**
   - 在 `main.go` 中调用 `setupDBEngine()`，内部使用 `model.NewDBEngine(global.DatabaseSetting)`。
   - `model.NewDBEngine`：
     - 根据配置拼接 MySQL DSN。
     - `gorm.Open(mysql.Open(dsn), &gorm.Config{})` 建立连接。
     - 设置连接池参数（`MaxIdleConns`、`MaxOpenConns`）。
     - `RunMode == "debug"` 时开启 SQL 日志。
     - 调用 `setupCallbacks(db)` 注册 GORM 回调（创建/更新/删除时自动维护时间戳和软删除字段）。
   - 最终将 `*gorm.DB` 保存到 `global.DBEngin`，供 DAO 层统一使用。

4. **日志写入**
   - 在 `main.go` 的 `setupLogger()` 中，使用 `pkg/logger` + `lumberjack` 创建带滚动策略的日志实例。
   - 日志文件路径来自 `global.AppSetting`（如：`storage/logs/app.log`）。
   - 最终写入 `global.Logger`，供各层通过 `global.Logger.Errorf` 等方法记录日志。

5. **响应处理 & 分页**
   - 响应工具：`pkg/app/app.go`。
     - `Response.ToResponse(data)`：正常返回，空数据自动转为 `{}`。
     - `Response.ToResponseList(list, totalRows)`：统一列表返回格式：`{list, pager}`。
     - `Response.ToErrorResponse(err *errcode.Error)`：统一错误返回格式。
   - 分页工具：`pkg/app/pagination.go`。
     - `GetPage(c)`、`GetPageSize(c)` 从查询参数中解析页码和每页数量，结合 `AppSetting.DefaultPageSize`、`AppSetting.MaxPageSize` 做限制。
     - `GetPageOffset(page, pageSize)` 用于 DAO 层构造 `Offset`。

---

## 二、分层结构与调用关系

目录简化如下（只列核心逻辑相关）：

- `internal/routers`：接口层（路由 + HTTP Handler）
- `internal/service`：业务层（Service，承接 Handler，封装业务规则）
- `internal/dao`：数据访问层（DAO，负责所有数据库读写）
- `internal/model`：数据模型层（Model，定义表结构和通用模型逻辑）
- `internal/middleware`：中间件（如参数校验错误翻译）

各层调用方向：

业务层 (Service)  
　　↓（调用）  
数据访问层 (DAO) —— 使用 `global.DBEngin` 操作数据库  
　　↓（调用）  
模型层 (Model) —— 基于 GORM 的表结构和通用逻辑  
　　↓（执行 SQL）  
数据库 (DB)



接口层 (Handler) 则位于最外层，接受 HTTP 请求、组装参数后调用 Service：

HTTP 请求  
　　↓  
Gin 路由 & Handler（`internal/routers/api/v1`）  
　　↓  
业务层 Service（`internal/service`）  
　　↓  
DAO（`internal/dao`）  
　　↓  
Model（`internal/model`）  
　　↓  
数据库

---

## 三、项目启动全流程（从 `main.go` 到 Gin 路由）

### 1. 进程启动顺序

Go 程序启动顺序为：

1. 执行所有包级别的 `init()` 函数。
2. 再执行 `main()` 函数。

本项目中，关键的初始化逻辑在 `main.go` 中： 

```go
func init() {
    err := setupSetting()   // 读取配置
    err = setupDBEngine()   // 初始化数据库
    err = setupLogger()     // 初始化日志
}

func main() {
    gin.SetMode(global.ServerSetting.RunMode)
    router := routers.NewRouters()

    s := &http.Server{
        Addr:           ":" + global.ServerSetting.HttpPort,
        Handler:        router,
        ReadTimeout:    global.ServerSetting.ReadTimeout,
        WriteTimeout:   global.ServerSetting.WriteTimeout,
        MaxHeaderBytes: 1 << 20,
    }

    if err := s.ListenAndServe(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
```

### 2. `setupSetting`：加载配置

- 调用 `setting.NewSetting()` 读取 `configs/config.yaml`。
- 使用 `ReadSection("Server", &global.ServerSetting)` 等方法填充配置结构体。
- 将 `ReadTimeout` / `WriteTimeout` 换算为 `time.Duration`（乘以 `time.Second`）。

### 3. `setupDBEngine`：初始化数据库

- 调用 `model.NewDBEngine(global.DatabaseSetting)`，内部：
  - 使用 `Database` 配置拼接 MySQL DSN。
  - `gorm.Open(mysql.Open(dsn), &gorm.Config{})` 创建连接。
  - 设置最大空闲连接数、最大打开连接数。
  - 如果 `RunMode == "debug"`，打开详细 SQL 日志。
  - 调用 `setupCallbacks(db)` 注册：
    - 创建前回调：自动维护 `CreatedOn`、`ModifiedOn`。
    - 更新前回调：自动更新 `ModifiedOn`。
    - 删除前回调：根据 `DeletedOn`、`IsDel` 字段执行软删除（`UPDATE ... SET deleted_on=?, is_del=1`）。
- 将得到的 `*gorm.DB` 赋值给 `global.DBEngin`，供 DAO 层使用。

### 4. `setupLogger`：初始化日志

- 使用 `logger.NewLogger(&lumberjack.Logger{...})` 创建带日志切割（按大小、时间）的 Logger。
- 日志路径由 `App.LogSavePath`、`App.LogFileName`、`App.LogFileExt` 组合而成。
- 最终写入 `global.Logger`。

### 5. `main()`：启动 HTTP 服务

- 根据 `global.ServerSetting.RunMode` 设置 Gin 模式（`debug` / `release`）。
- 调用 `routers.NewRouters()` 构造 Gin Engine（注册路由和中间件）。
- 构造 `http.Server`：端口、Handler、读写超时等均来自配置。
- `ListenAndServe()` 开启监听，服务正式对外提供 HTTP 接口。

---

## 四、路由与中间件（`internal/routers`）

### 1. 路由初始化：`routers.NewRouters()`

核心逻辑：

- 创建 Gin 引擎：`r := gin.New()`。
- 注册基础中间件：
  - `gin.Logger()`：请求日志。
  - `gin.Recovery()`：panic 恢复。
  - `middleware.Translations()`：参数校验错误的多语言翻译（通常是中文提示）。
- 注册 Swagger 路由：
  - `r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))`。
- 注册业务路由：
  - `apiv1 := r.Group("/api/v1")`
    - 标签模块：`/api/v1/tags` → `Tag` 相关 Handler。
    - 文章模块：`/api/v1/articles` → `Article` 相关 Handler。

### 2. 中间件：`middleware.Translations()`

- 负责注册验证器（validator）的错误信息翻译。
- 当 `app.BindAndValid` 在绑定和校验参数时出现错误，会通过该中间件注入的翻译器返回更友好的中文错误提示。

---

## 五、一次请求的完整调用链（以“获取标签列表”为例）

以：

> `GET /api/v1/tags?name=Go&page=1&page_size=10&state=1`

为例，完整流程如下。

### 1. 进入路由 & Handler

1. HTTP 请求到达 Gin 引擎。
2. 根据路由规则，匹配到 `apiv1.Group("/tags").GET("", tag.List)`。
3. 执行 `v1.Tag.List(c)`。

### 2. Handler 层：参数绑定与调用 Service

在 `v1.Tag.List` 中：

1. 定义请求参数结构体：`service.TagListRequest`（包含 `name`、`state` 等字段及校验规则）。
2. 使用 `app.BindAndValid(c, &param)` 做：
   - 将查询参数绑定到 `param`。
   - 使用 Gin 的 `binding` + 翻译中间件校验参数合法性。
3. 若校验失败：
   - 使用 `global.Logger.Errorf` 记录错误。
   - 使用 `response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))` 返回错误码和详细错误信息。
4. 校验通过后：
   - 创建 Service：`svc := service.New(c.Request.Context())`。
   - 构造分页对象：
     - `pager := app.Pager{Page: app.GetPage(c), PageSize: app.GetPageSize(c)}`。
   - 调用业务层：
     - `totalRows, err := svc.CountTag(&service.CountTagRequest{Name: param.Name, State: param.State})`。
     - `tags, err := svc.GetTagList(&param, &pager)`。
5. 若业务调用出错：
   - 记录日志，并使用对应的业务错误码（如 `ErrorCountTagFail`、`ErrorGetTagListFail`）返回。
6. 成功时：
   - 调用 `response.ToResponseList(tags, int(totalRows))` 返回数据，格式包含：
     - `list`: 标签数组。
     - `pager`: `{page, page_size, total_rows}`。

### 3. Service 层：封装业务逻辑

`internal/service/tag.go` 中：

- `CountTag(param *CountTagRequest)`：调用 `svc.dao.CountTag(param.Name, param.State)`。
- `GetTagList(param *TagListRequest, pager *app.Pager)`：调用 `svc.dao.GetTagList(param.Name, param.State, pager.Page, pager.PageSize)`。
- `CreateTag` / `UpdateTag` / `DeleteTag`：分别封装参数后调用 DAO 层的方法。

目前业务层较薄，但它是承载「复杂业务规则」的主要位置，未来可以在这里加入：权限校验、字段补全、操作记录等。

### 4. DAO 层：数据库访问封装

`internal/dao/tag.go` 中：

- `CountTag(name string, state uint8)`：
  - 构造 `model.Tag{Name: name, State: state}`。
  - 调用 `tag.Count(d.engine)`。
- `GetTagList(name string, state uint8, page, pageSize int)`：
  - 使用 `app.GetPageOffset(page, pageSize)` 计算偏移量。
  - 调用 `tag.List(d.engine, pageOffset, pageSize)`。
- `CreateTag` / `UpdateTag` / `DeleteTag`：
  - 通过构造 `model.Tag` 及通用 `Model` 字段（如 `CreatedBy`、`ModifiedBy`、`ID`），
  - 调用对应的 `Create`、`Update`、`Delete` 方法。

DAO 层的职责是 **屏蔽具体的 SQL/GORM 调用细节**，对上只暴露语义化的方法。

### 5. Model 层：GORM 模型与通用逻辑

以 `model.Tag` 为例：

- 结构体：
  - 组合了通用 `Model`：`ID`、`CreatedOn`、`DeletedOn`、`IsDel` 等字段。
  - 自身字段：`Name`、`State`。
- 方法：
  - `Count(db)`：根据 `Name`、`State`、`is_del=0` 统计数量。
  - `List(db, pageOffset, pageSize)`：分页查询列表。
  - `Create(db)`：直接 `db.Create(&t)`。
  - `Update(db, values)`：按 `id` + `is_del=0` 更新指定字段。
  - `Delete(db)`：按 `id` + `is_del=0` 删除记录，实质上会触发软删除回调。

通用 `Model` 提供的软删除和时间戳逻辑通过 GORM 回调全局生效，不需要每个模型单独维护。

### 6. 最终返回

- GORM 执行 SQL，从数据库中查出数据，经过 DAO → Service → Handler 层依次返回。
- Handler 使用 `app.Response` 统一封装为 JSON，返回给前端。

---

## 六、文章模块当前状态

- `internal/model/article.go`、`internal/model/article_tag.go` 已经定义了文章及文章-标签关系表结构。
- `internal/service/article.go` 定义了文章相关请求参数结构体（`CreateArticleRequest` 等）。
- `internal/routers/api/v1/article.go` 目前只是简单返回字符串的占位实现：
  - `Get` / `List` / `Create` / `Update` / `Delete` 还未接入 Service/DAO/Model 的完整逻辑。
- 后续若要完善文章模块，可以完全参考标签模块的分层模式：
  - 在 DAO 中实现文章的增删改查和列表统计。
  - 在 Service 中封装业务规则和参数处理。
  - 在 Handler 中完成参数绑定、调用 Service 并返回统一响应。

---

通过以上内容，本文件梳理了从 **进程启动 → 配置/数据库/日志初始化 → Gin 路由与中间件 → Handler → Service → DAO → Model → 数据库** 的完整调用链，便于后续维护和扩展。



go get -u gopkg.in/gomail.v2
Gomail 是一个用于发送电子邮件的简单又高效的第三方开源库，目前只支持使用 SMTP 服务器发送电子邮件，但是其 API 较为灵活，如果有其它的定制需求也可以轻易地借助其实现，这恰恰好符合我们的需求，因为目前我们只需要一个小而美的发送电子邮件的库就可以了。


curl.exe -X POST "http://127.0.0.1:8088/auth" -H "Content-Type: application/json" -d '{"app_key":"eddycjy","app_secret":"go-programming-tour-book"}'