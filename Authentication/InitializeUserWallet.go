package Authentication

import (
	"context"
	"database/sql"

	//"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func AssignUserWallet(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, userid string, chips int64) {

	// userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	userID := userid

	logger.Info("userID", userID)

	/*changeset := map[string]interface{}{
		"coins": float64(chips),
	}
	metadataWallet := map[string]interface{}{
		"game_result": "won",
	}*/
	/*if err := nk.WalletUpdate(ctx, userID, changeset, metadataWallet, true); err != nil {
		// Handle error.
		// logger.Println("Wallet update error -- AssignUserWallet -- ", err)
		// return err
	} else {
		logger.Info(" === WAllet Update ===")
	}
	//return nil */
}
