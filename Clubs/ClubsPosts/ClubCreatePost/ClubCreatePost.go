package ClubCreatePost

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubMiscFunctions"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubReadPosts"
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

func CreateClubPost(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club posts rpc called ===?%v", payload)
	props := &ClubClasses.ClubCreatePostClass{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	manager := false
	if group, _ := ClubMiscFunctions.CheckIfUserHasAGroup(ctx, logger, db, nk, props.UserId, props.ClubId); group != nil {
		//postdata := make(map[string]interface{})
		var state int32
		if groupUserList, _, err := nk.GroupUsersList(ctx, props.ClubId, 100, nil, ""); err != nil {
			logger.Error("Could not get user list for group: %s", err.Error())
			return "\"status\":\"Could not get user list for group\"", err
		} else {
			for _, member := range groupUserList {
				// States are => 0: Superadmin, 1: Admin, 2: Member, 3: Requested to join
				if member.User.Id == props.UserId {
					manager = ClubMiscFunctions.IfManager(ctx, logger, db, nk, payload)
					state = member.GetState().Value
					logger.Debug("User state %v", state)
				}
			}
		}
		if state == 2 || state == 3 {
			return "{\"status\":\"You connot create post\"}", nil
		}
		serverTime := time.Now().Unix()

		str, len, _ := ClubReadPosts.ReadClubPosts(ctx, logger, db, nk, payload)
		logger.Debug("str?=", str)
		logger.Debug("len?=", len)
		//logger.Debug("str?=",str)
		listRecords := &ClubClasses.PostClass{}
		if err1 := json.Unmarshal([]byte(str), &listRecords); err1 != nil {
			logger.Debug("Error occured during reading post")
			return "{\"status\":\"Error Reading Post ids\"}", nil
		} else {
			if len == 0 {
				value := make(map[string]interface{})
				// time:=strconv.FormatInt(serverTime, 10)
				value["postids"] = append(listRecords.PostsIDs, "post_"+props.ClubId+"_"+strconv.FormatInt(serverTime, 10))
				if manager {
					value["managerpostids"] = append(listRecords.PostsIDs, "post_"+props.ClubId+"_"+strconv.FormatInt(serverTime, 10))
				}
				post, _ := json.Marshal(value)
				postData := []*runtime.StorageWrite{
					{
						Collection: "Post_clubid_" + props.ClubId,
						Key:        "posts",
						//UserID:          props.UserId,
						Value:           string(post), // Value must be a valid encoded JSON object.
						PermissionRead:  2,
						PermissionWrite: 0,
					},
				}
				_, err := nk.StorageWrite(ctx, postData)
				if err != nil {
					logger.WithField("err", err).Error("Storage write error.")
					return "{\"status\":\"Error Creating Post\"}", nil
				}
				postseperate := make(map[string]interface{})
				postseperate["desc"] = props.Desc
				posting, _ := json.Marshal(postseperate)
				postData1 := []*runtime.StorageWrite{
					{
						Collection: "post_" + props.ClubId + "_" + strconv.FormatInt(serverTime, 10),
						Key:        "posts",
						//UserID:          props.UserId,
						Value:           string(posting), // Value must be a valid encoded JSON object.
						PermissionRead:  2,
						PermissionWrite: 0,
					},
				}
				_, err1 := nk.StorageWrite(ctx, postData1)
				if err1 != nil {
					logger.WithField("err", err).Error("Storage write error.")
					return "{\"status\":\"Error Creating Post\"}", nil
				}
			} else {
				//postppc := &ClubClasses.PostClass{}
				// for _, post := range listRecords.PostsIDs {
				// 	if err := json.Unmarshal([]byte(post.GetValue()), postppc); err != nil {
				// 		logger.Debug("===Unable to unmarshal props===?%v", err)
				// 		return "{\"status\":\"InValidProps\"}", err
				// 	}

				value := make(map[string]interface{})
				value["postids"] = append(listRecords.PostsIDs, "post_"+props.ClubId+"_"+strconv.FormatInt(serverTime, 10))
				if manager {
					value["managerpostids"] = append(listRecords.PostsIDs, "post_"+props.ClubId+"_"+strconv.FormatInt(serverTime, 10))
				}
				post, _ := json.Marshal(value)
				postData := []*runtime.StorageWrite{
					{
						Collection: "Post_clubid_" + props.ClubId,
						Key:        "posts",
						//UserID:          props.UserId,
						Value:           string(post), // Value must be a valid encoded JSON object.
						PermissionRead:  2,
						PermissionWrite: 0,
					},
				}
				_, err := nk.StorageWrite(ctx, postData)
				if err != nil {
					logger.WithField("err", err).Error("Storage write error.")
					return "{\"status\":\"Error Creating Post\"}", nil
				}
				postseperate := make(map[string]interface{})
				postseperate["desc"] = props.Desc
				posting, _ := json.Marshal(postseperate)
				postData1 := []*runtime.StorageWrite{
					{
						Collection: "post_" + props.ClubId + "_" + strconv.FormatInt(serverTime, 10),
						Key:        "posts",
						//UserID:          props.UserId,
						Value:           string(posting), // Value must be a valid encoded JSON object.
						PermissionRead:  2,
						PermissionWrite: 0,
					},
				}
				_, err1 := nk.StorageWrite(ctx, postData1)
				if err1 != nil {
					logger.WithField("err", err).Error("Storage write error.")
					return "{\"status\":\"Error Creating Post\"}", nil
				}

			}
		}

	} else {
		return "{\"status\":\"You are not in Club\"}", nil
	}

	return "{\"status\":\"Post Created\"}", nil
}
