package authcode

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	RestModels "restdoc-models/models"
)

func Detail(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	id := int64(0)

	var code RestModels.VerifyCode
	code.Id = id

	err := RestModels.GetOneVerifyCode(&code, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create verify code error"})
		return
	}

	idstr := strconv.FormatInt(id, 10)
	item := gin.H{
		"id":      idstr,
		"content": code.VerifyCode,
	}

	resp := gin.H{"data": item, "code": 0, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
