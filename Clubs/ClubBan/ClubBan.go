package ClubBan

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func BanClubMember(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	/*For now the best route would be to store the user IDs for banned users in metadata for the group
	Definitely the best place to hold the information but it could be useful to have a discussion on a github issue about adding a ban user option for groups
	We'd need to work out the best semantics for it though
	So open an issue when/if you get the chance*/
	
	// // groupID := "dcb891ea-a311-4681-9213-6741351c9994"
	// // userIds := []string{"9a51cf3a-2377-11eb-b713-e7d403afe081", "a042c19c-2377-11eb-b7c1-cfafae11cfbc"}
	// logger.Debug("===ClubMemberBan rpc called===")
	// props := &ClubClasses.ClubPropsIds{}
	// if err := json.Unmarshal([]byte(payload), &props); err != nil {
	// 	logger.Debug("===Unable to unmarshal props===?%v", err)
	// 	return "{\"status\":\"InValidProps\"}", err
	// }
	// if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {
	// 	logger.Error("Could not get user list for group: %s", err.Error())
	// 	return "\"status\":\"Could not get user list for group\"", err
	// } else {
	// 	for _, member := range groupUserList {
	// 		// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
	// 		if member.User.Id == props.UserId[0] {
	// 			callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	// 			if err := nk.ban(ctx, props.UserId); err != nil {
	// 				logger.WithField("err", err).Error("Group users kick error.")
	// 				return "{\"status\":\"Group users kick error\"}", err
	// 			}
	// 			return "{\"status\":\"User Kicked\"}", nil
	// 		}
	// 	}
	// }
	return "{\"status\":\"Rpc not ready yet\"}", nil
}
