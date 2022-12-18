# 接口API项目脚手架-博客后台

`Gin`框架从0到1搭建博客后台，梳理总结各个环节，具体包括：
- <a href="/golangv2/blog.html#_1-进行项目设计">1. 进行项目设计</a>
    * 项目结构搭建
    * 数据库设计
    * 创建 Model
    * 接口路由规划
    * 视图函数框架
- <a href="/golangv2/blog.html#_2-编写公共组件">2. 编写公共组件</a>
    * 错误码标准化
    * 配置管理
    * 数据库连接
    * 日志写入
    * 响应处理
- <a href="/golangv2/blog.html#_3-生成接口文档">3. 生成接口文档</a>
    * 给接口添加声明式注解
    * 扫描代码生成API文档
    * 在线渲染接口文档
- <a href="/golangv2/blog.html#_4-接口参数校验">4. 接口参数校验</a>
    * 接口校验规则
	* 进行接口校验
	* 国际化处理
- <a href="/golangv2/blog.html#_5-标签功能开发">5. 标签功能开发</a>
- <a href="/golangv2/blog.html#_6-上传文件服务">6. 上传文件服务</a>
	* 添加配置`config.yaml`
	* 通用功能`pkg/upload`
	* 业务逻辑`dao.go`
	* 处理函数`handlers.go`
- <a href="/golangv2/blog.html#_7-接口访问控制">7. 接口访问控制</a>
	* JWT 介绍
	* JWT 集成
	* 认证中间件
- <a href="/golangv2/blog.html#_8-常用应用中间件">8. 常用应用中间件</a>
	* 访问日志记录
	* 异常捕获处理
	* 服务信息存储
	* 接口限流控制
	* 统一超时控制
	* 注册中间件
- <a href="/golangv2/blog.html#_9-进行链路追踪">9. 进行链路追踪</a>
	* OpenTracing规范
	* 分布式链路追踪系统 Jaeger
	* 应用中注入追踪
- <a href="/golangv2/blog.html#_10-进行应用配置">10. 进行应用配置</a>
	* 配置解决方案
	* 编译不包括配置
	* 配置热更新
- <a href="/golangv2/blog.html#_11-进行应用编译">11. 进行应用编译</a>
	* 常用命令与参数
	* 编译信息写入
- <a href="/golangv2/blog.html#_12-优雅重启和停止">12. 优雅重启和停止</a>
	* 信号定义
	* 实现优雅重启和停止
	* 容器集群重启和停止


## 1. 进行项目设计

### 项目结构搭建
```
achilles
    |- apps      // 应用集合，项目核心业务逻辑，如果只是单个应用，可放置于internal目录下
        |- auth  // 用户管理
        |- blog  // 博客后台App
			|- dao.go
			|- handlers.go
			|- urls.go
        |- upload
    |- config   // 配置文件，如Server、App、Database
        |- config.yaml
    |- docs     // swag init 自动生成的接口文档
    |- global   // 全局变量，如配置、数据库连接、日志
    |- internal // 内部模块
        |- middleware   // 中间件
			|- access_log.go
			|- app_info.go
			|- context_timeout.go
			|- jwt.go
			|- limiter.go
			|- recovery.go
			|- tracer.go
			|- traslations.go
        |- model        // 模型层
        |- router       // 路由汇总
    |- pkg          // 项目公共模块
        |- app      // 响应格式，ToResponse、ToErrorResponse
        |- convert  // 类型转换，常用MustInt等工具函数
		|- email	// Email
        |- errcode  // 错误码，定义NewError，自定义业务错误码
		|- limiter	// 限流器
        |- logger   // 日志，定义NewLogger，结合初始化函数，得到全局变量
        |- setting  // 配置，定义NewSetting，结合初始化函数，得到全局变量
		|- tracer	// 追踪器
		|- upload	// 上传工具
		|- util		// 加密等工具
    |- storage
		|- uploads	// 静态文件
		|- logs		// 日志文件
    |- main.go
```

### 数据库设计

- 基本功能
    1. 标签管理：文章所归属的分类，也就是标签，文章可打多个标签;
    2. 文章管理：整个文章内容的管理，并且需要将文章和标签进行关联;

![](./assets/db-design.jpeg)

- 创建数据库 `blog`

```sql
create database if not exists blog charset utf8mb4 collate utf8mb4_general_ci;
```

