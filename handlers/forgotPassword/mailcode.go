package forgotPassword

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	//"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	Models "restdoc-models/models"
	"restdoc/config"
	"restdoc/internal/database/redis"
	"restdoc/internal/database/snowflake"
	"restdoc/utils"
)

const forgotKey = "forgot_"

func GetForgotPasswordMailCode(c *gin.Context) {

	ts := time.Now().Unix()

	var form resetForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": err.Error(), "data": gin.H{}})
		return
	}

	name := strings.TrimSpace(form.Name)
	email := strings.ToLower(strings.TrimSpace(form.Email))
	password := strings.TrimSpace(form.Password)

	//challenge := form.Challenge
	//validate := form.Validate
	//seccode := form.Seccode

	if email == "" {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "email field is empty.", "data": gin.H{}})
		return
	}

	if len(email) > 128 {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "email too long", "data": gin.H{}})
		return
	}

	if password == "" {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "password field required.", "data": gin.H{}})
		return
	}

	expire := config.DefaultConfig.DefaultExpire
	if expire == 0 {
		expire = 86400
	}

	/*
		checkd := geetestCheck(c, challenge, validate, seccode)
		if !checkd {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 3, "error": "验证码已失效", "message": "验证码已失效", "data": gin.H{}})
			return
		}
	*/

	key := fmt.Sprintf("%s-%s", forgotKey, email)
	verify_code := utils.LowerCaseRandStringRunes(6)

	ip := c.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}

	var user Models.User
	err = Models.GetUserByEmail(&user, email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "not found", "message": "not exist", "data": gin.H{"email": email}})
		return
	}

	if user.Name != name {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "not match", "message": "not match", "data": gin.H{"name": name}})
		return
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
	subject := "密码重置验证码是: " + verify_code
	err = utils.SendForgotPasswordEmail(email, from, verify_code, subject)
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
