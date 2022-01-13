package Authentication

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func EmailAuth(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateDeviceRequest) error {

	logger.Info("after device authenticate out value is  '%v'", out)

	//True only once when player created first time
	if out.Created {

	}
	return nil
}
