package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	Models "github.com/restdoc/restdoc-models"
)

func FormatSession(c *gin.Context) Models.Session {
	session := Models.Session{Id: "", User: "", Email: "", Valid: "false", Login: "false"}
	_session, ok := c.Get("session")
	if ok {
		s, ok := _session.(Models.Session)
		if ok {
			session = s
			return session
		} else {
			glog.Error("not get session")
		}
	}
	return session
}
