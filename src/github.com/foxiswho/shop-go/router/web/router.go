package web

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	"github.com/foxiswho/shop-go/middleware/captcha"
	"github.com/foxiswho/shop-go/middleware/staticbin"

	assets "github.com/foxiswho/shop-go/assets"
	. "github.com/foxiswho/shop-go/conf"
	"github.com/foxiswho/shop-go/middleware/opentracing"
	"github.com/foxiswho/shop-go/module/auth"
	"github.com/foxiswho/shop-go/module/cache"
	"github.com/foxiswho/shop-go/module/render"
	"github.com/foxiswho/shop-go/module/session"
	sauth "github.com/foxiswho/shop-go/service/user_service/auth"
	web_user "github.com/foxiswho/shop-go/router/web/user"
	web_index "github.com/foxiswho/shop-go/router/web/index"
	web_test "github.com/foxiswho/shop-go/router/web/test"
	"github.com/foxiswho/shop-go/router/base"
)

//---------
// Website Routers
//---------
func Routers() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Context自定义
	e.Use(base.NewBaseContext())
	// Customization
	if Conf.ReleaseMode {
		e.Debug = false
	}
	e.Logger.SetPrefix("web")
	e.Logger.SetLevel(GetLogLvl())

	// Session
	e.Use(session.Session())

	// CSRF
	e.Use(mw.CSRFWithConfig(mw.CSRFConfig{
		ContextKey:  "_csrf",
		TokenLookup: "form:_csrf",
	}))

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// 验证码，优先于静态资源
	e.Use(captcha.Captcha(captcha.Config{
		CaptchaPath: "/captcha/",
		SkipLogging: true,
	}))

	// 静态资源
	switch Conf.Static.Type {
	case BINDATA:
		e.Use(staticbin.Static(assets.Asset, staticbin.Options{
			Dir:         "/",
			SkipLogging: true,
		}))
	default:
		e.Static("/assets", "./assets")
	}

	// Gzip，在验证码、静态资源之后
	// 验证码、静态资源使用http.ServeContent()，与Gzip有冲突，Nginx报错，验证码无法访问
	e.Use(mw.GzipWithConfig(mw.GzipConfig{
		Level: 5,
	}))

	// OpenTracing
	if !Conf.Opentracing.Disable {
		e.Use(opentracing.OpenTracing("web"))
	}

	// 模板
	e.Renderer = render.LoadTemplates()
	e.Use(render.Render())

	// Cache
	e.Use(cache.Cache())

	// Auth
	//e.Use(auth.New(model.GenerateAnonymousUser))
	e.Use(auth.New(sauth.GenerateAnonymousUser))
	// Routers
	e.GET("/", base.Handler(web_index.HomeHandler))
	e.GET("/login", base.Handler(web_user.LoginHandler))
	e.GET("/register", base.Handler(web_user.RegisterHandler))
	e.GET("/logout", base.Handler(web_user.LogoutHandler))
	e.POST("/login", base.Handler(web_user.LoginPostHandler))
	e.POST("/register", base.Handler(web_user.RegisterPostHandler))



	user := e.Group("/user_service")
	user.Use(auth.LoginRequired())
	{
		user.GET("/:id", base.Handler(web_user.UserHandler))
	}

	about := e.Group("/about")
	about.Use(auth.LoginRequired())
	{
		about.GET("", base.Handler(web_index.AboutHandler))
	}
	test := e.Group("/test")
	{
		test.GET("/jwt/tester", base.Handler(web_test.JWTTesterHandler))
		test.GET("/ws", base.Handler(web_test.WsHandler))
		test.GET("/cache", base.Handler(web_test.CacheHandler))
		test.GET("/cookie", base.Handler(web_test.NewCookie().IndexHandler))
	}
	return e
}
