package Authentication

import (
	"context"
	"database/sql"

	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/speps/go-hashids"
)

func AuthDevice(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateDeviceRequest) error {

	logger.Info("after device authenticate out value is  '%v'", out)

	//True only once when player created first time
	if out.Created {
		userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

		//AuthenticatAssignUserWallet(ctx, logger, db, nk, userID, int64(1000000))
		AssignMetaData(ctx, logger, db, nk, userID)
		AssignUnlockCollection(ctx, logger, db, nk, userID)

		hd := hashids.NewData()
		hd.Salt = "this is my salt"
		h, _ := hashids.NewWithData(hd)
		timestamp := time.Now().UnixNano()
		t := int(timestamp)
		id, _ := h.Encode([]int{t})
		displayName := "Guest" + id

		AssignUserMatches(ctx, logger, nk, userID)

		nk.AccountUpdateId(ctx, userID, "", nil, displayName, "", "", "", "")
	}
	return nil
}
