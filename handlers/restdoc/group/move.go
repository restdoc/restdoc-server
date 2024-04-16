package restdocGroup

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"

	Models "github.com/restdoc/restdoc-models"
)

type moveListForm struct {
	AfterId  string `form:"after_id" `
	BeforeId string `form:"before_id" `
	GroupId  string `form:"group_id" `
}

func Move(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form moveListForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	before_id := strings.TrimSpace(form.BeforeId)
	after_id := strings.TrimSpace(form.AfterId)
	beforeId := int64(0)
	afterId := int64(0)
	if before_id != "" {
		bid, err := strconv.ParseInt(before_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong before_id"})
			return
		}
		beforeId = bid
	}

	if after_id != "" {
		aid, err := strconv.ParseInt(after_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong after_id"})
			return
		}
		afterId = aid
	}

	group_id := strings.TrimSpace(form.GroupId)
	groupId, err := strconv.ParseInt(group_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong list_id"})
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

	groupIds := []int64{groupId}
	if groupId > 0 {
		groupIds = append(groupIds, groupId)
	}

	if beforeId > 0 {
		groupIds = append(groupIds, beforeId)
	}

	if afterId > 0 {
		groupIds = append(groupIds, afterId)
	}

	var lists []Models.RestGroup

	err = Models.GetRestGroupsByIds(&lists, uid, groupIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get lists error"})
		return
	}

	var current Models.RestGroup
	var before Models.RestGroup
	var after Models.RestGroup
	var hasBefore = false
	var hasAfter = false
	newWeight, _ := new(decimal.Big).SetString("0")
	for _, list := range lists {
		switch list.Id {
		case beforeId:
			before = list
			beforeWeight, _ := new(decimal.Big).SetString(list.Weight)
			newWeight.Add(newWeight, beforeWeight)
			hasBefore = true
		case afterId:
			after = list
			afterWeight, _ := new(decimal.Big).SetString(list.Weight)
			newWeight.Add(newWeight, afterWeight)
			hasAfter = true
		case groupId:
			current = list
		default:
		}
	}

	if hasBefore && hasAfter {

		half, _ := new(decimal.Big).SetString("2")
		newWeight = newWeight.Quo(newWeight, half)

	} else {
		if hasBefore {
			minus1, _ := new(decimal.Big).SetString("1")
			newWeight = newWeight.Sub(newWeight, minus1)
		} else {
			if hasAfter {
				add1, _ := new(decimal.Big).SetString("1")
				newWeight = newWeight.Add(newWeight, add1)
			} else {
				newWeight, _ = new(decimal.Big).SetString(defaultWeight)
			}
		}
	}

	createdAt := int64(time.Now().Unix())

	current.Weight = newWeight.String()

	ts := time.Now().Unix()

	if current.Weight == "0" {
		current.Weight = defaultWeight
	}

	updates := map[string]interface{}{
		"weight":     current.Weight,
		"updated_at": ts,
	}

	if hasBefore || hasAfter {
		if hasBefore && hasAfter {
			if before.ProjectId != current.ProjectId || after.ProjectId != current.ProjectId {
				c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
				return
			}
			if before.ProjectId != after.ProjectId {
				c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
				return
			}
		} else {
			if hasBefore {
				if before.ProjectId != current.ProjectId {
					c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
					return
				}

			} else {
				if after.ProjectId != current.ProjectId {
					c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
					return
				}
			}
		}
	} else {
		current.Weight = defaultWeight
		updates["weight"] = defaultWeight
	}

	err = Models.UpdateRestGroup(&current, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update list weight error"})
		return
	}

	idstr := strconv.FormatInt(current.Id, 10)
	userInfo := gin.H{"id": userId}
	item := gin.H{
		"id":         idstr,
		"user":       userInfo,
		"created_at": createdAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
