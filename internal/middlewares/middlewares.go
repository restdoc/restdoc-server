package middlewares

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/amalfra/etag"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/consts"
	redispool "restdoc/internal/database/redis"
	"restdoc/logger"
)

const defaultCacheControl = "public, max-age=31536000"

var StaticBox fs.FS

// var etagMap map[string]string
var etagMap sync.Map

func Init() {
	initEtags()
}

func initEtags() {

	//etagMap = map[string]string{}
	box, _ := fs.Sub(StaticBox, ".")
	fs.WalkDir(box, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if path == "." {
			return nil
		}

		f, err := box.Open(path)
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(f)
		content := buf.String()

		fpath := fmt.Sprintf("/%s", strings.ToLower(path))
		if strings.HasSuffix(fpath, "index.html") {
			fpath = strings.Replace(fpath, "index.html", "", 1) // trim index.html
		}

		tag := etag.Generate(content, true)

		//etagMap[fpath] = tag
		etagMap.Store(fpath, tag)
		return nil
	})
}

func StackTrace(all bool) string {
	buf := make([]byte, 20480)
	number := runtime.Stack(buf, all)
	stack := string(buf[:number])
	return string(stack)
}

func checkSign(path string, request_source string, request_sign string, values url.Values) (bool, string) {
	// 检查签名是否正确
	fields, exists := config.DefaultConfig.SignFields[path]
	if !exists {
		return true, ""
	}

	/*
		if request_source == "h5" {
			// 部分接口需要兼容老版本的客户端，所以H5用自己的source，但是不做签名
			isIgnore := config.DefaultConfig.IgnoreSignConfig[path]
			if isIgnore {
				return true, ""
			} else {
				return false, ""
			}
		}
	*/

	if request_source == "" || request_sign == "" {
		return false, ""
	} else {
		raw := ""
		for _, field_name := range fields {
			field_value := ""

			field := string(field_name)
			if field == "source" {
				request_source := values.Get(field)
				field_value, exists = config.DefaultConfig.Sources[request_source]
				if !exists {
					return false, ""
				}
			} else {
				field_value = values.Get(field)
			}
			raw += field_value
		}
		hash := md5.New()
		hash.Write([]byte(raw))
		binhash := hash.Sum(nil)
		hex_hash := make([]byte, 32)
		hex.Encode(hex_hash, binhash)

		signature := string(hex_hash)
		if signature == request_sign {
			return true, signature
		} else {
			return false, signature
		}
	}
}

func Sign(c *gin.Context) {
	path := c.Request.URL.Path

	_, exists := config.DefaultConfig.SignFields[path]
	if !exists {
		c.Next()
		return
	}

	var reqSource string
	var reqSign string

	method := c.Request.Method
	switch method {
	case "GET":
		reqSource = c.Query("source")
		reqSign = c.Query("sign")
	case "POST":
		reqSource = c.PostForm("source")
		reqSign = c.PostForm("sign")
	default:
		reqSource = c.Query("source")
		reqSign = c.Query("sign")
	}

	if reqSource == "" {
		respData := gin.H{
			"code": 3,
			"msg":  "缺少参数source",
		}
		c.JSON(http.StatusOK, respData)
		c.AbortWithStatus(http.StatusOK)
		return
	}

	if reqSign == "" {
		respData := gin.H{
			"code": 3,
			"msg":  "缺少参数sign",
		}
		c.JSON(http.StatusOK, respData)
		c.AbortWithStatus(http.StatusOK)
		return
	}

	if ok, calcSign := checkSign(path, reqSource, reqSign, c.Request.URL.Query()); ok {
		c.Next()
		return
	} else {
		respData := gin.H{
			"code": 4,
			"msg":  "签名错误",
		}
		if config.DefaultConfig.Debug {
			respData["calc_sign"] = calcSign
		}
		c.JSON(http.StatusOK, respData)
		c.AbortWithStatus(http.StatusOK)
		return
	}
}

