package Authentication

import (
	"context"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func AssignUserMatches(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userid string) {
	userID := userid

	obj := IPM{InProgressMatches: []string{}} //for new users...

	jsonObject2, _ := json.Marshal(obj)
	objectIds2 := []*runtime.StorageWrite{
		&runtime.StorageWrite{
			Collection:      "UserMatches",
			Key:             "InProgress",
			UserID:          userID,
			Value:           string(jsonObject2), // Value must be a valid encoded JSON object.
			PermissionRead:  2,
			PermissionWrite: 0,
		},
	}
	if _, err := nk.StorageWrite(ctx, objectIds2); err != nil {
		// Handle error.
		// logger.Println("Error writing default decks in DB -- AssignUnlockCollection -- ", err)
		// return err
	}

	obj1 := EM{EndedMatches: []string{}} //for new users...
	jsonObject3, _ := json.Marshal(obj1)
	objectIds3 := []*runtime.StorageWrite{
		&runtime.StorageWrite{
			Collection:      "UserMatches",
			Key:             "Ended",
			UserID:          userID,
			Value:           string(jsonObject3), // Value must be a valid encoded JSON object.
			PermissionRead:  2,
			PermissionWrite: 0,
		},
	}
	if _, err := nk.StorageWrite(ctx, objectIds3); err != nil {
		// Handle error.
		// logger.Println("Error writing default decks in DB -- AssignUnlockCollection -- ", err)
		// return err
	}
}
