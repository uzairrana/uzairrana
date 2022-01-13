package ClubMiscFunctions

import (
	"context"
	"database/sql"
	"encoding/json"

	"arabicPoker.com/a/Clubs/ClubClasses"

	"github.com/heroiclabs/nakama-common/runtime"

	"github.com/heroiclabs/nakama-common/api"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
func CheckIfUserHasAGroup(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, userId string, groupId string) (*api.UserGroupList_UserGroup, bool) {

	if groups, _, err := nk.UserGroupsList(ctx, userId, 100, nil, ""); err != nil {
		logger.WithField("err", err).Error("User groups list error.")

	} else {
		for _, g := range groups { //can be done without indexing
			if g.Group.Id == groupId {
				return g, true
			}
		}
	}
	return nil, true
}
func IfManager(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) bool {
	logger.Debug("=== Club If manager rpc called ===?%v", payload)
	props := &ClubClasses.ClubMangersMetadata{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		//return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	if groups, _, err := nk.UserGroupsList(ctx, strSlice[0], 100, nil, ""); err != nil {
		logger.WithField("err", err).Error("User groups list error.")
		//return "{\"status\":\"User Groups List ERROR\"}", err
	} else {
		for i, g := range groups { //can be done without indexing
			if g.Group.Id == props.ClubId {
				if groups[i].Group.Metadata != "" {
					// data, _ := json.Marshal(groups[i].Group.Metadata)
					// logger.Debug("data?=%v", data)
					if err := json.Unmarshal([]byte(groups[i].Group.Metadata), &props); err != nil {
						logger.Debug("===Unable to unmarshal props===?%v", err)
						//return "{\"status\":\"InValidProps\"}", err
					}

				}
				logger.Debug("Metadata?=%v", props.Managers)
				logger.Debug("Metadata actual?=%v", groups[i].Group.Metadata)
				man := &ClubClasses.ClubManagers{}
				man.Managers = append(props.Managers, strSlice...)
				if Contains(props.Managers, strSlice[0]) {
					return true
				}
			}
		}
	}
	return false
}

// func GetManagersOfClubs(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) {

// 	var strSlice = []string{props.UserId}
// 	if groups, _, err := nk.UserGroupsList(ctx, strSlice[0], 100, nil, ""); err != nil {
// 		logger.WithField("err", err).Error("User groups list error.")
// 		return "{\"status\":\"User Groups List ERROR\"}", err
// 	} else {
// 		for i, g := range groups { //can be done without indexing
// 			if g.Group.Id == props.ClubId {
// 				if groups[i].Group.Metadata != "" {
// 					// data, _ := json.Marshal(groups[i].Group.Metadata)
// 					// logger.Debug("data?=%v", data)
// 					if err := json.Unmarshal([]byte(groups[i].Group.Metadata), &props); err != nil {
// 						logger.Debug("===Unable to unmarshal props===?%v", err)
// 						return "{\"status\":\"InValidProps\"}", err
// 					}

// 				}
// 				logger.Debug("Metadata?=%v", props.Managers)
// 				logger.Debug("Metadata actual?=%v", groups[i].Group.Metadata)
// 				man := &ClubClasses.ClubManagers{}
// 				man.Managers = append(props.Managers, strSlice...)
// 				if !ClubMiscFunctions.Contains(props.Managers, strSlice[0]) {
// 					//metadata := make(map[string]interface{})
// 					metadata := map[string]interface{}{ // Add whatever custom fields you want.
// 						"managers": man,
// 					}
// 					metadata["managers"] = append(props.Managers, strSlice...)
// 					if err := nk.GroupUpdate(ctx, props.ClubId, "", "", "", "", "", false, metadata, 100); err != nil {
// 						logger.WithField("err", err).Error("Group update error.")
// 						return "{\"status\":\"CLUB UPDATION FAILED\"}", err
// 					}

// 					return "{\"status\":\"CLUB SUCESSFULLY UPDATED FOR MANAGERS PROMOTE\"}", nil
// 				}
// 			}
// 		}
// 	}
// }
