package login

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/google/uuid"

	Models "restdoc-models/models"
	"restdoc/config"
	redispool "restdoc/internal/database/redis"
	"restdoc/utils"
)

type loginForm struct {
	Password string `form:"password" `
	Email    string `form:"email"`
}

func LoginPage(c *gin.Context) {

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	c.HTML(http.StatusOK, "Login", gin.H{"_timestamp": _timestamp, "year": year, "login": false, "saas": saasInfo})
	return
}

func Login(c *gin.Context) {

	now := time.Now()
	year := now.Format("2006")
	ts := now.Unix()

	if _session_id, err := c.Request.Cookie("session_id"); err == nil {
		session_id := _session_id.Value
		session, err := redispool.GetSession(session_id)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "data": gin.H{"cookie": true, "user": session}, "code": 0, "message": "OK"})
			return
		}
	}

	var login loginForm
	err := c.ShouldBind(&login)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	email := strings.TrimSpace(login.Email)
	password := strings.TrimSpace(login.Password)
	email = strings.ToLower(email)

	if email == "" {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "error": "email field required.", "data": gin.H{}, "code": 1, "message": "请输入邮箱"})
		return
	}

	uniq := int64(0)
	userPassword := ""
	uid := int64(0)
	_valid := false
	_name := ""
	_email := ""
	localeName := ""

	var user Models.User
	err = Models.GetUserByEmail(&user, email)
	if err != nil {
		glog.Error("get user by email error: ", err)
		c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "error": "user not exist.", "data": gin.H{}, "code": 2, "message": "帐号或密码错误"})
		return
	}

	uniq = user.Id
	userPassword = user.Password
	uid = user.Id
	_valid = user.Valid
	_name = user.Name
	_email = user.Email
	localeName = user.Locale
	isAdmin := true

	hashed := utils.GetHashedPassword(uniq, password)
	if hashed != userPassword {
		glog.Error("password is wrong error: ", err)
		c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "error": "wrong password.", "data": gin.H{}, "code": 2, "message": "帐号或密码错误"})
		return
	} else {
		session_id, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "error": err.Error(), "data": gin.H{}, "code": 1, "message": "生成uuid发生错误"})
			return
		} else {
			userId := strconv.FormatInt(uid, 10)
			valid := strconv.FormatBool(_valid)

			session := Models.Session{
				Id:     userId,
				User:   _name,
				Email:  _email,
				Valid:  valid,
				Locale: localeName,
				Login:  "true",
				Admin:  strconv.FormatBool(isAdmin),
			}
			err = redispool.SetSession(session_id.String(), session, 86400*7)
			if err == nil {
				httpOnly, secure := utils.GetCookieFlag(config.DefaultConfig.Debug, c)
				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie("session_id", session_id.String(), 86400*7, "/", "", secure, httpOnly)
			} else {
				c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "error": err.Error(), "data": gin.H{}, "code": 1, "message": "设置session失败"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"ts": ts, "year": year, "data": gin.H{}, "code": 0, "message": "成功登录"})
	}
	return
}
