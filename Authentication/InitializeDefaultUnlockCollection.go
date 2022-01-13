package Authentication

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
	//"github.com/heroiclabs/nakama/runtime"
)

func AssignUnlockCollection(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, userid string) {

	userID := userid

	var ObjectArray UnlockBundle
	var decksObject BundleDetails

	var collectionKeys = []string{"avatars", "decks", "tables", "gifts", "chatpacks"}

	for _, KeyOfCollection := range collectionKeys {
		ObjectArray.ARRAY = nil
		if KeyOfCollection == "decks" {

			decksObject = BundleDetails{BundleName: "defaultdeck"}
			ObjectArray.ARRAY = append(ObjectArray.ARRAY, decksObject)
		} else if KeyOfCollection == "tables" {

			decksObject = BundleDetails{BundleName: "defaulttable"}
			ObjectArray.ARRAY = append(ObjectArray.ARRAY, decksObject)

		} else if KeyOfCollection == "avatars" {

			decksObject = BundleDetails{BundleName: "defaultavatar"}
			ObjectArray.ARRAY = append(ObjectArray.ARRAY, decksObject)
		}

		jsonObject2, _ := json.Marshal(ObjectArray)
		objectIds2 := []*runtime.StorageWrite{
			&runtime.StorageWrite{
				Collection:      "UnlockItems",
				Key:             KeyOfCollection,
				UserID:          userID,
				Value:           string(jsonObject2), // Value must be a valid encoded JSON object.
				PermissionRead:  1,
				PermissionWrite: 1,
			},
		}

		if _, err := nk.StorageWrite(ctx, objectIds2); err != nil {
			// Handle error.
			// logger.Println("Error writing default decks in DB -- AssignUnlockCollection -- ", err)
			// return err
		}

	}

}
