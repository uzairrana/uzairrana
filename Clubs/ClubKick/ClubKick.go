package ClubKick

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func KickClubMember(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// groupID := "dcb891ea-a311-4681-9213-6741351c9994"
	// userIds := []string{"9a51cf3a-2377-11eb-b713-e7d403afe081", "a042c19c-2377-11eb-b7c1-cfafae11cfbc"}
	logger.Debug("===Clubkickmember rpc called===", payload)
	props := &ClubClasses.ClubPropsIds{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "\"status\":\"Could not get user list for group\"", err
	} else {
		logger.Debug("groupuserlist?=", groupUserList)
		for _, member := range groupUserList {
			// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
			if member.User.Id == strSlice[0] {
				callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
				if err := nk.GroupUsersKick(ctx, callerID, props.ClubId, strSlice); err != nil {
					logger.WithField("err", err).Error("Group users kick error.")
					return "{\"status\":\"Group users kick error\"}", err
				}
				return "{\"status\":\"User Kicked\"}", nil
			}
		}
	}
	return "{\"status\":\"User does not exists\"}", nil
}
