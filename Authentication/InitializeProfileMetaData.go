package Authentication

import (
	"context"
	"database/sql"

	
	"strconv"
	"github.com/heroiclabs/nakama-common/runtime"
)

func AssignMetaData(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, userid string) {

	// userID := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	userID := userid
	//freeChipsTime := time.Now().Unix() - 2
	matchesplayed := strconv.FormatInt(0, 10)

	logger.Info("No of matches  of  client :: %v " + matchesplayed)

	metadataAccount := make(map[string]interface{})
	//avatarName := GenerateRandomName.GenerateRandomAvatarName()

	//metadataAccount["avatars"] = "defaultAvatar"
	// metadataAccount["decks"] = "defaultdeck"
	// metadataAccount["tables"] = "defaulttable"
	//	metadataAccount["collectFreeChips"] = freeChipsTimeString
	metadataAccount["matchesplayed"] = matchesplayed
	//	metadataAccount["avatarName"] = "defaultavatar"

	if err := nk.AccountUpdateId(ctx, userID, "", metadataAccount, "", "", "", "", ""); err != nil {
		// logger.Println("Error in account metadata assignment -- AssignMetaData -- ", err)
		// return err
	} else {
		logger.Info(" === Account Metadata Update ===")
	}
	//return nil
}
