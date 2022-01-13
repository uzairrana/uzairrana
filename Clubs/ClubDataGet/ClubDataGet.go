package ClubDataGet

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func GetClubData(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== GetClubData rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {

		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}
	list, _, err := nk.GroupsList(ctx, "", "", nil, nil, 100, "")
	if err != nil {
		logger.WithField("err", err).Error("Group list error.")
		return "{\"status\":\"Error in finding Clubs\"}", err
	} else {
		for i, g := range list {
			if g.Id == props.ClubId {
				listofclubs, _ := json.Marshal(list[i])
				return "{\"clubinfo\":" + string(listofclubs) + "}", nil
				// logger.Info("ID %s - Name? =%b cursor?=%v", g.Id, g.Name, cursor)
			}
		}
	}
	// if len(list) > 0 {
	// 	listofclubs, _ := json.Marshal(list)
	// 	logger.Debug("listofclubs?=", "{\"list\":"+string(listofclubs)+"}")
	// 	return "{\"list\":" + string(listofclubs) + "}", nil
	// }
	return "{\"status\":\"No Club found\"}", nil
}
func CanUserJoinClubs(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== CanUserJoinClubs rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}
	groups, _, err := nk.UserGroupsList(ctx, props.UserId, 100, nil, "")
	if err != nil {
		logger.WithField("err", err).Error("User groups list error.")
		return "{\"status\":\"Error in finding Clubs\"}", err
	}
	for _, group := range groups {
		logger.Debug("User has state %d in group %s.", group.GetState(), group.GetGroup().Name)
	}
	can := true
	clubs := make(map[string]interface{})

	if len(groups) > 0 {
		if groups[0].GetState().Value != 3 {
			clubs["clubid"] = groups[0]
			can = false
		}
	}
	
	clubs["canjoin"] = can
	convClubs, _ := json.Marshal(clubs)
	return string(convClubs), nil
	//return "{\"clubinfo\":" + string(listofclubs) + ", \"canjoin\":\"" + can + "\"}", nil
}
