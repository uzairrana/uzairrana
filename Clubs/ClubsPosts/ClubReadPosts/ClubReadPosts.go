package ClubReadPosts

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func ReadClubPosts(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, int, error) {
	logger.Debug("=== Club Read posts rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", -1, err
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
	//if group, _ := ClubMiscFunctions.CheckIfUserHasAGroup(ctx, logger, db, nk, props.UserId, props.ClubId); group != nil {
	listRecords, _, err := nk.StorageList(ctx, "", "Post_clubid_"+props.ClubId, 100, "")
	if err != nil {
		logger.WithField("err", err).Error("Storage list error.")
	} else {
		if len(listRecords) > 0 {
			// payload, _ := json.Marshal(listRecords)
			postids := &ClubClasses.PostClass{}
			for _, r := range listRecords {
				logger.Info("read: %d, write: %d, value: %s", r.PermissionRead, r.PermissionWrite, r.Value)
				if err := json.Unmarshal([]byte(r.Value), &postids); err != nil {
					logger.WithField("err", err).Error("Unable to unmarshall postids.")
				} else {
					payload1, _ := json.Marshal(postids)
					return string(payload1), len(postids.PostsIDs), nil
				}
			}
		} else {
			return "{\"status\":\"No posts to show\"}", 0, err
		}

	}
	//}
	return "{\"status\":\"Error getting club posts\"}", -1, nil
	// records, err := nk.StorageRead(ctx, objectIds)
	// if err != nil {
	// 	logger.WithField("err", err).Error("Storage read error.")
	// 	return "{\"status\":\"Could not find that post\"}", err
	// } else {
	// 	// for _, record := range records {
	// 	// logger.Info("read: %d, write: %d, value: %s", record.PermissionRead, record.PermissionWrite, record.Value)
	// 	// }
	// 	obj, _ := json.Marshal(records)
	// 	return string(obj), err
	// }
}

// func ReadClubPostsID(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
// 	logger.Debug("=== Club Read posts rpc called ===?%v", payload)
// 	props := &ClubClasses.ClubJoinProps{}
// 	if err := json.Unmarshal([]byte(payload), &props); err != nil {
// 		logger.Debug("===Unable to unmarshal props===?%v", err)
// 		return "{\"status\":\"InValidProps\"}", err
// 	}
// 	//callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
// 	// objectIds := []*runtime.StorageRead{
// 	// 	{
// 	// 		Collection: "Post",
// 	// 		Key:        "post_" + props.ClubId,
// 	// 		UserID:     props.UserId,
// 	// 	},
// 	// }
// 	// Storage List
// 	//userID := props.UserId // Some user ID.
// 	if group, _ := ClubMiscFunctions.CheckIfUserHasAGroup(ctx, logger, db, nk, props.UserId, props.ClubId); group != nil {
// 		listRecords, _, err := nk.StorageList(ctx, "", "Post_clubid_"+props.ClubId, 100, "")
// 		if err != nil {
// 			logger.WithField("err", err).Error("Storage list error.")
// 		} else {
// 			if len(listRecords) > 0 {
// 				payload, _ := json.Marshal(listRecords)
// 				for _, r := range listRecords {
// 					logger.Info("read: %d, write: %d, value: %s", r.PermissionRead, r.PermissionWrite, r.Value)
// 				}
// 				return string(payload), nil
// 			} else {
// 				return "{\"status\":\"No posts to show\"}", err
// 			}

// 		}
// 	}
// 	return "{\"status\":\"Error getting club posts\"}", nil
// 	// records, err := nk.StorageRead(ctx, objectIds)
// 	// if err != nil {
// 	// 	logger.WithField("err", err).Error("Storage read error.")
// 	// 	return "{\"status\":\"Could not find that post\"}", err
// 	// } else {
// 	// 	// for _, record := range records {
// 	// 	// logger.Info("read: %d, write: %d, value: %s", record.PermissionRead, record.PermissionWrite, record.Value)
// 	// 	// }
// 	// 	obj, _ := json.Marshal(records)
// 	// 	return string(obj), err
// 	// }
// }
