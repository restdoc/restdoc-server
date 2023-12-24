package restdocHome

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Xuanwo/go-locale"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"

	Models "restdoc-models/models"
)

func Locale(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

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

	localeName := s.Locale
	var tag language.Tag

	if localeName != "" {
		t, err := language.Parse(localeName)
		if err == nil {
			tag = t
		}
	}

	if tag.String() == "" {
		t, err := locale.Detect()
		if err != nil {
			//log.Fatal(err)
			fmt.Println(err)
		} else {
			tag = t
		}
	}
	// Have fun with language.Tag!

	path := "/restdoc/zh-hans/"
	switch tag {
	case language.Chinese:
		path = "/restdoc/zh-hans/"
	case language.English:
		path = "/restdoc/en-us/"
	default:
	}

	c.Redirect(http.StatusTemporaryRedirect, path)
	return
}
