package teamuser

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	//"golang.org/x/crypto/bcrypt"

	Models "restdoc-models/models"
	"restdoc/config"
	"restdoc/internal/database/snowflake"
	"restdoc/utils"
)

var emailMatch = regexp.MustCompile(`[.@_a-zA-Z0-9]{1,128}`)

const maxRand = 9223372036854775807

type memberCreateForm struct {
	UserId string `form:"user_id" bind:"required"`
	TeamId string `form:"team_id" bind:"required"`
}

func MemberCreate(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form memberCreateForm
	err := c.ShouldBind(&form)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	teamId := strings.TrimSpace(form.TeamId)

	if teamId == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "team_id不能为空"})
		return
	}

	user_id := strings.TrimSpace(form.UserId)
	if user_id == "" {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "user_id不能为空"})
		return
	}

	tuId, err := strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 1, "message": "invalid user_id"})
		return
	}

	session, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "需要登录"})
		return
	}

	s, ok := session.(Models.Session)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "无效的session"})
		return
	}

	if s.Valid == "false" {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "请先验证邮箱"})
		return
	}

	/*
		userId := s.Id

		uid, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			glog.Error(err)
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
			return
		}
	*/

	var team Models.Team
	err = Models.GetOneTeam(&team, teamId)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "没有找到对应团队"})
		return
	}

	//todo check count limit to 10 for free version
	saasInfo := utils.ExtractSaaSInfo(c)
	if saasInfo.IsSaaS != "true" {
		total := Models.MembersCount()
		if total >= int64(config.DefaultConfig.Count) {
			c.JSON(http.StatusOK, gin.H{"id": "0", "code": 5, "error": "number of members is limited.Please upgrade.", "message": "number limited.", "ts": timestamp})
			return
		}
	}

	_type := int16(0)
	ts := time.Now().Unix()

	_newId, err := snowflake.Sf.NextID()
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"id": "0", "code": 5, "error": "generate mail id error", "message": "generate mail id error", "ts": timestamp})
		return
	}

	newId := int64(_newId)

	tid := team.Id

	//n, _ := rand.Int(rand.Reader, big.NewInt(maxRand))
	//uniq := uint64(n.Int64())

	_id := strconv.FormatInt(newId, 10)
	item := &Models.TeamUser{
		Id:       newId,
		TeamId:   tid,
		UserId:   tuId,
		Valid:    true,
		Type:     _type,
		CreateAt: ts,
		UpdateAt: ts,
	}
	err = Models.AddNewTeamUser(item)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "添加成员失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": _id, "code": 0, "message": "添加成员成功", "data": gin.H{}})
	return
}
