package teamuser

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"gorm.io/gorm"

	Models "restdoc-models/models"
)

func MemberInfo(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()
	duId := c.Param("id")

	id, err := strconv.ParseInt(duId, 10, 64)
	if err != nil {
		glog.Error("parse id error ", err)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "invalid id"})
		return
	}

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

	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		glog.Error("parse user id error ", err)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	var tu Models.TeamUser

	err = Models.GetOneTeamUser(&tu, id)

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			glog.Error("get domain user inbo error ", err)
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 5, "message": "获取用户信息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 5, "message": "域名不存在"})
	} else {
		if tu.Id != uid {
			c.JSON(http.StatusOK, gin.H{"error": "Access denied", "code": 5, "message": "访问受限"})
			return
		}

		info := map[string]interface{}{
			"id": tu.Id,
		}
		c.JSON(http.StatusOK, gin.H{"data": info, "code": 0, "message": "OK"})
	}
	return
}
