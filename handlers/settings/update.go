package settings

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	redispool "restdoc/internal/database/redis"
	Models "restdoc-models/models"
)

type settingsUpdateForm struct {
	Locale string `form:"locale" `
}

func Update(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form settingsUpdateForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	locale := strings.TrimSpace(form.Locale)

	if locale == "" {
		c.JSON(http.StatusOK, gin.H{"error": "empty locale", "code": 1, "message": "参数错误"})
		return
	}

	_, err = language.Parse(locale)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	session_id := ""
	if _session_id, err := c.Request.Cookie("session_id"); err == nil {
		session_id = _session_id.Value
	}

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

	var du Models.TeamUser
	err = Models.GetOneTeamUser(&du, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 5, "message": "get user error"})
		return
	}

	du.Locale = locale

	err = Models.UpdateTeamUserLocale(&du, locale)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 5, "message": "update locale error"})
		return
	}

	s.Locale = locale
	err = redispool.SetSession(session_id, s, 3600*24*7)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "标签缓存失败"})
		return
	}

	updatedAt := int64(time.Now().Unix())
	item := gin.H{
		"user_id":    userId,
		"locale":     locale,
		"updated_at": updatedAt,
	}

	results := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, results)
	return
}
