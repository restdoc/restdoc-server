package restdocProject

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	Models "github.com/restdoc/restdoc-models"
)

type projectUpdateForm struct {
	Id        string `form:"id" binding:"required"`
	Name      string `form:"name" `
	Color     string `form:"color" `
	Icon      string `form:"icon" `
	NameColor string `form:"name_color" `
	IconColor string `form:"icon_color" `
}

func Update(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form projectUpdateForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	project_id := strings.TrimSpace(form.Id)
	id, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong card id "})
		return
	}

	updatedAt := int64(time.Now().Unix())
	updates := map[string]interface{}{"updated_at": updatedAt}

	name := strings.TrimSpace(form.Name)
	icon := strings.TrimSpace(form.Icon)
	nameColor := strings.TrimSpace(form.NameColor)
	iconColor := strings.TrimSpace(form.IconColor)
	color := strings.TrimSpace(form.Color)

	if name != "" {
		updates["name"] = name
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

	var pr Models.RestProject
	err = Models.GetOneRestProject(&pr, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get project error"})
		return
	}

	if iconColor != "" && iconColor != pr.IconColor {
		updates["icon_color"] = iconColor
	}

	if icon != "" && icon != pr.Icon {
		updates["icon"] = icon
	}

	if nameColor != "" && nameColor != pr.NameColor {
		updates["name_color"] = nameColor
	}

	if color != "" && color != pr.Color {
		updates["color"] = color
	}

	err = Models.UpdateRestProject(&pr, updates)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "create songlist error"})
		return
	}

	u := &Models.User{Id: uid}

	item := gin.H{
		"id":         id,
		"user":       u,
		"name":       pr.Name,
		"color":      pr.Color,
		"icon":       pr.Icon,
		"icon_color": pr.IconColor,
		"name_color": pr.NameColor,
		"project_id": pr.Id,
		"updated_at": updatedAt,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
