package ClubUserList

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func UserListClub(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== ClubUserList rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {

		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}

	groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, "")
	if err != nil {
		logger.WithField("err", err).Error("Group users list error.")
		return "{\"status\":\"Error Occured while getting user from club\"}", err
	}

	for _, member := range groupUserList {
		// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
		logger.Info("Member username %s has state %d", member.GetUser().Username, member.GetState())
	}
	listofusersrequests, _ := json.Marshal(groupUserList)
	return "{\"list\":" + string(listofusersrequests) + "}", nil

}