- 创建库表
```sql
create table blog_article(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    -- deleted_at datetime null comment '删除时间',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    title varchar(100) default '' comment '文章标题',
    `desc` varchar(255) default '' comment '文章简述',
    cover_image_url varchar(255) default '' comment '封面图片地址',
    content longtext comment '文章内容',
    state tinyint(3) unsigned DEFAULT 1 comment '状态 0 为禁用、1 为启用'
)engine=innodb charset=utf8mb4 comment '文章表';


create table blog_tag(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    -- deleted_at datetime null comment '删除时间',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',

    name varchar(100) default '' comment '标签名称',
    state tinyint(3) unsigned DEFAULT 1 comment '状态 0 为禁用、1 为启用'
)engine=innodb charset=utf8mb4 comment '标签表';


create table blog_article_tag(
    id bigint(20) unsigned not null primary key auto_increment comment '主键',
    created_at datetime not null default now() comment '创建时间',
    created_by varchar(64) default '' comment '创建人',
    updated_at datetime not null default now() comment '修改时间',
    updated_by varchar(64) default '' comment '修改人',
    -- deleted_at datetime null comment '删除时间',
    is_deleted tinyint(3) default 0 comment '是否删除 0未删除 1已删除',
    
    article_id bigint(20) unsigned not null comment '文章 ID',
    tag_id bigint(20) unsigned not null comment '标签 ID'
)engine=innodb charset=utf8mb4 comment '文章标签关联表';
```

### 创建 Model
```go
// internal/model/blog.go
package model
import "gorm.io/plugin/soft_delete"

type Model struct {
	ID        uint64    `gorm:"primarykey;" json:"id"` // 主键
	CreatedAt time.Time `json:"-"`                     // 创建时间
	CreatedBy string    `json:"-"`                     // 创建人
	UpdatedAt time.Time `json:"-"`                     // 修改时间
	UpdatedBy string    `json:"-"`                     // 修改人
	// DeletedAt gorm.DeletedAt `json:"-"`                     // 删除时间 库表中对应字段 deleted_at datetime类型
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-" ` // 删除时间 库表中对应字段 is_deleted bool
}

type Article struct {
	*Model
	Title         string `json:"title"`           // 文章标题
	Desc          string `json:"desc"`            // 文章简述
	Content       string `json:"content"`         // 文章内容
	CoverImageUrl string `json:"cover_image_url"` // 封面图片地址
	State         uint8  `json:"state"`           // 状态 0 为禁用、1 为启用
}

func (a Article) TableName() string {
	return "blog_article"
}

type Tag struct {
	*Model
	Name  string `json:"name"`  // 标签名称
	State uint8  `json:"state"` // 状态 0 为禁用、1 为启用
}

func (t Tag) TableName() string {
	return "blog_tag"
}

type ArticleTag struct {
	*Model
	TagID     uint64 `json:"tag_id"`     // 标签 ID
	ArticleID uint64 `json:"article_id"` // 文章 ID
}

func (a ArticleTag) TableName() string {
	return "blog_article_tag"
}
```

### 接口路由规划
```go
package blog

import "github.com/gin-gonic/gin"

var article = NewArticle()
var tag = NewTag()

func Routers(e *gin.Engine) {
	apiv1 := e.Group("/api/v1")
	{
		apiv1.POST("/tags", tag.Create)       // 新增标签
		apiv1.GET("/tags", tag.List)          // 获取标签列表
		apiv1.PUT("/tags/:id", tag.Update)    // 更新指定标签
		apiv1.DELETE("/tags/:id", tag.Delete) // 删除指定标签

		apiv1.POST("/articles", article.Create)           // 新增文章
		apiv1.DELETE("/articles/:id", article.Delete)     // 删除指定文章
		apiv1.PUT("/articles/:id", article.Update)        // 更新指定文章
		apiv1.PATCH("/articles/:id/state", article.Patch) // 禁用或启用文章
		apiv1.GET("/articles/:id", article.Get)           // 获取指定文章
		apiv1.GET("/articles", article.List)              // 获取文章列表
	}
}
```

### 视图函数框架

```go
package blog

// apps/blog/tag.go
type Tag struct{}

func NewTag() Tag {
	return Tag{}
}

func (t Tag) Get(c *gin.Context)    {}
func (t Tag) List(c *gin.Context)   {}
func (t Tag) Create(c *gin.Context) {}
func (t Tag) Update(c *gin.Context) {}
func (t Tag) Delete(c *gin.Context) {}

