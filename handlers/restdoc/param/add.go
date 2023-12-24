package restdocParam

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/gin-gonic/gin"

	RestModels "restdoc-models/models"
	"restdoc/internal/database/snowflake"
	"restdoc/utils"
)

const defaultWeight = "420000000"

type paramForm struct {
	Name         string `form:"name" binding:"required"`
	APIId        string `form:"api_id" binding:"required"`
	DefaultValue string `form:"default_value" `
	Required     string `form:"required" `
	Type         string `form:"type"`
	Enabled      string `form:"enabled" `
	AfterId      string `form:"after_id" `
	BeforeId     string `form:"before_id" `
}

func Create(c *gin.Context) {

	now := time.Now()
	timestamp := now.Unix()

	var form paramForm
	err := c.ShouldBind(&form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error(), "code": 1, "message": "参数错误"})
		return
	}

	s := utils.FormatSession(c)
	_, err = strconv.ParseInt(s.Id, 10, 64)
	if err != nil {
		//glog.Error("wrong session user id:", s.Id)
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	required := false
	_required := strings.TrimSpace(form.Required)
	if _required == "true" {
		required = true
	}

	defaultValue := ""
	_default := strings.TrimSpace(form.DefaultValue)
	if _default != "" {
		defaultValue = _default
	}

	_type := strings.TrimSpace(form.Type)
	paramType := RestModels.PARAM_GET
	switch _type {
	case "get":
		paramType = RestModels.PARAM_GET
	case "form", "post_form":
		paramType = RestModels.PARAM_POST_FORM
	case "header":
		paramType = RestModels.PARAM_HEADER
	default:
	}
	api_id := strings.TrimSpace(form.APIId)
	aid, err := strconv.ParseInt(api_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong group_id in url"})
		return
	}

	//get api
	var api RestModels.RestAPI
	err = RestModels.GetOneRestAPI(&api, aid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get api error"})
		return
	}

	var pr RestModels.RestProject
	err = RestModels.GetOneRestProject(&pr, api.ProjectId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "get project error"})
		return
	}

	permitted := utils.CheckPermission(pr.TeamId, s.TeamId)
	if !permitted {
		c.JSON(http.StatusOK, gin.H{"err": "not permitted", "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "permission denied"})
		return

	}

	name := strings.TrimSpace(form.Name)

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"error": "参数缺失", "code": 1, "message": "缺少name参数"})
		return
	}

	afterId := int64(0)
	after_id := strings.TrimSpace(form.AfterId)
	if after_id != "" {
		afId, err := strconv.ParseInt(after_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "参数错误", "code": 1, "message": "无效的after_id"})
			return
		}
		afterId = afId
	}

	beforeId := int64(0)
	before_id := strings.TrimSpace(form.BeforeId)
	if before_id != "" {
		bfId, err := strconv.ParseInt(before_id, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"error": "参数错误", "code": 1, "message": "无效的before_id"})
			return
		}
		beforeId = bfId
	}

	userId := s.Id
	_, err = strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "wrong user_id in session"})
		return
	}

	_id, err := snowflake.Sf.NextID()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 1, "message": "could not generate id"})
		return
	}

	id := int64(_id)

	weight, _ := new(decimal.Big).SetString(defaultWeight)
	if afterId > 0 {

		if beforeId > 0 {

			pmIds := []int64{beforeId, afterId}
			var pms []RestModels.RestParam

			err = RestModels.GetRestParamsByIds(&pms, pmIds)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "get cards error"})
				return
			}

			if len(pmIds) != 2 {
				c.JSON(http.StatusOK, gin.H{"err": "wrong params numbers", "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "wrong params numbers"})
				return
			}

			weight, _ = new(decimal.Big).SetString("0")
			for _, param := range pms {
				w := param.Weight
				if w != "" {
					paramWeight, _ := new(decimal.Big).SetString(w)
					weight.Add(weight, paramWeight)
				}
			}
			half, _ := new(decimal.Big).SetString("2")
			weight = weight.Quo(weight, half)

		} else {

			var prevParam RestModels.RestParam
			err = RestModels.GetOneRestParam(&prevParam, afterId)
			if err == nil {
				_weight := prevParam.Weight
				if _weight == "0" {
					_weight = defaultWeight
				}
				prevWeight, _ := new(decimal.Big).SetString(_weight)
				weight.Add(weight, prevWeight)
			}

		}
	}

	createdAt := int64(time.Now().Unix())

	var pm RestModels.RestParam
	pm.Id = id
	pm.ApiId = aid
	pm.Name = name
	pm.Title = ""
	pm.Default = defaultValue
	pm.Required = required
	pm.Weight = weight.String()
	pm.Status = 0
	pm.Type = paramType
	pm.CreatedAt = createdAt

	err = RestModels.AddNewRestParam(&pm)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error(), "timestamp": timestamp, "data": gin.H{}, "code": 5, "message": "could not get group"})
		return

	}

	idstr := strconv.FormatInt(id, 10)
	apiId := strconv.FormatInt(aid, 10)
	projectId := strconv.FormatInt(api.ProjectId, 10)
	item := gin.H{
		"id":         idstr,
		"user":       gin.H{"id": userId},
		"name":       name,
		"api_id":     apiId,
		"project_id": projectId,
		"created_at": createdAt,
		"status":     pm.Status,
	}

	resp := gin.H{"data": gin.H{"detail": item}, "code": 0, "ts": timestamp, "msg": "OK"}
	c.JSON(http.StatusOK, resp)
	return
}
