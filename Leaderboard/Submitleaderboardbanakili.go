package Leaderboard

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func LeaderboardBanakili(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("=== LeaderboardBanakili    Rpc CAlled ===?%v", payload)
	props := &LeaderboardSubmitProps{}
	if err := json.Unmarshal([]byte(payload), &props); err != nil {

		logger.Debug("Unmarshal failed?=%v", props)
		return "{\"status\":\"InValidProps\"}", err
	}
	// Leaderboard Record Write
	// id := "4ec4f126-3f9d-11e7-84ef-b7c182b36521"
	// ownerID := "4c2ae592-b2a7-445e-98ec-697694478b1c" //same rhny ge
	// username := "02ebb2c8"
	// score := int64(10)
	// subscore := int64(0)

	metadata := map[string]interface{}{}
	logger.Debug("=== LeaderboardBanakili  props.Scoreprops.Scoreprops.Scoreprops.Score  Rpc CAlled ===?%v", props.Score)
	if record, err := nk.LeaderboardRecordWrite(ctx, props.LeaderboardId, props.UserId, props.UserName, int64(props.Score), 0, metadata, nil); err != nil {
		logger.WithField("err", err).Error("Leaderboard record write error.")
		return "{\"status\":\"Leaderboard record write error.\"}", err
	} else {
		logger.Debug("record submitted   ?=%v", record)
		return "{\"status\":\"record submitted successfully\"}", err

	}

}