// apps/blog/article.go
type Article struct{}

func NewArticle() Article {
	return Article{}
}

func (a Article) Get(c *gin.Context)    {}
func (a Article) List(c *gin.Context)   {}
func (a Article) Create(c *gin.Context) {}
func (a Article) Update(c *gin.Context) {}
func (a Article) Delete(c *gin.Context) {}
func (a Article) Patch(c *gin.Context)  {}
```


## 2. 编写公共组件

基础组件，或公共组件，保证应用程序的标准化

- 错误码标准化
- 配置管理
- 数据库连接
- 日志写入
- 响应处理

### 错误码标准化

与客户端进行交互两点：
- 一个是正确响应下的结果集返回;
- 一个是错误响应的错误码和消息体返回;

客户端会调用不同的服务端，需要保证标准化的错误码格式

预定义项目中一些公共错误码, 引导和规范大家使用
```go
package errcode

// pkg/errcode/common_code.go
var (
	Success                   = NewError(0, "成功")
	ServerError               = NewError(10000000, "服务内部错误")
	InvalidParams             = NewError(10000001, "入参错误")
	NotFound                  = NewError(10000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败, 找不到对应的 AppKey 和 AppSecret")
	UnauthorizedTokenError    = NewError(10000004, "鉴权失败, Token 错误")
	UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败, Token 超时")
	UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败, Token 生成失败")
	TooManyRequests           = NewError(10000007, "请求过多")
)
```

编写常用的错误处理公共方法，标准化错误输出，自定义实现`error`
```go
package errcode

// pkg/errcode/errcode.go
type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e // 值拷贝
	newError.details = []string{}
	newError.details = append(newError.details, details...)
	return &newError
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
```

### 配置管理

配置读取和更新，封装第三方库`viper`

![](./assets/config.jpeg)

在项目目录下的 `config` 目录新建 `config.yaml` 文件，写入以下配置：
```yaml
Server:
  RunMode: debug
  HttpPort: 8090
  ReadTimeout: 60
  WriteTimeout: 60
App:
  DefaultPageSize: 10
  MaxPageSize: 100
  LogSavePaht: storage/logs
  LogFileName: app
  LogFileExt: .log
Database:
  DBType: mysql
  Username: root        # 填写你的数据库账号
  Password: 123456      # 填写你的数据库密码
  Host: 127.0.0.1:3306
  DBName: achilles      # 填写你的数据库名称
  TablePrefix: blog_
  Charset: utf8mb4
  ParseTime: True
  MaxIdleConns: 10
  MaxOpenConns: 30
```

分层思想构建程序结构，其他程序可以调用全局`global.*Setting`实例即可
```go
// config/config.yaml 配置文件

// global/setting.go 全局配置结构体实例变量
package global

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
)

// setting/setting.go 封装 viper.New 为 NewSetting
package setting

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.AddConfigPath("config/")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")

	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp}, nil
}

// setting/section.go 自定义结构体
package setting

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
......

// main.go 定义初始化配置函数
func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
}
```

### 数据库连接
基本结构参考 配置管理，封装第三方库`gorm`

```go
// internal/model/model.go NewDBEngine 用于初始化连接
package model

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	dsn := fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
```

### 日志写入
基本结构参考 配置管理，封装第三方库`zap`

```go
// pkg/logger/logger.go NewLogger 用户初始化日志
func NewLogger(appSetting *setting.AppSettingS) *zap.Logger {
	// 1. Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式 2022-12-08T18:24:07.979+0800
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写 "level":"INFO"
	encoderConfig.TimeKey = "timestamp"                     // "timestamp":"2022-12-08T18:41:35.596+0800"
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// 2. WriterSyncer 同时输出到 lumberJackLogger 和控制台
	logLocation := fmt.Sprintf("%s%s%s",
		appSetting.LogSavePath,
		appSetting.LogFileName,
		appSetting.LogFileExt)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logLocation, // 日志文件的位置
		MaxSize:    1,           // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 5,           // 保留旧文件的最大个数
		MaxAge:     30,          // 保留旧文件的最大天数
		Compress:   false,       // 是否压缩/归档旧文件
	}
	writeSyncer := zapcore.AddSync(io.MultiWriter(lumberJackLogger, os.Stdout))

	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	// 调用函数信息 如 "caller":"day15/main.go:27", zap.AddCallerSkip(1)用于额外封装一层场景
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}
```

### 响应处理

通用的响应对象，自定义自定义结构体 `Response` 封装 `gin.Context`

```go
package app

