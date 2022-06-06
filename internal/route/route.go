package route

import (
	"embed"
	"flag"
	"io/fs"
	"io/ioutil"
	"net/http"
	"time"

	//"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
	//"github.com/markbates/pkger"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/nanmu42/gzip"
	ginglog "github.com/szuecs/gin-glog"

	"restdoc/handlers/about"
	//"restdoc/handlers/changePassword"
	"restdoc/handlers/authcode"
	"restdoc/handlers/forgotPassword"
	"restdoc/handlers/home"
	"restdoc/handlers/login"
	"restdoc/handlers/logout"
	"restdoc/handlers/privacy"
	"restdoc/handlers/restdoc/api"
	"restdoc/handlers/restdoc/endpoint"
	"restdoc/handlers/restdoc/group"
	"restdoc/handlers/restdoc/home"
	"restdoc/handlers/restdoc/project"
	"restdoc/handlers/settings"
	"restdoc/handlers/signup"
	"restdoc/handlers/team"
	"restdoc/handlers/teamuser"
	"restdoc/handlers/terms"
	"restdoc/handlers/upload"
	"restdoc/handlers/user"
	"restdoc/internal/middlewares"
	"restdoc/render"
)

var StaticBox fs.FS
var TemplateBox embed.FS

func InitRouter() *gin.Engine {
	// 初始化路由
	flag.Parse()
	r := gin.Default()

	//store := cookie.NewStore([]byte("secret"))

	handler := gzip.NewHandler(gzip.Config{
		// gzip compression level to use
		CompressionLevel: 2,
		// minimum content length to trigger gzip, the unit is in byte.
		MinContentLength: 2048,
		// RequestFilter decide whether or not to compress response judging by request.
		// Filters are applied in the sequence here.
		RequestFilter: []gzip.RequestFilter{
			gzip.NewCommonRequestFilter(),
			gzip.DefaultExtensionFilter(),
		},
		// ResponseHeaderFilter decide whether or not to compress response
		// judging by response header
		ResponseHeaderFilter: []gzip.ResponseHeaderFilter{
			gzip.DefaultContentTypeFilter(),
		},
	})

	// 使用Sentry
	// r.Use(sentry.Recovery(raven.DefaultClient, false))
	// r.Use(middlewares.GlobalRecover)
	//r.Use(sentry.Recovery(raven.DefaultClient, false))

	r.Use(sentrygin.New(sentrygin.Options{}))
	r.Use(ginglog.Logger(10 * time.Second))
	r.Use(middlewares.Cors)
	r.Use(middlewares.CacheControl)
	r.Use(middlewares.GetSession)
	r.Use(middlewares.JsonP())
	r.Use(handler.Gin)

	//r.Use(sessions.Sessions("mysession", store))

	render.TemplateBox = TemplateBox
	render.InitRender()
	//rder := render.Render(box)
	r.HTMLRender = &render.Render

	r.GET("/", home.Home)
	r.POST("/", home.Home)
	r.GET("/login", login.LoginPage)
	r.POST("/login", login.Login)
	r.GET("/signup", signup.SignUpPage)
	r.GET("/signupemail", signup.SignupMailPage)
	r.GET("/forgotpasswordemail", forgotPassword.ForgotPasswordMailPage)
	r.POST("/signup", signup.SignUp)
	r.POST("/getmailcode", signup.GetMailCode)
	r.GET("/forgotpassword", forgotPassword.ForgotPasswordPage)
	r.POST("/resetpassword/user", forgotPassword.ResetPassword)
	r.POST("/forgotpassword/user", forgotPassword.GetForgotPasswordMailCode)
	r.POST("/logout", logout.Logout)

	r.GET("/about", about.About)
	r.GET("/terms", terms.Terms)
	r.GET("/privacy-policy", privacy.Privacy)

	r.POST("/code/add", authcode.Add)
	r.GET("/code/get", authcode.Detail)

	//r.POST("/gt/validate", gt.ValidateGeetest)

	r.POST("/api/file/upload", upload.Upload)

	r.GET("/api/restdoc/project/list", restdocProject.List)
	r.GET("/api/restdoc/project/detail/:id", restdocProject.Detail)
	r.POST("/api/restdoc/project/create", restdocProject.Add)
	r.POST("/api/restdoc/project/update", restdocProject.Update)
	r.POST("/api/restdoc/project/delete", restdocProject.Delete)

	r.POST("/api/restdoc/group/create", restdocGroup.Add)
	r.POST("/api/restdoc/group/move", restdocGroup.Move)

	r.GET("/api/restdoc/api/list", restdocApi.List)
	r.POST("/api/restdoc/api/create", restdocApi.Add)
	r.POST("/api/restdoc/api/update", restdocApi.Update)
	r.GET("/api/restdoc/api/detail/:id", restdocApi.Detail)

	r.POST("/api/restdoc/endpoint/update", restdocEndpoint.Update)
	r.POST("/api/restdoc/endpoint/create", restdocEndpoint.Create)

	r.GET("/api/user/info", user.UserInfo)

	r.GET("/restdoc", restdocHome.Locale)
	r.GET("/restdoc/", restdocHome.Locale)

	r.GET("/team", team.Page)

	r.GET("/api/team/list", team.List)
	r.POST("/api/team/create", team.Create)

	r.GET("/teamuser/:id", teamuser.MemberPage)
	r.POST("/teamuser/:id/create", teamuser.MemberCreate)
	r.GET("/teamuser/:id/list", teamuser.MemberList)
	r.GET("/teamuser/:id/detail/:id", teamuser.DetailPage)
	r.GET("/teamuser/:id/info/:id", teamuser.MemberInfo)
	r.GET("/teamuser/:id/delete/:id", teamuser.DeletePage)
	r.POST("/teamuser/:id/delete", teamuser.Delete)

	r.GET("/api/settings/detail", settings.Info)
	r.POST("/api/settings/update", settings.Update)

	staticBox, _ := fs.Sub(StaticBox, "static")
	faviconBox, _ := fs.Sub(StaticBox, "static")
	//r.StaticFS("/static", m.Middleware(http.FileServer(http.FS(staticBox))))
	//r.GET("/static", gin.WrapH(m.Middleware(http.FS(staticBox))))
	r.StaticFS("/static", http.FS(staticBox))

	wellknownBox, _ := fs.Sub(StaticBox, "wellknown")
	r.StaticFS("/.well-known", http.FS(wellknownBox))

	locales := []string{"zh-hans", "en-US"}
	for _, locale := range locales {
		restdocBox, _ := fs.Sub(StaticBox, "restdoc/"+locale)
		r.StaticFS("/restdoc/"+locale, http.FS(restdocBox))
	}

	r.GET("/favicon.ico", func(c *gin.Context) {
		filename := "image/favicon.ico"
		f, err := faviconBox.Open(filename)
		if err != nil {
			c.String(404, "page not found")
			return
		}

		iconData, err := ioutil.ReadAll(f)
		if err != nil {
			c.String(404, "page not found")
			return
		}

		c.Header("Content-Type", "image/vnd.microsoft.icon")
		c.String(200, string(iconData))
	})

	//r.StaticFS("/mail", http.Dir("./mail"))
	//r.StaticFS("/mail", pkger.Dir("/mail"))
	//r.StaticFS("/mail", http.Dir("../frontend/dist/restdoc/"))
	return r
}