func GlobalRecover(c *gin.Context) {
	defer func(c *gin.Context) {
		if err := recover(); err != nil {
			stack := StackTrace(false)
			sentry.CaptureMessage(stack)

			logger.Info("异常处理, msg:%s 出错日志 => %s", stack, err)
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  err,
			})
			c.Abort()
		}
	}(c)
	c.Next()
}

func Cors(c *gin.Context) {

	requestOrigin := c.Request.Header.Get("origin")
	origin := config.DefaultConfig.AllowOrigin
	if origin == "" {
		origin = requestOrigin
	}

	path := strings.ToLower(c.Request.URL.Path)
	if strings.HasPrefix(path, "/api") {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "content-type,json")
		//c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains");
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	}
	c.Next()
}

func setCacheHeader(c *gin.Context, path string) {
	if _tag, ok := etagMap.Load(path); ok {
		tag := _tag.(string)
		c.Header("Cache-Control", defaultCacheControl)
		c.Header("Etag", tag)
	}
}

func CacheControl(c *gin.Context) {
	path := strings.ToLower(c.Request.URL.Path)
	if strings.HasSuffix(path, "/") {
		setCacheHeader(c, path)
	} else {
		if strings.HasPrefix(path, "/static") {
			setCacheHeader(c, path)
		} else {
			var extension = filepath.Ext(path)
			switch extension {
			case ".css", ".js", ".png", ".jpeg", ".jpg", ".gif", ".pdf", ".ico", ".woff", ".woff2", ".svg":
				setCacheHeader(c, path)
			default:
			}
		}
	}
	c.Next()
}

func handleSession(c *gin.Context, need_permission bool) {

	needJsonResponse := c.Request.Header.Get("json") == "true"
	if _session_id, err := c.Request.Cookie(consts.CookieKey); err == nil {
		session_id := _session_id.Value
		session, err := redispool.GetSession(session_id)
		if err != nil {
			if need_permission {

				//c.SetSameSite(http.SameSiteLaxMode)
				//c.SetCookie("session_id", "", 10, "/", "", secure, httpOnly)
				if needJsonResponse {
					respData := gin.H{
						"code": 403,
						"msg":  "need login",
					}
					c.JSON(http.StatusOK, respData)
					c.AbortWithStatus(http.StatusOK)
					return
				} else {
					c.Redirect(http.StatusTemporaryRedirect, "/login")
					c.AbortWithStatus(http.StatusTemporaryRedirect)
					return
				}

			}
		} else {
			c.Set("session", session)
		}
		c.Next()
	} else {
		if need_permission {
			if needJsonResponse {
				respData := gin.H{
					"code": 403,
					"msg":  "need login",
				}
				c.JSON(http.StatusOK, respData)
				c.AbortWithStatus(http.StatusOK)
				return
			} else {
				c.Redirect(http.StatusTemporaryRedirect, "/login")
				c.AbortWithStatus(http.StatusTemporaryRedirect)
				return
			}

		} else {
			c.Next()
		}
	}
}

func GetSession(c *gin.Context) {

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	path := strings.TrimSpace(strings.ToLower(c.Request.URL.Path))

	need_permission := false

	basenames := strings.Split(path, ".")
	if len(basenames) > 0 {
		filetype := basenames[len(basenames)-1]
		switch filetype {
		case "css", "js", "png", "jpg", "jpeg", "gif", "html", "woff", "pdf", "woff2":
			need_permission = false
			return
		default:
		}
	}

	arr := strings.Split(path, "/")
	firstPath := ""
	secondPath := ""
	if len(arr) > 1 {
		firstPath = arr[1]
		if len(arr) > 2 {
			secondPath = arr[2]
		}
	} else {
		firstPath = path
	}

	switch firstPath {
	case "", "price", "extension":
		need_permission = false
		handleSession(c, need_permission) // 需要加载session
		return
	case "team", "teamuser", "user", "api", "kanban", "restdoc":
		need_permission = true
		handleSession(c, need_permission)
		return
	}
	switch secondPath {

	}

}
