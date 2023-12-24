package restdocEndpoint

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
)

type deleteEndpointForm struct {
	Id string `form:"id" `
}

func Delete(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form deleteEndpointForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	endpoint_id := strings.TrimSpace(form.Id)
	endpointId, err := strconv.ParseInt(endpoint_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong endpoint id"})
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

	endpoint := Models.RestEndpoint{Id: endpointId}

	updatedAt := int64(time.Now().Unix())

	updates := map[string]interface{}{
		"status":     Models.RestAPIDeleted,
		"updated_at": updatedAt,
	}

	err = Models.UpdateRestEndpoint(&endpoint, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update card weight error"})
		return
	}

	eid := strconv.FormatInt(endpoint.Id, 10)

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         eid,
		"user":       u,
		"created_at": updatedAt,
	}

	//todo: delete card detail

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
