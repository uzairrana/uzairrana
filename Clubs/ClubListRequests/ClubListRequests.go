package ClubListRequests

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func ListClubRequests(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===ClubListRequests rpc called===", payload)
	props := &ClubClasses.ClubProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	i := 3
	state := &i
	if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, state, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "{\"status\":\"Unable to get Requests List\"}", err
	} else {
		for _, member := range groupUserList {
			// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
			logger.Debug("Member?%v MemberState?%v", member, member.GetState())
		}
		listofusersrequests, _ := json.Marshal(groupUserList)
		logger.Debug("return?=", string(listofusersrequests))
		return "{\"list\":" + string(listofusersrequests) + "}", nil
	}
}
