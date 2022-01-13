package ClubSearch

import (
	"context"
	"encoding/json"

	"arabicPoker.com/a/Clubs/ClubClasses"

	"github.com/heroiclabs/nakama-common/runtime"
)

func SearchClubByName(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, payload string) (string, error) {
	//groupName := "MyClub"
	//langTag := "en"

	// open := true
	logger.Debug("===CreateClub rpc called===")
	props := &ClubClasses.ClubProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {
		logger.Debug("===Unable to unmarshal props===?%v", err)
		return "{\"status\":\"InValidProps\"}", err
	}

	list, _, err := nk.GroupsList(ctx, props.Name, "", nil, nil, 100, "")
	if err != nil {
		logger.WithField("err", err).Error("Group list error.")
		return "{\"status\":\"Error\"}", err
	} else {
		for _, g := range list {
			logger.Info("ID %s - can enter? %b", g.Id, g.Open.Value)
		}
		if len(list) == 0 {
			return "{\"status\":\"No Club Found\"}", err
		}
		group, _ := json.Marshal(list)
		return "{\"list\":" + string(group) + "}", nil
	}
}
