package utils

import (
	Models "restdoc-models/models"
	"sort"
	"strconv"

	"github.com/ericlagergren/decimal"
)

func FormatEndpoints(endpoints []Models.RestEndpoint) map[int64][]map[string]interface{} {

	sortRestEndpoints(endpoints)

	endpointMaps := map[int64][]map[string]interface{}{}

	for i := range endpoints {
		endpoint := endpoints[i]
		if endpoint.Status == Models.RestAPIDeleted {
			continue
		}

		id := strconv.FormatInt(endpoint.Id, 10)
		projectId := endpoint.ProjectId
		if _, ok := endpointMaps[endpoint.ProjectId]; !ok {
			arr := []map[string]interface{}{}
			item := map[string]interface{}{"id": id, "name": endpoint.Name, "value": endpoint.Value}
			arr = append(arr, item)
			endpointMaps[projectId] = arr
		} else {
			item := map[string]interface{}{"id": id, "name": endpoint.Name, "value": endpoint.Value}
			endpointMaps[projectId] = append(endpointMaps[projectId], item)
		}
	}
	return endpointMaps
}

func sortRestEndpoints(restEndpoints []Models.RestEndpoint) {
	sort.Slice(restEndpoints, func(i, j int) bool {
		first := restEndpoints[i]
		second := restEndpoints[j]
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
