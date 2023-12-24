package restdocProject

import (
	"net/http"
	"strconv"
	"time"

	//"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	//"gorm.io/gorm"

	Models "restdoc-models/models"
	"restdoc/utils"
)

func Detail(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	s := utils.FormatSession(c)
	_, err := strconv.ParseInt(s.Id, 10, 64)
	if err != nil {
		glog.Error("wrong session user id:", s.Id)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	project_id := c.Param("id")

	if project_id == "" {
		c.JSON(http.StatusOK, gin.H{"err": "id missing", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "id missing"})
		return
	}

	var project Models.RestProject

	projectId, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "invalid id", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "invalid id"})
		return
	}

	err = Models.GetOneRestProject(&project, projectId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get project error"})
		return
	}

	//principalId := strconv.FormatInt(api.PrincipalId, 10)

	//get endpoints

	projectIds := []int64{projectId}
	var endpoints []Models.RestEndpoint
	err = Models.GetRestEndpointsByProjectIds(&endpoints, projectIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get endpoint error"})
		return
	}

	endpointMaps := utils.FormatEndpoints(endpoints)

	ends, ok := endpointMaps[projectId]
	if !ok {
		ends = []map[string]interface{}{}
	}

	item := gin.H{
		"id":         project_id,
		"name":       project.Name,
		"status":     project.Status,
		"project_id": project_id,
		"endpoints":  ends,
	}

	result := gin.H{"detail": item}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取成功", "data": result})
	return
}
