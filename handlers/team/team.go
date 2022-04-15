package team

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/utils"
)

func Page(c *gin.Context) {

	s := utils.FormatSession(c)

	if s.Admin != "true" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	c.HTML(http.StatusOK, "Team", gin.H{"_timestamp": _timestamp, "year": year, "login": false, "s": s, "saas": saasInfo})
	return
}
