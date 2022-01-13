package ClubDelete

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubLeaderboard"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubDeletePost"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubReadPosts"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func DeleteClub(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Delete rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {

		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}
	if err := nk.GroupDelete(ctx, props.ClubId); err != nil {
		logger.WithField("err", err).Error("Group delete error.")
		return "{\"status\":\"Group Deletion Error\"}", err
	} else {

		callerID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
		props.UserId = callerID
		payload1, _ := json.Marshal(props)
		str, _, _ := ClubReadPosts.ReadClubPosts(ctx, logger, db, nk, string(payload1))
		//manger := ClubMiscFunctions.IfManager(ctx, logger, db, nk, string(payload1))

		postids := &ClubClasses.PostClass{}
		if err12 := json.Unmarshal([]byte(str), &postids); err != nil {
			logger.Debug("===Unable to unmarshal props===?%v", err12)
			return "{\"status\":\"InValidProps Error read post props\"}", err12
		} else {
			for _, r := range postids.PostsIDs {
				//logger.Info("read: %d, write: %d, value: %s", r.PermissionRead, r.PermissionWrite, r.Value)
				propsforpost := &ClubClasses.ClubPosts{}
				propsforpost.ClubId = props.ClubId
				propsforpost.PostId = r
				payload2, _ := json.Marshal(propsforpost)
				str2, _ := ClubDeletePost.DeleteClubPost(ctx, logger, db, nk, string(payload2))
				logger.Debug("post status?=", str2)

			}
		}
		postData := []*runtime.StorageDelete{
			{
				Collection: "Post_clubid_" + props.ClubId,
				Key:        "posts",
			},
		}
		err123 := nk.StorageDelete(ctx, postData)
		if err123 != nil {
			logger.WithField("err", err).Error("Storage Delete error.")
			return "{\"status\":\"Error Deleting club Post database\"}", nil
		}
		logger.WithField("Successfull deletion clubID?%v", props.ClubId)
		err := ClubLeaderboard.DeleteLeaderboard(ctx, logger, db, nk, props.ClubId)
		if err != nil {
			logger.WithField("err", err).Error("Leaderboard delete error.")
		} else {
			logger.Debug("Leaderboard deleted successfully id?=%v", props.ClubId)
		}
		return "{\"status\":\"ClubDeleted\"}", nil
	}
}
