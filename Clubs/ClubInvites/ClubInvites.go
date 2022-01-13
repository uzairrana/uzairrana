package ClubInvites

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InviteClub(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Invite rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {

		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}
	clubID := props.ClubId
	userID := props.UserId

	if groupUserList, _, err := nk.GroupUsersList(ctx, clubID, 100, nil, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "{\"status\":\"Could not get user list for group\"}", err
	} else {
		notificationContent := map[string]interface{}{
			"invite": true,
		}
		for _, member := range groupUserList {
			if member.GetState().Value != 3 && member.User.Id == ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string) {
				// if err := nk.GroupUserJoin(ctx, clubID, userID, ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)); err != nil {
				// 	logger.WithField("err", err).Error("Group user join error.")
				// 	return "{\"status\":\"Group User Join Error\"}", err
				// }
				nk.NotificationSend(ctx, props.UserId, "Invite sent to Join", notificationContent, 54, "", true)
				return string("{\"status\":\"Invited Successfully\", \"userid\":\"") + userID + "\"}", nil
			}
		}
		logger.Debug("i am here")
		return "{\"status\":\"You do not have permission to invite others\"}", nil
	}
}
