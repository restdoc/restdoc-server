package utils

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	isd "github.com/jbenet/go-is-domain"

	"restdoc/config"
)

const saasDomain = "restdoc.com"

type SaaSInfo struct {
	SaasDomain string
	HostDomain string
	IsSaaS     string
}

func init() {
	var _ = reflect.TypeOf(SaaSInfo{})
}

func ExtractSaaSInfo(c *gin.Context) SaaSInfo {

	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	hostdomain := host
	isSaaS := false

	switch hostdomain {
	case "www.restdoc.com", "restdoc.com":
		if config.DefaultConfig.SaaSDomain == saasDomain {
			isSaaS = true
		}
	default:
		host = "restdoc.com"
	}

	isSaaS = true

	hostdomain = scheme + "://" + host
	isSaaSMode := strconv.FormatBool(isSaaS)

	info := SaaSInfo{HostDomain: hostdomain, IsSaaS: isSaaSMode}
	return info
}

func IsSaaS(domain string, host string) bool {
	if !isd.IsDomain(host) {
		return false
	}
	return domain == saasDomain
}

func GetCookieFlag(debug bool, c *gin.Context) (bool, bool) {

	secure := true
	if config.DefaultConfig.Debug == true {
		if c.Request.TLS == nil {
			secure = false
		}
	}
	httpOnly := true
	return httpOnly, secure
}
