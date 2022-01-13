package ClubDeleteAllPosts

import (
	"context"
	"database/sql"
	"encoding/json"

	"arabicPoker.com/a/Clubs/ClubClasses"
	"arabicPoker.com/a/Clubs/ClubMiscFunctions"

	"github.com/heroiclabs/nakama-common/runtime"
)

type managerposts struct {
	ManagerPosts []string `json:"managerpostids"`
}

func DeleteAllPosts(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== Club Delete All posts rpc called ===?%v", payload)
	props := &ClubClasses.ClubJoinProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}
	listRecords, _, err := nk.StorageList(ctx, "", "Post_clubid_"+props.ClubId, 100, "")
	if err != nil {
		logger.WithField("err", err).Error("Storage list error.")
	} else {
		if len(listRecords) > 0 {
			// payload, _ := json.Marshal(listRecords)
			postids := &ClubClasses.PostClass{}
			for _, r := range listRecords {
				logger.Info("read: %d, write: %d, value: %s", r.PermissionRead, r.PermissionWrite, r.Value)
				if err := json.Unmarshal([]byte(r.Value), &postids); err != nil {
					logger.WithField("err", err).Error("Unable to unmarshall postids.")
				} else {
					managerpostsids := &managerposts{}
					if err := json.Unmarshal([]byte(r.Value), &managerpostsids); err != nil {
						logger.WithField("err", err).Error("Unable to unmarshall postids.")
					}
					for _, str := range postids.PostsIDs {
						if ClubMiscFunctions.Contains(managerpostsids.ManagerPosts, str) {
							managerpostsids.ManagerPosts = ClubMiscFunctions.RemoveIndex(managerpostsids.ManagerPosts, ClubMiscFunctions.Find(managerpostsids.ManagerPosts, str))
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
								//return "{\"status\":\"Error reading Post while deleting\"}", nil
							}
							objectIds := []*runtime.StorageDelete{
								{
									Collection: str,
									Key:        "posts",
									//UserID:     props.UserId,
								},
							}
							err1 := nk.StorageDelete(ctx, objectIds)
							if err1 != nil {
								logger.WithField("err", err).Error("Storage delete error.")
								//return "{\"status\":\"Error Deleting Post\"}", err
							}
						} else {
							postids.PostsIDs = ClubMiscFunctions.RemoveIndex(postids.PostsIDs, ClubMiscFunctions.Find(postids.PostsIDs, str))
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
									Collection: str,
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
					}
				}
			}
		} else {
			return "{\"status\":\"No posts to delete\"}", err
		}
	}
	return "{\"status\":\"No posts to delete\"}", err
}
