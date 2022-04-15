package team

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/utils"
)

func DetailPage(c *gin.Context) {

	s := utils.FormatSession(c)

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	c.HTML(http.StatusOK, "DomainDetail", gin.H{"_timestamp": _timestamp, "year": year, "s": s, "saas": saasInfo})
	return
}
