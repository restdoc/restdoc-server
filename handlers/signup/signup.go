package signup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"github.com/gin-contrib/sessions"
	"github.com/Xuanwo/go-locale"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/google/uuid"

	//"github.com/xxtea/xxtea-go/xxtea"

	"restdoc/config"
	"restdoc/consts"
	redispool "restdoc/internal/database/redis"
	"restdoc/internal/database/snowflake"
	"restdoc/third/geetest"
	"restdoc/utils"

	Models "github.com/restdoc/restdoc-models"
)

const (
	captchaID  = "d8bb7647cfee4b47da69607847326b22"
	privateKey = "ab9d414eb8022531f6fb83ded2d7aea0"
	maxRand    = 9223372036854775807
)

type signupForm struct {
	Password    string `form:"password" `
	Email       string `form:"email"`
	BackupEmail string `form:"backupemail"`
	Name        string `form:"name"`
	Admin       string `"form:"admin"`
	Company     string `form:"company"`
	Code        string `form:"code"`
	Challenge   string `form:"geetest_challenge"`
	Validate    string `form:"geetest_validate"`
	Seccode     string `form:"geetest_seccode"`
}

func SignUpPage(c *gin.Context) {

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")
	s := utils.FormatSession(c)

	saasInfo := utils.ExtractSaaSInfo(c)

	c.HTML(http.StatusOK, "Signup", gin.H{"_timestamp": _timestamp, "year": year, "s": s, "saas": saasInfo})
	return
}

func SignupMailPage(c *gin.Context) {
	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")
	s := utils.FormatSession(c)

	saasInfo := utils.ExtractSaaSInfo(c)
	verifyCode := "234234"
	subject := "邮箱注册验证码为：" + verifyCode

	c.HTML(http.StatusOK, "SignupEmail", gin.H{"_timestamp": _timestamp, "year": year, "s": s, "saas": saasInfo, "verifyCode": verifyCode, "subject": subject})
	return
}

func checkField() {

}

func RegisterGeetest(c *gin.Context) {
	ts := time.Now().Unix()
	geetest := geetest.NewGeetestLib(captchaID, privateKey, 2*time.Second)
	//status, response := geetest.PreProcess("", "")
	_, response := geetest.PreProcess("", "")

	/*
		session := sessions.Default(c)
		v, _ := session.Get("geetest").(map[string]interface{})
		if v != nil {
			fmt.Println(v)
		} else {
			v = map[string]interface{}{}
		}
		//v.Set("geetest_status", status)
		v["geetest_status"] = status
		session.Set("geetest", v)
	*/

	var resp map[string]interface{}
	err := json.Unmarshal(response, &resp)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"ts": ts})
		return
	} else {
		c.JSON(http.StatusOK, resp)
		return
	}
}

func SignUp(c *gin.Context) {

	ts := time.Now().Unix()

	var signup signupForm
	err := c.ShouldBind(&signup)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": err.Error(), "data": gin.H{}})
		return
	}

	email := strings.ToLower(strings.TrimSpace(signup.Email))

	name := strings.TrimSpace(signup.Name)
	company := strings.TrimSpace(signup.Company)
	password := strings.TrimSpace(signup.Password)
	code := strings.TrimSpace(signup.Code)
	//challenge := signup.Challenge
	//validate := signup.Validate
	//seccode := signup.Seccode
	isAdmin := true

	if len(code) != 6 {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "code field is invalid.", "data": gin.H{}})
		return
	}

	if isAdmin {
		if name == "" || email == "" {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "name or email field is empty.", "data": gin.H{}})
			return
		}
	}

	if len(email) > 128 {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "email too long", "data": gin.H{}})
		return
	}

	if password == "" {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "password field required.", "data": gin.H{}})
		return
	}

	var verifyItem Models.VerifyCode
	fmt.Println(email)
	fmt.Println(code)
	err = Models.GetVerifyCodeByEmailAndCode(&verifyItem, email, code)
	if err != nil {
		// if not exist
		// if error
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "error": "verify_code invalid", "data": gin.H{}})
		return
	}

	expire := config.DefaultConfig.DefaultExpire
	if expire == 0 {
		expire = 86400
	}
	expire = 0

	ip := c.Request.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(c.Request.RemoteAddr, ":")[0]
	}

	//ipInt := utils.InetAtoN(ip)

	_newId, err := snowflake.Sf.NextID()
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": "generate id error", "message": "generate id error", "data": gin.H{}})
		return
	}
	newId := int64(_newId)

	uniq := newId

	hashedPassword := utils.GetHashedPassword(uniq, password)

	user := &Models.User{
		Id:       newId,
		Name:     name,
		Password: hashedPassword,
		Email:    email,
		Company:  company,
		Valid:    true,
		CreateAt: ts,
		UpdateAt: ts,
	}

	valid := "true"

	localeName := ""
	tag, err := locale.Detect()
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
	} else {
		localeName = tag.String()
	}

	err = Models.AddNewUser(user)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": err.Error(), "data": gin.H{}})
		return
	} else {

		//add default team

		_tid, err := snowflake.Sf.NextID()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": err.Error(), "data": gin.H{}})
			return
		}

		tid := int64(_tid)
		var t Models.Team
		t.Id = tid
		t.UserId = user.Id
		t.Name = fmt.Sprintf("%s%s", user.Name, "的团队")
		t.CreateAt = ts
		t.UpdateAt = ts
		t.Type = 0
		t.Valid = true

		err = Models.AddNewTeam(&t)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": err.Error(), "data": gin.H{}})
			return
		}

		//delete verify code
		vid := strconv.FormatInt(verifyItem.Id, 10)
		_ = Models.DeleteVerifyCode(&verifyItem, vid)

		session_id, err := uuid.NewUUID()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": "generate uuid error", "data": gin.H{}})
			return
		}

		userId := strconv.FormatInt(user.Id, 10)

		session := Models.Session{
			Id:     userId,
			User:   user.Name,
			Email:  user.Email,
			Locale: localeName,
			Valid:  valid,
			Admin:  strconv.FormatBool(isAdmin),
		}

		err = redispool.SetSession(session_id.String(), session, expire)
		if err == nil {
			httpOnly, secure := utils.GetCookieFlag(config.DefaultConfig.Debug, c)

			c.SetCookie(consts.CookieKey, session_id.String(), expire, "/", "", secure, httpOnly)

			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 0, "message": "OK", "data": gin.H{}, "mailed": true})
			return
		} else {
			glog.Error(err)
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "error": err.Error(), "message": err.Error(), "data": gin.H{}})
			return
		}

	}
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
