package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"

	"arabicPoker.com/a/Leaderboard"

	// "arabicPoker.com/a/Clubs/ClubCreate"
	// "arabicPoker.com/a/Clubs/ClubDelete"
	// "arabicPoker.com/a/Clubs/ClubJoin"
	// "arabicPoker.com/a/Clubs/ClubListRequests"
	// "arabicPoker.com/a/Clubs/ClubSearch"

	"arabicPoker.com/a/Authentication"
	// "arabicPoker.com/a/Clubs/ActionLogs"
	"arabicPoker.com/a/Clubs/ClubAdd"
	"arabicPoker.com/a/Clubs/ClubCreate"
	"arabicPoker.com/a/Clubs/ClubDataGet"
	"arabicPoker.com/a/Clubs/ClubDelete"
	"arabicPoker.com/a/Clubs/ClubInvites"
	"arabicPoker.com/a/Clubs/ClubJoin"
	"arabicPoker.com/a/Clubs/ClubKick"
	"arabicPoker.com/a/Clubs/ClubLeave"
	"arabicPoker.com/a/Clubs/ClubListRequests"
	"arabicPoker.com/a/Clubs/ClubLists"
	"arabicPoker.com/a/Clubs/ClubMiscFunctions"
	"arabicPoker.com/a/Clubs/ClubPermissions"
	"arabicPoker.com/a/Clubs/ClubSearch"

	// "arabicPoker.com/a/Clubs/ClubTournaments"
	"arabicPoker.com/a/Clubs/ClubUpdate"
	"arabicPoker.com/a/Clubs/ClubUserList"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubCreatePost"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubDeleteAllPosts"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubDeletePost"
	"arabicPoker.com/a/Clubs/ClubsPosts/ClubReadPosts"

	"arabicPoker.com/a/tarneeb"

	//"arabicPoker.com/a/Authentication"
	"arabicPoker.com/a/banakil"
	"arabicPoker.com/a/fourHundred"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	//"packagess.com/Authentication"
	//"packagess.com/CHK"
	//"packagess.com/chess"
	//"packagess.com/tarneeb"
)

var state ServerState
var wg sync.WaitGroup

type ServerState struct {
	requests   int
	flag       bool
	prevHitVal int
	newHitVal  int
}
type FromMatchSearch struct {
	Gamecode string   `json:"gamecode"`
	MatchId  []string `json:"MatchId"`
	//Status           string  `json:"status"`
	Min        string `json:"min"`
	MaxPlayers string `json:"maxPlayers"`
	HitValue   int    `json:"HitValue"`
	Label      string `json:"label"`
}

type ForgetPwd struct {
	Email string `json:"email"`
}

//////////////////////////////////////////////////////////////////////////////////

type Expire struct {
	// dailyReward int64 `json:"dailyreward"`
	GiftNames     string    `json:"giftName"`
	ExpireDate    int64     `json:"expiredate"`
	received_time time.Time `json:"received_time"`
}

const (
	rpcIdCanClaimDailyReward = "canclaimdailyreward_go"
	rpcIdClaimDailyReward    = "claimdailyreward_go"
)

var (
	errInternalError  = runtime.NewError("internal server error", 13) // INTERNAL
	errMarshal        = runtime.NewError("cannot marshal type", 13)   // INTERNAL
	errNoInputAllowed = runtime.NewError("no input allowed", 3)       // INVALID_ARGUMENT
	errNoUserIdFound  = runtime.NewError("no user ID in context", 3)  // INVALID_ARGUMENT
	errUnmarshal      = runtime.NewError("cannot unmarshal type", 13) // INTERNAL
)

type dailyReward struct {
	LastClaimUnix int64    `json:"last_claim_unix"` // The last time the user claimed the reward in UNIX time.
	GiftName      string   `json:"giftName"`
	ExpireDate    int64    `json:"expiredate"`
	Gifts         []string `json:"Gifts"`
	Expiredates   []int64  `json:"Expiredates"`
}
type gift_expire struct {
	Gifts       []string `json:"Gifts"`
	Expiredates []int64  `json:"Expiredates"`
}

var randumgift gift_expire

