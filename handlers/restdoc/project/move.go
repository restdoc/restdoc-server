package restdocProject

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"

	Models "restdoc-models/models"
	"restdoc/utils"
)

type moveForm struct {
	AfterId   string `form:"after_id" `
	BeforeId  string `form:"before_id" `
	ProjectId string `form:"project_id" `
}

func Move(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form moveForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "code": 1, "message": "参数错误"})
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

	project_id := strings.TrimSpace(form.ProjectId)
	projectId, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong project_id"})
		return
	}

	s := utils.FormatSession(c)

	userId := s.Id
	_, err = strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	projectIds := []int64{projectId}
	if projectId > 0 {
		projectIds = append(projectIds, projectId)
	}

	if beforeId > 0 {
		projectIds = append(projectIds, beforeId)
	}

	if afterId > 0 {
		projectIds = append(projectIds, afterId)
	}

	var projects []Models.RestProject

	err = Models.GetRestProjectsByIds(&projects, projectIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get lists error"})
		return
	}

	//todo: check is self own

	var current Models.RestProject
	var before Models.RestProject
	var hasBefore = false
	var hasAfter = false
	newWeight, _ := new(decimal.Big).SetString("0")
	for _, project := range projects {
		switch project.Id {
		case beforeId:
			before = project
			beforeWeight, _ := new(decimal.Big).SetString(project.Weight)
			newWeight.Add(newWeight, beforeWeight)
			hasBefore = true
		case afterId:
			//after = project
			afterWeight, _ := new(decimal.Big).SetString(project.Weight)
			newWeight.Add(newWeight, afterWeight)
			hasAfter = true
		case projectId:
			current = project
		default:
		}
	}

	if hasBefore && hasAfter {

		half, _ := new(decimal.Big).SetString("2")
		newWeight = newWeight.Quo(newWeight, half)
		if newWeight.String() == defaultWeight {
			add1, _ := new(decimal.Big).SetString("1")
			newWeight = newWeight.Add(newWeight, add1)
		}

		if newWeight.String() == before.Weight {
			add1, _ := new(decimal.Big).SetString("1")
			newWeight = newWeight.Add(newWeight, add1)
		}

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
		}
	} else {
		current.Weight = defaultWeight
		updates["weight"] = defaultWeight
	}

	err = Models.UpdateProject(&current, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update list weight error"})
		return
	}

	_id := strconv.FormatInt(current.Id, 10)
	userInfo := gin.H{"id": userId}
	item := gin.H{
		"id":         _id,
		"user":       userInfo,
		"created_at": createdAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
