package ClubLists

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func ListClubs(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== ListClubs rpc called ===?%v", payload)
	list, cursor, err := nk.GroupsList(ctx, "", "", nil, nil, 100, "")
	if err != nil {
		logger.WithField("err", err).Error("Group list error.")
	} else {
		for _, g := range list {
			logger.Info("ID %s - Name? =%b cursor?=%v", g.Id, g.Name, cursor)
		}
	}
	if len(list) > 0 {
		listofclubs, _ := json.Marshal(list)
		//logger.Debug("listofclubs?=", "{\"list\":"+string(listofclubs)+"}")
		return "{\"list\":" + string(listofclubs) + "}", nil
	}
	return "{\"status\":\"No Club found\"}", nil
}
