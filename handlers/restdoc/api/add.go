package restdocApi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"

	"restdoc/internal/database/snowflake"

	RestModels "github.com/restdoc/restdoc-models"
)

const defaultWeight = "420000000"

type cardForm struct {
	Name    string `form:"name" binding:"required"`
	GroupId string `form:"group_id" binding:"required"`
	//
	AfterId  string `form:"after_id" `
	BeforeId string `form:"before_id" `
}

func Add(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form cardForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	group_id := strings.TrimSpace(form.GroupId)
	gid, err := strconv.ParseInt(group_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong group_id in url"})
		return
	}

	name := strings.TrimSpace(form.Name)

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少name参数"})
		return
	}

	afterId := int64(0)
	after_id := strings.TrimSpace(form.AfterId)
	if after_id != "" {
		afId, err := strconv.ParseInt(after_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "参数错误", "code": 1, "message": "无效的after_id"})
			return
		}
		afterId = afId
	}

	beforeId := int64(0)
	before_id := strings.TrimSpace(form.BeforeId)
	if before_id != "" {
		bfId, err := strconv.ParseInt(before_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "参数错误", "code": 1, "message": "无效的before_id"})
			return
		}
		beforeId = bfId
	}

	session, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "Maybe not login"})
		return
	}

	s, ok := session.(RestModels.Session)
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

	_id, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
		return
	}

	id := int64(_id)

	fmt.Println(userId)
	fmt.Println(gid)
	var group RestModels.RestGroup
	err = RestModels.GetOneRestGroup(&group, gid, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "could not get group"})
		return

	}

	weight, _ := new(decimal.Big).SetString(defaultWeight)
	if afterId > 0 {

		if beforeId > 0 {

			cardIds := []int64{beforeId, afterId}
			var cards []RestModels.RestAPI

			err = RestModels.GetRestAPIsByIds(&cards, cardIds)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get cards error"})
				return
			}

			if len(cards) != 2 {
				c.JSON(http.StatusOK, gin.H{"err": "wrong card numbers", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "wrong cards numbers"})
				return
			}

			weight, _ = new(decimal.Big).SetString("0")
			for _, card := range cards {
				w := card.Weight
				if w != "" {
					cardWeight, _ := new(decimal.Big).SetString(w)
					weight.Add(weight, cardWeight)
				}
			}
			half, _ := new(decimal.Big).SetString("2")
			weight = weight.Quo(weight, half)

		} else {

			var prevCard RestModels.RestAPI
			err = RestModels.GetOneRestAPI(&prevCard, afterId)
			if err == nil {
				_weight := prevCard.Weight
				if _weight == "0" {
					_weight = defaultWeight
				}
				prevWeight, _ := new(decimal.Big).SetString(_weight)
				weight.Add(weight, prevWeight)
			}

		}
	}

	createdAt := int64(time.Now().Unix())

	var cd RestModels.RestAPI
	cd.Id = id
	cd.Name = name
	cd.ProjectId = group.ProjectId
	cd.GroupId = group.Id
	cd.CreatorId = uid
	cd.Weight = weight.String()
	cd.Status = 0
	cd.CreatedAt = createdAt
	err = RestModels.AddNewRestAPI(&cd)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
		return
	}

	idstr := strconv.FormatInt(id, 10)
	listId := strconv.FormatInt(group.Id, 10)
	projectId := strconv.FormatInt(group.ProjectId, 10)
	item := gin.H{
		"id":         idstr,
		"user":       gin.H{"id": userId},
		"name":       name,
		"list_id":    listId,
		"project_id": projectId,
		"created_at": createdAt,
		"status":     cd.Status,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
