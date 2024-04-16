package restdocGroup

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"

	"restdoc/internal/database/snowflake"

	Models "github.com/restdoc/restdoc-models"
)

const defaultWeight = "420000000"

type listForm struct {
	Name      string `form:"name" binding:"required"`
	ProjectId string `form:"project_id" binding:"required"`
	AfterId   string `form:"after_id" `
	BeforeId  string `form:"before_id" `
}

func Add(c *gin.Context) {

	//todo after_id sort weight

	now := time.Now()
	timestamp := now.Unix()

	var form listForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	project_id := strings.TrimSpace(form.ProjectId)
	pid, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong project_id in url"})
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

	_id, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
		return
	}

	id := int64(_id)

	weight, _ := new(decimal.Big).SetString(defaultWeight)
	if afterId > 0 {

		if beforeId > 0 {
			listIds := []int64{beforeId, afterId}
			var lists []Models.RestGroup

			err = Models.GetRestGroupsByIds(&lists, uid, listIds)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get lists error"})
				return
			}

			if len(lists) != 2 {
				c.JSON(http.StatusOK, gin.H{"err": "wrong list numbers", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "wrong list numbers"})
				return
			}

			weight, _ = new(decimal.Big).SetString("0")
			for _, list := range lists {
				w := list.Weight
				if w != "" {
					cardWeight, _ := new(decimal.Big).SetString(w)
					weight.Add(weight, cardWeight)
				}
			}
			half, _ := new(decimal.Big).SetString("2")
			weight = weight.Quo(weight, half)

		} else {

			var prevList Models.RestGroup
			err = Models.GetOneRestGroup(&prevList, afterId, uid)
			if err == nil {
				_weight := prevList.Weight
				if _weight == "0" {
					_weight = defaultWeight
				}
				prevWeight, _ := new(decimal.Big).SetString(_weight)
				weight.Add(weight, prevWeight)
			}

		}
	}

	createdAt := int64(time.Now().Unix())

	var ls Models.RestGroup
	ls.Id = id
	ls.UserId = uid
	ls.Name = name
	ls.ProjectId = pid
	//ls.Status = 0
	ls.CreatedAt = createdAt
	ls.Weight = weight.String()
	err = Models.AddNewRestGroup(&ls)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
		return
	}

	//u := &Models.User{ID: uid}
	idstr := strconv.FormatInt(id, 10)
	item := gin.H{
		"id":         idstr,
		"user":       gin.H{"id": userId},
		"name":       name,
		"created_at": createdAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
