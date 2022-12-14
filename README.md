# 项目脚手架-博客后台

`Gin`框架实现博客后台，具体包括：
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
    * 国际化处理
    * 接口校验
- <a href="/golangv2/blog.html#_5-标签功能开发">5. 标签功能开发</a>
- <a href="/golangv2/blog.html#_12-优雅重启和停止">12. 优雅重启和停止</a>
	* 信号定义
	* 实现优雅重启和停止
	* 容器集群重启和停止
   

## 1. 进行项目设计

### 项目结构搭建
```js
achilles
    |- apps      // 应用集合，项目核心业务逻辑，如果只是单个应用，可放置于internal目录下
        |- blog  // 博客后台App
        |- orm ......
    |- config   // 配置文件，如Server、App、Database
        |- config.yaml
    |- docs     // swag init 自动生成的接口文档
    |- global   // 全局变量，如配置、数据库连接、日志
    |- internal // 内部模块
        |- middleware   // 中间件
        |- model        // 模型层
        |- router       // 路由汇总
    |- pkg          // 项目公共模块
        |- app      // 响应格式，ToResponse、ToErrorResponse
        |- convert  // 类型转换，常用MustInt等工具函数
        |- errcode  // 错误码，定义NewError，自定义业务错误码
        |- logger   // 日志，定义NewLogger，结合初始化函数，得到全局变量
        |- setting  // 配置，定义NewSetting，结合初始化函数，得到全局变量
    |- storage
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

编写常用的错误处理公共方法，标准化错误输出
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

配置读取和更新

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
基本结构参考 配置管理

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
基本结构参考 配置管理

```go
// pkg/logger/logger.go NewLogger 用户初始化日志
package logger
func NewLogger(appSetting *setting.AppSettingS) *zap.Logger
```

### 响应处理
可以通用响应对象，自定义自定义结构体 `Response` 包裹 `gin.Context`，就可以自定义

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



## 4. 接口参数校验

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

### 接口校验

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


## 5. 标签功能开发

```go
// blog/tag_dao.go 搭建基本逻辑
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

## 7. 接口访问控制


## 8. 应用中间件

## 9. 进行链路追踪

[https://www.jaegertracing.io/](https://www.jaegertracing.io/)

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

[http://127.0.0.1:16686](http://127.0.0.1:16686)

```go
go get -u github.com/opentracing/opentracing-go
go get -u github.com/uber/jaeger-client-go
```


## 10. 进行应用配置

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

### 解决方案
1. 命令行参数
2. 系统环境变量
3. 打包进二进制文件中 安装`go-bindata`库，但不支持热更新和监听
4. 集中式配置中心等


### 配置热更新

```go
go get -u golang.org/x/sys/...
go get -u github.com/fsnotify/fsnotify
```


## 11. 进行应用编译


### 使用Idflags设置编译信息


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






