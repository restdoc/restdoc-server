package teamuser

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/utils"
)

func MemberPage(c *gin.Context) {

	s := utils.FormatSession(c)
	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)

	c.HTML(http.StatusOK, "Member", gin.H{"_timestamp": _timestamp, "year": year, "login": false, "s": s, "saas": saasInfo})
	return
}
