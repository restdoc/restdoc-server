package utils

import (
	Models "restdoc-models/models"
	"strconv"
)

func FormatEndpoints(endpoints []Models.RestEndpoint) map[int64][]map[string]interface{} {
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
