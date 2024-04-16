package team

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"restdoc/config"
	"restdoc/utils"

	Models "github.com/restdoc/restdoc-models"
)

type teamDeleteForm struct {
	TeamId string `form:"id" bind:"required"`
}

func TeamDeletePage(c *gin.Context) {

	s := utils.FormatSession(c)
	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	c.HTML(http.StatusOK, "TeamDelete", gin.H{"_timestamp": _timestamp, "year": year, "login": false, "s": s, "saas": saasInfo})
	return
}

func TeamDelete(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form teamDeleteForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	id := strings.TrimSpace(form.TeamId)
	if id == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "id不能为空"})
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
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "invalid team_id"})
		return
	}

	var team Models.Team
	err = Models.DeleteTeam(&team, tid, uid)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "查询错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "域名删除失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "域名删除成功"})
	}
	return
}
