package restdocEndpoint

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "github.com/restdoc/restdoc-models"
)

type endpointUpdateForm struct {
	Id    string `form:"id" binding:"required"`
	Name  string `form:"name" `
	Value string `form:"value" `
}

func Update(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form endpointUpdateForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	endpoint_id := strings.TrimSpace(form.Id)
	id, err := strconv.ParseInt(endpoint_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong card id "})
		return
	}

	updatedAt := int64(time.Now().Unix())
	updates := map[string]interface{}{"updated_at": updatedAt}

	name := strings.TrimSpace(form.Name)
	value := strings.TrimSpace(form.Value)

	if name == "" && value == "" {
		c.JSON(http.StatusOK, gin.H{"err": "param is empty", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "name and value are empty."})
		return
	}

	if name != "" {
		updates["name"] = name
	}

	if value != "" {
		updates["value"] = value
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

	var re Models.RestEndpoint
	err = Models.GetOneRestEndpoint(&re, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get endpoint error"})
		return
	}

	err = Models.UpdateRestEndpoint(&re, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update endpoint error"})
		return
	}

	project_id := strconv.FormatInt(re.ProjectId, 10)

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         id,
		"user":       u,
		"name":       name,
		"project_id": project_id,
		"updated_at": updatedAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
