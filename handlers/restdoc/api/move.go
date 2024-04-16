package restdocApi

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"

	Models "github.com/restdoc/restdoc-models"
)

type moveCardForm struct {
	AfterId    string `form:"after_id" `
	BeforeId   string `form:"before_id" `
	CardId     string `form:"card_id" `
	NewGroupId string `form:"new_group_id" `
}

func Move(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form moveCardForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	before_id := strings.TrimSpace(form.BeforeId)
	after_id := strings.TrimSpace(form.AfterId)
	new_group_id := strings.TrimSpace(form.NewGroupId)
	beforeId := int64(0)
	afterId := int64(0)
	newGroupId := int64(0)
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

	card_id := strings.TrimSpace(form.CardId)
	cardId, err := strconv.ParseInt(card_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong card_id"})
		return
	}

	if beforeId == 0 && afterId == 0 {

		if new_group_id != "" {
			gid, err := strconv.ParseInt(new_group_id, 10, 64)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong new_group_id"})
				return
			}
			newGroupId = gid
		}
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

	cardIds := []int64{cardId}
	if cardId > 0 {
		cardIds = append(cardIds, cardId)
	}

	if beforeId > 0 {
		cardIds = append(cardIds, beforeId)
	}

	if afterId > 0 {
		cardIds = append(cardIds, afterId)
	}

	var cards []Models.RestAPI

	err = Models.GetRestAPIsByIds(&cards, cardIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get cards error"})
		return
	}

	var current Models.RestAPI
	var before Models.RestAPI
	var after Models.RestAPI
	var hasBefore = false
	var hasAfter = false
	newWeight, _ := new(decimal.Big).SetString("0")
	for _, card := range cards {
		switch card.Id {
		case beforeId:
			before = card
			beforeWeight, _ := new(decimal.Big).SetString(card.Weight)
			newWeight.Add(newWeight, beforeWeight)
			hasBefore = true
		case afterId:
			after = card
			afterWeight, _ := new(decimal.Big).SetString(card.Weight)
			newWeight.Add(newWeight, afterWeight)
			hasAfter = true
		case cardId:
			current = card
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
			if before.GroupId != after.GroupId {
				c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
				return
			} else {
				updates["group_id"] = before.GroupId
			}
		} else {
			if hasBefore {
				if before.ProjectId != current.ProjectId {
					c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
					return
				}
				updates["group_id"] = before.GroupId

			} else {
				if after.ProjectId != current.ProjectId {
					c.JSON(http.StatusOK, gin.H{"err": "project_id not same", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "project_id not same"})
					return
				}
				updates["group_id"] = after.GroupId
			}
		}
	} else {
		current.Weight = defaultWeight
		current.GroupId = newGroupId
		updates["weight"] = defaultWeight
		updates["group_id"] = newGroupId
	}

	err = Models.UpdateRestAPI(&current, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update card weight error"})
		return
	}

	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         current.Id,
		"user":       u,
		"created_at": createdAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
