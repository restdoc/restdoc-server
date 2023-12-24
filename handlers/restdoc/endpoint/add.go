package restdocEndpoint

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
	"restdoc/internal/database/snowflake"
)

const defaultWeight = "420000000"

type Endpoint struct {
	Id    string `json:id`
	Name  string `json:name`
	Value string `json:value`
}

type endpointForm struct {
	ProjectId string `form:"project_id" binding:"required"`
	Name      string `form:"name" binding:"required"`
	Value     string `form:"value" binding:"required"`
}

func Create(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form endpointForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	project_id := strings.TrimSpace(form.ProjectId)
	name := strings.TrimSpace(form.Name)
	value := strings.TrimSpace(form.Value)

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少name参数"})
		return
	}

	if value == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少value参数"})
		return
	}

	if project_id == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少project_id参数"})
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

	projectId, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong project_id"})
		return
	}

	userId := s.Id
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	_id, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
		return
	}

	id := int64(_id)

	var end = Models.RestEndpoint{}

	_newId, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
		return
	}

	newId := int64(_newId)
	end.Id = newId
	end.ProjectId = projectId
	end.Name = name
	end.Value = value
	end.Status = int16(0)
	end.Weight = defaultWeight
	end.CreatedAt = timestamp
	end.UpdatedAt = timestamp

	Models.AddNewRestEndpoint(&end)

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         id,
		"user":       u,
		"name":       name,
		"value":      name,
		"project_id": project_id,
	}

	//todo add default group

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
