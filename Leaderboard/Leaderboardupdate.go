package Leaderboard

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func LeaderboardBanakilisubmited(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	// props := &Leaderboard.LeaderboardSubmitProps{}
	// props.Score = 302
	// props.UserId = "1fe84f38-1ffa-4e05-a837-f862e841b3e9"
	// props.UserName = "WAjVCKCVWl"
	// props.LeaderboardId = "banakilleader_board"
	// payload1, _ := json.Marshal(props)
	rpcReturn, err := LeaderboardBanakili(ctx, logger, db, nk, string(payload))

	return rpcReturn, err
}
