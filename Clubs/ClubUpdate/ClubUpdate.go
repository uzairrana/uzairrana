package ClubUpdate

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubMiscFunctions"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func UpdateClubInfo(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Update rpc called ===?%v", payload)
	props := &ClubClasses.ClubProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}

	metadata := make(map[string]interface{})
	metadata["managers"] = make([]string, 0)
	if err := nk.GroupUpdate(ctx, props.ClubId, props.Name, props.OwnerID, "", props.Desc, props.ClubPictureUrl, false, metadata, 100); err != nil {
		logger.WithField("err", err).Error("Group update error.")
		return "{\"status\":\"CLUB UPDATION FAILED\"}", nil
	}
	return "{\"status\":\"CLUB SUCESSFULLY UPDATED\"}", nil

}
func UpdateClubInfoManagersPromote(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Update Managers function called ===?%v", payload)
	props := &ClubClasses.ClubMangersMetadata{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	if groups, _, err := nk.UserGroupsList(ctx, strSlice[0], 100, nil, ""); err != nil {
		logger.WithField("err", err).Error("User groups list error.")
		return "{\"status\":\"User Groups List ERROR\"}", err
	} else {
		for i, g := range groups { //can be done without indexing
			if g.Group.Id == props.ClubId {
				if groups[i].Group.Metadata != "" {
					// data, _ := json.Marshal(groups[i].Group.Metadata)
					// logger.Debug("data?=%v", data)
					if err := json.Unmarshal([]byte(groups[i].Group.Metadata), &props); err != nil {
						logger.Debug("===Unable to unmarshal props===?%v", err)
						return "{\"status\":\"InValidProps\"}", err
					}

				}
				logger.Debug("Metadata?=%v", props.Managers)
				logger.Debug("Metadata actual?=%v", groups[i].Group.Metadata)
				man := &ClubClasses.ClubManagers{}
				man.Managers = append(props.Managers, strSlice...)
				if !ClubMiscFunctions.Contains(props.Managers, strSlice[0]) {
					//metadata := make(map[string]interface{})
					metadata := map[string]interface{}{ // Add whatever custom fields you want.
						"managers": man,
					}
					metadata["managers"] = append(props.Managers, strSlice...)
					if err := nk.GroupUpdate(ctx, props.ClubId, "", "", "", "", "", false, metadata, 100); err != nil {
						logger.WithField("err", err).Error("Group update error.")
						return "{\"status\":\"CLUB UPDATION FAILED\"}", err
					}

					return "{\"status\":\"CLUB SUCESSFULLY UPDATED FOR MANAGERS PROMOTE\"}", nil
				}
			}
		}
	}
	return "{\"status\":\"MANAGERS ALREADY EXISTS\"}", nil
}
func UpdateClubInfoManagersDemote(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Update Managers Demote function called ===?%v", payload)
	props := &ClubClasses.ClubMangersMetadata{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	if groups, _, err := nk.UserGroupsList(ctx, strSlice[0], 100, nil, ""); err != nil {
		logger.WithField("err", err).Error("User groups list error.")
		return "{\"status\":\"User Groups List ERROR\"}", err
	} else {
		for i, g := range groups {
			if g.Group.Id == props.ClubId {
				if groups[i].Group.Metadata != "" {
					if err := json.Unmarshal([]byte(groups[i].Group.Metadata), &props); err != nil {
						logger.Debug("===Unable to unmarshal props===?%v", err)
						return "{\"status\":\"InValidProps\"}", err
					}

				}
				logger.Debug("Metadata?=%v", props.Managers)
				logger.Debug("Metadata actual?=%v", groups[i].Group.Metadata)
				if ClubMiscFunctions.Contains(props.Managers, strSlice[0]) {
					metadata := make(map[string]interface{})
					metadata["managers"] = ClubMiscFunctions.RemoveIndex(props.Managers, ClubMiscFunctions.Find(props.Managers, strSlice[0]))
					if err := nk.GroupUpdate(ctx, props.ClubId, "", "", "", "", "", false, metadata, 100); err != nil {
						logger.WithField("err", err).Error("Group update error.")
						return "{\"status\":\"CLUB UPDATION FAILED\"}", nil
					}
					//logger.Debug("in if contains before returning")
					return "{\"status\":\"CLUB SUCESSFULLY UPDATED FOR MANAGERS DEMOTE\"}", nil
				}
			}
		}
	}
	return "{\"status\":\"MANAGER EXISTS BUT COULD NOT DEMOTE\"}", nil
}
