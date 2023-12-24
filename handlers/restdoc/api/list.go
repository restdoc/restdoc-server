package restdocApi

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	Models "restdoc-models/models"
	"restdoc/utils"
)

func List(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	s := utils.FormatSession(c)
	_, err := strconv.ParseInt(s.Id, 10, 64)
	if err != nil {
		glog.Error("wrong session user id:", s.Id)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	project_id := c.Query("project_id")
	pid, err := strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "invalid project_id", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong project_id in url"})
		return
	}

	var pr Models.RestProject
	err = Models.GetOneRestProject(&pr, pid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "get project error", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": err.Error()})
		return
	}

	userId := s.Id
	uid, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	fmt.Println(pr.TeamId)
	fmt.Println(s.TeamId)
	var teams []Models.Team
	err = Models.GetTeamsByUserId(&teams, uid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get teams error"})
		return
	}
	teamIds := []string{}
	for _, team := range teams {
		tid := strconv.FormatInt(team.Id, 10)
		teamIds = append(teamIds, tid)
	}
	teamIdString := strings.Join(teamIds, ",")

	permitted := utils.CheckPermission(pr.TeamId, teamIdString)
	if !permitted {
		c.JSON(http.StatusOK, gin.H{"err": "not permitted error", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "no permittion"})
		return
	}

	var apis []Models.RestAPI

	//add team id todo
	//teamIds := []int64{uid}
	err = Models.GetRestAPIsByProjectId(&apis, pid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "查询用户失败"})
		return
	}

	listInfo := map[int64][]int64{}
	//lists := map[int64]*Models.KanbanList{}
	cardsInfo := map[int64]*Models.RestAPI{}

	groupIdMap := map[int64]bool{}

	if len(apis) == 0 {

		var groups []Models.RestGroup
		err = Models.GetRestGroupsByProjectId(&groups, pid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "获取列表失败"})
			return
		}

		sortRestGroup(groups)
		formatGroups := []map[string]interface{}{}
		for i := range groups {
			group := groups[i]
			groupIdStr := strconv.FormatInt(group.Id, 10)
			groupName := group.Name
			projectId := strconv.FormatInt(group.ProjectId, 10)

			formatList := map[string]interface{}{
				"info": gin.H{"id": groupIdStr, "name": groupName, "project_id": projectId},
				"apis": []interface{}{},
			}
			formatGroups = append(formatGroups, formatList)
		}
		results := gin.H{"groups": formatGroups}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取列表成功", "data": results})
		return
	}

	for i := range apis {
		item := apis[i]
		id := item.Id
		if item.Status == Models.RestAPIDeleted {
			continue
		}
		//idStr := strconv.FormatUint(id, 10)
		//creatorId := strconv.FormatUint(item.CreatorId, 10)
		groupId := item.GroupId
		//listIdStr := strconv.FormatUint(listId, 10)
		//listWeight := item.Weight
		_, ok := groupIdMap[groupId]
		if !ok {
			//l := &Models.KanbanList{Id: listId}
			//lists[listId] = l
			groupIdMap[groupId] = true
		}

		cardsInfo[id] = &item

		existList, ok := listInfo[groupId]
		if ok {
			existList = append(existList, id)
			listInfo[groupId] = existList
		} else {
			existList = []int64{id}
			listInfo[groupId] = existList
		}

		cardsInfo[id] = &item
	}

	restGroups := []Models.RestGroup{}
	err = Models.GetRestGroupsByProjectId(&restGroups, pid)
	if err != nil {

	}

	sortRestGroup(restGroups)

	groupIds := []int64{}
	for _, v := range restGroups {
		groupIds = append(groupIds, v.Id)
	}
	if len(groupIds) == 0 {
		formatLists := []map[string]interface{}{}
		results := gin.H{"list": formatLists}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取列表成功", "data": results})
		return
	}

	formatLists := []map[string]interface{}{}
	for _, group := range restGroups {
		groupId := group.Id
		groupIdStr := strconv.FormatInt(group.Id, 10)
		groupName := group.Name
		projectId := strconv.FormatInt(group.ProjectId, 10)

		apis := []*Models.RestAPI{}
		formatedAPIs := []map[string]interface{}{}

		apiIds, ok := listInfo[groupId]
		if ok {
			for _, v := range apiIds {
				api := cardsInfo[v]
				apis = append(apis, api)
			}

			sort.Slice(apis, func(i, j int) bool {

				first := apis[i]
				second := apis[j]
				firstWeight := first.Weight
				secondWeight := second.Weight
				firstV, ok := new(decimal.Big).SetString(firstWeight)
				if !ok {
					return false
				}
				secondV, ok := new(decimal.Big).SetString(secondWeight)
				if !ok {
					return false
				}
				result := firstV.Cmp(secondV)
				return result < 0
			})

			for i, _ := range apis {
				card := apis[i]
				id := strconv.FormatInt(card.Id, 10)
				name := card.Name
				status := card.Status
				method := utils.FormatMethod(card.Method)
				postType := ""
				//url := fmt.Sprintf("%s%s", card.BaseUrl, card.Path)
				path := card.Path
				if status == Models.RestAPIDeleted || status == Models.RestAPIForeverDeleted {
					continue
				}
				formData := []map[string]interface{}{}

				formDataItem := map[string]interface{}{
					"key":         "x",
					"value":       "y",
					"enabled":     true,
					"required":    true,
					"description": " desc x",
				}
				formData = append(formData, formDataItem)

				item := gin.H{
					"id":        id,
					"name":      name,
					"status":    status,
					"method":    method,
					"path":      path,
					"params":    []map[string]interface{}{},
					"form_data": formData,
					"post_type": postType,
				}
				formatedAPIs = append(formatedAPIs, item)
			}
		}

		list := map[string]interface{}{
			"info": gin.H{"id": groupIdStr, "name": groupName, "project_id": projectId},
			"apis": formatedAPIs,
		}
		formatLists = append(formatLists, list)
	}
	/*

		project := map[string]interface{}{
					"name":       item.Name,
					"id":         id,
					"creator_id": creatorId,
					"team_id":    teamId,
				}
				projects = append(projects, project)
	*/

	projectIds := []int64{pid}
	var endpoints []Models.RestEndpoint
	err = Models.GetRestEndpointsByProjectIds(&endpoints, projectIds)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get endpoint error"})
		return
	}

	endpointMaps := utils.FormatEndpoints(endpoints)

	ends, ok := endpointMaps[pid]
	if !ok {
		ends = []map[string]interface{}{}
	}

	results := gin.H{"groups": formatLists, "endpoints": ends}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取列表成功", "data": results})
	return
}

func sortRestGroup(kanbanLists []Models.RestGroup) {
	sort.Slice(kanbanLists, func(i, j int) bool {
		first := kanbanLists[i]
		second := kanbanLists[j]
		firstWeight := first.Weight
		secondWeight := second.Weight
		firstV, ok := new(decimal.Big).SetString(firstWeight)
		if !ok {
			return false
		}
		secondV, ok := new(decimal.Big).SetString(secondWeight)
		if !ok {
			return false
		}
		result := firstV.Cmp(secondV)
		return result < 0
	})
}
