package ClubLeave

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/heroiclabs/nakama-common/runtime"
)

func LeaveClubMember(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// groupID := "dcb891ea-a311-4681-9213-6741351c9994"
	// userIds := []string{"9a51cf3a-2377-11eb-b713-e7d403afe081", "a042c19c-2377-11eb-b7c1-cfafae11cfbc"}
	logger.Debug("===Clubleaverpc called===", payload)
	props := &ClubClasses.ClubLeaveProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	// groupID := "dcb891ea-a311-4681-9213-6741351c9994"
	// userID := "9a51cf3a-2377-11eb-b713-e7d403afe081"
	// username := "myusername"
	if err := nk.GroupUserLeave(ctx, props.ClubId, ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string), ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)); err != nil {
		logger.WithField("err", err).Error("Group user leave error.")
		return "{\"status\":\"Error Leaving Club\"}", err
	}
	return "{\"status\":\"User Left the club\"}", nil

}






















































