package restdocProject

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/internal/database/snowflake"

	Models "github.com/restdoc/restdoc-models"
)

const defaultWeight = "420000000"

type Endpoint struct {
	Name  string `json:name`
	Value string `json:value`
}

type projectForm struct {
	Name      string `form:"name" binding:"required"`
	TeamId    string `form:"team_id" `
	Endpoints string `form:"endpoints"`
}

func Add(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form projectForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	name := strings.TrimSpace(form.Name)
	team_id := strings.TrimSpace(form.TeamId)
	endpointsBody := strings.TrimSpace(form.Endpoints)

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少name参数"})
		return
	}

	if team_id == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少team_id参数"})
		return
	}

	if endpointsBody == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少endpoints参数"})
		return
	}

	raw, err := base64.StdEncoding.DecodeString(endpointsBody)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "parse endpoints error", "code": 1, "message": "parse endpoints error"})
		return
	}

	var endpoints []Endpoint

	err = json.Unmarshal(raw, &endpoints)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "parse endpoints error", "code": 1, "message": "endpoints格式错误"})
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

	_id, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
		return
	}

	id := int64(_id)

	fmt.Println(endpoints)
	//
	var ends []Models.RestEndpoint
	for i := range endpoints {
		endpoint := endpoints[i]
		var end = Models.RestEndpoint{}

		_newId, err := snowflake.Sf.NextID()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
			return
		}

		newId := int64(_newId)
		end.Id = newId
		end.ProjectId = id
		end.Name = endpoint.Name
		end.Value = endpoint.Value
		end.Status = int16(0)
		end.Weight = defaultWeight
		end.CreatedAt = timestamp
		end.UpdatedAt = timestamp
		ends = append(ends, end)
	}

	Models.AddRestEndpoints(&ends)

	teamId := uid
	if team_id != "" {
		teamId, err = strconv.ParseInt(team_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong team_id param"})
			return
		}
	}

	createdAt := int64(time.Now().Unix())

	var pr Models.RestProject
	pr.Id = id
	pr.Name = name
	pr.CreatorId = uid
	pr.TeamId = teamId
	pr.CreatedAt = createdAt
	err = Models.AddNewRestProject(&pr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
		return
	}

	_gid, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate group id"})
		return
	}

	gid := int64(_gid)

	var gp Models.RestGroup
	gp.Id = gid
	gp.UserId = uid
	gp.ProjectId = id
	gp.Name = "default"
	gp.Weight = defaultWeight
	gp.CreatedAt = createdAt
	gp.UpdatedAt = createdAt
	gp.Type = 0

	err = Models.AddNewRestGroup(&gp)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
		return
	}

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         id,
		"user":       u,
		"name":       name,
		"created_at": createdAt,
	}

	//todo add default group

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
