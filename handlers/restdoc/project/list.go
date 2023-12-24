package restdocProject

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	Models "restdoc-models/models"
	"restdoc/utils"
)

func List(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	s := utils.FormatSession(c)
	fmt.Println(s.Id)
	uid, err := strconv.ParseInt(s.Id, 10, 64)
	if err != nil {
		glog.Error("wrong session user id:", s.Id)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	var teams []Models.Team

	err = Models.GetTeamsByUserId(&teams, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get team error"})
		return
	}

	teamIds := []int64{}
	for i := range teams {
		team := teams[i]
		teamId := team.Id
		teamIds = append(teamIds, teamId)
	}

	var list []Models.RestProject

	//add team id todo
	//get teamids by user id

	err = Models.GetHomeRestProjects(&list, teamIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "查询用户失败"})
		return
	}

	projectIds := []int64{}
	for i := range list {
		item := list[i]
		projectIds = append(projectIds, item.Id)
	}

	//get endpoints by project_id
	var endpoints []Models.RestEndpoint
	err = Models.GetRestEndpointsByProjectIds(&endpoints, projectIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get endpoint error"})
		return
	}

	endpointMaps := utils.FormatEndpoints(endpoints)

	projects := []map[string]interface{}{}
	for i := range list {
		item := list[i]
		creatorId := strconv.FormatInt(item.CreatorId, 10)
		teamId := strconv.FormatInt(item.TeamId, 10)
		projectId := item.Id
		id := strconv.FormatInt(projectId, 10)

		endpoints, ok := endpointMaps[projectId]
		if !ok {
			continue
		}

		project := map[string]interface{}{
			"name":       item.Name,
			"id":         id,
			"creator_id": creatorId,
			"team_id":    teamId,
			"endpoints":  endpoints,
		}
		projects = append(projects, project)
	}
	results := gin.H{"list": projects}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取列表成功", "data": results})
	return
}
