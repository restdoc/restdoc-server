package teamuser

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"restdoc/config"

	Models "github.com/restdoc/restdoc-models"
)

func MemberList(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	id := c.Param("id")

	_, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "Maybe not login"})
		return
	}

	//check is admin

	/*
		s, ok := session.(Models.Session)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "Invalid session."})
			return
		}
	*/

	/*
		userId := s.Id
		uid, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
			return
		}
	*/

	var team Models.Team
	err := Models.GetOneTeam(&team, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "没有找到对应域名"})
		return
	}

	var members []Models.TeamUser
	err = Models.GetTeamUsersByTeamId(&members, id)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "查询错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "获取成员列表失败"})
		return
	} else {
		results := []interface{}{}
		max := config.DefaultConfig.Count
		for i := range members {
			if i > max {
				break
			}
			member := members[i]
			url := fmt.Sprintf("/member/%s/detail/%d", member.TeamId, member.Id)
			_type := member.Type
			mid := strconv.FormatInt(member.Id, 10)
			item := map[string]interface{}{
				"id":       mid,
				"url":      url,
				"valid":    member.Valid,
				"type":     _type,
				"password": "********",
			}
			results = append(results, item)
		}
		c.JSON(http.StatusOK, gin.H{"data": results, "code": 0, "message": "OK"})
		return
	}
	return
}
