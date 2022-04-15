package team

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"

	//"golang.org/x/crypto/bcrypt"

	RestModels "restdoc-models/models"
	"restdoc/internal/database/snowflake"
	"restdoc/utils"
)

type teamCreateForm struct {
	Name string `form:"name" bind:"required"`
	Type int    `form:"type" bind:"required"`
}

func Create(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form teamCreateForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}
	name := strings.TrimSpace(form.Name)

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "name不能为空"})
		return
	}

	session, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "需要登录"})
		return
	}

	s, ok := session.(RestModels.Session)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "无效的session"})
		return
	}

	if s.Valid == "false" {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "请先验证邮箱"})
		return
	}

	if s.Admin != "true" {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "您没有域名添加权限"})
		return
	}

	userId := s.Id

	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	saasInfo := utils.ExtractSaaSInfo(c)

	if saasInfo.IsSaaS != "true" {
		//

	}

	_type := int16(form.Type)
	if _type != 0 && _type != 1 {
		_type = 0
	}

	ts := time.Now().Unix()

	_id, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"id": 0, "code": 1, "error": "generate team id error.", "ts": timestamp})
		return
	}

	id := int64(_id)

	teamItem := &RestModels.Team{
		Id:       id,
		Name:     name,
		UserId:   uid,
		CreateAt: ts,
		UpdateAt: ts,
		Type:     _type,
	}
	err = RestModels.AddNewTeam(teamItem)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "team create error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "团队创建成功", "data": gin.H{}})
	return
}
