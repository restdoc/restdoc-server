package settings

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"golang.org/x/text/language"

	Models "github.com/restdoc/restdoc-models"
)

func Info(c *gin.Context) {

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

	_, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		glog.Error("parse user id error ", err)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	locale := s.Locale
	_, err = language.Parse(locale)
	if err != nil {
		locale = ""
	}

	info := map[string]interface{}{
		"locale": locale,
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "code": 0, "message": "OK"})
}