import (
	"achilles/pkg/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 包装一层，就可以自定义方法，配合New函数
type Response struct {
	Ctx *gin.Context
}

// 基于列表返回，定义小类型
type Pager struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	TotalRows int `json:"total_rows"`
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, gin.H{
		"code": errcode.Success.Code(),
		"msg":  errcode.Success.Msg(),
		"data": data,
	})
}

func (r *Response) ToResponseList(list interface{}, totalRows int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"code": errcode.Success.Code(),
		"list": list,
		"pager": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}
```

## 3. 生成接口文档

从功能使用上来讲，OpenAPI 规范能够帮助我们描述一个 API 的基本信息
- 有关该 API 的描述;
- 可用路径(/资源);
- 在每个路径上的可用操作(GET/POST/PUT/PATCH/DELETE...);
- 每个操作的输入/输出格式;

[Swagger](https://swaggo.github.io/swaggo.io/declarative_comments_format/) 基于标准的 OpenAPI 规范进行设计，生成统一标准的接口文档，基本步骤
1. 给接口添加声明式注解
2. 扫描代码生成API文档
3. 在线渲染接口文档


### 安装依赖

```go
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get -u github.com/alecthomas/template
```

### 生成 swag 命令
```sh
# 具体路径结合个人环境
cd /Users/nining/go/pkg/mod/github.com/swaggo/swag@v1.8.8/cmd/swag
go install

vim /etc/profile
export PATH=/Users/nining/go/bin:$PATH
source /etc/profile

swag -v
swag version v1.8.8
```

### 生成接口文档

1. 给接口添加声明式注解
2. 扫描代码生成API文档
```go
// 项目主目录下执行
swag int
```
3. 在线渲染接口文档

```go
// internal/router/router.go
import (
	_ "achilles/docs"  // 很重要

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)
url := ginSwagger.URL("http://127.0.0.1:8090/swagger/doc.json")
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
```

### 访问接口文档
[http://127.0.0.1:8090/swagger/index.html](http://127.0.0.1:8090/swagger/index.html)

![](./assets/swagger.png)

## 4. 接口参数校验

### 接口校验规则

`Gin`框架，其内部的模型绑定和验证默认使用的是 [go-playground/validator](https://github.com/go-playground/validator) 来进行参数绑定和校验

业务接口校验，就是将校验规则写在对应的结构体的字段标签上
|  标签   | 含义  |
|  ----  | ----  |
| required	  | 必填 |
| gt  | 大于 |
| gte  | 大于等于 |
| lt  | 小于 |
| lte  | 小于等于 |
| gt  | 大于 |
| gt  | 大于 |
|min	|最小值|
|max	|最大值|
|oneof	|参数集内的其中之一|
|len	|长度要求与 len 给定的一致|

```go
// apps/blog/dao_tag.go
type TagListCountRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State *uint8 `form:"state"`
}