//////////////////////////////////////////////////////////////////////////////////
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	//logger.Info("@@@@!!!!!!!loaded")
	Leaderboard_Clubs(ctx, logger, db, nk)
	Leaderboard_Clubs_LastWeek(ctx, logger, db, nk)
	banakil_leaderboard(ctx, logger, db, nk)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	id := "banakilleader_board"
	authoritative := false
	sortOrder := "desc"
	operator := "best"
	resetSchedule := "0 0 * * 1"
	metadata := map[string]interface{}{
		"weather_conditions": "rain",
	}

	if err := initializer.RegisterRpc("verify_email", verifyEmail); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("verify_email_code", verifyCode); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("current_time", getCurrentDateTime); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("pwd_email_verify", verifyEmail); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("pwd_verify_email_code", verifyCode); err != nil {
		return err
	}
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	if err := initializer.RegisterRpc("leaderboardBanakilisubmited", LeaderboardBanakilisubmited); err != nil {
		return err
	}

	// if err := initializer.RegisterRpc("rpcIdClaimDailyReward", RpcClaimDailyReward); err != nil {
	// 	return err
	// }
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	// Register as match handler, this call should be in InitModule.
	if err := initializer.RegisterMatch("TarneebGame", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		logger.Info("TarneebGame RPC Hit")
		type Tarneeb = tarneeb.Tarneeb
		return &Tarneeb{}, nil
	}); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterMatch("400Game", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		logger.Info("400Game RPC Hit")
		type FourHundred = fourHundred.FourHundred
		return &FourHundred{}, nil
	}); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterMatch("BanakilGame", func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
		logger.Info("BanakilGame RPC Hit")
		type Banakil = banakil.Banakil
		return &Banakil{}, nil
	}); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("createMatch", CreateMatchRPC); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("fourHundredCreateMatch", FourHundredCreateMatchRPC); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("BanakilGame", BanakilCreateMatchRPC); err != nil {
		return err
	}

	// Register as matchmaker matched hook, this call should be in InitModule.
	if err := initializer.RegisterMatchmakerMatched(MakeMatch); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("listing", MatchListing); err != nil {
		return err
	}
	// rpcdeletegift
	// RpcCanClaimDailyReward
	if err := initializer.RegisterRpc("getDelayValue", HitDelay); err != nil {
		return err
	}

	///	Faizii
	if err := initializer.RegisterRpc("clubcreaterpc", RpcClaimDailyReward); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubdeleterpc", clubdeleterpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clublistrequestrpc", clublistrequestrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubsearchrpc", clubsearchrpc); err != nil {
		return err
	}
	// if err := initializer.RegisterRpc("delect", rpcdeletegift); err != nil {
	// 	return err
	// }

	if err := initializer.RegisterRpc("clubjoinrpc", clubjoinrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubupdaterpc", clubupdaterpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubpromoterpc", clubpromoterpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubdemoterpc", clubdemoterpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubadduserrpc", clubadduserrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubkickmemberrpc", clubkickmemberrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clublistrpc", clublistrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubinviterpc", clubinviterpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubpromotemanagerrpc", clubpromotemanagerrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubdemotemanagerrpc", clubdemotemanagerrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubcreatepost", clubcreatepost); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubreadpost", clubreadpost); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubdeletepost", clubdeletepost); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubgetdatarpc", clubgetdatarpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubuserlistrpc", clubuserlistrpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubuserleaverpc", clubuserleaverpc); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubusercanjoin", clubusercanjoin); err != nil {
		return err
	}
	//
	if err := initializer.RegisterRpc("clubifmanager", clubifmanager); err != nil {
		return err
	}
	if err := initializer.RegisterRpc("clubdeleteallposts", clubdeleteallposts); err != nil {
		return err
	}
	if err := nk.LeaderboardCreate(ctx, id, authoritative, sortOrder, operator, resetSchedule, metadata); err != nil {
		logger.WithField("err", err).Error("Leaderboard create error.")
	}
	// if err := initializer.RegisterRpc("clubcreatetournament", clubcreatetournament); err != nil {
	// 	return err
	// }
	// if err := initializer.RegisterRpc("clubdeletetournament", clubdeletetournament); err != nil {
	// 	return err
	// }
	// if err := initializer.RegisterRpc("clubgettournamentbyid", clubgettournamentbyid); err != nil {
	// 	return err
	// }
	// if err := initializer.RegisterRpc("clubjointournament", clubjointournament); err != nil {
	// 	return err
	// }
	// if err := initializer.RegisterRpc("clubsubmittournamentscore", clubsubmittournamentscore); err != nil {
	// 	return err
	// }
	// if err := initializer.RegisterRpc("getservertime", getservertime); err != nil {
	// 	return err
	// }
	// if err := initializer.RegisterRpc("getactoinlogs", getactoinlogs); err != nil {
	// 	return err
	// }
	//

	return nil
}

