package ClubJoin

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func JoinClub(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Join rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {

		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}
	clubID := props.ClubId
	//userID := props.UserId
	if err := nk.GroupUserJoin(ctx, clubID, ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string), ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)); err != nil {
		logger.WithField("err", err).Error("Group user join error.")
		return "{\"status\":\"Group User Join Error\"}", err
	}
	
	if groupUserList, _, err := nk.GroupUsersList(ctx, clubID, 100, nil, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "{\"status\":\"Could not get user list for group\"}", err
	} else {
		notificationContent := map[string]interface{}{
			"clubid": clubID,
		}
		for _, member := range groupUserList {
			if member.GetState().Value != 3 {
				nk.NotificationSend(ctx, member.GetUser().Id, "A Club Member Requested to Join", notificationContent, 53, "", true)
			}
		}
	}

	return string("{\"status\":\"Requested to join\", \"clubID\":\"") + clubID + "\"}", nil

}