// 与 TagListRequest 中 State取值不同
type TagCreateRequest struct {
	Name  string `form:"name" binding:"required,min=1,max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

// 全量更新
type TagUpdateRequest struct {
	ID    uint64 `form:"id" binding:"required,gte=1"`
	Name  string `form:"name" binding:"min=1,max=100"`
	State *uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

type TagDeleteRequest struct {
	ID uint64 `form:"id" binding:"required,gte=1"`
}
```

### 进行接口校验

封装通用的校验函数 `BindAndValid`
```go
package app

// pkg/app.form.go 封装BindAndValid
type ValidError struct {
	Key     string
	Message string
}

type ValidErrors []*ValidError

func (v *ValidError) Error() string {
	return v.Message
}

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}

	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		verrs, ok := err.(val.ValidationErrors)
		if !ok {
			return false, errs
		}

		for key, value := range verrs.Translate(trans) {
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value,
			})
		}

		return false, errs
	}

	return true, nil
}
```

### 国际化处理

定义中间件，统一处理

```go
package middleware

// internal/middleware/translations.go
func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			case "en":
				_ = en_translations.RegisterDefaultTranslations(v, trans)
			default:
				_ = zh_translations.RegisterDefaultTranslations(v, trans)
			}
			c.Set("trans", trans)
		}

		c.Next()

	}
}
```


## 5. 标签功能开发

当业务线逻辑较多时，需要规划代码结构
```go
|- apps
	|- blog  // 博客后台App
		|- dao_article.go
		|- dao_tag.go
		|- dao_xxx.go
		|- handlers_article.go
		|- handlers_tag.go
		|- handlers_xxx.go
		|- urls.go
```
- `dao_xxx.go`封装数据库处理
- `handlers_xxx.go`视图处理函数，进行参数校验、调用`dao_xxx.go`中方法和进行相应返回
```go
// apps/blog/dao_tag.go 搭建基本逻辑
package blog

import (
	"achilles/global"
	"achilles/internal/model"
	"achilles/pkg/app"
	"fmt"
)

type TagListCountRequest struct {
	Name  string `form:"name" binding:"max=100"`
	State *uint8 `form:"state"`
}

// 与 TagListRequest 中 State取值不同
type TagCreateRequest struct {
	Name  string `form:"name" binding:"required,min=1,max=100"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

// 全量更新
type TagUpdateRequest struct {
	ID    uint64 `form:"id" binding:"required,gte=1"`
	Name  string `form:"name" binding:"min=1,max=100"`
	State *uint8 `form:"state,default=1" binding:"oneof=0 1"`
}

type TagDeleteRequest struct {
	ID uint64 `form:"id" binding:"required,gte=1"`
}

func GetTagList(req TagListCountRequest, page, pageSize int) ([]*model.Tag, error) {
	var tags []*model.Tag
	db := global.DBEngine.Model(&model.Tag{})
	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.State != nil {
		db = db.Where("state = ?", req.State)
	}

	fmt.Printf("\nGetTagList Req:%#v\n", req)

	pageOffset := app.GetPageOffset(page, pageSize)
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}

	if err := db.Debug().Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func GetTagCount(req TagListCountRequest) (int, error) {
	var count int64
	db := global.DBEngine.Model(&model.Tag{})
	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if req.State != nil {
		db = db.Where("state = ?", req.State)
	}

	fmt.Printf("\nGetTagCount Req:%#v\n", req)

	if err := db.Debug().Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func CreateTag(req TagCreateRequest) error {
	fmt.Printf("\nCreateTag Req:%#v\n", req)

	return global.DBEngine.Debug().Create(&model.Tag{Name: req.Name, State: req.State}).Error
}

func UpdateTag(req TagUpdateRequest) error {
	fmt.Printf("\nUpdateTag Req:%#v\n", req)

	db := global.DBEngine.Model(&model.Tag{})
	db = db.Debug().Where("id = ?", req.ID)
	return db.Updates(map[string]interface{}{"name": req.Name, "state": req.State}).Error
}

func DeleteTag(req TagDeleteRequest) error {
	fmt.Printf("\nDeleteTag Req:%#v\n", req)

	return global.DBEngine.Delete(&model.Tag{}, req.ID).Error
}
```

## 6. 上传文件服务

实现单独上传文件接口服务，该接口返回文件地址，结合其他接口完成业务服务。另外实际业务中后端可以接口`S3`等在线存储服务

### 添加配置`config.yaml`
`configs/config.yaml`配置文件新增配置项，完善配置读取和初始化相关代码

```yaml
App:
  ...
  UploadSavePath: storage/uploads
  UploadServerUrl: http://127.0.0.1:8090/static
  UploadImageMaxSize: 5  # MB
  UploadImageAllowExts:
    - .jpg
    - .jpeg
    - .png
```

### 通用功能`pkg/upload`
```go
// pkg/uplaod/file.go

```
### 业务逻辑`dao.go`
```go
func UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported")
	}
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximum file limit")
	}

	uploadSavePath := upload.GetSavePath()
	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory")
		}
	}
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions")
	}

	dst := uploadSavePath + "/" + fileName
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
```

### 处理函数`handlers.go`
```go
type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	fileInfo, err := UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Error("UploadFile Error", zap.String("error", err.Error()))
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
```

## 7. 接口访问控制

常见的两种 API 访问控制方案
- OAuth 2.0: 授权机制的指导标准，常用于使用第三方登陆的情况
- JWT: 常用于前后端分离的情况，能够非常便捷的给 API 接口提供安全鉴权


### JWT 介绍
JWT全称JSON Web Token是一种跨域认证解决方案，属于一个开放的标准，规定了一种Token实现方式，常用于前后端分离项目

推荐阅读[阮一峰JSON Web Token 入门教程](https://www.ruanyifeng.com/blog/2018/07/json_web_token-tutorial.html)


JSON Web 令牌（JWT）是由紧凑的形式三部分组成，这些部分由点 “.“ 分隔，组成为 ”xxxxx.yyyyy.zzzzz“ 的格式，三个部分分别代表的意义如下：
![](./assets/jwt-format.jpeg)


### JWT 集成
- 安装依赖
```go
go get github.com/golang-jwt/jwt
```

- 生成token与解析token
```go
// pkg/app/jwt.go
package app

