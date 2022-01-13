package ClubPostUpdate

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubMiscFunctions"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func UpdateClubPost(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Read posts rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	//callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	// objectIds := []*runtime.StorageRead{
	// 	{
	// 		Collection: "Post",
	// 		Key:        "post_" + props.ClubId,
	// 		UserID:     props.UserId,
	// 	},
	// }
	// Storage List
	//userID := props.UserId // Some user ID.
	if group, _ := ClubMiscFunctions.CheckIfUserHasAGroup(ctx, logger, db, nk, props.UserId, props.ClubId); group != nil {
		listRecords, _, err := nk.StorageList(ctx, "", "Post_clubid_"+props.ClubId, 100, "")
		if err != nil {
			logger.WithField("err", err).Error("Storage list error.")
		} else {
			if len(listRecords) > 0 {
				payload, _ := json.Marshal(listRecords)
				for _, r := range listRecords {
					logger.Info("read: %d, write: %d, value: %s", r.PermissionRead, r.PermissionWrite, r.Value)
				}
				return string(payload), nil
			} else {
				return "{\"status\":\"No posts to show\"}", err
			}

		}
	}
	return "{\"status\":\"Error getting club posts\"}", nil
}
