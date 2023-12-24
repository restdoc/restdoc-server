package restdocProject

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
)

type deleteProjectForm struct {
	ProjectId string `form:"id" `
}

func Delete(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form deleteProjectForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	project_id := strings.TrimSpace(form.ProjectId)
	projectId, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong project id"})
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

	//
	//delete lists
	//delete cards
	//delete project

	var apis []Models.RestAPI
	err = Models.DeleteRestAPIsByProjectId(apis, projectId, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "delete cards in project error"})
		return

	}

	var groups []Models.RestGroup
	err = Models.DeleteRestGroupsByProjectId(groups, projectId, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "delete lists in project error"})
		return

	}

	current := Models.RestProject{Id: projectId}

	err = Models.DeleteRestProject(&current, projectId, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update card weight error"})
		return
	}

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":   current.Id,
		"user": u,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