import (
	"achilles/global"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	AccountID int64 `json:"account_id"`
	jwt.StandardClaims
}

func GetJWTSecret() []byte {
	return []byte(global.JWTSetting.Secret)
}

func GenerateToken(accountId int64) (string, error) {
	// accountIdStr := strconv.FormatInt(accountId, 10)
	nowTime := time.Now()
	expireTime := nowTime.Add(global.JWTSetting.Expire)
	claims := Claims{
		// AccountID: util.EncodeMD5(accountIdStr),
		AccountID: accountId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    global.JWTSetting.Issuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(GetJWTSecret())
	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
```

### 认证中间件

```go
// internal/middleware/jwt.go
package middleware

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//@Param token header string true "Token"
		var (
			token string
			ecode = errcode.Success
		)
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}

		response := app.NewResponse(c)
		if token == "" {
			response.ToErrorResponse(errcode.UnauthorizedTokenLack)
			c.Abort()
			return
		}

		claims, err := app.ParseToken(token)
		if err != nil {
			switch err.(*jwt.ValidationError).Errors {
			case jwt.ValidationErrorExpired:
				ecode = errcode.UnauthorizedTokenTimeout
			default:
				ecode = errcode.UnauthorizedTokenError
			}
			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		existed, account := auth.GetAccountById(claims.AccountID)
		if existed {
			c.Set("account", account)
		}

		c.Next()
	}
}
```


## 8. 常用应用中间件

中间件，对访问统一处理

### 访问日志记录
```go
// interal/middleware/access_log.go
package middleware

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 截流多写一次
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer = bodyWriter重写接入一下
		bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWriter

		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		global.Logger.Info("Access Log",
			zap.String("method", c.Request.Method),
			zap.Int("status_code", bodyWriter.Status()),
			zap.Int64("begin_time", beginTime),
			zap.Int64("endTime", endTime),
			zap.String("request", c.Request.PostForm.Encode()),
			zap.String("response", bodyWriter.body.String()),
		)
	}
}
```
自定义结构体`AccessLogWriter`，实现双写`Write`，就可以读取到响应体`body.String()`

### 异常捕获处理

自定义异常捕获处理，实现邮件提醒

`Gomail` 是一个小而美，用于发送电子邮件的简单又高效的第三方开源库

- 安装依赖
```go
go get -u gopkg.in/gomail.v2
```
- 封装邮箱工具包
```go
// pkg/email.email.go
package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Email struct {
	*SMTPInfo
}

type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}
```
- 实现中间件
```go
// interal/middleware/recovery.go

package middleware

func Recovery() gin.HandlerFunc {
	defailtMailer := email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	})
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				global.Logger.Error("panic recover err", zap.Any("error", err))

				err := defailtMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("异常抛出，发生时间: %d", time.Now().Unix()),
					fmt.Sprintf("错误信息: %v", err),
				)
				if err != nil {
					global.Logger.Panic("mail.SendMail Error%v", zap.Any("error", err))
				}

				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()

		c.Next()
	}
}
```

### 服务信息存储

进程内上下文传递内部信息

```go
// interal/middleware/app_info.go

package middleware

import "github.com/gin-gonic/gin"

func AppInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("app_name", "achilles")
		c.Set("app_version", "v1.0.0")
		c.Next()
	}
}
```
### 接口限流控制

针对应用接口进行限流控，`ratelimit` 提供了大量接口方法实现限流

- 安装依赖
```go
go get -u github.com/juju/ratelimit
```

- 编写中间件
```go
// internal/middleware/limiter.go
func RateLimiter(l limiter.LimiterIface) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := l.Key(c)
		if bucket, ok := l.GetBucket(key); ok {
			count := bucket.TakeAvailable(1)
			if count == 0 {
				response := app.NewResponse(c)
				response.ToErrorResponse(errcode.TooManyRequests)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
```


### 统一超时控制

统一的在应用程序中针对所有请求都进行一个最基本的超时时间控制
![](./assets/app-context-deadline.jpeg)

```go
// internal/middleware/context_timeout.go
func ContextTimeout(t time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
```
在使用时，只要传入`c.Request.Context()`的请求，统一接受超时时间控制
```go
r.GET("/timeout", func(c *gin.Context) {

	global.DBEngine.Raw("select sleep(100);")

	_, err := ctxhttp.Get(c.Request.Context(), http.DefaultClient, "https://www.google.com/")
	if err != nil {
		// 会中断主程序，慎用
		// log.Fatalf("ctxhttp.Get err: %v", err)
		global.Logger.Error("context deadline exceeded")
	}
})
```

### 注册中间件

注意注册中间的顺序

```go
var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

func NewRouter() *gin.Engine {
	// 初始化路由
	gin.SetMode(global.ServerSetting.RunMode)
	r := gin.New()

	r.Use(middleware.Tracing())

	if global.ServerSetting.RunMode == "debug" {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(middleware.AccessLog())
		r.Use(middleware.Recovery())
	}

	r.Use(middleware.RateLimiter(methodLimiters))
	r.Use(middleware.ContextTimeout(global.AppSetting.DefaultContextTimeout))
	r.Use(middleware.Translations())

	// r.Use(middleware.JWT())

	return r
}
```


## 9. 进行链路追踪

随着项目不断迭代，会涉及到很多接口，这些接口可能分布式部署的，既存在着多份副本，又存在着相互调用，并且各自的调用中还包含大量的SQL、HTTP、Redis以及应用的逻辑调用。使用分布式链路追踪系统，可以解决多程序部署在多环境下调用链路的"观察"

### OpenTracing规范

追踪模型的思维模型起源于 Google 的 Dapper, a Large-Scale Distributed Systems Tracing Infrastructure论文，有一系列约定的术语

|  术语   | 含义  | 概述 |
|  ----  | ----  |---- |
| Trace	 | 追踪   | 代表事务或流程在(分布式)系统执行过程 |
| Span	 | 跨度   | 代表事务中的工作单元，多个Span完成一个Trace |
| SpanContext	 | 跨度上下文   | 代表事务的跟踪信息 |

### 分布式链路追踪系统 Jaeger

[https://www.jaegertracing.io/](https://www.jaegertracing.io/)

- 安装
```sh
docker run -d  --name jaeger \
-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
-p 5775:5775/udp \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778	\
-p 16686:16686	\
-p 14268:14268	\
-p 9411:9411	\
jaegertracing/all-in-one:latest
```

打开浏览器，访问 [http://127.0.0.1:16686](http://127.0.0.1:16686)

### 应用中注入追踪

- 安装第三方库
```go
go get -u github.com/opentracing/opentracing-go
go get -u github.com/uber/jaeger-client-go
```
- 编写 tracer 工具
```go
// pkg/tracer/tracer.go
package tracer

import (
	"io"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func NewJaegerTracer(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHostPort,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
```
- 全局配置与初始化
```go
var (
	Tracer opentracing.Tracer
)

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer(
		"achilles",
		"127.0.0.1:6831",
	)
	if err != nil {
		return err
	}

	global.Tracer = jaegerTracer
	return nil
}
```
- 编写中间件
```go
// internal/middleware/tracer.go
package middleware

import (
	"achilles/global"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func Tracing() func(c *gin.Context) {
	return func(c *gin.Context) {
		var newCtx context.Context
		var span opentracing.Span

		spanCtx, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(c.Request.Header),
		)
		if err != nil {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				c.Request.URL.Path,
			)
		} else {
			span, newCtx = opentracing.StartSpanFromContextWithTracer(
				c.Request.Context(),
				global.Tracer,
				c.Request.URL.Path,
				opentracing.ChildOf(spanCtx),
				opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
			)
		}
		defer span.Finish()
		c.Request = c.Request.WithContext(newCtx)

		c.Next()

	}
}


// 追踪中间件在所有中间件之前生效
r.Use(middleware.Tracing())
```

- 日志追踪
```go
// internal/middleware/tracer.go 中间件增加
....
var traceID string
var spanID string
spanContext := span.Context()
switch spanContext.(type) {
case jaeger.SpanContext:
	jaegerContext := spanContext.(jaeger.SpanContext)
	traceID = jaegerContext.TraceID().String()
	spanID = jaegerContext.SpanID().String()
}
c.Set("X-Trace-ID", traceID)
c.Set("X-Span-ID", spanID)

c.Request = c.Request.WithContext(newCtx)
c.Next()

// internal/middleware/access_log.go 中间件增加
zap.String("trace_id", c.GetString("X-Trace-ID")),
zap.String("span_id", c.GetString("X-Span-ID")),
```
- SQL追踪


## 10. 进行应用配置

### 编译不包括配置
配置文件，即非.go文件的文件类型，不会被打包到二进制文件中，所以需要解决配置文件读取问题

- 如何获取相对路径
```go
func main(){
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	fmt.Println(path)
}

// 预期一致
go build . && ./demo
/Users/nining/learning-go/demo/demo

// 预期不一致，临时目录
go run main.go
/var/folders/wj/j8p_hk852pj_xxv_7wgplpc80000gn/T/go-build3222903838/b001/exe/main
```

### 配置解决方案
1. 命令行参数
```go
// main.go中定义包内全局变量，优先处理的逻辑
var (
	port    string
	runMode string
	config  string
)
// init 中调用 setupFlag() setupSetting()
func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "config/", "指定要使用的配置文件路径")
	flag.Parse()
	return nil
}

func setupSetting()error{
	......
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
}

// 启动时就可以加上参数
/achilles -config=config/ -port=8001 -mode=debug
```
2. 系统环境变量
```go
// 将配置文件存放在系统自带的环境变量中
os.Getenv("ENV")

优先级：命令行 > 环境变量 > 手动文件配置
```
3. 打包进二进制文件中 安装`go-bindata`库，但不支持热更新和监听
4. 集中式配置中心等

### 配置热更新

- 安装第三方库
```go
go get -u golang.org/x/sys/...
go get -u github.com/fsnotify/fsnotify
```
- 协程监听实现
```go
// pkg/setting/setting.go
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
```


## 11. 进行应用编译

### 常用命令与参数

```go
go help
	build       compile packages and dependencies
	install     compile and install packages and dependencies
	run         compile and run Go program
	...
```
- 交叉编译
```
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -o blog-service .
```

### 编译信息写入

使用Idflags设置编译信息

```go
// 定义局部全局变量
var (
	isVersion    bool
	buildTime    string
	buildVersion string
)
// 初始化 -version参数
func setupFlag() error {
	.....
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()
	return nil
}
// 使用 -version 仅产看编译信息
func main() {
	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		return
	}
}

