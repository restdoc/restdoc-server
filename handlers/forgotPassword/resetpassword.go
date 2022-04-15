package forgotPassword

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	Models "restdoc-models/models"
	"restdoc/config"
	redispool "restdoc/internal/database/redis"
	"restdoc/third/geetest"
	"restdoc/utils"
)

const (
	captchaID  = "d8bb7647cfee4b47da69607847326b22"
	privateKey = "ab9d414eb8022531f6fb83ded2d7aea0"
	maxRand    = 9223372036854775807
)

const personalDomain = "feetmail.com"

type resetForm struct {
	Email       string `form:"email" `
	Password    string `form:"password" `
	BackupEmail string `form:"backupemail"`
	Name        string `form:"name"`
	Admin       string `"form:"admin"`
	Company     string `form:"company"`
	Code        string `form:"code"`
	Challenge   string `form:"geetest_challenge"`
	Validate    string `form:"geetest_validate"`
	Seccode     string `form:"geetest_seccode"`
}

func geetestCheck(c *gin.Context, challenge string, validate string, seccode string) bool {
	gt := geetest.NewGeetestLib(captchaID, privateKey, 2*time.Second)
	res := make(map[string]interface{})

	//todo get status from session
	status := 1
	var geetestRes bool
	if status == 1 {
		geetestRes = gt.SuccessValidate(challenge, validate, seccode, "", "")
	} else {
		geetestRes = gt.FailbackValidate(challenge, validate, seccode)
	}
	if geetestRes {
		res["code"] = 0
		res["msg"] = "Success"
	} else {
		res["code"] = -100
		res["msg"] = "Failed"
	}
	//response, _ := json.Marshal(res)
	return geetestRes
}

func ForgotPasswordPage(c *gin.Context) {

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	//todo
	c.HTML(http.StatusOK, "ForgotPassword", gin.H{"_timestamp": _timestamp, "year": year, "saas": saasInfo})
	return
}

func ForgotPasswordMailPage(c *gin.Context) {
	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")
	s := utils.FormatSession(c)

	saasInfo := utils.ExtractSaaSInfo(c)
	verifyCode := "234234"
	subject := "密码重置验证码为：" + verifyCode

	c.HTML(http.StatusOK, "ForgotPasswordEmail", gin.H{"_timestamp": _timestamp, "year": year, "s": s, "saas": saasInfo, "verifyCode": verifyCode, "subject": subject})
	return

}

func ResetPassword(c *gin.Context) {

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	ts := now.Unix()
	year := now.Format("2006")

	var reset resetForm
	err := c.ShouldBind(&reset)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": err.Error(), "error": err.Error()})
		return
	}
	email := strings.ToLower(strings.TrimSpace(reset.Email))
	code := strings.TrimSpace(reset.Code)
	password := strings.TrimSpace(reset.Password)

	if email == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "email field is empty", "error": "email field is empty"})
		return
	}

	if password == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "password field is empty", "error": "password field is empty"})
		return
	}

	if code == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "code field is empty", "error": "code field is empty"})
		return
	}

	cacheKey := fmt.Sprintf("resetPassword-%s", email)
	exist, err := redispool.GetResendState(cacheKey)
	if err != nil {
		//todo
	}

	if exist == "true" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "error": "密码重置太频繁", "message": "密码重置太频繁", "_timestamp": _timestamp, "year": year})
		return
	}

	var verifyItem Models.VerifyCode
	err = Models.GetVerifyCodeByEmailAndCode(&verifyItem, email, code)
	if err != nil {
		// if not exist
		// if error
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{"ts": _timestamp, "code": 1, "error": "verify_code invalid", "message": "verify_code invalid", "data": gin.H{}})
		return
	}

	//todo
	if (ts - verifyItem.CreateAt) > 3600 {
		//expire
		_ = redispool.DeleteResetpasswordState(cacheKey) // 删除key
		//delete verify code from mysql
		vid := strconv.FormatInt(verifyItem.Id, 10)
		_ = Models.DeleteVerifyCode(&verifyItem, vid)

		c.JSON(http.StatusOK, gin.H{"ts": _timestamp, "code": 1, "error": "verify_code expired", "message": "verify_code expired", "data": gin.H{}})
		return
	}
	//get verify code
	//if has verify code   update user set new password = x where email = y
	//

	var user Models.User
	err = Models.GetUserByEmail(&user, email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5, "message": "Get user error", "error": "Get user error"})
		return
	}

	uniq := user.Id
	hashedPassword := utils.GetHashedPassword(uniq, password)
	err = Models.UpdatePassword(&user, email, hashedPassword)
	if err == nil {

		_ = redispool.DeleteResetpasswordState(cacheKey) // 删除key
		//delete verify code from mysql
		vid := strconv.FormatInt(verifyItem.Id, 10)
		_ = Models.DeleteVerifyCode(&verifyItem, vid)

		c.JSON(http.StatusOK, gin.H{"code": 0, "_timestamp": _timestamp, "year": year, "data": gin.H{"changed": true}})
	} else {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": err.Error(), "error": err.Error(), "_timestamp": _timestamp, "year": year})
	}

}
