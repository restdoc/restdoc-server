package authcode

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	RestModels "github.com/restdoc/restdoc-models"
)

type cardForm struct {
	Payload string `form:"payload" binding:"required"`
}

func Add(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form cardForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	payload := strings.TrimSpace(form.Payload)
	if payload == "" {
		c.JSON(http.StatusOK, gin.H{"err": "payload is empty", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "payload is empty"})
		return
	}

	id := int64(0)

	var code RestModels.VerifyCode
	code.Id = id
	code.Email = "x@x.com"
	code.VerifyCode = payload
	code.IP = 0
	code.UpdateAt = timestamp
	code.CreateAt = timestamp
	code.Type = 0

	err = RestModels.UpdateVerifyCode(&code, code.Email, code.VerifyCode)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create verify code error"})
		return
	}

	idstr := strconv.FormatInt(id, 10)
	item := gin.H{
		"id": idstr,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
