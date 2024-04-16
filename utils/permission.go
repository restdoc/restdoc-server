package utils

import (
	"strconv"
	"strings"

	Models "github.com/restdoc/restdoc-models"
)

func CheckPermission(teamId int64, permittedTeamIds string) bool {
	team_id := strconv.FormatInt(teamId, 10)
	permitted := false
	arr := strings.Split(permittedTeamIds, ",")
	for _, tid := range arr {
		if tid == team_id {
			permitted = true
			break
		}
	}
	return permitted
}

func GetTeamIds(uid int64) (string, error) {
	teamIds := []string{}
	var teams []Models.Team
	err := Models.GetTeamsByUserId(&teams, uid)
	if err != nil {
		return "", err
	}
	for _, team := range teams {
		teamId := strconv.FormatInt(team.Id, 10)
		teamIds = append(teamIds, teamId)
	}
	ids := strings.Join(teamIds, ",")
	return ids, nil
}
