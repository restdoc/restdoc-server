package signup

import (
	"net/http"
	"strings"
	"time"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"restdoc/config"
	redispool "restdoc/internal/database/redis"
	"restdoc/internal/database/snowflake"
	"restdoc/utils"

	Models "github.com/restdoc/restdoc-models"
)

func GetMailCode(c *gin.Context) {

	ts := time.Now().Unix()

	var signup signupForm
	err := c.ShouldBind(&signup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": err.Error(), "data": gin.H{}})
		return
	}

	email := strings.ToLower(strings.TrimSpace(signup.Email))

	//challenge := signup.Challenge
	//validate := signup.Validate
	//seccode := signup.Seccode

	if email == "" {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "email field is empty.", "data": gin.H{}})
		return
	}

	if len(email) > 128 {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "email too long", "data": gin.H{}})
		return
	}

	/*
		checkd := geetestCheck(c, challenge, validate, seccode)
		if !checkd {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 3, "error": "验证码已失效", "message": "验证码已失效", "data": gin.H{}})
			return
		}
	*/

	expire := config.DefaultConfig.DefaultExpire
	if expire == 0 {
		expire = 86400
	}

	key := email
	verify_code := utils.RandNumberRunes(6)

	ip := c.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}

	var verifyItem Models.VerifyCode
	err = Models.GetVerifyCodeByEmail(&verifyItem, email)
	if err == nil {
		err = Models.UpdateVerifyCode(&verifyItem, email, verify_code)
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": "update code error", "message": err.Error(), "data": gin.H{}})
			return
		}
	} else {
		_id, err := snowflake.Sf.NextID()
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": "generate id error", "message": "generate id error", "data": gin.H{}})
			return
		}
		newId := int64(_id)

		ipInt := utils.InetAtoN(ip)
		verifyItem = Models.VerifyCode{
			Id:         newId,
			Email:      email,
			VerifyCode: verify_code,
			IP:         int64(ipInt),
			CreateAt:   ts,
			UpdateAt:   ts,
		}

		err = Models.AddNewVerifyCode(&verifyItem)
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": err.Error(), "data": gin.H{}})
			return
		}
	}

	from := config.DefaultConfig.FromUser
	subject := "邮箱注册验证码是: " + verify_code
	err = utils.SendSignupEmail("", email, from, verify_code, subject)
	if err == nil {
		err = redispool.SetResendState(key, 3600)
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 0, "message": "OK", "data": gin.H{}, "mailed": true})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 0, "message": "OK", "data": gin.H{}, "mailed": true})
		return
	} else {

		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": err.Error(), "data": gin.H{}, "mailed": false})
		return
	}
}
