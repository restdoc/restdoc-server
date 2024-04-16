package restdocApi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	//"gorm.io/gorm"

	Models "github.com/restdoc/restdoc-models"
)

const dateFormat = "2006-01-02"

type cardUpdateForm struct {
	Id     string `form:"id" binding:"required"`
	Name   string `form:"name" `
	Status string `form:"status" `
	Path   string `form:"path" `
	Method string `form:"method" `
	Color  string `form:"color" `
	Desc   string `form:"desc" `
}

func Update(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form cardUpdateForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	api_id := strings.TrimSpace(form.Id)
	id, err := strconv.ParseInt(api_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong card id "})
		return
	}

	updatedAt := int64(time.Now().Unix())
	updates := map[string]interface{}{}

	name := strings.TrimSpace(form.Name)

	if name != "" {
		updates["name"] = name
	}

	statusId := uint16(0)
	status := strings.TrimSpace(form.Status)
	if status != "" {
		st, err := strconv.ParseInt(status, 10, 64)
		if err == nil {
			statusId = uint16(st)
			updates["status"] = statusId
		}

	}

	path := strings.TrimSpace(form.Path)

	if path != "" {
		updates["path"] = path
	}

	color := strings.TrimSpace(form.Color)
	if color != "" {
		_color := strings.Replace(color, "#", "", -1)
		fmt.Println("_color", _color)
		parsedColor, err := strconv.ParseInt(_color, 16, 32)

		if err == nil {
			updates["color"] = int32(parsedColor)
		} else {
			fmt.Println(err)
		}
	}

	desc := form.Desc

	method := strings.TrimSpace(form.Method)

	if method != "" {
		switch strings.ToLower(method) {
		case "get":
			updates["method"] = Models.METHOD_GET
		case "post":
			updates["method"] = Models.METHOD_POST
		case "option":
			updates["method"] = Models.METHOD_OPTION
		default:
		}
	}

	fmt.Println(updates)
	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "name, status, color and desc are empty."})
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

	aid, err := strconv.ParseInt(api_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "invalid card id"})
		return
	}

	var api Models.RestAPI
	err = Models.GetOneRestAPI(&api, aid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get card error"})
		return
	}

	if len(updates) > 0 {
		updates["updated_at"] = updatedAt
	}

	if len(updates) == 0 {
		if desc != "" {
			/*
				var detail Models.RestAPIDetail
				err = Models.GetOneRestAPIDetail(&detail, api_id)
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						detail.Id = cid
						detail.Desc = desc
						err = Models.AddNewRestAPIDetail(&detail)
						if err != nil {
							c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "failed to add card detail"})
							return

						}
					}

				} else {
					detail.Id = cid
					updates = map[string]interface{}{"desc": desc}
					err = Models.UpdateCardDetail(&detail, updates)
					if err != nil {
						c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
						return
					}
				}
			*/
		}

	} else {
		err = Models.UpdateRestAPI(&api, updates)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
			return
		}
	}

	project_id := strconv.FormatInt(api.ProjectId, 10)
	u := &Models.User{Id: uid}
	item := gin.H{
		"id":         id,
		"user":       u,
		"name":       name,
		"path":       path,
		"method":     method,
		"color":      color,
		"project_id": project_id,
		"updated_at": updatedAt,
	}

	if desc != "" {
		item["desc"] = desc
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
