package logout

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"restdoc/config"
	redispool "restdoc/internal/database/redis"
	"restdoc/utils"
)

func Logout(c *gin.Context) {

	needJsonResponse := c.Request.Header.Get("json") == "true"
	ts := time.Now().Unix()

	_session_id, err := c.Request.Cookie("session_id")
	if err != nil {
		if needJsonResponse {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "message": err.Error(), "data": gin.H{}})
			return
		}
		return
	}

	if _session_id == nil {
		if needJsonResponse {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 1, "message": "no session id", "data": gin.H{}})
			return
		}
		return
	}

	session_id := _session_id.Value
	_, err = redispool.GetSession(session_id)
	if err != nil {
		if needJsonResponse {
			c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 5, "message": err.Error(), "data": gin.H{}})
			return
		}
	}

	err = redispool.DeleteSession(session_id)
	if err == nil {
		httpOnly, secure := utils.GetCookieFlag(config.DefaultConfig.Debug, c)

		c.SetCookie("session_id", "", 0, "/", "", secure, httpOnly)
	}

	if needJsonResponse {
		c.JSON(http.StatusOK, gin.H{"ts": ts, "code": 0, "message": "OK", "data": gin.H{}})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
