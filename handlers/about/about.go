package about

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/utils"
)

func About(c *gin.Context) {

	_timestamp := config.DefaultConfig.VersionTimestamp
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)
	c.HTML(http.StatusOK, "About", gin.H{"_timestamp": _timestamp, "year": year, "saas": saasInfo})
	return
}
