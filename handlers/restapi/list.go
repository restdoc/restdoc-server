package restapi

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	Models "restdoc-models/models"
)

func List(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	session, ok := c.Get("session")
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "Maybe not login"})
		return
	}

	s, ok := session.(Models.Session)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"ts": timestamp, "data": gin.H{}, "code": 1, "message": "Invalid session."})
		return
	}
	userId := s.Id
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "ts": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	var domains []Models.Domain
	err = Models.GetDomainsByUserId(&domains, uid, 0, 100)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "查询错误"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "domain create error"})
		return
	} else {
		results := []interface{}{}
		for i := range domains {
			domain := domains[i]
			url := "/domain/detail/" + domain.Domain
			_type := domain.Type
			item := map[string]interface{}{"domain": domain.Domain, "url": url, "valid": domain.Valid, "type": _type}
			results = append(results, item)
		}
		c.JSON(http.StatusOK, gin.H{"data": results, "code": 0, "message": "OK"})
		return
	}
	return
}
