package team

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"gorm.io/gorm"

	Models "restdoc-models/models"
)

func TeamInfo(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()
	id := c.Param("id")

	session, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "Maybe not login"})
		return
	}

	s, ok := session.(Models.Session)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "Invalid session."})
		return
	}
	userId := s.Id

	_, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		glog.Error("parse user id error ", err)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	var t Models.Team
	err = Models.GetOneTeam(&t, id)

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			glog.Error("get domain by name error ", err)
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "获取域名信息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "域名不存在"})
	} else {

		//defaultSmtpPassword := domain.DefaultSmtpPassword
		info := map[string]interface{}{
			"id":    t.Id,
			"name":  t.Name,
			"type":  t.Type,
			"valid": t.Valid,
		}
		c.JSON(http.StatusOK, gin.H{"data": info, "code": 0, "message": "OK"})
	}
	return
}
