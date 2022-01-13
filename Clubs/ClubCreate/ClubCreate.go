package ClubCreate

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubLeaderboard"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func CreateClub(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===CreateClub rpc called===")
	props := &ClubClasses.ClubProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	// metadata := map[string]interface{}{ // Add whatever custom fields you want.
	// 	"ownerid":      props.OwnerID,
	// 	"name":          props.Name,
	// 	"level": "0"

	// }
	metadata := map[string]interface{}{}
	userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if group, err := nk.GroupCreate(ctx, userID, props.Name, userID, "en", props.Desc, props.ClubPictureUrl, false, metadata, 100); err == nil {
		logger.Debug("===Group Created Successfully=== ", group.Id)
		errL := ClubLeaderboard.CreateLeaderboard(ctx, logger, db, nk, group.Id)
		if errL != nil {
			logger.Debug("===Unable to create leaderboard for the club===?%v", err)
		} else {
			group.AvatarUrl = props.ClubPictureUrl
			logger.Debug("===Leaderboard created for clubid?===?%v", group.Id)
			postdata := make(map[string]interface{})
			postdata["postids"] = make([]string, 0)
			post, _ := json.Marshal(postdata)
			objectIDs := []*runtime.StorageWrite{
				{
					Collection: "Post_clubid_" + group.Id,
					Key:        "posts",
					//UserID:          userID,
					Value:           string(post), // Value must be a valid encoded JSON object.
					PermissionRead:  2,
					PermissionWrite: 0,
				},
			}
			_, err := nk.StorageWrite(ctx, objectIDs)
			if err != nil {
				logger.WithField("err", err).Error("Storage write error.")
			}
		}
		return string("{\"status\":\"Club Created\", \"clubID\":\"") + group.Id + "\"}", nil
	} else {
		return "{\"status\":\"CLUB CREATION FAILED\"}", nil
	}

}