// func getservertime(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	now := time.Now()
// 	startTime := now.Second()

// 	// serverTime := time.Now().Unix()
// 	// // payload = strconv.FormatInt(serverTime, 10)
// 	// startTime = t1.Second()
// 	payloadtime := map[string]interface{}{
// 		"time": startTime,
// 	}
// 	payloada, err := json.Marshal(payloadtime)
// 	return string(payloada), err
// 	// logger.Debug("getservertime by server :: %v " + payload)
// 	// return "{\"time\":" + startTime + "}", nil
// 	// return payload, nil
// }
func Leaderboard_Clubs_LastWeek(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (string, error) {
	id := "clubleaderboard_lastweek"
	authoritative := false
	sort := "desc"
	operator := "incr"
	reset := "0 0 14 * *"
	metadata := map[string]interface{}{
		"coins": int64(0),
	}

	if err := nk.LeaderboardCreate(ctx, id, authoritative, sort, operator, reset, metadata); err != nil {
		logger.Info("Cannot create Leaderboard", err)
		return "error clubleaderboard", err
		// Handle error.
	} else {

		logger.Info("created Leaderboard")
		return "Leaderboard created", err
	}
}
func Leaderboard_Clubs(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (string, error) {
	id := "clubleaderboard"
	authoritative := false
	sort := "desc"
	operator := "incr"
	reset := "0 0 * * 7"
	metadata := map[string]interface{}{
		"coins": int64(0),
	}

	if err := nk.LeaderboardCreate(ctx, id, authoritative, sort, operator, reset, metadata); err != nil {
		logger.Info("Cannot create Leaderboard", err)
		return "error clubleaderboard", err
		// Handle error.
	} else {

		logger.Info("created Leaderboard")
		return "Leaderboard created", err
	}
}
func banakil_leaderboard(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (string, error) {
	id := "banakilleader_board"
	authoritative := false
	sort := "desc"
	operator := "incr"
	reset := "0 0 * * 7"
	metadata := map[string]interface{}{
		"coins": int64(0),
	}

	if err := nk.LeaderboardCreate(ctx, id, authoritative, sort, operator, reset, metadata); err != nil {
		logger.Info("Cannot create Leaderboard", err)
		return "error clubleaderboard", err
		// Handle error.
	} else {

		logger.Info("created Leaderboard")
		// Leaderboard Record Write

		return "Leaderboard created", err
	}
}

// func getactoinlogs(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, _, err := ActionLogs.ReadLogEvent(ctx, logger, db, nk, payload)
// 	return rpcReturn, err
// }
func getservertime(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	serverTime := time.Now().Unix()
	payload = strconv.FormatInt(serverTime, 10)

	logger.Debug("send servertime:? %v " + payload)
	return "{\"time\":" + payload + "}", nil
}

// func clubcreatetournament(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, err := ClubTournaments.CreateClubTournaments(ctx, logger, db, nk, payload)
// 	return rpcReturn, err
// }
// func clubdeletetournament(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, err := ClubTournaments.DeleteClubTournament(ctx, logger, db, nk, payload)
// 	return rpcReturn, err
// }
// func clubgettournamentbyid(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, err := ClubTournaments.GetTournnamentByID(ctx, logger, db, nk, payload)
// 	return rpcReturn, err
// }
// func clubjointournament(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, err := ClubTournaments.JoinClubTournament(ctx, logger, db, nk, payload)
// 	return rpcReturn, err
// }
// func clubsubmittournamentscore(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, err := ClubTournaments.SubmitTournamentScore(ctx, logger, db, nk, payload)
// 	return rpcReturn, err
// }
func clubdeleteallposts(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubDeleteAllPosts.DeleteAllPosts(ctx, logger, db, nk, payload)
	return rpcReturn, err
}

func LeaderboardBanakilisubmited(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	// props := &Leaderboard.LeaderboardSubmitProps{}
	// props.Score = 302
	// props.UserId = "1fe84f38-1ffa-4e05-a837-f862e841b3e9"
	// props.UserName = "WAjVCKCVWl"
	// props.LeaderboardId = "banakilleader_board"
	// payload1, _ := json.Marshal(props)
	rpcReturn, err := Leaderboard.LeaderboardBanakili(ctx, logger, db, nk, string(payload))

	return rpcReturn, err
}
func clubifmanager(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	Return := ClubMiscFunctions.IfManager(ctx, logger, db, nk, payload)
	man := make(map[string]interface{})
	man["ifmanager"] = Return
	rpcReturn, _ := json.Marshal(man)
	return string(rpcReturn), nil
}

func clubusercanjoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubDataGet.CanUserJoinClubs(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubuserleaverpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubLeave.LeaveClubMember(ctx, logger, db, nk, payload)
	return rpcReturn, err
}

func clubuserlistrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubUserList.UserListClub(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubgetdatarpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubDataGet.GetClubData(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubdeletepost(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubDeletePost.DeleteClubPost(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubreadpost(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, _, err := ClubReadPosts.ReadClubPosts(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubcreatepost(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubCreatePost.CreateClubPost(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubdemotemanagerrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubPermissions.PermissionClubDemoteManager(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubpromotemanagerrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubPermissions.PermissionClubPromoteManager(ctx, logger, db, nk, payload)
	return rpcReturn, err
}

func clubinviterpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubInvites.InviteClub(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clublistrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubLists.ListClubs(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubkickmemberrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubKick.KickClubMember(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubadduserrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubAdd.AddUserToClub(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubdemoterpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubPermissions.PermissionClubDemote(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubpromoterpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubPermissions.PermissionClubPromote(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubupdaterpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubUpdate.UpdateClubInfo(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubdeleterpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubDelete.DeleteClub(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubjoinrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubJoin.JoinClub(ctx, logger, db, nk, payload)
	return rpcReturn, err
}
func clubsearchrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubSearch.SearchClubByName(ctx, logger, nk, payload)
	return rpcReturn, err
}

// func clubsearchbycountry(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

// 	rpcReturn, err := ClubSearch.SearchByCountry(ctx, logger, nk, payload)
// 	return rpcReturn, err
// }

func clubcreaterpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubCreate.CreateClub(ctx, logger, db, nk, payload)

	return rpcReturn, err
}
func clublistrequestrpc(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	rpcReturn, err := ClubListRequests.ListClubRequests(ctx, logger, db, nk, payload)

	return rpcReturn, err
}
func verifyEmail(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	from := "dsfds@gmail.com"

	var eml ForgetPwd
	if err := json.Unmarshal([]byte(payload), &eml); err != nil {
		return "", err
	}
	logger.Info("ressseettt  emmaaaiill iss", eml.Email)

	from = eml.Email

	var id int
	const selectSQL = `
SELECT email FROM users WHERE email = $1;
`
	err := db.QueryRowContext(ctx, selectSQL, from).Scan(&id)

	var payloadObj string
	//emailVerification := false
	if err == sql.ErrNoRows {

		logger.Info("error case is.....: ", sql.ErrNoRows)
		payloadObj = "no Email exist"

	} else {

		logger.Info("no error case is.....: ")

		payloadObj, err = Authentication.VerifyEmail(ctx, logger, db, nk, payload)

		logger.Info("email verfied---- :")
		if err != nil {
			logger.Info("err111111111 :-p :: %v", err)
			return "", err
		}
	}

	return payloadObj, nil
}

func verifyCode(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	payloadObj, err := Authentication.VerifyCode(ctx, logger, db, nk, payload)

	if err != nil {
		logger.Info("err :-p :: %v", err)
		return "", err
	}

	logger.Info("aaa1111aaaaaaaaaaaaaaaaaaaaaa :-p :: %v", payloadObj)
	return payloadObj, nil
}

func getCurrentDateTime(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	payloadObj, err := Authentication.GetCurrentDateTime(ctx, logger, db, nk)

	if err != nil {
		logger.Info("err :-p :: %v", err)
		return "", err
	}
	return payloadObj, nil
}

func CreateMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		return "", err
	}
	logger.Info("sssssssssssssss", params)

	modulename := "TarneebGame" // Name with which match handler was registered in InitModule, see example above.
	if matchId, err := nk.MatchCreate(ctx, modulename, params); err != nil {
		return "", err
	} else {

		if err != nil {
			return "", err
		}
		return matchId, nil
	}
}

func FourHundredCreateMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		return "", err
	}
	logger.Info("sssssssssssssss", params)

	modulename := "400Game" // Name with which match handler was registered in InitModule, see example above.
	if matchId, err := nk.MatchCreate(ctx, modulename, params); err != nil {
		return "", err
	} else {

		if err != nil {
			return "", err
		}
		return matchId, nil
	}
}

func BanakilCreateMatchRPC(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		return "", err
	}
	logger.Info("sssssssssssssss", params)

	modulename := "BanakilGame" // Name with which match handler was registered in InitModule, see example above.
	if matchId, err := nk.MatchCreate(ctx, modulename, params); err != nil {
		return "", err
	} else {

		if err != nil {
			return "", err
		}
		return matchId, nil
	}
}

func MakeMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	for _, e := range entries {
		logger.Info("Matched user '%s' named '%s'", e.GetPresence().GetUserId(), e.GetPresence().GetUsername())
		for k, v := range e.GetProperties() {
			logger.Info("Matched on '%s' value '%v'", k, v)
		}
	}

	matchId, err := nk.MatchCreate(ctx, "chess", map[string]interface{}{"invited": entries})
	if err != nil {
		return "", err
	}

	return matchId, nil
}

func MatchListing(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	state.requests++

	logger.Info("TarneebMatchListing CHECK MATCH LISTING PARAMS******* :: %v", payload)

	var frmMatch FromMatchSearch
	if err := json.Unmarshal([]byte(payload), &frmMatch); err != nil {
		return "", err
	}

	logger.Info("stateMain.requests is:: ", state.requests)
	logger.Info("_______")
	//acc, _ := nk.AccountGetId(ctx, "user-id here")
	logger.Info("____++++_____acc is:: ", frmMatch.Label)

	if state.flag {
		state.flag = false
		time.Sleep(3 * time.Second)

	} else {
		state.flag = true
	}

	var s string

	s, _ = chkMatchList(ctx, logger, nk, frmMatch.Label)

	return s, nil
}

func chkMatchList(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, lab string) (string, error) {

	//lab = "TarneebGame"
	label := lab
	logger.Info("module name is:: %v", label)

	min_size := 0
	max_size := 4
	if matches, err := nk.MatchList(ctx, 10, true, label, &min_size, &max_size, "*"); err != nil {
		logger.Info("Error is:: %v", err)

		return "", err
	} else {

		logger.Info("matchess is:: %v", matches)

		logger.Info("NooooooooOOOOOOooooError is:: %v", len(matches))
		obj := &FromMatchSearch{}
		if len(matches) > 0 {
			//logger.Info("ssssssiiiiizzzeee is:: %v", matches[0].HandlerName)
			//logger.Info("Geeetttt  ssssssiiiiizzzeee is:: %v", matches[0].GetSize())

			sort.Slice(matches, func(i, j int) bool {
				return matches[i].Size < matches[j].Size
			})

			//var matchIds []string
			for _, match := range matches {

				//logger.Info("ssssssiiiiizzzeee is:: %v", match.Size)
				//	logger.Info("Geeetttt  sssssiizeee is:: %v", match.GetSize())
				//logger.Info("Match id %s", match.GetMatchId())

				logger.Info("match.GetLabel().GetValue() is:: %v", match.GetLabel().GetValue())

				if match.GetSize() < 4 && match.GetLabel().GetValue() == lab {

					obj.MatchId = append(obj.MatchId, match.GetMatchId())
					//return match.GetMatchId(), nil

				}
			}

			if len(obj.MatchId) > 0 {

				logger.Info("Runing Matches iss:: %v", obj.MatchId)

				if Json_MatchIdList, err := json.Marshal(obj); err != nil {
					logger.Info("Error is: ", err)
				} else {

					return string(Json_MatchIdList), nil

				}
			}

			//onlive if all matches with with complete players
			modulename := lab // Name with which match handler was registered in InitModule, see example above.
			if matchId, err := nk.MatchCreate(ctx, modulename, nil); err != nil {
				return "", err
			} else {

				logger.Info("________________________Inner Maaatch created id : ", matchId)
				obj.MatchId = append(obj.MatchId, matchId)

				if Json_MatchIdList, err := json.Marshal(obj); err != nil {
					logger.Info("Error is: ", err)
				} else {

					return string(Json_MatchIdList), nil
				}
			}
			//onlive if all matches with with complete players

			//return matches[0].MatchId, nil
			//jsonList, _ := json.Marshal(matches)
			//return string(jsonList), nil
		}

		modulename := lab // Name with which match handler was registered in InitModule, see example above.
		if matchId, err := nk.MatchCreate(ctx, modulename, nil); err != nil {
			logger.Info("Outer Maaatch creat error : ", err)

			return "", err
		} else {

			logger.Info("__________________________________Outer Maaatch created id : ", matchId)
			obj.MatchId = append(obj.MatchId, matchId)

			if Json_MatchIdList, err := json.Marshal(obj); err != nil {
				logger.Info("Error is: ", err)
			} else {

				return string(Json_MatchIdList), nil
			}
		}

	}

	return "", nil
}

func HitDelay(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	flag := true
	min := 1
	max := 6
	logger.Info("prev is: ", state.prevHitVal)
	logger.Info("Flag is: ", flag)

	for flag {

		rand.Seed(time.Now().UnixNano())

		state.newHitVal = rand.Intn(max-min+1) + min
		logger.Info("new is: ", state.newHitVal)

		if state.newHitVal != state.prevHitVal {
			flag = false
		}

	}

	state.prevHitVal = state.newHitVal

	FrmMatchHit := &FromMatchSearch{}
	FrmMatchHit.HitValue = state.prevHitVal
	logger.Info("Prev is: ", state.prevHitVal)

	if Json_FrmMatchHit, err := json.Marshal(FrmMatchHit); err != nil {
		logger.Info("Error is: ", err)
	} else {

		return string(Json_FrmMatchHit), nil
	}

	return "", nil
}

func RpcCanClaimDailyReward(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var resp struct {
		CanClaimDailyReward bool `json:"canClaimDailyReward"`
	}

	dailyReward, _, err := getLastDailyRewardObject(ctx, logger, nk, payload)
	if err != nil {
		logger.Error("Error getting daily reward: %v", err)
		return "", errInternalError
	}

	resp.CanClaimDailyReward = canUserClaimDailyReward(dailyReward)

	out, err := json.Marshal(resp)
	if err != nil {
		logger.Error("Marshal error: %v", err)
		return "", errMarshal
	}

	logger.Debug("rpcCanClaimDailyReward resp: %v", string(out))

	return string(out), nil
}
func getLastDailyRewardObject(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, payload string) (dailyReward, *api.StorageObject, error) {
	var d dailyReward
	d.LastClaimUnix = 0
	userID := "a1b38492-a40f-41be-bc60-1c4dcabbc7a2"
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return d, nil, errNoUserIdFound
	}
	payload = ""
	if len(payload) > 0 {
		return d, nil, errNoInputAllowed
	}

	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: "reward",
		Key:        "daily",
		UserID:     userID,
	}})

	if err != nil {
		logger.Error("StorageRead error: %v", err)
		return d, nil, errInternalError
	}

	fmt.Println("len of object    ", len(objects))
	var o *api.StorageObject
	if len(objects) == 0 {
		fmt.Println("cccccccccccccccccccccccccccccccccccccccccccccc")
	} else {
		for _, object := range objects {
			switch object.GetKey() {
			case "daily":
				if err := json.Unmarshal([]byte(object.GetValue()), &d); err != nil {
					logger.Error("Unmarshal error: %v", err)
					return d, nil, errUnmarshal
				}

				return d, object, nil
			}
		}
	}

	return d, o, nil
}
func canUserClaimDailyReward(d dailyReward) bool {
	t := time.Now()
	// midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	midnight := t.Add(time.Minute * 1)
	return time.Unix(d.LastClaimUnix, 0).Before(midnight)

}
func RpcClaimDailyReward(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errNoUserIdFound
	}

	var resp struct {
		GetGift string `json:"gift"`
	}
	// resp.gift = int64(0)
	Gift := []string{"Strawberry ", "Coffee_cups", "100_J_Bucks", "One_Rose", "Star", "One_Day_J_Hukim", "Slice_of_Cake", "Lemon", "Tea", "150_J_Bucks", "Shisha_Hook", "Blue_Eye", "Teddy_Bear", "Cappuccino", "Harte", "Fox_Face", "Unicorn", "Ice_Cream", "Bear", "Lion", "Banana", "Lollipop"}
	number := rand.Int() % len(Gift)
	dailyReward, dailyRewardObject, err := getLastDailyRewardObject(ctx, logger, nk, payload)
	if err != nil {
		logger.Error("Error getting daily reward: %v", err)
		return "", errInternalError
	}
	resp.GetGift = Gift[number]
	// now := time.Now()

	//	t1 := now.Add(time.Hour * 720)
	today := time.Now()
	// yesterday := today.AddDate(0, 0, 30)
	yesterday := today.Add(time.Minute * 2)

	t1 := (yesterday.Unix())
	// t2 := (yesterday.Unix())
	//var randumgift GiftExpire
	randumgift.Gifts = append(randumgift.Gifts, resp.GetGift)
	randumgift.Expiredates = append(randumgift.Expiredates, t1)
	fmt.Println("randumgift_________________________________.Gift", randumgift.Gifts)
	fmt.Println("randumgift___________ccc______________________.Gift", randumgift.Expiredates)
	Json_randumgift, _ := json.Marshal(randumgift)
	s1 := string([]byte(Json_randumgift))

	fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&(((((((((((((((((((((((((((((((((((((&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", dailyReward, "___", dailyRewardObject, " ", s1)
	if canUserClaimDailyReward(dailyReward) {
		fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&(((((((((((((((((((((((((((((((((((((&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&", canUserClaimDailyReward(dailyReward))

		obj := Expire{
			GiftNames:  resp.GetGift,
			ExpireDate: t1,
		}
		Json_FrmMatchHit, _ := json.Marshal(obj)
		s := string([]byte(Json_FrmMatchHit))
		dailyReward.LastClaimUnix = time.Now().Unix()
		dailyReward.GiftName = resp.GetGift
		dailyReward.ExpireDate = t1

		object, err := json.Marshal(dailyReward)
		s = string([]byte(object))
		fmt.Println(s)
		if err != nil {
			logger.Error("Marshal error: %v", err)
			return "", errInternalError
		}

		var randumgift1 gift_expire
		objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
			Collection: "reward",
			Key:        "daily",
			UserID:     userID,
		}})

		if err != nil {
			logger.Error("StorageWrite error: %v", err)
			return "", errInternalError
		} else {

			logger.Error("object is______: %v", objects)
			if len(objects) == 0 {
				dailyReward.Gifts = append(dailyReward.Gifts, resp.GetGift)
				dailyReward.Expiredates = append(dailyReward.Expiredates, t1)
				// listRecords, nextCursor, err := nk.StorageList(ctx, userID, "reward", 10, "")
				Json_FrmMatchHit, _ = json.Marshal(dailyReward)
				s1 := string([]byte(Json_FrmMatchHit))
				_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
					{
						Collection:      "reward",
						Key:             "daily",
						PermissionRead:  1,
						PermissionWrite: 0, // No client write.
						Value:           string(s1),
						// Version:         version,
						UserID: userID,
					}})

			} else {
				for _, object := range objects {
					switch object.GetKey() {
					case "daily":
						if err := json.Unmarshal([]byte(object.GetValue()), &randumgift1); err != nil {
							logger.Error("Unmarshal error: %v", err)

						} else {
							logger.Error("old data____: %v", obj)
							dailyReward.Gifts = append(dailyReward.Gifts, resp.GetGift)
							dailyReward.Expiredates = append(dailyReward.Expiredates, t1)
							listRecords, nextCursor, err := nk.StorageList(ctx, userID, "reward", 10, "")
							// for i
							fmt.Println("qwertyuiojhgfdcvbnmnbvcftgv vbhuhbnjikmko", listRecords, nextCursor, err)
						}

					}
				}
			}
			//	err  = json.Unmarshal(objects, &obj)

		}
		Json_FrmMatchHit, _ = json.Marshal(dailyReward)
		s1 := string([]byte(Json_FrmMatchHit))
		_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
			{
				Collection:      "reward",
				Key:             "daily",
				PermissionRead:  1,
				PermissionWrite: 0, // No client write.
				Value:           string(s1),
				// Version:         version,
				UserID: userID,
			}})

		randumgift.Gifts = nil
		randumgift.Expiredates = nil

		if err != nil {
			logger.Error("StorageWrite error: %v", err)
			return "", errInternalError
		}
	}

	out, err := json.Marshal(resp)
	if err != nil {
		logger.Error("Marshal error: %v", err)
		return "", errMarshal
	}

	logger.Debug("rpcClaimDailyRewcccccccccccccccccard-resp: %v", string(out))

	return string(out), nil
}

func rpcdeletegift(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var keys dailyReward
	userID := "a1b38492-a40f-41be-bc60-1c4dcabbc7a2"

	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: "reward",
		Key:        "daily",
		UserID:     userID,
	}})
	crnttime := time.Now()
	tUnix := crnttime.Unix()

	var o *api.StorageObject
	fmt.Println("oooooooooooooooooooooooooozzzzzooooo", o)
	for _, object := range objects {
		switch object.GetKey() {
		case "daily":
			if err := json.Unmarshal([]byte(object.GetValue()), &keys); err != nil {
				logger.Error("Unmarshal error: %v", err)
				return "", err

			} else {
				var list []int

				for i, explist := range keys.Expiredates {
					if explist <= tUnix {
						list = append(list, i)

					}
				}
				fmt.Println("liiiiiiiiiiiiiiiisssaaa/.,amnhbgavfgbnmssssssssssssssssssttttttttttt", list)
				fmt.Println("liiiiiiiiiiiiiiiisssaaa/.,amnhbgavlemnnnnssssttttttttttt", len(list))

				// var Expiredates []int
				if len(keys.Expiredates) > 2 {
					for _, romoveindex := range list {

						copy(keys.Expiredates[romoveindex:], keys.Expiredates[romoveindex+1:])

						keys.Expiredates[len(keys.Expiredates)-1] = 0
						keys.Expiredates = keys.Expiredates[:len(keys.Expiredates)-1]
						// keys.Expiredates = append(keys.Expiredates[:romoveindex], keys.Expiredates[romoveindex+1])
						// keys.Gifts = append(keys.Gifts[:romoveindex], keys.Gifts[romoveindex+1])

						copy(keys.Gifts[romoveindex:], keys.Gifts[romoveindex+1:])

						keys.Gifts[len(keys.Gifts)-1] = " "
						keys.Gifts = keys.Gifts[:len(keys.Gifts)-1]

						fmt.Println("keys.Expiredates     ", keys.Expiredates)
						fmt.Println("keys.Gifts          ", keys.Gifts)

					}
				} else {
					list = nil
					keys.Expiredates = nil
					keys.Gifts = nil
				}

				fmt.Println("liiiiiiiiiiiiiiiisssssssssssssssssssttttttttttt", keys.Expiredates)
				fmt.Println("liiiiiiiiiiiiiiiissseys.Gifteys.Gifteys.Giftssssssssssssssssttttttttttt", keys.Gifts)

				//after updating a again writing into a database

				Json_FrmMatchHit, _ := json.Marshal(keys)
				s1 := string([]byte(Json_FrmMatchHit))
				_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{
					{
						Collection:      "reward",
						Key:             "daily",
						PermissionRead:  1,
						PermissionWrite: 0, // No client write.
						Value:           string(s1),
						// Version:         version,
						UserID: userID,
					}})
				fmt.Println("errrrrroooooooooorrr     1     ", err)

				return string(s1), err
				// list = nil
			}

		}
	}
	fmt.Println("errrrrroooooooooorrr     2    ", err)

	return "", err
}

func main() {

}