// 编译参数
go build -ldflags "-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVersion=v1.0.0" -o blog-service

// 测试结果
% ./blog-service -version
build_time: 2022-xx-xx,17:10:00
build_version: v1.0.0
```


## 12. 优雅重启和停止

程序需要不间断地更新和发布，可以通过信号量方式来解决

### 信号定义

信号是UNIX、类UNIX，以及其他POSIX兼容的操作系统中进程间通信的一种有限制的方式

当一个信号发送给一个进程时，操作系统中断了进程正常的控制流程。此时，任何非原子操作都被中断，如果进程定义了信号的处理函数，则被执行；后轴执行默认的函数处理

- `kill -l` 查看系统支持的信号
- `ctrl + c` 发送SIGINT信号, 希望进程结束的通知
- `kill -9 pid` 发送SIGKILL信号, 强制中断进程


### 实现优雅重启和停止

```go
func main() {

	s := &http.Server{......}
	// 新增优雅删除逻辑
	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	// the channel used with signal.Notify should be buffered (staticcheck)
	// 接受 syscall.SIGINT和syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 最大时间控制, 用户通知该服务它有5秒的时间来处理原有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server force to shutdown:", err)
	}

	log.Println("Server exiting...")
}
```

- 伪造一个需要长时间查询的接口 `/timeout  global.DBEngine.Raw("select sleep(100);")`
- 执行`ctrl + c` 发送SIGINT信号, 希望进程结束，过5秒后，强制结束

```go
Shutting down server...
2022/xx/xx 18:01:02 Server force to shutdown:context deadline exceeded
```
- 如果没有要处理的旧请求，`ctrl + c` 不需要等待，直接退出


### 容器集群重启和停止

- Kubernetes 向 Pod 发送 SIGTERM 信号，告诉容器它很快就会被关闭，让应用程序停止接收新的请求;
- Kubernetes 等待一会，称为优雅终止宽限期(默认30s)，如果容器仍在运行，发送 SIGKILL 强制删除;