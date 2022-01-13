package ClubLeaderboard

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func CreateLeaderboard(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, clubid string) error {
	logger.Debug("===CreateLeaderboard function called===")
	id := clubid
	authoritative := true
	sort := "desc"
	operator := "incr"
	reset := "0 0 * * 0"
	metadata := map[string]interface{}{}

	if err := nk.LeaderboardCreate(ctx, id, authoritative, sort, operator, reset, metadata); err != nil {
		logger.WithField("err", err).Error("Leaderboard creation error.")
		return err
	}
	logger.Debug("leaderboard created id?%v", id)
	return nil

}
func DeleteLeaderboard(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, clubid string) error {
	logger.Debug("===Deleteleaderboard function called===")
	if err := nk.LeaderboardDelete(ctx, clubid); err != nil {
		logger.WithField("err", err).Error("Leaderboard delete error.")
		return err
	}
	logger.Debug("leaderboard Deleted Successfully id?%v", clubid)
	return nil
}
