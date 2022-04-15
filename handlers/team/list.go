package team

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	Models "restdoc-models/models"

	"restdoc/utils"
)

func List(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

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

	var teams []Models.Team

	isSaaS := false

	saasInfo := utils.ExtractSaaSInfo(c)
	if saasInfo.IsSaaS == "true" {
		isSaaS = true
	}

	err = Models.GetTeamsByUserId(&teams, uid)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "查询错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "domain create error"})
		return
	} else {
		results := []interface{}{}
		if isSaaS {

			for i := range teams {
				team := teams[i]
				url := fmt.Sprintf("/team/detail/%d", team.Id)
				tid := strconv.FormatInt(team.Id, 10)
				_type := team.Type
				item := map[string]interface{}{"id": tid, "name": team.Name, "url": url, "valid": team.Valid, "type": _type}
				results = append(results, item)
			}

		} else {

			for i := range teams {
				team := teams[i]
				tid := strconv.FormatInt(team.Id, 10)
				url := fmt.Sprintf("/team/detail/%d", team.Id)
				_type := team.Type
				item := map[string]interface{}{"id": tid, "name": team.Name, "url": url, "valid": team.Valid, "type": _type}
				results = append(results, item)
			}

		}

		data := gin.H{"list": results}

		c.JSON(http.StatusOK, gin.H{"data": data, "code": 0, "message": "OK"})
		return
	}
	return
}
