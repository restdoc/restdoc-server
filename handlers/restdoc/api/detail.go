package restdocApi

import (
	"net/http"
	"strconv"
	"time"

	//"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	//"gorm.io/gorm"

	"restdoc/utils"

	Models "github.com/restdoc/restdoc-models"
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

	api_id := c.Param("id")

	if api_id == "" {
		c.JSON(http.StatusOK, gin.H{"err": "id missing", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "id missing"})
		return
	}

	var api Models.RestAPI

	apiId, err := strconv.ParseInt(api_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "invalid id", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "invalid id"})
		return
	}

	err = Models.GetOneRestAPI(&api, apiId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "查询用户失败"})
		return
	}

	groupId := strconv.FormatInt(api.GroupId, 10)
	projectId := strconv.FormatInt(api.ProjectId, 10)
	//principalId := strconv.FormatInt(api.PrincipalId, 10)

	desc := ""

	//get params
	var params []Models.RestParam
	err = Models.GetRestParamsByAPIId(&params, apiId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get params error"})
		return
	}

	getParams := []map[string]interface{}{}
	postFormParams := []map[string]interface{}{}
	headerParams := []map[string]interface{}{}
	for _, param := range params {
		_id := strconv.FormatInt(param.Id, 10)
		if param.Status == Models.RestAPIDeleted {
			continue
		}
		item := map[string]interface{}{
			"id":       _id,
			"key":      param.Name,
			"value":    "",
			"required": param.Required,
			"enabled":  true,
			"desc":     param.Title,
		}
		switch param.Type {
		case Models.PARAM_GET:
			getParams = append(getParams, item)
		case Models.PARAM_POST_FORM:
			postFormParams = append(postFormParams, item)
		case Models.PARAM_HEADER:
			headerParams = append(headerParams, item)
		}

	}

	item := gin.H{
		"id":            api_id,
		"name":          api.Name,
		"status":        api.Status,
		"project_id":    projectId,
		"group_id":      groupId,
		"desc":          desc,
		"get_params":    getParams,
		"form_params":   postFormParams,
		"header_params": headerParams,
	}

	//get profile_Images from mysql
	profileImages := gin.H{}
	result := gin.H{
		"detail": item, "profile_images": profileImages}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "获取成功", "data": result})
	return
}
