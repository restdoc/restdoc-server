package restdocGroup

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "github.com/restdoc/restdoc-models"
)

type deleteForm struct {
	Id string `form:"id" `
}

func Delete(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form deleteForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	list_id := strings.TrimSpace(form.Id)
	listId, err := strconv.ParseInt(list_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong list id"})
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

	fmt.Println(uid)
	fmt.Println(listId)

	//
	//delete lists
	//delete cards
	//delete project

	/*
		err = Models.DeleteKanbanCardsByProjectId(listId, uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "delete cards in project error"})
			return

		}

		err = Models.DeleteKanbanListsByProjectId(listId, uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "delete lists in project error"})
			return

		}

		current := Models.KanbanProject{Id: projectId}

		err = Models.DeleteProject(&current, projectId, uid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "update card weight error"})
			return
		}

		u := &Models.User{ID: uid}
		item := gin.H{
			"id":   current.Id,
			"user": u,
		}

		resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	*/
	resp := gin.H{}
	c.JSON(http.StatusOK, resp)
	return
}
