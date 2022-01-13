package ClubAdd

import (
	"context"
	"database/sql"
	"encoding/json"

	"arabicPoker.com/a/Clubs/ClubClasses"

	"github.com/heroiclabs/nakama-common/runtime"
)

func AddUserToClub(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===ClubUserAdd rpc called===", payload)
	props := &ClubClasses.ClubPropsIds{}

	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	logger.Debug("clubid?=", props.ClubId)
	logger.Debug("userid?=", props.UserId)

	logger.Debug("strSlice?=", strSlice)
	callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if err := nk.GroupUsersAdd(ctx, callerID, props.ClubId, strSlice); err != nil {
		logger.WithField("err", err).Error("Group users add error.")
		return "{\"status\":\"Group users add error\"}", err
	}
	return "{\"status\":\"User added successfully\"}", nil

}
