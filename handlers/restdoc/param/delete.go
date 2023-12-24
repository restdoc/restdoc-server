package restdocParam

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
)

type deleteParamForm struct {
	Id string `form:"id" `
}

func Delete(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form deleteParamForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	param_id := strings.TrimSpace(form.Id)
	paramId, err := strconv.ParseInt(param_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong param id"})
		return
	}

	session, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "Maybe not login"})
		return
	}

	s, ok := session.(Models.Session)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "Invalid session."})
		return
	}

	userId := s.Id
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	//check permission

	param := Models.RestParam{Id: paramId}

	updatedAt := int64(time.Now().Unix())

	updates := map[string]interface{}{
		"status":     Models.RestAPIDeleted,
		"updated_at": updatedAt,
	}

	err = Models.UpdateParam(&param, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update card weight error"})
		return
	}

	pid := strconv.FormatInt(param.Id, 10)

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         pid,
		"user":       u,
		"created_at": updatedAt,
	}

	//todo: delete card detail

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
