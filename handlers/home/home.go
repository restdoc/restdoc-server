package home

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/utils"
)

func Home(c *gin.Context) {

	s := utils.FormatSession(c)

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	c.HTML(http.StatusOK, "Home", gin.H{"_timestamp": _timestamp, "year": year, "s": s, "saas": saasInfo})
	return
}

func Mail(c *gin.Context) {

	s := utils.FormatSession(c)

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)

	c.HTML(http.StatusOK, "Mail", gin.H{"_timestamp": _timestamp, "year": year, "s": s, "saas": saasInfo})
	return
}
