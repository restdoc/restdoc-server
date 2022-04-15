package privacy

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	"restdoc/utils"
)

func Privacy(c *gin.Context) {

	_timestamp := config.DefaultConfig.VersionTimestamp

	s := utils.FormatSession(c)
	now := time.Now()
	year := now.Format("2006")

	saasInfo := utils.ExtractSaaSInfo(c)

	c.HTML(http.StatusOK, "Privacy", gin.H{"_timestamp": _timestamp, "year": year, "login": false, "s": s, "saas": saasInfo})
	return
}
