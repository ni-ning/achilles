package router

import (
	"achilles/apps/auth"
	"achilles/apps/blog"
	"achilles/apps/upload"
	_ "achilles/docs"
	"achilles/global"
	"achilles/internal/middleware"
	"achilles/pkg/limiter"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context/ctxhttp"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

func includeRouter(e *gin.Engine) {
	// 加载多个APP的路由配置
	// routers.Include(blog.Routers, orm.Routers, params.Routers, todo.Routers)
	Include(blog.Routers, auth.Routers, upload.Routers)

	for _, opt := range options {
		opt(e)
	}
}

var methodLimiters = limiter.NewMethodLimiter().AddBuckets(limiter.LimiterBucketRule{
	Key:          "/auth",
	FillInterval: time.Second,
	Capacity:     10,
	Quantum:      10,
})

// 初始化
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
	r.Use(middleware.ContextTimeout(10 * time.Second))
	r.Use(middleware.Translations())

	// r.Use(middleware.JWT())

	// 注入路由表项
	includeRouter(r)

	url := ginSwagger.URL("http://127.0.0.1:8090/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 测试超时
	r.GET("/timeout", func(c *gin.Context) {

		global.DBEngine.Raw("select sleep(100);")

		_, err := ctxhttp.Get(c.Request.Context(), http.DefaultClient, "https://www.google.com/")
		if err != nil {
			// 会中断主程序，慎用
			// log.Fatalf("ctxhttp.Get err: %v", err)
			global.Logger.Error("context deadline exceeded")
		}
	})

	return r
}
