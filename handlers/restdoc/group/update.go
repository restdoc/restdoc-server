package restdocGroup

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
)

type listUpdateForm struct {
	Id     string `form:"id" binding:"required"`
	Name   string `form:"name" `
	Status string `form:"status" `
}

func Update(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form listUpdateForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	list_id := strings.TrimSpace(form.Id)
	id, err := strconv.ParseInt(list_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong card id "})
		return
	}

	updatedAt := int64(time.Now().Unix())
	updates := map[string]interface{}{"updated_at": updatedAt}

	name := strings.TrimSpace(form.Name)

	if name != "" {
		updates["name"] = name
	}

	statusId := uint8(0)
	status := strings.TrimSpace(form.Status)
	if status != "" {
		st, err := strconv.ParseInt(status, 10, 64)
		if err == nil {
			statusId = uint8(st)
			updates["status"] = statusId
		}

	}

	if name == "" && status == "" {
		c.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "status and name are empty."})
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

	var rg Models.RestGroup
	err = Models.GetOneRestGroup(&rg, id, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get list error"})
		return
	}

	err = Models.UpdateRestGroup(&rg, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
		return
	}

	project_id := strconv.FormatInt(rg.ProjectId, 10)
	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         id,
		"user":       u,
		"name":       name,
		"group_id":   rg.Id,
		"project_id": project_id,
		"updated_at": updatedAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
