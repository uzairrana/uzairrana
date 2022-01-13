package ClubDeletePost

import (
	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubMiscFunctions"
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

type managerposts struct {
	ManagerPosts []string `json:"managerpostids"`
}

func DeleteClubPost(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Delete post rpc called ===?%v", payload)
	props := &ClubClasses.ClubPosts{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	//if group, _ := ClubMiscFunctions.CheckIfUserHasAGroup(ctx, logger, db, nk, props.UserId, props.ClubId); group != nil {
	listRecords, _, err := nk.StorageList(ctx, "", "Post_clubid_"+props.ClubId, 100, "")
	if err != nil {
		logger.WithField("err", err).Error("Storage list error.")
	} else {
		// payload, _ := json.Marshal(listRecords)
		postids := &ClubClasses.PostClass{}
		for _, r := range listRecords {
			logger.Info("read: %d, write: %d, value: %s", r.PermissionRead, r.PermissionWrite, r.Value)
			if err := json.Unmarshal([]byte(r.Value), &postids); err != nil {
				logger.WithField("err", err).Error("Unable to unmarshall postids.")
			} else {
				logger.Debug("postid?=", postids)
				if ClubMiscFunctions.Contains(postids.PostsIDs, props.PostId) {
					manager := ClubMiscFunctions.IfManager(ctx, logger, db, nk, payload)
					if manager {
						managerpostsids := &managerposts{}
						if err := json.Unmarshal([]byte(r.Value), &managerpostsids); err != nil {
							logger.WithField("err", err).Error("Unable to unmarshall postids.")
						}
						managerpostsids.ManagerPosts = ClubMiscFunctions.RemoveIndex(managerpostsids.ManagerPosts, ClubMiscFunctions.Find(managerpostsids.ManagerPosts, props.PostId))
						value := make(map[string]interface{})
						//value["postids"] = append(listRecords.PostsIDs, "post_"+props.ClubId+"_"+strconv.FormatInt(serverTime, 10))
						value["managerpostids"] = managerpostsids.ManagerPosts
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
							return "{\"status\":\"Error reading Post while deleting\"}", nil
						}
						objectIds := []*runtime.StorageDelete{
							{
								Collection: props.PostId,
								Key:        "posts",
								//UserID:     props.UserId,
							},
						}
						err1 := nk.StorageDelete(ctx, objectIds)
						if err1 != nil {
							logger.WithField("err", err).Error("Storage delete error.")
							//return "{\"status\":\"Error Deleting Post\"}", err
						}

					}
					//	logger.Debug("club post delete true if first")
					postids.PostsIDs = ClubMiscFunctions.RemoveIndex(postids.PostsIDs, ClubMiscFunctions.Find(postids.PostsIDs, props.PostId))
					value := make(map[string]interface{})
					//value["postids"] = append(listRecords.PostsIDs, "post_"+props.ClubId+"_"+strconv.FormatInt(serverTime, 10))
					value["postids"] = postids.PostsIDs
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
						return "{\"status\":\"Error reading Post while deleting\"}", nil
					}
					objectIds := []*runtime.StorageDelete{
						{
							Collection: props.PostId,
							Key:        "posts",
							//UserID:     props.UserId,
						},
					}
					err1 := nk.StorageDelete(ctx, objectIds)
					if err1 != nil {
						logger.WithField("err", err).Error("Storage delete error.")
						return "{\"status\":\"Error Deleting Post\"}", err
					} else {
						return "{\"status\":\"Post Deleted\"}", err

					}

				}
				//payload1, _ := json.Marshal(postids)
				//return string(payload1), len(postids.PostsIDs), nil
			}
		}

	}
	//}
	return "{\"status\":\"Post not deleted\"}", nil

}
