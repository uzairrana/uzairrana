package ClubPermissions

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubUpdate"
	"context"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/heroiclabs/nakama-common/runtime"
)

func PermissionClubPromote(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===ClubPermission Promote rpc called===")
	props := &ClubClasses.ClubPropsIds{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	convInt, _ := strconv.ParseInt(props.Permission, 0, 32)
	if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "\"status\":\"Could not get user list for group\"", err
	} else {
		for _, member := range groupUserList {
			// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
			if member.User.Id == strSlice[0] {
				state := member.GetState().Value
				//logger.Debug("Before User state %v", state)
				if state <= int32(convInt) {
					return "{\"status\":\"Invalid Permission \"}", err
				}
				for i := state; i > int32(convInt); i-- {
					callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
					if err := nk.GroupUsersPromote(ctx, callerID, props.ClubId, strSlice); err != nil {
						logger.WithField("err", err).Error("Group users promote error.")
						return "{\"status\":\"User does not exits\"}", err
					}
				}
				//logger.Debug("after User state %v", member.GetState().Value)
				content := map[string]interface{}{}
				nk.NotificationSend(ctx, member.GetUser().Id, "You are Promoted", content, 55, "", true)
				return "{\"status\":\"User Promoted\"}", err
				// nk.NotificationSend(ctx, member.GetUser().Id, "A Club Member Requested to Join", notificationContent, 53, "", true)
			}
		}
		return "{\"status\":\"User does not exits\"}", err
	}

}
func PermissionClubDemote(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===ClubPermission Demote rpc called===")
	props := &ClubClasses.ClubPropsIds{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	convInt, _ := strconv.ParseInt(props.Permission, 0, 32)
	if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {

		logger.Error("Could not get user list for group: %s", err.Error())
		return "\"status\":\"Could not get user list for group\"", err
	} else {

		for _, member := range groupUserList {
			// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
			if member.User.Id == strSlice[0] {
				state := member.GetState().Value
				//logger.Debug("Before User state%v", state)
				if state >= int32(convInt) {
					return "{\"status\":\"Invalid Permission\"}", err
				}
				for i := state; i < int32(convInt); i++ {
					callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
					if err := nk.GroupUsersDemote(ctx, callerID, props.ClubId, strSlice); err != nil {
						logger.WithField("err", err).Error("Group users Demote error.")
					}

				}
				//logger.Debug("after User state%v", member.GetState().Value)
				content := map[string]interface{}{}
				nk.NotificationSend(ctx, member.GetUser().Id, "You are Demoted", content, 53, "", true)
				return "{\"status\":\"User Demoted\"}", err
				// nk.NotificationSend(ctx, member.GetUser().Id, "A Club Member Requested to Join", notificationContent, 53, "", true)
			}
		}
		return "{\"status\":\"User does not exits\"}", err
	}

}
func PermissionClubDemoteManager(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===ClubPermission Demote rpc called===")
	props := &ClubClasses.ClubPropsIds{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	//convInt, _ := strconv.ParseInt(props.Permission, 0, 32)
	if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "\"status\":\"Could not get user list for group\"", err
	} else {

		for _, member := range groupUserList {
			// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
			if member.User.Id == strSlice[0] {
				state := member.GetState().Value
				//logger.Debug("Before User state%v", state)
				if state == 1 {
					callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
					if err := nk.GroupUsersDemote(ctx, callerID, props.ClubId, strSlice); err != nil {
						logger.WithField("err", err).Error("Group users Demote error.")
						return "{\"status\":\"Error demoting a manager\"}", err
					}
					if rtr, err := ClubUpdate.UpdateClubInfoManagersDemote(ctx, logger, db, nk, payload); err == nil {
						logger.Debug("UpdateClubInfoManager rtr?=%v", rtr)
						content := map[string]interface{}{}
						nk.NotificationSend(ctx, member.GetUser().Id, "You are Demoted from a manager post", content, 56, "", true)
						return "{\"status\":\"User Demoted\"}", err
					}
					return "{\"status\":\"User Demotion error\"}", err
				}
			}
		}
		return "{\"status\":\"User does not exits\"}", err
	}

}
func PermissionClubPromoteManager(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("===PermissionClubPromoteManager  rpc called===")
	props := &ClubClasses.ClubPropsIds{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	var strSlice = []string{props.UserId}
	//convInt, _ := strconv.ParseInt(props.Permission, 0, 32)
	if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {
		logger.Error("Could not get user list for group: %s", err.Error())
		return "\"status\":\"Could not get user list for group\"", err
	} else {
		for _, member := range groupUserList {
			// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
			if member.User.Id == strSlice[0] {
				state := member.GetState().Value
				logger.Debug("Before User state %v", state)
				if state == 2 {
					callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
					logger.Debug("before member.user.state?=", member.State.Value)
					if err := nk.GroupUsersPromote(ctx, callerID, props.ClubId, strSlice); err != nil {
						logger.WithField("err", err).Error("Group users promote error.")
						return "{\"status\":\"User does not exits\"}", err
					}
					logger.Debug("after+ member.user.state?=", member.State.Value)
					//logger.Debug("after User state %v", member.GetState().Value)
					if rtr, err := ClubUpdate.UpdateClubInfoManagersPromote(ctx, logger, db, nk, payload); err != nil {
						logger.Debug("UpdateClubInfoManager rtr?=%v", rtr)
					}
					content := map[string]interface{}{}
					nk.NotificationSend(ctx, member.GetUser().Id, "You are Promoted to manager", content, 54, "", true)
					return "{\"status\":\"User Promoted\"}", err
				} else if state == 1 {
					if rtr, err := ClubUpdate.UpdateClubInfoManagersPromote(ctx, logger, db, nk, payload); err == nil {
						logger.Debug("UpdateClubInfoManager rtr?=%v", rtr)
						return "{\"status\":\"User is already a manager\"}", err
					}
					//logger.Debug("UpdateClubInfoManager rtr?=%v", rtr)
				}
			}
		}
		logger.Debug("I am here")
		return "{\"status\":\"User does not exits\"}", err
	}

}
