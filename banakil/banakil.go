// Copyright 2020 The Nakama Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package banakil

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"arabicPoker.com/a/Leaderboard"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type DuringMatchForDummy struct{}

type Team struct {
	teamScore  int
	member1    string
	member2    string
	teamBid    int
	roundScore int
}

type Player struct {
	userId         string
	displayName    string
	userName       string
	placeHolderUrl string

	teamMate     string
	teamNum      string
	points       int
	seatPosition int
	cards        []string
	presence     runtime.Presence
	score        int
	// cardWithValue  map[string]float64 //uzair
	opcode           int64
	cardWithValue    map[string]*mapping
	mainList         []SetStruct
	takeCardfromPile bool
	putCardinPile    bool
	firstGodown      bool
	scoreCards       []GoDownSet
	meldScore        float32
	roundScore       float32
	handScore        float32
	totalScore       float32
}

type mapping struct {
	valueofcards float64
	status       bool
}
type MatchState struct {
	debug         bool
	joinedPlayers int
	userNames     []string
	players       map[string]*Player
	teams         map[string]*Team

	deck               [106]string
	users              []*api.User
	hostUserName       string
	hostDisplayName    string
	hostSignal         bool
	otherPlayers       bool
	playersLimit       int
	playersLimitSignal bool
	sittingArangement  [4]string
	hostOpponent       string
	dealer             string
	bidFlag            bool

	tableCardsTeam1 []GoDownSet
	tableCardsTeam2 []GoDownSet

	restricgtedCardOfDiscard string

	firstTurnsCounter int
	firstTurn         string

	cardWithValue map[string]*mapping

	highestBid int

	botOpcode int64

	allowCarddata bool

	matchExitCounter int
	matchExitBool    bool

	stockpile []string

	discardpile       []string
	temporaryHit      int
	matchFirstMeld    bool
	fourtyBonus       bool
	userFourtyBonus   string
	oneTim106         bool
	roundScoreCounter int
}

type SetStruct struct {
	Set []int "Set"
}

type MainStruct struct {
	//mainList []*SetStruct
}

type GetTeamMate struct {
	X string `json:"x"`
	Y bool   `json:"y"`
}

type SetList struct {
	Set []SetStruct `json:"Sets"`
}
type Before_GamePlay struct {
	FirstPlayer      bool            `json:"FirstPlayer"`
	PlayersJoined    int             `json:"PlayersJoined"`
	HostUserName     string          `json:"HostUserName"`
	HostDisplayName  string          `json:"HostDisplayName"`
	JoinedPlayers    []Name_UserName `json:"JoinedPlayers"`
	Winner1          string          `json:"Winner1"`
	Winner2          string          `json:"Winner2"`
	started          bool            `json:"Started"`
	SetTeamMate      string          `json:"SetTeamMate"`
	TeamMateResponse GetTeamMate     `json:"TeamMateResponse"`

	SittingArangement []string `json:"SittingArangement"`

	Dealer string `json:"Dealer"`
	Turn   string `josn:"Turn"`

	PlayerUserName string `json:"PlayerUserName"`

	TakeCardfromPile bool `json:"TakeCardfromPile"`

	Cards               []string `json:"Cards"`
	ThrownCard          string   `json:"ThrownCard"`
	Chat                string   `json:"Chat"`
	HostRoundScore1     float32  `json:"HostRoundScore1"`
	HostOppoRoundScore2 float32  `json:"HostOppoRoundScore2"`
	Team1               float32  `json:"Team1"`
	Team2               float32  `json:"Team2"`
	Exit                string   `json:"Exit"`

	RemoveCards string `json:"RemoveCards"`

	HostTeam         int `json:"HostTeam"`
	HostOpponentTeam int `json:"HostOpponentTeam"`

	Sets               []SetList             `json:"Sets"`
	StockPile          []string              `json:"StockPile"`
	DiscardPile        []string              `json:"DiscardPile"`
	GoDownCardsData    GoDownCardsDataStruct `json:"GoDownCardsData"`
	CardName           string                `json:"cardname"`
	PutCardinPile      bool                  `json:"PutCardinPile"`
	FirstGodown        bool                  `json:"FirstGodown"`
	StockPileIndex     string                `json:"StockPileIndex"`
	DiscardPileObj     DiscardPile           `json:"discardpilecard"`
	TableSequenceCards Melds                 `json:"TableSequenceCards"`
	CardIncreased      int                   `json:"CardIncreased"`
}

type Name_UserName struct {
	Name      string
	UserName  string
	AvatarUrl string
	TeamNum   string
	ID        string
}
type Melds struct {
	TablecardsSequencelist []string `json:"TablecardsSequencelist"`
	CardName               string   `json:"cardName"`
	AllowMeld              bool     `json:"AllowMeld"`
	UserName               string   `json:"UserName"`
	CardIndex              int      `json:"cardIndex"`
	CurrentListIndex       int      `json:"CurrentListIndex"`
	NotMatchedcardName     string   `json:"NotMatchedcardName"`
}

type HandCardLisdt struct {
	HandCard GetHandList `json:"HandCard"`
}
type GetHandList struct {
	CardList []string `json:"CardList"`
	UserName string   `json:"userName"`
	ScoreCal bool     `json:"ScoreCal"`
}

type GoDownCardsDataStruct struct {
	GoDownCardsMain []GoDownSet `json:"GoDownCardsMain"`
	UserName        string      `json:"userName"`
	HandEmptyOrNot  bool        `json:"HandEmptyOrNot"`
}

type GoDownSet struct {
	Set []string `json:"strngArr"`
}

type DiscardPile struct {
	DiscardList []string `json:"disCardList"`
	UserName    string   `json:"userName"`
}

type Banakil struct{}

func (b *Banakil) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {

	debug := false

	state := &MatchState{
		debug:              debug,
		hostSignal:         true,
		joinedPlayers:      0,
		players:            make(map[string]*Player),
		teams:              make(map[string]*Team),
		deck:               [106]string{"joker_A_3_1", "joker_A_3_0", "club_A_1", "club_2_1", "club_3_1", "club_4_1", "club_5_1", "club_6_1", "club_7_1", "club_8_1", "club_9_1", "club_10_1", "club_J_1", "club_Q_1", "club_K_1", "diamond_A_1", "diamond_2_1", "diamond_3_1", "diamond_4_1", "diamond_5_1", "diamond_6_1", "diamond_7_1", "diamond_8_1", "diamond_9_1", "diamond_10_1", "diamond_J_1", "diamond_Q_1", "diamond_K_1", "heart_A_1", "heart_2_1", "heart_3_1", "heart_4_1", "heart_5_1", "heart_6_1", "heart_7_1", "heart_8_1", "heart_9_1", "heart_10_1", "heart_J_1", "heart_Q_1", "heart_K_1", "spade_A_1", "spade_2_1", "spade_3_1", "spade_4_1", "spade_5_1", "spade_6_1", "spade_7_1", "spade_8_1", "spade_9_1", "spade_10_1", "spade_J_1", "spade_Q_1", "spade_K_1", "spade_K_1", "club_2_0", "club_3_0", "club_4_0", "club_5_0", "club_6_0", "club_7_0", "club_8_0", "club_9_0", "club_10_0", "club_J_0", "club_Q_0", "club_K_0", "diamond_A_0", "diamond_2_0", "diamond_3_0", "diamond_4_0", "diamond_5_0", "diamond_6_0", "diamond_7_0", "diamond_8_0", "diamond_9_0", "diamond_10_0", "diamond_J_0", "diamond_Q_0", "diamond_K_0", "heart_A_0", "heart_2_0", "heart_3_0", "heart_4_0", "heart_5_0", "heart_6_0", "heart_7_0", "heart_8_0", "heart_9_0", "heart_10_0", "heart_J_0", "heart_Q_0", "heart_K_0", "spade_A_0", "spade_2_0", "spade_3_0", "spade_4_0", "spade_5_0", "spade_6_0", "spade_7_0", "spade_8_0", "spade_9_0", "spade_10_0", "spade_J_0", "spade_Q_0", "spade_K_0"},
		playersLimit:       4,
		playersLimitSignal: true,
		bidFlag:            false,
		highestBid:         0,
		cardWithValue:      make(map[string]*mapping), //uzair
		allowCarddata:      true,
		sittingArangement:  [4]string{"BAtXtNDFgZ", "BAtXtNDFgZ", "BAtXtNDFgZ", "BAtXtNDFgZ"},
		firstTurnsCounter:  0,
		oneTim106:          true,
	}

	if state.debug {
		logger.Info("match init, starting with debug: %v", state.debug)
	}
	tickRate := 1

	shuffle(state.deck, state)
	label := "BanakilGame"
	return state, tickRate, label
}

func (b *Banakil) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	state := stateMain.(*MatchState)
	var statement string
	if stateMain.(*MatchState).debug {
		logger.Info("match join attempt username %v user_id %v session_id %v node %v with metadata %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId(), metadata)
	}

	logger.Info("___@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@", state.joinedPlayers)

	state.joinedPlayers++
	if state.joinedPlayers <= state.playersLimit {

		setUserNameList(presence, state, logger)
		setPlayerData(presence, ctx, state, logger, nk)

		logger.Info("___________Match Join Attempt, starting: ", state.joinedPlayers)
		if state.joinedPlayers == 1 {
			state.hostSignal = true

		} else {
			state.otherPlayers = true
		}

		logger.Info("_____***************______Match Join Attempt, starting: ", presence.GetUsername())

		if state.joinedPlayers == 1 {
			state.hostUserName = presence.GetUsername()
		}

		return state, true, statement
	} else {
		state.joinedPlayers--

		statement = "Player capacity reached"
	}
	logger.Info("_____/////////////////////////////______Match Join Attempt, starting: ", presence.GetUsername())

	return state, false, statement
}

func (b *Banakil) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, presences []runtime.Presence) interface{} {
	state := stateMain.(*MatchState)
	if stateMain.(*MatchState).debug {
		for _, presence := range presences {
			logger.Info("match join username %v user_id %v session_id %v node %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId())

		}
	}

	logger.Info("___________Match Join, starting", state.joinedPlayers)

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)
	if !ok {

		logger.Info("presence isssss: %s", ok, userID)
		// return "", errors.New("Invalid context")
	}

	//logger.Info("*********************************-", state.joinedPlayers)

	return state
}

func (b *Banakil) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, presences []runtime.Presence) interface{} {
	state := stateMain.(*MatchState)

	for _, presence := range presences {
		logger.Info("match leave username %v user_id %v session_id %v node %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId())
	}

	state.matchExitCounter--
	state.matchExitBool = true

	logger.Info("___________Match Leave, starting ", state.matchExitCounter)

	return state
}

func setUserNameList(presence runtime.Presence, stateMain interface{}, logger runtime.Logger) {
	userName := presence.GetUsername()
	state := stateMain.(*MatchState)
	flag := true
	if len(state.userNames) > 0 {
		for _, uzrNam := range state.userNames {

			if uzrNam == userName {

				//flag = false
				return
			}
		}
		if flag {
			state.userNames = append(state.userNames, userName)
			//state.joinedPlayers++
			state.matchExitCounter++

		}
	} else {
		state.userNames = append(state.userNames, userName)
		//state.joinedPlayers++
		state.matchExitCounter++
	}
	logger.Info("___________setUserNameList, state.userNames: ", state.userNames)

}

func setPlayerData(presence runtime.Presence, ctx context.Context, stateMain interface{}, logger runtime.Logger, nk runtime.NakamaModule) {
	state := stateMain.(*MatchState)
	userName := presence.GetUsername()
	arr := []string{}
	arr1 := []SetStruct{}
	var emptyMap map[string]*mapping
	emptyMap = make(map[string]*mapping)

	flag := true
	state.users, _ = nk.UsersGetUsername(ctx, state.userNames)
	for _, uzr := range state.players {

		if uzr.userName == userName {

			return
		}
	}
	if flag {

		for _, uzr := range state.users {

			if uzr.GetUsername() == userName && uzr.GetUsername() != "" {
				//	uzr.GetAvatarUrl

				player := &Player{uzr.GetId(), uzr.GetDisplayName(), uzr.GetUsername(), uzr.GetAvatarUrl(), "", "", 0, 0, arr, presence, 0, 0, emptyMap, arr1, true, false, false, nil, 0, 0, 0, 0}

				state.players[userName] = player
				logger.Info("___________ setPlayerData, state.Players: ", state.players[userName].seatPosition)

				return
			} else {
				logger.Info("userName is Empty: ", uzr.GetUsername())

			}
		}

	}
	logger.Info("___________ setPlayerData, state.Players: ", state.players[userName].seatPosition)

}

func (b *Banakil) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, messages []runtime.MatchData) interface{} {
	state := stateMain.(*MatchState)
	if state.debug {
		logger.Info("match loop match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match loop match_id %v message count %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), len(messages))
	}

	if state.matchExitCounter <= 0 && state.matchExitBool {
		logger.Info("_____going to kil________-", state.matchExitCounter)
		logger.Info("______going to kil_______-", state.matchExitBool)

		return nil
	}
	HostPlyr := &Before_GamePlay{}
	// dataFromBot := &Before_GamePlay{}
	var dataFromBot_byte []byte

	if state.joinedPlayers == 4 {
		//	logger.Info("_____________-", state.joinedPlayers)

	}

	if state.joinedPlayers == 1 && state.hostSignal {

		logger.Info("____0000000_______Match , starting", state.hostUserName)

		HostPlyr.FirstPlayer = true
		if Json_started, err := json.Marshal(HostPlyr); err != nil {
			logger.Info("Error is: ", err)

		} else {
			logger.Info("state.hostSignal is: ", state.hostSignal)

			state.hostSignal = false
			logger.Info("____1000_______BroadcastMessage , starting", state.joinedPlayers)

			presence := []runtime.Presence{state.players[state.hostUserName].presence}

			dispatcher.BroadcastMessage(100, Json_started, presence, nil, true)

		}
	} else if state.joinedPlayers < state.playersLimit && state.otherPlayers {
		OtherPlyr := &Before_GamePlay{}

		OtherPlyr.FirstPlayer = false
		if Json_started, err := json.Marshal(OtherPlyr); err != nil {
			logger.Info("Error is: ", err)

		} else {
			state.otherPlayers = false
			dispatcher.BroadcastMessage(100, Json_started, nil, nil, true)

		}

	}

	//logger.Info("state.joinedPlayers______________: ", state.joinedPlayers)

	if state.joinedPlayers == state.playersLimit && state.playersLimitSignal {

		JoiningData := &Before_GamePlay{}
		Nam_UserNam := Name_UserName{}
		for _, plyr := range state.players {
			Nam_UserNam.UserName = plyr.userName
			Nam_UserNam.Name = plyr.displayName
			Nam_UserNam.AvatarUrl = plyr.placeHolderUrl
			Nam_UserNam.ID = plyr.userId

			JoiningData.JoinedPlayers = append(JoiningData.JoinedPlayers, Nam_UserNam)
		}

		//JoiningData = &Before_GamePlay{state.joinedPlayers, state.hostUserName, JoiningData.JoinedPlayers, false, "", }
		JoiningData.PlayersJoined = state.joinedPlayers
		JoiningData.HostUserName = state.hostUserName
		JoiningData.HostDisplayName = state.hostDisplayName

		JoiningData.started = false

		if JsonJoiningData, err := json.Marshal(JoiningData); err != nil {
			logger.Info("Error is: ", err)

		} else {

			dispatcher.BroadcastMessage(101, JsonJoiningData, nil, nil, true)
			state.playersLimitSignal = false
			logger.Info("hittteedd 1111110000OOOOOoooooo11111 is: ", state.hostUserName)

			if state.hostUserName != "" {

			}

		}

	}

	if tick >= 10 {
		//	return nil
		//	logger.Info("tick is: ", tick)

	}

	var data []byte
	var opcd int64
	for _, msg := range messages {

		logger.Info("msg username", msg.GetUsername())

		logger.Info("msg status", msg.GetStatus())

		logger.Info("msg opcode", msg.GetOpCode())

		logger.Info("msg Data", string(msg.GetData()))
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		logger.Info("######", string(msg.GetData()))

		fmt.Println("")
		fmt.Println("")
		fmt.Println("")

		decoded := &DuringMatchForDummy{}

		if err := json.Unmarshal(msg.GetData(), &decoded); err != nil {
			logger.Info("Error is: ", err)

		}

		data = msg.GetData()
		opcd = msg.GetOpCode()

		logger.Info("=opcd After decode=", opcd)

		//	logger.Info("=msg After decode=", data)

		var opCOde int64
		if state.botOpcode == 0 {
			opCOde = opcd //msg.GetOpCode()
		} else {
			opCOde = state.botOpcode //msg.GetOpCode()
			state.botOpcode = 0
			data = dataFromBot_byte

		}
		//	logger.Info("=Data After Opcode=", data)

		switch opCOde {

		case 400:

			host := &Before_GamePlay{}

			if err := json.Unmarshal(data, &host); err != nil {
				logger.Info("Error is: ", err)
			} else {
				logger.Info("Data is: ", host)

				state.hostUserName = host.HostUserName
				state.hostDisplayName = host.HostDisplayName

			}

			//opCOde = 700
			logger.Info("400 Data is________________________: ", state.hostUserName)

			// opCOde = 0
			// opcd = 0
		case 401:

			logger.Info("401 401 401___________________: ", state.hostUserName)

			Started := &Before_GamePlay{}
			Started.started = true
			if Json_started, err := json.Marshal(Started); err != nil {
				logger.Info("Error is: ", err)

			} else {

				logger.Info("Data is________________________: ", state.hostUserName)

				dispatcher.BroadcastMessage(102, Json_started, nil, nil, true)
				cardWithValue(state, logger)

			}
			// opCOde = 0
			// opcd = 0
		case 402:
			logger.Info("402 is 402 ________________________: ", state.hostUserName)

			TeamMate := &Before_GamePlay{}
			if err := json.Unmarshal(data, &TeamMate); err != nil {
				logger.Info("Error is: ", err)
			} else {

				if TeamMate.SetTeamMate != "" {

					logger.Info("402 is state.players[TeamMate.PlayerUserName].presence_______: ", state.players[TeamMate.SetTeamMate].presence)

					presence := []runtime.Presence{state.players[TeamMate.SetTeamMate].presence}

					if Json_TeamMateReq, err := json.Marshal(TeamMate.SetTeamMate); err != nil {
						logger.Info("Error is: ", err)
					} else {
						logger.Info("402 sseendd 103 ________________________: ", state.hostUserName)

						dispatcher.BroadcastMessage(103, Json_TeamMateReq, presence, nil, true)

						//IF Requested Person IS NOT AVAILABLE

					}
				}

			}

		case 403:

			// After team mate slction
			logger.Info("403    TeamMate.TeamMateResponse.y___________________: ")

			TeamMate := &Before_GamePlay{}
			if err := json.Unmarshal(data, &TeamMate); err != nil {
				logger.Info("Error is: ", err)
			} else {

				TeamMate.SetTeamMate = TeamMate.TeamMateResponse.X

				if TeamMate.SetTeamMate != "" {

					logger.Info("403 TeamMate.TeamMateResponse.y___________________: ", TeamMate.TeamMateResponse.Y)
					logger.Info("403 TeamMate.TeamMateResponse.y___________________: ", TeamMate.TeamMateResponse.X)

					if TeamMate.TeamMateResponse.Y {
						TeamMate.SetTeamMate = TeamMate.TeamMateResponse.X

						for _, plyr := range state.players {
							if plyr.userName == state.hostUserName {

								state.players[plyr.userName].teamMate = TeamMate.SetTeamMate
							} else if plyr.userName == TeamMate.SetTeamMate {
								state.players[TeamMate.SetTeamMate].teamMate = state.hostUserName
							} else {
								for _, p := range state.players {

									if p.userName != state.hostUserName && p.userName != TeamMate.SetTeamMate && p.userName != plyr.userName {

										state.players[p.userName].teamMate = plyr.userName
										state.players[plyr.userName].teamMate = p.userName

									}
								}
							}
						}

						// for _, playercard := range state.players {
						// 	logger.Info("&&&&&&&&&&&&&&&&&", playercard)
						// }

						arrangePositions(state, logger)
						B4GamePlay := &Before_GamePlay{}

						sittingArrangement := &Before_GamePlay{}
						Nam_UserNam := Name_UserName{}
						for _, seat := range state.sittingArangement {

							Nam_UserNam.UserName = state.players[seat].userName
							Nam_UserNam.Name = state.players[seat].displayName
							Nam_UserNam.AvatarUrl = state.players[seat].placeHolderUrl
							Nam_UserNam.TeamNum = state.players[seat].teamNum
							Nam_UserNam.ID = state.players[seat].userId

							sittingArrangement.JoinedPlayers = append(sittingArrangement.JoinedPlayers, Nam_UserNam)
						}

						if Json_sittingArrangement, err := json.Marshal(sittingArrangement); err != nil {
							logger.Info("Error is: ", err)
						} else {
							dispatcher.BroadcastMessage(104, Json_sittingArrangement, nil, nil, true)

						}
						dealerSelection(state, logger)

						B4GamePlay.Dealer = state.dealer
						if Json_Dealer, err := json.Marshal(B4GamePlay); err != nil {
							logger.Info("Error is: ", err)
						} else {
							dispatcher.BroadcastMessage(105, Json_Dealer, nil, nil, true)
						}
						///////////////////////////////////////////////////////////////Next will statrt from here///////////////////////////
						afterSitting(state, logger, B4GamePlay, dispatcher)
						// if Json_stockpile, err := json.Marshal(B4GamePlay); err != nil {
						// 	logger.Info("Error is: ", err)

						// } else {

						// 	dispatcher.BroadcastMessage(201, Json_stockpile, nil, nil, true)

						// }

					} else {

						logger.Info("rejected Resending the 102 ___: ", state.hostUserName)

						Started := &Before_GamePlay{}
						Started.started = true
						if Json_started, err := json.Marshal(Started); err != nil {
							logger.Info("Error is: ", err)

						} else {
							logger.Info("to show again team list_________: ", state.hostUserName)

							dispatcher.BroadcastMessage(102, Json_started, nil, nil, true)

						}
					}
				}

			}

		case 404:
			HandCardList := &HandCardLisdt{}

			if err := json.Unmarshal(data, &HandCardList); err != nil {
				logger.Info("Error is: ", err)
			} else {
				logger.Info("44400444 is: ", HandCardList.HandCard)

				state.players[HandCardList.HandCard.UserName].mainList = nil
				presence := []runtime.Presence{state.players[HandCardList.HandCard.UserName].presence}

				logger.Info("*********state.players[HandCardList.HandCard.UserName].takeCardfromPile : ", state.players[HandCardList.HandCard.UserName].takeCardfromPile)

				if len(HandCardList.HandCard.CardList) <= 0 && !state.players[HandCardList.HandCard.UserName].takeCardfromPile || HandCardList.HandCard.ScoreCal { //|| state.temporaryHit >= 1 {
					logger.Info("***************************************************************************** : ")

					ScoreAlert := &Before_GamePlay{}

					if !HandCardList.HandCard.ScoreCal {
						Broadcast(ScoreAlert, int64(216), logger, nil, dispatcher)

					} else {
						if scoreCalc_404(state, HandCardList, logger, nil, dispatcher) {

							B4GamePlay := &Before_GamePlay{}

							//FOR CHECK
							B4GamePlay.Dealer = TurnSequence(state.players[state.dealer].seatPosition, state, logger)
							state.dealer = B4GamePlay.Dealer
							if Json_Dealer, err := json.Marshal(B4GamePlay); err != nil {
								logger.Info("Error is: ", err)
							} else {
								dispatcher.BroadcastMessage(105, Json_Dealer, nil, nil, true)
							}

							if state.oneTim106 {
								afterSitting(state, logger, B4GamePlay, dispatcher)

								state.oneTim106 = false
							}
						}

					}

					//	FOR CHECK

				} else {

					if len(HandCardList.HandCard.CardList) <= 0 {
						logger.Info("400004444 returning: ", err)
						submitLeadboardscore(state, logger, ctx, db, nk)
						return nil
					}
					sendSets(HandCardList.HandCard.CardList, HandCardList.HandCard.UserName, state, logger, presence, dispatcher)
					state.oneTim106 = true
				}

				//Scoring while ending hand

			}

		case 405:

			GoDownResp := &Before_GamePlay{}
			state.temporaryHit++
			if err := json.Unmarshal(data, &GoDownResp); err != nil {
				logger.Info("Error is: ", err)
			} else {
				logger.Info("4000555 is: ", GoDownResp.GoDownCardsData)

				goDown_405(state, GoDownResp, logger, nil, dispatcher, ctx, db, nk)

			}
		case 406:

			DiscardPile_obj := &Before_GamePlay{}

			//logger.Info("Error is: ", err)
			logger.Info("__________CardFOR Discard pile ____ 44000666666")

			if err := json.Unmarshal(data, &DiscardPile_obj); err != nil {
				logger.Info("Error is: ", err)
			} else {

				logger.Info("__________Response Discard pile ____ 4400066666", DiscardPile_obj.DiscardPileObj)

				allow := false
				if DiscardPile_obj.DiscardPileObj.UserName != "" {

					state.players[DiscardPile_obj.DiscardPileObj.UserName].takeCardfromPile = false
					state.players[DiscardPile_obj.DiscardPileObj.UserName].putCardinPile = true
					DiscardPile_obj.Turn = DiscardPile_obj.DiscardPileObj.UserName
					DiscardPile_obj.CardIncreased = len(DiscardPile_obj.DiscardPileObj.DiscardList)

					if len(DiscardPile_obj.DiscardPileObj.DiscardList) == len(state.discardpile) {
						for i, discardCard := range state.discardpile {
							logger.Info("Discard card removing: ", discardCard)
							remove(state.discardpile, i)
							allow = true
						}

					} else if len(DiscardPile_obj.DiscardPileObj.DiscardList) > 1 && len(DiscardPile_obj.DiscardPileObj.DiscardList) < len(state.discardpile) {

						for _, client_discardCard := range DiscardPile_obj.DiscardPileObj.DiscardList {

							for j, discardCard := range state.discardpile {
								//logger.Info("Discard card removing: ", discardCard)
								if client_discardCard == discardCard {
									remove(state.discardpile, j)
									allow = true

									break

								}
							}

						}

					} else if len(DiscardPile_obj.DiscardPileObj.DiscardList) == 1 {
						for i, discardCard := range state.discardpile {
							//logger.Info("Discard card removing: ", discardCard)
							if DiscardPile_obj.DiscardPileObj.DiscardList[0] == discardCard {

								state.restricgtedCardOfDiscard = DiscardPile_obj.DiscardPileObj.DiscardList[0]
								remove(state.discardpile, i)
								allow = true

							}

						}
					}

					if allow {

						//append card in player card list after the demo

						//DiscardPile_obj.DiscardPileObj.DiscardList = nil
						DiscardPile_obj.TakeCardfromPile = state.players[DiscardPile_obj.DiscardPileObj.UserName].takeCardfromPile
						DiscardPile_obj.PutCardinPile = state.players[DiscardPile_obj.DiscardPileObj.UserName].putCardinPile

						//DiscardPile_obj.DiscardPileObj.UserName = ""

						Broadcast(DiscardPile_obj, int64(212), logger, nil, dispatcher)
					}
				}
			}

		case 407:

			CardfromstockPile := &Before_GamePlay{}

			//logger.Info("Error is: ", err)
			logger.Info("____440007777 ____ ") //, state.stockpile)

			if err := json.Unmarshal(data, &CardfromstockPile); err != nil {
				logger.Info("Error is: ", err)
			} else {

				logger.Info("__________Response:  ____ ", CardfromstockPile)

				stockPileIndex, _ := strconv.Atoi(CardfromstockPile.StockPileIndex)
				logger.Info("__________len(state.stockpile)", len(state.stockpile))

				logger.Info("____400777_____state.stockpile)", state.stockpile)
				if len(state.stockpile) > 0 && state.stockpile[0] == "" {

					ScoreAlert := &Before_GamePlay{}

					Broadcast(ScoreAlert, int64(216), logger, nil, dispatcher)

				} else if stockPileIndex < len(state.stockpile) {

					for i, card := range state.stockpile {
						fmt.Println(i, ": ", card)
					}
					logger.Info("____400777_____state.stockpile[stockPileIndex])", state.stockpile[stockPileIndex])

					if state.stockpile[stockPileIndex] == CardfromstockPile.CardName {
						state.players[CardfromstockPile.PlayerUserName].takeCardfromPile = false
						state.players[CardfromstockPile.PlayerUserName].putCardinPile = true

						remove(state.stockpile, stockPileIndex)

						//CardfromstockPile.StockPileIndex = ""
						//CardfromstockPile.CardName = ""
						CardfromstockPile.Turn = CardfromstockPile.PlayerUserName
						CardfromstockPile.TakeCardfromPile = state.players[CardfromstockPile.PlayerUserName].takeCardfromPile
						CardfromstockPile.PutCardinPile = state.players[CardfromstockPile.PlayerUserName].putCardinPile
						//CardfromstockPile.StockPile = nil
						//CardfromstockPile.StockPile = append(CardfromstockPile.StockPile, state.stockpile...)

						logger.Info("____400777_____After remove)", CardfromstockPile.StockPile)

						Broadcast(CardfromstockPile, int64(213), logger, nil, dispatcher)
						//CardfromstockPile.StockPile = nil
						CardfromstockPile.CardName = ""
						CardfromstockPile.CardIncreased = 1
						Broadcast(CardfromstockPile, int64(212), logger, nil, dispatcher)

					}

				}

			}

		case 408:

			CardForDiscardPile := &Before_GamePlay{}

			//logger.Info("Error is: ", err)
			logger.Info("__________Card in Discard pile ____ 440008888")

			if err := json.Unmarshal(data, &CardForDiscardPile); err != nil {
				logger.Info("Error is: ", err)
			} else {

				logger.Info("__________Card in Discard pile ____ 440008888", CardForDiscardPile.CardName)
				logger.Info("_________state.restricgtedCardOfDiscard ", state.restricgtedCardOfDiscard)

				if CardForDiscardPile.CardName != state.restricgtedCardOfDiscard {

					if state.players[CardForDiscardPile.PlayerUserName].putCardinPile == true {
						state.discardpile = append(state.discardpile, CardForDiscardPile.CardName)

						state.players[CardForDiscardPile.PlayerUserName].putCardinPile = false
						state.players[CardForDiscardPile.PlayerUserName].takeCardfromPile = true

						Broadcast(CardForDiscardPile, int64(211), logger, nil, dispatcher)

						B4GamePlay := &Before_GamePlay{}

						B4GamePlay.DiscardPile = nil
						B4GamePlay.PlayerUserName = CardForDiscardPile.PlayerUserName
						B4GamePlay.CardName = CardForDiscardPile.CardName
						B4GamePlay.Turn = TurnSequence(state.players[CardForDiscardPile.PlayerUserName].seatPosition, state, logger)
						state.players[B4GamePlay.Turn].takeCardfromPile = true
						state.players[B4GamePlay.Turn].putCardinPile = false
						B4GamePlay.TakeCardfromPile = state.players[B4GamePlay.Turn].takeCardfromPile
						B4GamePlay.PutCardinPile = state.players[B4GamePlay.Turn].putCardinPile
						B4GamePlay.FirstGodown = state.players[B4GamePlay.Turn].firstGodown
						Broadcast(B4GamePlay, int64(212), logger, nil, dispatcher)
					}

				}

			}

		case 409:

			CardForDiscardPile := &Before_GamePlay{}
			SendCardForDiscardPile := &Before_GamePlay{}
			//logger.Info("Error is: ", err)
			logger.Info("__________Meeellds ____ 44000999999999")

			if err := json.Unmarshal(data, &CardForDiscardPile); err != nil {
				logger.Info("Error is: ", err)
			} else {

				logger.Info("__________MEEEllllldddsss ____ 44400099", CardForDiscardPile.TableSequenceCards)

				if CardForDiscardPile.TableSequenceCards.TablecardsSequencelist != nil {

					if state.players[CardForDiscardPile.TableSequenceCards.UserName].firstGodown {
						if len(CardForDiscardPile.TableSequenceCards.TablecardsSequencelist) <= 0 {
							logger.Info("40000999 returning: ", err)

							return nil
						}
						SendCardForDiscardPile.TableSequenceCards.TablecardsSequencelist, SendCardForDiscardPile.TableSequenceCards.CardName = major(CardForDiscardPile.TableSequenceCards.TablecardsSequencelist, CardForDiscardPile.TableSequenceCards.CardIndex, CardForDiscardPile.TableSequenceCards.CardName)
						SendCardForDiscardPile.TableSequenceCards.AllowMeld = state.players[CardForDiscardPile.TableSequenceCards.UserName].firstGodown
						SendCardForDiscardPile.TableSequenceCards.CurrentListIndex = CardForDiscardPile.TableSequenceCards.CurrentListIndex
						SendCardForDiscardPile.TableSequenceCards.UserName = CardForDiscardPile.TableSequenceCards.UserName
						SendCardForDiscardPile.TableSequenceCards.CardIndex = CardForDiscardPile.TableSequenceCards.CardIndex
						SendCardForDiscardPile.TableSequenceCards.NotMatchedcardName = CardForDiscardPile.TableSequenceCards.CardName
						logger.Info("______Return____MEEEllllldddsss ____ 44400099", SendCardForDiscardPile.TableSequenceCards)
						Broadcast(SendCardForDiscardPile, int64(215), logger, nil, dispatcher)

						if len(SendCardForDiscardPile.TableSequenceCards.TablecardsSequencelist) > 0 {

							var set GoDownSet
							set.Set = append(set.Set, CardForDiscardPile.TableSequenceCards.CardName)
							state.players[CardForDiscardPile.TableSequenceCards.UserName].scoreCards = append(state.players[CardForDiscardPile.TableSequenceCards.UserName].scoreCards, set)

							if SendCardForDiscardPile.TableSequenceCards.CardName != "" {
								// just remove card from SendCardForDiscardPile.TableSequenceCards.CardName

								for _, setObj := range state.players[CardForDiscardPile.TableSequenceCards.UserName].scoreCards {

									if len(setObj.Set) == 1 {
										for cardIndex, card := range setObj.Set {
											if card == SendCardForDiscardPile.TableSequenceCards.CardName {
												// Remove it

												setObj.Set[cardIndex] = ""
											}
										}

									}

								}

							}

						}
						//	logger.Info("____Return______MEEEllllldddsss ____ 44400099", CardForDiscardPile.TableSequenceCards)

					} else {
						SendCardForDiscardPile.TableSequenceCards.TablecardsSequencelist = nil
						SendCardForDiscardPile.TableSequenceCards.CardName = ""
						SendCardForDiscardPile.TableSequenceCards.AllowMeld = state.players[CardForDiscardPile.TableSequenceCards.UserName].firstGodown
						SendCardForDiscardPile.TableSequenceCards.CurrentListIndex = CardForDiscardPile.TableSequenceCards.CurrentListIndex
						SendCardForDiscardPile.TableSequenceCards.UserName = CardForDiscardPile.TableSequenceCards.UserName
						SendCardForDiscardPile.TableSequenceCards.CardIndex = CardForDiscardPile.TableSequenceCards.CardIndex
						SendCardForDiscardPile.TableSequenceCards.NotMatchedcardName = CardForDiscardPile.TableSequenceCards.CardName

						Broadcast(SendCardForDiscardPile, int64(215), logger, nil, dispatcher)

					}

				}

			}
		case 410:
			Cht := &Before_GamePlay{}

			if err := json.Unmarshal(data, &Cht); err != nil {
				logger.Info("Error is: ", err)
			} else {
				logger.Info("Data is: ", Cht)
			}

			if Json_chat, err := json.Marshal(Cht); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(219, Json_chat, nil, nil, true)
			}

		}
	}

	return state
}

func melds_409(stateMain interface{}, HandCardList *HandCardLisdt, logger runtime.Logger, presence []runtime.Presence, dispatcher runtime.MatchDispatcher) {
	//state := stateMain.(*MatchState)

}

func sendGoDownCards(obj Before_GamePlay, stateMain interface{}, logger runtime.Logger) {

	// for _, setNum := range obj.GoDownData.GoDownCardSetNum {

	// 	// for _,  cardstate.players[obj.GoDownData.UserName].mainList {

	// 	// }

	// }

}

func TurnSequence(currentPosition int, stateMain interface{}, logger runtime.Logger) string {
	state := stateMain.(*MatchState)

	if (currentPosition + 1) < len(state.sittingArangement) {

		if state.sittingArangement[currentPosition+1] == state.dealer {
			//state.bidFlag = false
		}
		return state.sittingArangement[currentPosition+1]
	} else {

		if state.sittingArangement[0] == state.dealer {
			//state.bidFlag = false
		}
		return state.sittingArangement[0]
	}

}

func cardWithValue(stateMain interface{}, logger runtime.Logger) {

	state := stateMain.(*MatchState)
	for _, plyr := range state.players {

		state.players[plyr.userName].cardWithValue = make(map[string]*mapping)

	}

	state.cardWithValue = nil
	state.cardWithValue = make(map[string]*mapping)
	state.cardWithValue = map[string]*mapping{"club_A": {1.5, false}, "club_2": {2, false}, "club_3": {0.5, false}, "club_4": {0.5, false}, "club_5": {0.5, false}, "club_6": {0.5, false}, "club_7": {1, false},
		"club_8": {1, false}, "club_9": {1, false}, "club_10": {1, false}, "club_J": {1, false}, "club_Q": {1, false}, "club_K": {1, false}, "diamond_A": {1.5, false},
		"diamond": {2, false}, "diamond_3": {0.5, false}, "diamond_4": {0.5, false}, "diamond_5": {0.5, false}, "diamond_6": {0.5, false}, "diamond_7": {1, false},
		"diamond_8": {1, false}, "diamond_9": {1, false}, "diamond_10": {1, false}, "diamond_J": {1, false}, "diamond_Q": {1, false}, "diamond_K": {1, false},
		"heart_A": {1.5, false}, "heart_2": {2, false}, "heart_3": {0.5, false}, "heart_4": {0.5, false}, "heart_5": {0.5, false}, "heart_6": {0.5, false}, "heart_7": {1, false},
		"heart_8": {1, false}, "heart_9": {1, false}, "heart_10": {1, false}, "heart_J": {1, false}, "heart_Q": {1, false}, "heart_K": {1, false}, "spade_A": {1.5, false},
		"spade_2": {2, false}, "spade_3": {0.5, false}, "spade_4": {0.5, false}, "spade_5": {0.5, false}, "spade_6": {0.5, false}, "spade_7": {1, false}, "spade_8": {1, false},
		"spade_9": {1, false}, "spade_10": {1, false}, "spade_J": {1, false}, "spade_Q": {1, false}, "spade_K": {1, false}, "joker_A_3_1": {4, false}, "joker_A_3_2": {4, false}}

}

func afterSitting(stateMain interface{}, logger runtime.Logger, B4GamePlay *Before_GamePlay, dispatcher runtime.MatchDispatcher) {
	state := stateMain.(*MatchState)
	for _, card := range state.deck {
		//	state.deck[i] = card
		fmt.Println("shuffle0____", card)

	}
	//state.deck = nil
	state.deck = [106]string{"joker_A_3_1", "joker_A_3_0", "club_A_1", "club_2_1", "club_3_1", "club_4_1", "club_5_1", "club_6_1", "club_7_1", "club_8_1", "club_9_1", "club_10_1", "club_J_1", "club_Q_1", "club_K_1", "diamond_A_1", "diamond_2_1", "diamond_3_1", "diamond_4_1", "diamond_5_1", "diamond_6_1", "diamond_7_1", "diamond_8_1", "diamond_9_1", "diamond_10_1", "diamond_J_1", "diamond_Q_1", "diamond_K_1", "heart_A_1", "heart_2_1", "heart_3_1", "heart_4_1", "heart_5_1", "heart_6_1", "heart_7_1", "heart_8_1", "heart_9_1", "heart_10_1", "heart_J_1", "heart_Q_1", "heart_K_1", "spade_A_1", "spade_2_1", "spade_3_1", "spade_4_1", "spade_5_1", "spade_6_1", "spade_7_1", "spade_8_1", "spade_9_1", "spade_10_1", "spade_J_1", "spade_Q_1", "spade_K_1", "spade_K_1", "club_2_0", "club_3_0", "club_4_0", "club_5_0", "club_6_0", "club_7_0", "club_8_0", "club_9_0", "club_10_0", "club_J_0", "club_Q_0", "club_K_0", "diamond_A_0", "diamond_2_0", "diamond_3_0", "diamond_4_0", "diamond_5_0", "diamond_6_0", "diamond_7_0", "diamond_8_0", "diamond_9_0", "diamond_10_0", "diamond_J_0", "diamond_Q_0", "diamond_K_0", "heart_A_0", "heart_2_0", "heart_3_0", "heart_4_0", "heart_5_0", "heart_6_0", "heart_7_0", "heart_8_0", "heart_9_0", "heart_10_0", "heart_J_0", "heart_Q_0", "heart_K_0", "spade_A_0", "spade_2_0", "spade_3_0", "spade_4_0", "spade_5_0", "spade_6_0", "spade_7_0", "spade_8_0", "spade_9_0", "spade_10_0", "spade_J_0", "spade_Q_0", "spade_K_0"}

	shuffle(state.deck, state)
	cardDealing(state, logger)
	for _, plyr := range state.players {
		logger.Info("Sending Card  : ")

		//var handList = []string{"club_A_1", "diamond_A_1", "heart_A_1", "club_5_1", "club_A_1", "diamond_A_1", "heart_A_1", "spade_A_1", "club_5_1", "club_5_1", "club_A_1", "club_5_1", "heart_A_1", "club_A_1", "diamond_A_1", "club_5_1", "spade_A_1", "club_A_1" /*"diamond_A", "heart_A", "club_A", "diamond_A", "joker_A_3", "joker_A_3" , "club_A", "joker_A_3", "heart_A", "spade_A", "joker_A_3", "club_5", "club_A", "club_5", "heart_A", "club_A", "joker_A_3", "joker_A_3", "spade_A", "club_A", "diamond_A", "heart_A", "club_3", "diamond_3", "heart_3", "club_5", "club_3", "diamond_3", "heart_3", "spade_3", "club_5", "club_5", "club_3", "club_5", "heart_3", "club_3", "diamond_3", "club_5", "spade_3", "club_3", "diamond_3", "heart_3", "club_3", "diamond_3", "joker_A_3", "joker_A_3", "club_3", "joker_A_3", "heart_3", "spade_3", "joker_A_3", "club_5", "club_3", "club_5", "heart_3", "club_3", "joker_A_3", "joker_A_3", "spade_3", "club_3", "diamond_3", "heart_3"*/}
		//	var handList = []string{"spade_4_0", "diamond_Q_1", "spade_2_1", "club_5_0", "heart_J_1", "club_K_1", "heart_8_0", "spade_K_1", "heart_10_0", "diamond_7_0", "heart_10_1", "joker_A_3_0", "club_3_0", "diamond_3_0", "diamond_3_1", "spade_3_1", "spade_A_1", "diamond_2_0"}
		B4GamePlay.Cards = nil
		for _, Card := range state.players[plyr.userName].cards { //handList {

			B4GamePlay.Cards = append(B4GamePlay.Cards, Card)
			//logger.Info("card Player is: ", plyr)

		}
		//logger.Info("card Player is: ", plyr)
		logger.Info("JsoooOOOOoon Player cards are:______ ", B4GamePlay.Cards)

		B4GamePlay.PlayerUserName = plyr.userName
		if Json_PlayerCards, err := json.Marshal(B4GamePlay); err != nil {
			logger.Info("Error is: ", err)
		} else {

			logger.Info("Sending Card is: ")

			presence := []runtime.Presence{plyr.presence}
			dispatcher.BroadcastMessage(106, Json_PlayerCards, presence, nil, true)

			sendSets( /*handList */ state.players[plyr.userName].cards, plyr.userName, state, logger, presence, dispatcher)

			logger.Info("After Ck()()()  : ", state.players[plyr.userName].mainList)

		}
		B4GamePlay.Cards = nil
	}

	B4GamePlay.StockPile = append(B4GamePlay.StockPile, state.stockpile...)
	if Json_StockPile, err := json.Marshal(B4GamePlay); err != nil {
		logger.Info("Error is: ", err)
	} else {

		logger.Info("Stockpile Card is: ", B4GamePlay.StockPile)

		dispatcher.BroadcastMessage(208, Json_StockPile, nil, nil, true)

		B4GamePlay.DiscardPile = nil
		B4GamePlay.StockPile = nil

		logger.Info("state.discardpile Card is: ", state.discardpile)

		B4GamePlay.DiscardPile = append(B4GamePlay.DiscardPile, state.discardpile...)

		Broadcast(B4GamePlay, int64(209), logger, nil, dispatcher)

		B4GamePlay.DiscardPile = nil
		B4GamePlay.Turn = TurnSequence(state.players[state.firstTurn].seatPosition, state, logger)
		state.firstTurn = B4GamePlay.Turn
		B4GamePlay.TakeCardfromPile = state.players[B4GamePlay.Turn].takeCardfromPile
		B4GamePlay.PutCardinPile = state.players[B4GamePlay.Turn].putCardinPile
		B4GamePlay.FirstGodown = state.players[B4GamePlay.Turn].firstGodown
		Broadcast(B4GamePlay, int64(212), logger, nil, dispatcher)
	}

}
func scoreCalc_404(stateMain interface{}, HandCardList *HandCardLisdt, logger runtime.Logger, presence []runtime.Presence, dispatcher runtime.MatchDispatcher) bool {
	state := stateMain.(*MatchState)
	cardsPoint := map[string]float32{"joker_A_3_0": 4.0, "joker_A_3_1": 4.0, "club_2_0": 2.0, "club_3_0": 0.5, "club_4_0": 0.5, "club_5_0": 0.5,
		"club_6_0": 0.5, "club_7_0": 1.0, "club_8_0": 1.0, "club_9_0": 1.0, "club_10_0": 1.0, "club_J_0": 1.0, "club_Q_0": 1.0, "club_K_0": 1.0, "club_A_0": 1.5,
		"diamond_2_0": 2.0, "diamond_3_0": 0.5, "diamond_4_0": 0.5, "diamond_5_0": 0.5, "diamond_6_0": 0.5, "diamond_7_0": 1.0, "diamond_8_0": 1.0,
		"diamond_9_0": 1.0, "diamond_10_0": 1.0, "diamond_J_0": 1.0, "diamond_Q_0": 1.0, "diamond_K_0": 1.0, "diamond_A_0": 1.5, "heart_2_0": 2.0,
		"heart_3_0": 0.5, "heart_4_0": 0.5, "heart_5_0": 0.5, "heart_6_0": 0.5, "heart_7_0": 1.0, "heart_8_0": 1.0, "heart_9_0": 1.0, "heart_10_0": 1.0,
		"heart_J_0": 1.5, "heart_Q_0": 1.5, "heart_K_0": 1.5, "heart_A_0": 2.0, "spade_2_0": 2.0, "spade_3_0": 0.5, "spade_4_0": 0.5, "spade_5_0": 0.5,
		"spade_6_0": 0.5, "spade_7_0": 1.0, "spade_8_0": 1.0, "spade_9_0": 1.0, "spade_10_0": 1.0, "spade_J_0": 1.0, "spade_Q_0": 1.0, "spade_K_0": 1.0, "spade_A_0": 1.5,
		"club_2_1": 2.0, "club_3_1": 0.5, "club_4_1": 0.5, "club_5_1": 0.5,
		"club_6_1": 0.5, "club_7_1": 1.0, "club_8_1": 1.0, "club_9_1": 1.0, "club_10_1": 1.0, "club_J_1": 1.0, "club_Q_1": 1.0, "club_K_1": 1.0, "club_A_1": 1.5,
		"diamond_2_1": 2, "diamond_3_1": 0.5, "diamond_4_1": 0.5, "diamond_5_1": 0.5, "diamond_6_1": 0.5, "diamond_7_1": 1.0, "diamond_8_1": 1.0,
		"diamond_9_1": 1.0, "diamond_10_1": 1.0, "diamond_J_1": 1.0, "diamond_Q_1": 1.0, "diamond_K_1": 1.0, "diamond_A_1": 1.5, "heart_2_1": 2,
		"heart_3_1": 0.5, "heart_4_1": 0.5, "heart_5_1": 0.5, "heart_6_1": 0.5, "heart_7_1": 1.0, "heart_8_1": 1.0, "heart_9_1": 1.0, "heart_10_1": 1.0,
		"heart_J_1": 1.0, "heart_Q_1": 1.0, "heart_K_1": 1.0, "heart_A_1": 1.5, "spade_2_1": 2, "spade_3_1": 0.5, "spade_4_1": 0.5, "spade_5_1": 0.5,
		"spade_6_1": 0.5, "spade_7_1": 1.0, "spade_8_1": 1.0, "spade_9_1": 1.0, "spade_10_1": 1.0, "spade_J_1": 1.0, "spade_Q_1": 1.0, "spade_K_1": 1.5, "spade_A_1": 2, "": 0}

	state.players[HandCardList.HandCard.UserName].meldScore = 0
	var score float32
	//if len(HandCardList.HandCard.CardList) <= 0 {

	for _, setObj := range state.players[HandCardList.HandCard.UserName].scoreCards {
		for _, card := range setObj.Set {
			score = score + cardsPoint[card]
		}
	}

	state.players[HandCardList.HandCard.UserName].meldScore = score
	logger.Info("Caculated Score of melds is____: ", state.players[HandCardList.HandCard.UserName].meldScore)

	var handCardListScore float32
	for _, handCard := range HandCardList.HandCard.CardList {

		handCardListScore = handCardListScore + cardsPoint[handCard]

	}
	state.players[HandCardList.HandCard.UserName].handScore = handCardListScore
	logger.Info("Caculated Score of HandList is____: ", state.players[HandCardList.HandCard.UserName].handScore)

	state.players[HandCardList.HandCard.UserName].roundScore = state.players[HandCardList.HandCard.UserName].meldScore - state.players[HandCardList.HandCard.UserName].handScore
	logger.Info("Caculated Score of Round is____: ", state.players[HandCardList.HandCard.UserName].roundScore)

	if len(HandCardList.HandCard.CardList) == 0 && state.matchFirstMeld {
		state.players[HandCardList.HandCard.UserName].roundScore = state.players[HandCardList.HandCard.UserName].roundScore + 20
		logger.Info("Adding 20 into Round score is____: ", state.players[HandCardList.HandCard.UserName].roundScore)

	} else if !state.players[HandCardList.HandCard.UserName].firstGodown {
		state.players[HandCardList.HandCard.UserName].roundScore = state.players[HandCardList.HandCard.UserName].roundScore - 20
		logger.Info("Subtracting 20 into Round score is____: ", state.players[HandCardList.HandCard.UserName].roundScore)

	} else if state.fourtyBonus && HandCardList.HandCard.UserName == state.userFourtyBonus {
		state.players[HandCardList.HandCard.UserName].roundScore = state.players[HandCardList.HandCard.UserName].roundScore + 40
		logger.Info("Adding 40 into Round score is____: ", state.players[HandCardList.HandCard.UserName].roundScore)

	}

	state.players[HandCardList.HandCard.UserName].totalScore = state.players[HandCardList.HandCard.UserName].totalScore + state.players[HandCardList.HandCard.UserName].roundScore
	logger.Info("Caculated Total Score is____: ", state.players[HandCardList.HandCard.UserName].totalScore)
	state.roundScoreCounter++

	SendTotalScore := &Before_GamePlay{}
	logger.Info("state.roundScoreCounter is__", state.roundScoreCounter)
	if state.roundScoreCounter == 4 {
		state.roundScoreCounter = 0
		if state.players[state.hostUserName].totalScore > state.players[state.players[state.hostUserName].teamMate].totalScore {
			SendTotalScore.HostRoundScore1 = state.players[state.hostUserName].totalScore
		} else {
			SendTotalScore.HostRoundScore1 = state.players[state.players[state.hostUserName].teamMate].totalScore
		}

		if state.players[state.hostOpponent].totalScore > state.players[state.players[state.hostOpponent].teamMate].totalScore {
			SendTotalScore.HostOppoRoundScore2 = state.players[state.hostOpponent].totalScore
		} else {
			SendTotalScore.HostOppoRoundScore2 = state.players[state.players[state.hostOpponent].teamMate].totalScore
		}

		Broadcast(SendTotalScore, int64(217), logger, nil, dispatcher)

		if state.players[HandCardList.HandCard.UserName].totalScore == 230 {

			var highestScore float32
			highestScoreUser := ""
			var finalhostOpponentTeamScore float32
			var finalhostTeamScore float32

			for _, plyr := range state.players {

				if plyr.totalScore > highestScore {
					highestScore = plyr.totalScore
					highestScoreUser = plyr.userName
				}
			}

			logger.Info("Highest Score is____: ", highestScore)

			EndScore := &Before_GamePlay{}

			if state.hostUserName == highestScoreUser || state.players[state.hostUserName].teamMate == highestScoreUser {
				finalhostTeamScore = highestScore
				EndScore.Winner1 = highestScoreUser
				EndScore.Winner2 = state.players[highestScoreUser].teamMate

				hostOpponentTeam := state.players[state.hostOpponent].totalScore + state.players[state.players[state.hostOpponent].teamMate].totalScore
				finalhostOpponentTeamScore = hostOpponentTeam - state.players[state.players[highestScoreUser].teamMate].totalScore

				EndScore.Team1 = finalhostTeamScore
				EndScore.Team2 = finalhostOpponentTeamScore

			} else {
				finalhostOpponentTeamScore = highestScore
				EndScore.Winner1 = highestScoreUser
				EndScore.Winner2 = state.players[highestScoreUser].teamMate

				hostTeam := state.players[state.hostUserName].totalScore + state.players[state.players[state.hostUserName].teamMate].totalScore
				finalhostTeamScore = hostTeam - state.players[state.players[highestScoreUser].teamMate].totalScore

				EndScore.Team1 = finalhostOpponentTeamScore
				EndScore.Team2 = finalhostTeamScore
			}

			logger.Info("Highest Score is____: ", finalhostTeamScore)
			logger.Info("hostOpponentTeamScore Score is____: ", finalhostOpponentTeamScore)

			Broadcast(EndScore, int64(218), logger, nil, dispatcher)

			return false

		}

	}

	return true
}

func goDown_405(stateMain interface{}, GoDownResp *Before_GamePlay, logger runtime.Logger, presence []runtime.Presence, dispatcher runtime.MatchDispatcher, ctx context.Context, db *sql.DB, nk runtime.NakamaModule) {
	state := stateMain.(*MatchState)

	if GoDownResp.GoDownCardsData.HandEmptyOrNot && !state.matchFirstMeld {
		state.fourtyBonus = true
		state.userFourtyBonus = GoDownResp.GoDownCardsData.UserName
		state.players[GoDownResp.GoDownCardsData.UserName].firstGodown = true
	} else {
		state.matchFirstMeld = true
		state.players[GoDownResp.GoDownCardsData.UserName].firstGodown = true
	}
	submitLeadboardscore(state, logger, ctx, db, nk)

	Broadcast(GoDownResp, int64(210), logger, nil, dispatcher)

	state.matchFirstMeld = true
	if state.players[GoDownResp.GoDownCardsData.UserName].teamNum == "team1" {
		for _, set := range GoDownResp.GoDownCardsData.GoDownCardsMain {

			state.tableCardsTeam1 = append(state.tableCardsTeam1, set)
			state.players[GoDownResp.GoDownCardsData.UserName].scoreCards = append(state.players[GoDownResp.GoDownCardsData.UserName].scoreCards, set)

		}

		logger.Info("size of Team 1 table cards are____: ", len(state.tableCardsTeam1))

		logger.Info("Team 1 table cards are____: ", state.tableCardsTeam1)

	} else if state.players[GoDownResp.GoDownCardsData.UserName].teamNum == "team2" {
		for _, set := range GoDownResp.GoDownCardsData.GoDownCardsMain {

			state.tableCardsTeam2 = append(state.tableCardsTeam2, set)
			state.players[GoDownResp.GoDownCardsData.UserName].scoreCards = append(state.players[GoDownResp.GoDownCardsData.UserName].scoreCards, set)

		}

		logger.Info("size of Team 2 table cards are____: ", len(state.tableCardsTeam2))

		logger.Info("Team 2 table cards are____: ", state.tableCardsTeam2)
	}

}
func Broadcast(obj *Before_GamePlay, opcode int64, logger runtime.Logger, presence []runtime.Presence, dispatcher runtime.MatchDispatcher) {

	if Json, err := json.Marshal(obj); err != nil {
		logger.Info("Error is: ", err)
	} else {

		logger.Info("_____Successfully  Broadcasted____ ", opcode)

		dispatcher.BroadcastMessage(opcode, Json, presence, nil, true)
	}

}

func sendSets(handList []string, userName string, stateMain interface{}, logger runtime.Logger, presence []runtime.Presence, dispatcher runtime.MatchDispatcher) {
	state := stateMain.(*MatchState)

	logger.Info(" _____Before ThREEESSaAsoos is: ", handList)

	threesAsoos("_A", handList /*state.players[userName].cards,*/, userName, state, logger)
	logger.Info(" _____After ThREEESSaAsoos is: ", handList)

	ranking(handList /*state.players[plyr.userName].cards*/, userName, state, logger)
	logger.Info(" _____After Raaanking is: ", handList)

	reverseRanking(handList /*state.players[plyr.userName].cards*/, userName, state, logger)

	var Obj SetList

	for _, Set := range state.players[userName].mainList {

		var List SetStruct
		for _, Card := range Set.Set {

			List.Set = append(List.Set, Card)
		}
		Obj.Set = append(Obj.Set, List)
	}

	logger.Info(" ____Total sets are: ", state.players[userName].mainList)

	if Json_Sets, err := json.Marshal(Obj); err != nil {
		logger.Info("Error is: ", err)
	} else {

		dispatcher.BroadcastMessage(207, Json_Sets, presence, nil, true)
	}

	if len(state.players[userName].mainList) > 0 {

	} else {

		logger.Info("No set exist in hand: ")

	}

}

func arrangePositions(stateMain interface{}, logger runtime.Logger) {
	state := stateMain.(*MatchState)

	i := 0

	for _, plyr := range state.players {

		logger.Info(" i  __________ ", plyr.displayName)
		logger.Info(" =+++++++++++++++++", plyr.userName)

	}

	state.sittingArangement[0] = state.hostUserName
	state.players[state.hostUserName].seatPosition = 0
	state.players[state.hostUserName].teamNum = "team1"

	state.sittingArangement[2] = state.players[state.hostUserName].teamMate
	state.players[state.players[state.hostUserName].teamMate].seatPosition = 2
	state.players[state.players[state.hostUserName].teamMate].teamNum = "team1"

	state.teams["team1"] = &Team{0, state.sittingArangement[0], state.sittingArangement[2], 0, 0}

	for _, plyr := range state.players {

		if i == 0 {

			logger.Info(" i = 0 &&&&&&& ", plyr)

		}

		if plyr.userName != state.hostUserName && plyr.userName != state.players[state.hostUserName].teamMate {

			logger.Info(" if plyr.userName != state.hostUserName && plyr.userName &&&&&&& ", plyr)

			state.sittingArangement[1] = plyr.userName
			state.players[plyr.userName].seatPosition = 1
			state.players[plyr.userName].teamNum = "team2"

			state.sittingArangement[3] = state.players[plyr.userName].teamMate
			state.players[state.players[plyr.userName].teamMate].seatPosition = 3
			state.players[state.players[plyr.userName].teamMate].teamNum = "team2"

			state.teams["team2"] = &Team{0, state.sittingArangement[1], state.sittingArangement[3], 0, 0}

			state.hostOpponent = plyr.userName
			break

		}
		i++
	}

	logger.Info("___________ arrangePositions, : ", state.sittingArangement)

}

func dealerSelection(stateMain interface{}, logger runtime.Logger) {

	state := stateMain.(*MatchState)
	logger.Info("len(state.sittingArangement)   is: ", len(state.sittingArangement))
	// a := rand.Intn(max-min) + min
	max := 2 //len(state.sittingArangement)
	r := rand.Intn(max-0) + 0

	logger.Info("rand   is: ", r)

	state.dealer = state.sittingArangement[r]
	state.firstTurn = state.dealer
	logger.Info("Dealer is: ", state.dealer)
}

func cardDealing(stateMain interface{}, logger runtime.Logger) {

	logger.Info("cardDealing is: ")
	state := stateMain.(*MatchState)

	for _, plyr := range state.players {
		state.players[plyr.userName].cards = nil
	}

	i := 3
	//for _, card := range state.deck {
	for j := 0; j < 72; j++ {
		if i < len(state.sittingArangement) {
			state.players[state.sittingArangement[i]].cards = append(state.players[state.sittingArangement[i]].cards, state.deck[j])
			fmt.Println(state.deck[j])
			state.players[state.sittingArangement[i]].cardWithValue[state.deck[j]] = state.cardWithValue[state.deck[j]]

		}
		if i == 0 {
			i = 3
		} else {
			i--
		}
		// if len(state.players[state.sittingArangement[i]].cards) == 18 {
		// 	break
		// }

	}

	for _, plyr := range state.players {
		logger.Info("sSaroor Nalka, null card_______", state.players[plyr.userName].cards)
		//	state.players[plyr.userName].cards = nil
	}
	state.stockpile = nil
	state.stockpile = state.deck[72:]
	state.discardpile = nil
	state.discardpile = append(state.discardpile, state.stockpile[len(state.stockpile)-1])

	logger.Info("state.stockpile[len(state.stockpile)-1]_______", state.stockpile[len(state.stockpile)-1])

	logger.Info("len(state.stockpile) before________", len(state.stockpile))
	logger.Info("state.stockpile before________", state.stockpile)

	index := len(state.stockpile) - 1
	remove(state.stockpile, index)
	logger.Info("state.stockpile) After_______", state.stockpile)

	logger.Info("len(state.stockpile) After_______", len(state.stockpile))

}

func shuffle(deck [106]string, stateMain interface{}) [106]string {
	state := stateMain.(*MatchState)

	for _, card := range deck {
		//	state.deck[i] = card
		fmt.Println("shuffle1____", card)

	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	for i, card := range deck {
		state.deck[i] = card
		//fmt.Println("shuffle2____", card)

	}
	return deck
}

func (b *Banakil) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	if state.(*MatchState).debug {
		logger.Info("match terminate match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match terminate match_id %v grace seconds %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), graceSeconds)
	}

	return state
}

func reverseRanking(handList []string, userName string, stateMain interface{}, logger runtime.Logger) {

	state := stateMain.(*MatchState)

	var seqList []string
	var seqListInt []int

	wildCard := false
	jokerCounter := 0
	twoCounter := 0
	currentSuit := ""
	wildCardCounter := 0

	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 14, "club_4_0": 13, "club_5_0": 12,
		"club_6_0": 11, "club_7_0": 10, "club_8_0": 9, "club_9_0": 8, "club_10_0": 7, "club_J_0": 6, "club_Q_0": 5, "club_K_0": 4, "club_A_0": 3,
		"diamond_2_0": -1, "diamond_3_0": 14, "diamond_4_0": 13, "diamond_5_0": 12, "diamond_6_0": 11, "diamond_7_0": 10, "diamond_8_0": 9,
		"diamond_9_0": 8, "diamond_10_0": 7, "diamond_J_0": 6, "diamond_Q_0": 5, "diamond_K_0": 4, "diamond_A_0": 3, "heart_2_0": -1,
		"heart_3_0": 14, "heart_4_0": 13, "heart_5_0": 12, "heart_6_0": 11, "heart_7_0": 10, "heart_8_0": 9, "heart_9_0": 8, "heart_10_0": 7,
		"heart_J_0": 6, "heart_Q_0": 5, "heart_K_0": 4, "heart_A_0": 3, "spade_2_0": -1, "spade_3_0": 14, "spade_4_0": 13, "spade_5_0": 12,
		"spade_6_0": 11, "spade_7_0": 10, "spade_8_0": 9, "spade_9_0": 8, "spade_10_0": 7, "spade_J_0": 6, "spade_Q_0": 5, "spade_K_0": 4, "spade_A_0": 3,
		"club_2_1": -1, "club_3_1": 14, "club_4_1": 13, "club_5_1": 12,
		"club_6_1": 11, "club_7_1": 10, "club_8_1": 9, "club_9_1": 8, "club_10_1": 7, "club_J_1": 6, "club_Q_1": 5, "club_K_1": 4, "club_A_1": 3,
		"diamond_2_1": -1, "diamond_3_1": 14, "diamond_4_1": 13, "diamond_5_1": 12, "diamond_6_1": 11, "diamond_7_1": 10, "diamond_8_1": 9,
		"diamond_9_1": 8, "diamond_10_1": 7, "diamond_J_1": 6, "diamond_Q_1": 5, "diamond_K_1": 4, "diamond_A_1": 3, "heart_2_1": -1,
		"heart_3_1": 14, "heart_4_1": 13, "heart_5_1": 12, "heart_6_1": 11, "heart_7_1": 10, "heart_8_1": 9, "heart_9_1": 8, "heart_10_1": 7,
		"heart_J_1": 6, "heart_Q_1": 5, "heart_K_1": 4, "heart_A_1": 3, "spade_2_1": -1, "spade_3_1": 14, "spade_4_1": 13, "spade_5_1": 12,
		"spade_6_1": 11, "spade_7_1": 10, "spade_8_1": 9, "spade_9_1": 8, "spade_10_1": 7, "spade_J_1": 6, "spade_Q_1": 5, "spade_K_1": 4, "spade_A_1": 3}

	logger.Info("In ranking_______________________", handList)

	for i := 0; i < len(handList); i++ {

		if strings.Contains(handList[i], "joker") {
			wildCardCounter = jokerCounter

		} else if strings.Contains(handList[i], "_2") {
			wildCardCounter = twoCounter
		}

		if cardsValue[handList[i]] == -1 && wildCardCounter < 1 {

			wildCardSuit := ""

			if strings.Contains(handList[i], "heart") {
				wildCardSuit = "heart"
			} else if strings.Contains(handList[i], "spade") {
				wildCardSuit = "spade"
			} else if strings.Contains(handList[i], "diamond") {
				wildCardSuit = "diamond"
			} else if strings.Contains(handList[i], "club") {
				wildCardSuit = "club"
			} else if strings.Contains(handList[i], "joker") {
				wildCardSuit = "joker"
			}
			logger.Info("1 WildCard Suit: ", wildCardSuit)

			if len(seqList) > 0 && i < len(handList) {
				if cardsValue[seqList[len(seqList)-1]]+1 <= 14 /*&& !strings.Contains(handList[i+1], wildCardSuit)*/ {
					cardsValue[handList[i]] = cardsValue[seqList[len(seqList)-1]] + 1

				}
			} else if len(seqList) == 0 {
				if (i + 1) < len(handList) {

					wildCardSuit := ""

					if strings.Contains(handList[i], "heart") {
						wildCardSuit = "heart"
					} else if strings.Contains(handList[i], "spade") {
						wildCardSuit = "spade"
					} else if strings.Contains(handList[i], "diamond") {
						wildCardSuit = "diamond"
					} else if strings.Contains(handList[i], "club") {
						wildCardSuit = "club"
					} else if strings.Contains(handList[i], "joker") {
						wildCardSuit = "joker"
					}
					logger.Info("2 WildCard Suit: ", wildCardSuit)
					if cardsValue[handList[i+1]] == -1 /*&& !strings.Contains(handList[i+1], wildCardSuit)*/ {
						if i+2 < len(handList) {

							if cardsValue[handList[i+2]] >= 5 {
								cardsValue[handList[i]] = (cardsValue[handList[i+2]] - 2)

							} else if cardsValue[handList[i+2]] == 4 {
								cardsValue[handList[i]] = -1
								continue

							} else if cardsValue[handList[i+2]] == -1 {
								continue
							}
						} else {
							//verify later just setting due to out of bound error
							continue
						}

					} else if cardsValue[handList[i+1]] != 3 {

						cardsValue[handList[i]] = (cardsValue[handList[i+1]] - 1)

					} else if cardsValue[handList[i+1]] == 3 {

						fmt.Println("Next card is: ", handList[i+1])
						fmt.Println("Next value is: ", cardsValue[handList[i+1]])

						//cardsValue[handList[i]] = (cardsValue[handList[i+1]] - 1)
						cardsValue[handList[i]] = (cardsValue[handList[i+1]] + 1) // it was working ("	cardsValue[handList[i]] = (cardsValue[handList[i+1]] - 1)") except for the 0 index onn wild card

						fmt.Println("cardsValue[handList[i]] is: ", cardsValue[handList[i]])

					}
				}
			}

			wildCard = true
		} else if cardsValue[handList[i]] == -1 {

			if len(seqList) >= 3 {
				fmt.Println("seq is: ", seqList)
				i--

				/////////////For main list
				setObj := SetStruct{}

				for _, tempSaved := range seqListInt {

					setObj.Set = append(setObj.Set, tempSaved)
				}

				state.players[userName].mainList = append(state.players[userName].mainList, setObj)

				logger.Info("Set is: ", state.players[userName].mainList)
				logger.Info("")

				fmt.Println("")

				/////////////For main list

			} else if len(seqList) < 3 && wildCardCounter < 1 && (strings.Contains(handList[i], "joker") || strings.Contains(handList[i], "_2")) {

				seqList = nil
				seqListInt = nil
				wildCard = false
				jokerCounter = 0
				twoCounter = 0
				currentSuit = ""
				cardsValue["joker_A_3_1"] = -1
				cardsValue["joker_A_3_0"] = -1 // reset the values of _2 later..

				cardsValue["club_2_0"] = -1
				cardsValue["diamond_2_0"] = -1
				cardsValue["heart_2_0"] = -1
				cardsValue["spade_2_0"] = -1
				cardsValue["club_2_1"] = -1
				cardsValue["diamond_2_1"] = -1
				cardsValue["heart_2_1"] = -1
				cardsValue["spade_2_1"] = -1

				i = i - 2
			} else if len(seqList) < 3 && wildCardCounter <= 1 && (strings.Contains(handList[i], "joker") || strings.Contains(handList[i], "_2")) {

				seqList = nil
				wildCard = false
				jokerCounter = 0
				twoCounter = 0
				currentSuit = ""
				cardsValue["joker_A_3_1"] = -1
				cardsValue["joker_A_3_0"] = -1

				cardsValue["club_2_0"] = -1
				cardsValue["diamond_2_0"] = -1
				cardsValue["heart_2_0"] = -1
				cardsValue["spade_2_0"] = -1
				cardsValue["club_2_1"] = -1
				cardsValue["diamond_2_1"] = -1
				cardsValue["heart_2_1"] = -1
				cardsValue["spade_2_1"] = -1

				if strings.Contains(handList[i], "joker") {

					if i > 0 {
						if strings.Contains(handList[i-1], "joker") {
							i--
						} else {
							i = i - 2
						}
					}
				} else if strings.Contains(handList[i], "_2") {

					if i > 0 {
						if strings.Contains(handList[i-1], "_2") {
							i--
						} else {
							i = i - 2
						}
					}
				} else {
					i = i - 2
				}
				//i = i - 2
			} else {
				i--
			}
			seqList = nil
			seqListInt = nil
			wildCard = false
			jokerCounter = 0
			twoCounter = 0
			currentSuit = ""
			cardsValue["joker_A_3_1"] = -1
			cardsValue["joker_A_3_0"] = -1

			cardsValue["club_2_0"] = -1
			cardsValue["diamond_2_0"] = -1
			cardsValue["heart_2_0"] = -1
			cardsValue["spade_2_0"] = -1
			cardsValue["club_2_1"] = -1
			cardsValue["diamond_2_1"] = -1
			cardsValue["heart_2_1"] = -1
			cardsValue["spade_2_1"] = -1

			continue

		}

		if currentSuit == "" && !wildCard {
			if strings.Contains(handList[i], "heart") {
				currentSuit = "heart"
			} else if strings.Contains(handList[i], "spade") {
				currentSuit = "spade"
			} else if strings.Contains(handList[i], "diamond") {
				currentSuit = "diamond"
			} else if strings.Contains(handList[i], "club") {
				currentSuit = "club"
			}

		}

		if strings.Contains(handList[i], currentSuit) || strings.Contains(handList[i], "joker") || strings.Contains(handList[i], "_2") {
			if len(seqList) == 0 {
				seqList = append(seqList, handList[i])
				seqListInt = append(seqListInt, i)

				if wildCard {
					if strings.Contains(handList[i], "joker") {

						jokerCounter++

					} else if strings.Contains(handList[i], "_2") {
						twoCounter++
					}
					wildCard = false
				}
			} else if cardsValue[handList[i]] == cardsValue[seqList[len(seqList)-1]]+1 {

				seqList = append(seqList, handList[i])
				seqListInt = append(seqListInt, i)

				if wildCard {
					if strings.Contains(handList[i], "joker") {
						jokerCounter++

					} else if strings.Contains(handList[i], "_2") {
						twoCounter++
					}
					wildCard = false
				}
			} else {
				if len(seqList) >= 3 {
					fmt.Println("seq is: ", seqList)

					/////////////For main list
					setObj := SetStruct{}

					for _, tempSaved := range seqListInt {

						setObj.Set = append(setObj.Set, tempSaved)
					}

					state.players[userName].mainList = append(state.players[userName].mainList, setObj)

					logger.Info("Set is: ", state.players[userName].mainList)
					logger.Info("")

					fmt.Println("")

					/////////////For main list
				}

				if (strings.Contains(handList[i-1], "joker") || strings.Contains(handList[i-1], "_2")) && len(seqList) < 3 && cardsValue[handList[i]] != 3 {

					cardsValue[handList[i-1]] = -1
					i--
				}

				seqList = nil
				seqListInt = nil
				wildCard = false
				jokerCounter = 0
				twoCounter = 0
				currentSuit = ""
				cardsValue["joker_A_3_1"] = -1
				cardsValue["joker_A_3_0"] = -1

				i--

			}
		} else {
			if len(seqList) >= 3 {
				fmt.Println("seq is: ", seqList)
				/////////////For main list
				setObj := SetStruct{}

				for _, tempSaved := range seqListInt {

					setObj.Set = append(setObj.Set, tempSaved)
				}

				state.players[userName].mainList = append(state.players[userName].mainList, setObj)

				logger.Info("Set is: ", state.players[userName].mainList)
				logger.Info("")

				fmt.Println("")

				/////////////For main list
			}

			if (strings.Contains(handList[i-1], "joker") || strings.Contains(handList[i-1], "_2")) && len(seqList) < 3 && cardsValue[handList[i]] != 3 {

				cardsValue[handList[i-1]] = -1
				i--
			}
			seqList = nil
			seqListInt = nil
			wildCard = false
			jokerCounter = 0
			twoCounter = 0
			currentSuit = ""
			cardsValue["joker_A_3_1"] = -1
			cardsValue["joker_A_3_0"] = -1

			i--

		}

		if i == len(handList)-1 {

			if len(seqList) >= 3 {
				fmt.Println("_____________LAAAASSST________________seq is: ", seqList)
				/////////////For main list
				setObj := SetStruct{}

				for _, tempSaved := range seqListInt {

					setObj.Set = append(setObj.Set, tempSaved)
				}

				state.players[userName].mainList = append(state.players[userName].mainList, setObj)

				logger.Info("Set is: ", state.players[userName].mainList)
				logger.Info("")

				fmt.Println("")

				/////////////For main list

			}
		}

	}
}

func ranking(handList []string, userName string, stateMain interface{}, logger runtime.Logger) {

	state := stateMain.(*MatchState)

	var seqList []string
	var seqListInt []int

	wildCard := false
	jokerCounter := 0
	twoCounter := 0
	currentSuit := ""
	wildCardCounter := 0

	// var handList = []string{"spade_2_0", "club_4", "spade_2_1", "club_2_1", "club_2_0", "club_7", "club_8",
	// 	"club_9", "spade_2_0", "diamond_3", "spade_8", "joker_A_3_0",

	// 	"spade_10", "spade_J", "heart_J", "heart_J", "heart_J", "heart_Q",
	// 	"heart_K", "heart_A", "joker_A_3_1", "heart_K", "heart_A",

	// 	"joker_A_3_1", "heart_A", "heart_Q", "heart_K", "heart_A",
	// 	"spade_2_0", "club_A", "joker_A_3_1", "spade_A", "club_2_1", "club_3",

	// 	"club_4", "club_5", "club_6", "club_7", "club_8", "club_9", "club_10",
	// 	"heart_9", "heart_10", "heart_J", "heart_K", "joker_A_3_0",

	// 	"heart_A", "heart_K", "heart_Q", "heart_K", "heart_A", "heart_9",
	// 	"heart_10", "heart_J", "heart_Q", "heart_K", "joker_A_3_1", "heart_A",

	// 	"heart_9", "spade_2_0", "joker_A_3_0", "joker_A_3_1", "spade_2_1",
	// 	"heart_9", "heart_10", "heart_J", "heart_Q", "heart_K", "joker_A_3_1",
	// 	"heart_A", "heart_9", "heart_10", "heart_J",

	// 	"heart_Q", "heart_K", "joker_A_3_1", "heart_A", "heart_9",
	// 	"heart_10", "heart_J", "heart_Q", "heart_K", "joker_A_3_1",

	// 	"heart_A", "heart_9", "heart_10", "heart_J", "heart_Q",
	// 	"heart_K", "joker_A_3_1", "heart_A",
	// }

	cardsValue := map[string]int{"*": -2, "joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 3, "club_4_0": 4, "club_5_0": 5,
		"club_6_0": 6, "club_7_0": 7, "club_8_0": 8, "club_9_0": 9, "club_10_0": 10, "club_J_0": 11, "club_Q_0": 12, "club_K_0": 13, "club_A_0": 14,
		"diamond_2_0": -1, "diamond_3_0": 3, "diamond_4_0": 4, "diamond_5_0": 5, "diamond_6_0": 6, "diamond_7_0": 7, "diamond_8_0": 8,
		"diamond_9_0": 9, "diamond_10_0": 10, "diamond_J_0": 11, "diamond_Q_0": 12, "diamond_K_0": 13, "diamond_A_0": 14, "heart_2_0": -1,
		"heart_3_0": 3, "heart_4_0": 4, "heart_5_0": 5, "heart_6_0": 6, "heart_7_0": 7, "heart_8_0": 8, "heart_9_0": 9, "heart_10_0": 10,
		"heart_J_0": 11, "heart_Q_0": 12, "heart_K_0": 13, "heart_A_0": 14, "spade_2_0": -1, "spade_3_0": 3, "spade_4_0": 4, "spade_5_0": 5,
		"spade_6_0": 6, "spade_7_0": 7, "spade_8_0": 8, "spade_9_0": 9, "spade_10_0": 10, "spade_J_0": 11, "spade_Q_0": 12, "spade_K_0": 13, "spade_A_0": 14,
		"club_2_1": -1, "club_3_1": 3, "club_4_1": 4, "club_5_1": 5,
		"club_6_1": 6, "club_7_1": 7, "club_8_1": 8, "club_9_1": 9, "club_10_1": 10, "club_J_1": 11, "club_Q_1": 12, "club_K_1": 13, "club_A_1": 14,
		"diamond_2_1": -1, "diamond_3_1": 3, "diamond_4_1": 4, "diamond_5_1": 5, "diamond_6_1": 6, "diamond_7_1": 7, "diamond_8_1": 8,
		"diamond_9_1": 9, "diamond_10_1": 10, "diamond_J_1": 11, "diamond_Q_1": 12, "diamond_K_1": 13, "diamond_A_1": 14, "heart_2_1": -1,
		"heart_3_1": 3, "heart_4_1": 4, "heart_5_1": 5, "heart_6_1": 6, "heart_7_1": 7, "heart_8_1": 8, "heart_9_1": 9, "heart_10_1": 10,
		"heart_J_1": 11, "heart_Q_1": 12, "heart_K_1": 13, "heart_A_1": 14, "spade_2_1": -1, "spade_3_1": 3, "spade_4_1": 4, "spade_5_1": 5,
		"spade_6_1": 6, "spade_7_1": 7, "spade_8_1": 8, "spade_9_1": 9, "spade_10_1": 10, "spade_J_1": 11, "spade_Q_1": 12, "spade_K_1": 13, "spade_A_1": 14}

	logger.Info("In ranking_______________________", handList)

	for i := 0; i < len(handList); i++ {
		//logger.Info("CARD IS_______1______", handList[i])

		//logger.Info(" _____***()()()SS: ", handList)

		if strings.Contains(handList[i], "joker") {
			wildCardCounter = jokerCounter

		} else if strings.Contains(handList[i], "_2") {
			wildCardCounter = twoCounter
		}

		if cardsValue[handList[i]] == -1 && wildCardCounter < 1 {

			wildCardSuit := ""

			if strings.Contains(handList[i], "heart") {
				wildCardSuit = "heart"
			} else if strings.Contains(handList[i], "spade") {
				wildCardSuit = "spade"
			} else if strings.Contains(handList[i], "diamond") {
				wildCardSuit = "diamond"
			} else if strings.Contains(handList[i], "club") {
				wildCardSuit = "club"
			} else if strings.Contains(handList[i], "joker") {
				wildCardSuit = "joker"
			}
			logger.Info("3 WildCard Suit: ", wildCardSuit)

			if len(seqList) > 0 && i < len(handList) {
				if cardsValue[seqList[len(seqList)-1]]+1 <= 14 /*&& !strings.Contains(handList[i+1], wildCardSuit)*/ {
					cardsValue[handList[i]] = cardsValue[seqList[len(seqList)-1]] + 1

				}
			} else if len(seqList) == 0 {
				if (i + 1) < len(handList) {

					wildCardSuit := ""

					if strings.Contains(handList[i], "heart") {
						wildCardSuit = "heart"
					} else if strings.Contains(handList[i], "spade") {
						wildCardSuit = "spade"
					} else if strings.Contains(handList[i], "diamond") {
						wildCardSuit = "diamond"
					} else if strings.Contains(handList[i], "club") {
						wildCardSuit = "club"
					} else if strings.Contains(handList[i], "joker") {
						wildCardSuit = "joker"
					}
					logger.Info("4 WildCard Suit: ", wildCardSuit)
					if cardsValue[handList[i+1]] == -1 /*&& !strings.Contains(handList[i+1], wildCardSuit)*/ {
						if i+2 < len(handList) {

							if cardsValue[handList[i+2]] >= 5 {
								cardsValue[handList[i]] = (cardsValue[handList[i+2]] - 2)

							} else if cardsValue[handList[i+2]] == 4 {
								cardsValue[handList[i]] = -1
								continue

							} else if cardsValue[handList[i+2]] == -1 {
								continue
							}
						} else {
							//verify later just setting due to out of bound error
							continue
						}

					} else if cardsValue[handList[i+1]] != 3 {

						cardsValue[handList[i]] = (cardsValue[handList[i+1]] - 1)

					} else if cardsValue[handList[i+1]] == 3 {

						fmt.Println("Next card is: ", handList[i+1])
						fmt.Println("Next value is: ", cardsValue[handList[i+1]])

						cardsValue[handList[i]] = (cardsValue[handList[i+1]] - 1)
						fmt.Println("cardsValue[handList[i]] is: ", cardsValue[handList[i]])

					}
				}
			}

			wildCard = true
		} else if cardsValue[handList[i]] == -1 {

			if len(seqList) >= 3 {
				fmt.Println("seq is: ", seqList)
				i--

				/////////////For main list
				setObj := SetStruct{}

				for _, tempSaved := range seqListInt {

					setObj.Set = append(setObj.Set, tempSaved)
					handList[tempSaved] = "*"

				}

				state.players[userName].mainList = append(state.players[userName].mainList, setObj)

				logger.Info("Set is: ", state.players[userName].mainList)
				logger.Info("")

				fmt.Println("")

				/////////////For main list

			} else if len(seqList) < 3 && wildCardCounter < 1 && (strings.Contains(handList[i], "joker") || strings.Contains(handList[i], "_2")) {

				seqList = nil
				seqListInt = nil
				wildCard = false
				jokerCounter = 0
				twoCounter = 0
				currentSuit = ""
				cardsValue["joker_A_3_1"] = -1
				cardsValue["joker_A_3_0"] = -1 // reset the values of _2 later..

				cardsValue["club_2_0"] = -1
				cardsValue["diamond_2_0"] = -1
				cardsValue["heart_2_0"] = -1
				cardsValue["spade_2_0"] = -1
				cardsValue["club_2_1"] = -1
				cardsValue["diamond_2_1"] = -1
				cardsValue["heart_2_1"] = -1
				cardsValue["spade_2_1"] = -1

				i = i - 2
			} else if len(seqList) < 3 && wildCardCounter <= 1 && (strings.Contains(handList[i], "joker") || strings.Contains(handList[i], "_2")) {

				seqList = nil
				wildCard = false
				jokerCounter = 0
				twoCounter = 0
				currentSuit = ""
				cardsValue["joker_A_3_1"] = -1
				cardsValue["joker_A_3_0"] = -1

				cardsValue["club_2_0"] = -1
				cardsValue["diamond_2_0"] = -1
				cardsValue["heart_2_0"] = -1
				cardsValue["spade_2_0"] = -1
				cardsValue["club_2_1"] = -1
				cardsValue["diamond_2_1"] = -1
				cardsValue["heart_2_1"] = -1
				cardsValue["spade_2_1"] = -1

				if strings.Contains(handList[i], "joker") {

					if i > 0 {
						if strings.Contains(handList[i-1], "joker") {
							i--
						} else {
							i = i - 2
						}
					}
				} else if strings.Contains(handList[i], "_2") {

					if i > 0 {
						if strings.Contains(handList[i-1], "_2") {
							i--
						} else {
							i = i - 2
						}
					}
				} else {
					i = i - 2
				}
				//i = i - 2
			} else {
				i--
			}
			seqList = nil
			seqListInt = nil
			wildCard = false
			jokerCounter = 0
			twoCounter = 0
			currentSuit = ""
			cardsValue["joker_A_3_1"] = -1
			cardsValue["joker_A_3_0"] = -1

			cardsValue["club_2_0"] = -1
			cardsValue["diamond_2_0"] = -1
			cardsValue["heart_2_0"] = -1
			cardsValue["spade_2_0"] = -1
			cardsValue["club_2_1"] = -1
			cardsValue["diamond_2_1"] = -1
			cardsValue["heart_2_1"] = -1
			cardsValue["spade_2_1"] = -1

			continue

		}

		if currentSuit == "" && !wildCard {
			if strings.Contains(handList[i], "heart") {
				currentSuit = "heart"
			} else if strings.Contains(handList[i], "spade") {
				currentSuit = "spade"
			} else if strings.Contains(handList[i], "diamond") {
				currentSuit = "diamond"
			} else if strings.Contains(handList[i], "club") {
				currentSuit = "club"
			}

		}

		if strings.Contains(handList[i], currentSuit) || strings.Contains(handList[i], "joker") || strings.Contains(handList[i], "_2") {
			if len(seqList) == 0 {
				if (cardsValue[handList[i]]) != -2 {
					seqList = append(seqList, handList[i])
					seqListInt = append(seqListInt, i)

					if wildCard {
						if strings.Contains(handList[i], "joker") {

							jokerCounter++

						} else if strings.Contains(handList[i], "_2") {
							twoCounter++
						}
						wildCard = false
					}
				}
			} else if cardsValue[handList[i]] == cardsValue[seqList[len(seqList)-1]]+1 {

				seqList = append(seqList, handList[i])
				seqListInt = append(seqListInt, i)

				if wildCard {
					if strings.Contains(handList[i], "joker") {
						jokerCounter++

					} else if strings.Contains(handList[i], "_2") {
						twoCounter++
					}
					wildCard = false
				}
			} else {
				if len(seqList) >= 3 {
					fmt.Println("seq is: ", seqList)

					/////////////For main list
					setObj := SetStruct{}

					for _, tempSaved := range seqListInt {

						setObj.Set = append(setObj.Set, tempSaved)
						handList[tempSaved] = "*"

					}

					state.players[userName].mainList = append(state.players[userName].mainList, setObj)

					logger.Info("Set is: ", state.players[userName].mainList)
					logger.Info("")

					fmt.Println("")

					/////////////For main list
				}

				if (strings.Contains(handList[i-1], "joker") || strings.Contains(handList[i-1], "_2")) && len(seqList) < 3 && cardsValue[handList[i]] != 3 {

					cardsValue[handList[i-1]] = -1
					i--
				}

				seqList = nil
				seqListInt = nil
				wildCard = false
				jokerCounter = 0
				twoCounter = 0
				currentSuit = ""
				cardsValue["joker_A_3_1"] = -1
				cardsValue["joker_A_3_0"] = -1

				i--

			}
		} else {
			if len(seqList) >= 3 {
				fmt.Println("seq is: ", seqList)
				/////////////For main list
				setObj := SetStruct{}

				for _, tempSaved := range seqListInt {

					setObj.Set = append(setObj.Set, tempSaved)
					handList[tempSaved] = "*"

				}

				state.players[userName].mainList = append(state.players[userName].mainList, setObj)

				logger.Info("Set is: ", state.players[userName].mainList)
				logger.Info("")

				fmt.Println("")

				/////////////For main list
			}

			if (strings.Contains(handList[i-1], "joker") || strings.Contains(handList[i-1], "_2")) && len(seqList) < 3 && cardsValue[handList[i]] != 3 {

				cardsValue[handList[i-1]] = -1
				i--
			}
			seqList = nil
			seqListInt = nil
			wildCard = false
			jokerCounter = 0
			twoCounter = 0
			currentSuit = ""
			cardsValue["joker_A_3_1"] = -1
			cardsValue["joker_A_3_0"] = -1

			i--

		}

		if i == len(handList)-1 {

			if len(seqList) >= 3 {
				fmt.Println("_____________LAAAASSST________________seq is: ", seqList)
				/////////////For main list
				setObj := SetStruct{}

				for _, tempSaved := range seqListInt {

					setObj.Set = append(setObj.Set, tempSaved)
					handList[tempSaved] = "*"
				}

				state.players[userName].mainList = append(state.players[userName].mainList, setObj)

				logger.Info("Set is: ", state.players[userName].mainList)
				logger.Info("")

				fmt.Println("")

				/////////////For main list

			}
		}

	}
}

func threesAsoos(num string, handList []string, userName string, stateMain interface{}, logger runtime.Logger) {
	state := stateMain.(*MatchState)

	//cardMap := map[string]int{"club_A": 1, "club_2": 2, "club_3": 3, "club_4": 4, "club_5": 5, "club_6": 6, "club_7": 7, "club_8": 8, "club_9": 9, "club_10": 10, "club_J": 11, "club_Q": 12, "club_K": 13, "diamond_A": 14, "diamond_2": 15, "diamond_3": 16, "diamond_4": 17, "diamond_5": 18, "diamond_6": 19, "diamond_7": 20, "diamond_8": 21, "diamond_9": 22, "diamond_10": 23, "diamond_J": 24, "diamond_Q": 25, "diamond_K": 26, "heart_A": 27, "heart_2": 28, "heart_3": 29, "heart_4": 30, "heart_5": 31, "heart_6": 32, "heart_7": 33, "heart_8": 34, "heart_9": 35, "heart_10": 36, "heart_J": 37, "heart_Q": 38, "heart_K": 39, "spade_A": 40, "spade_2": 41, "spade_3": 42, "spade_4": 43, "spade_5": 44, "spade_6": 45, "spade_7": 46, "spade_8": 47, "spade_9": 48, "spade_10": 49, "spade_J": 50, "spade_Q": 51, "spade_K": 52, "joker": 53}

	cardMap := map[string]int{"club_A_1": 1, "club_2_1": 2, "club_3_1": 3, "club_4_1": 4, "club_5_1": 5, "club_6_1": 6, "club_7_1": 7, "club_8_1": 8, "club_9_1": 9, "club_10_1": 10,
		"club_J_1": 11, "club_Q_1": 12, "club_K_1": 13, "diamond_A_1": 14, "diamond_2_1": 15, "diamond_3_1": 16, "diamond_4_1": 17, "diamond_5_1": 18, "diamond_6_1": 19, "diamond_7_1": 20,
		"diamond_8_1": 21, "diamond_9_1": 22, "diamond_10_1": 23, "diamond_J_1": 24, "diamond_Q_1": 25, "diamond_K_1": 26, "heart_A_1": 27, "heart_2_1": 28, "heart_3_1": 29, "heart_4_1": 30,
		"heart_5_1": 31, "heart_6_1": 32, "heart_7_1": 33, "heart_8_1": 34, "heart_9_1": 35, "heart_10_1": 36, "heart_J_1": 37, "heart_Q_1": 38, "heart_K_1": 39, "spade_A_1": 40, "spade_2_1": 41,
		"spade_3_1": 42, "spade_4_1": 43, "spade_5_1": 44, "spade_6_1": 45, "spade_7_1": 46, "spade_8_1": 47, "spade_9_1": 48, "spade_10_1": 49, "spade_J_1": 50, "spade_Q_1": 51, "spade_K_1": 52,
		"joker_A_3_1": 53, "club_A_0": 2, "club_2_0": 2, "club_3_0": 3, "club_4_0": 4, "club_5_0": 5, "club_6_0": 6, "club_7_0": 7, "club_8_0": 8, "club_9_0": 9, "club_10_0": 10, "club_J_0": 11,
		"club_Q_0": 12, "club_K_0": 13, "diamond_A_0": 14, "diamond_2_0": 15, "diamond_3_0": 16, "diamond_4_0": 17, "diamond_5_0": 18, "diamond_6_0": 19, "diamond_7_0": 20, "diamond_8_0": 21,
		"diamond_9_0": 22, "diamond_10_0": 23, "diamond_J_0": 24, "diamond_Q_0": 25, "diamond_K_0": 26, "heart_A_0": 27, "heart_2_0": 28, "heart_3_0": 29, "heart_4_0": 30, "heart_5_0": 31,
		"heart_6_0": 32, "heart_7_0": 33, "heart_8_0": 34, "heart_9_0": 35, "heart_10_0": 36, "heart_J_0": 37,
		"heart_Q_0": 38, "heart_K_0": 39, "spade_A_0": 40, "spade_2_0": 41, "spade_3_0": 42, "spade_4_0": 43, "spade_5_0": 44, "spade_6_0": 45, "spade_7_0": 46, "spade_8_0": 47, "spade_9_0": 48,
		"spade_10_0": 49, "spade_J_0": 50, "spade_Q_0": 51, "spade_K_0": 52, "joker_A_3_0": 53}

	var threesList []int
	for i := 0; i < len(handList); i++ {

		//fmt.Println("CARD IS: ", handList[i])
		if (i + 2) < len(handList) {
			if strings.Contains(handList[i], "joker") {

				if strings.Contains(handList[i+1], "_A") {
					num = "_A"
				} else if strings.Contains(handList[i+1], "_3") {
					num = "_3"
				}

			} else if strings.Contains(handList[i], "_3") {
				num = "_3"
			} else if strings.Contains(handList[i], "_A") {
				num = "_A"
			}
		}

		if (i + 2) < len(handList) {
			if strings.Contains(handList[i], num) {
				if strings.Contains(handList[i+1], num) {

					if strings.Contains(handList[i+2], num) {

						//////  if all cards are according to pattern

						if cardMap[handList[i]] != cardMap[handList[i+1]] {
							if cardMap[handList[i+2]] != cardMap[handList[i]] && cardMap[handList[i+2]] != cardMap[handList[i+1]] {

								// Now save in the Threeslist
								threesList = append(threesList, i)
								//handList[i] = "*"
								threesList = append(threesList, i+1)
								//handList[i+1] = "*"
								threesList = append(threesList, i+2)
								//handList[i+2] = "*"

								removeCard := false
								// Now save in the Threeslist
								if (i + 3) < len(handList) {
									if strings.Contains(handList[i+3], num) {
										// checkCard := true
										// for _,savedCard := range threesList {}
										if cardMap[handList[i+3]] != cardMap[handList[i]] && cardMap[handList[i+3]] != cardMap[handList[i+1]] && cardMap[handList[i+3]] != cardMap[handList[i+2]] {
											threesList = append(threesList, i+3)

											removeCard = true

											//i = i + 3
										}

									} //else {
									// 	i = i + 3
									// }
								}

								if removeCard {
									handList[i] = "*"
									handList[i+1] = "*"
									handList[i+2] = "*"
									handList[i+3] = "*"
								} else {

									handList[i] = "*"
									handList[i+1] = "*"
									handList[i+2] = "*"
								}
								fmt.Println("Set is: ", threesList)
								// fmt.Println("Length is: ", len(threesList))
								// fmt.Println("I is: ", i)

								setObj := SetStruct{}

								for _, tempSaved := range threesList {

									setObj.Set = append(setObj.Set, tempSaved)
								}

								state.players[userName].mainList = append(state.players[userName].mainList, setObj)

								logger.Info("Set is: ", state.players[userName].mainList)
								logger.Info("")

								fmt.Println("")
								if len(threesList) == 3 {
									i = i + 2
								} else if len(threesList) == 4 {
									i = i + 3
								}
								threesList = nil

							}
						}

						//////  if all cards are according to pattern
						///////
					} else {
						//	i = i + 2
					}

				} else {
					//i = i + 1
				}
			}
		}

	}

}

func remove(slice []string, s int) {

	slice[s] = ""
	//return append(slice[:s], slice[s+1:]...)
}

// func check(userName string, mainList []SetStruct, stateMain interface{}, logger runtime.Logger, dispatcher runtime.MatchDispatcher) {
// 	state := stateMain.(*MatchState)
// 	crntPtrn := ""

// 	var handList = []string{"joker_1", "heart_A", "club_A", "spade_A", "spade_5", "club_A", "joker_2", "spade_A", "club_A", "spade_A", "diamond_A", "heart_A", "spade_A", "heart_A", "diamond_A", "spade_3", "heart_3", "diamond_3"}

// 	for _, card := range handList /* state.players[userName].cards*/ {
// 		//var handList = []string{"club_A", "spade_A", "heart_A", "diamond_A", "spade_5", "club_A", "spade_A", "heart_A", "club_6", "spade_A", "heart_A", "diamond_A"}

// 		thisCard := "*"
// 		if strings.Contains(card, "_A") || strings.Contains(card, "joker") {
// 			logger.Info(" card is : ", card)

// 			thisCard = "_A"
// 		} else if strings.Contains(card, "_3") || strings.Contains(card, "joker") {

// 			thisCard = "_3"
// 		} else {
// 			//crntPtrn = ""
// 		}

// 		if crntPtrn == "" {

// 			if strings.Contains(card, "_A") || strings.Contains(card, "joker") {
// 				//	fmt.Println(" card is : ", card)

// 				crntPtrn = "_A"
// 			} else if strings.Contains(card, "_3") || strings.Contains(card, "joker") {

// 				crntPtrn = "_3"
// 			} else {

// 			}

// 		}

// 		if crntPtrn == thisCard {
// 			//ftn for _A
// 			_A(card, userName, state.players[userName].mainList, state, logger)
// 		} else {

// 			state.asoos = false
// 			if state.threesCounter == 3 {

// 				setObj := SetStruct{}

// 				for _, tempSaved := range state.tempThrees {

// 					setObj.Set = append(setObj.Set, tempSaved)
// 				}

// 				if len(setObj.Set) == 3 {
// 					fmt.Println("|||_checkedCard is : ", setObj)
// 					state.players[userName].mainList = append(state.players[userName].mainList, setObj)
// 					//state.threesCounter = 0
// 					//state.tempThrees = nil
// 					state.asoos = true
// 				}

// 			}
// 			state.tempThrees = nil
// 			state.threesCounter = 0
// 			crntPtrn = ""
// 		}
// 	}

// 	///output

// }

// func _A(card string, userName string, mainList []SetStruct, stateMain interface{}, logger runtime.Logger) {
// 	state := stateMain.(*MatchState)
// 	var sign string
// 	checkedCard := true
// 	logger.Info(" card is : ", card)

// 	if strings.Contains(card, "club") {
// 		sign = "club"
// 	} else if strings.Contains(card, "diamond") {
// 		sign = "diamond"
// 	} else if strings.Contains(card, "heart") {
// 		sign = "heart"
// 	} else if strings.Contains(card, "spade") {
// 		sign = "spade"
// 	} else if strings.Contains(card, "joker") {
// 		sign = "joker"
// 	}

// 	if len(state.players[userName].mainList) > 0 && (len(state.players[userName].mainList[len(state.players[userName].mainList)-1].Set)) == 3 && state.asoos {

// 		logger.Info(" AsoosAsoosAsoos : ")

// 		for _, savedCard := range state.players[userName].mainList[len(state.players[userName].mainList)-1].Set {

// 			if strings.Contains(savedCard, sign) {
// 				fmt.Println(" savedCard : ", savedCard)

// 				checkedCard = false

// 				state.threesCounter = 0
// 				state.tempThrees = nil

// 				state.tempThrees = append(state.tempThrees, card)
// 				state.threesCounter++
// 				break
// 			}
// 		}

// 		if checkedCard {

// 			state.players[userName].mainList[len(state.players[userName].mainList)-1].Set = append(state.players[userName].mainList[len(state.players[userName].mainList)-1].Set, card)

// 		}

// 		state.asoos = false

// 	} else {

// 		for _, tempSaved := range state.tempThrees {

// 			if tempSaved == card {
// 				checkedCard = false
// 				break
// 			}
// 		}

// 		if checkedCard {
// 			state.tempThrees = append(state.tempThrees, card)
// 			state.threesCounter++
// 		} else {
// 			state.threesCounter = 0
// 			state.tempThrees = nil

// 			// to restart pattern from current
// 			state.tempThrees = append(state.tempThrees, card)
// 			state.threesCounter++
// 		}

// 		if state.threesCounter == 3 {

// 			setObj := SetStruct{}

// 			for _, tempSaved := range state.tempThrees {

// 				setObj.Set = append(setObj.Set, tempSaved)
// 			}

// 			if len(setObj.Set) == 3 {
// 				logger.Info("********_checkedCard is : ", setObj)
// 				state.players[userName].mainList = append(state.players[userName].mainList, setObj)
// 				state.threesCounter = 0
// 				state.tempThrees = nil
// 				state.asoos = true
// 			}
// 		}

// 	}

// }

func major(godown []string, CardIndex int, checkedcard string) ([]string, string) {

	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 3, "club_4_0": 4, "club_5_0": 5,
		"club_6_0": 6, "club_7_0": 7, "club_8_0": 8, "club_9_0": 9, "club_10_0": 10, "club_J_0": 11, "club_Q_0": 12, "club_K_0": 13, "club_A_0": 14,
		"diamond_2_0": -1, "diamond_3_0": 3, "diamond_4_0": 4, "diamond_5_0": 5, "diamond_6_0": 6, "diamond_7_0": 7, "diamond_8_0": 8,
		"diamond_9_0": 9, "diamond_10_0": 10, "diamond_J_0": 11, "diamond_Q_0": 12, "diamond_K_0": 13, "diamond_A_0": 14, "heart_2_0": -1,
		"heart_3_0": 3, "heart_4_0": 4, "heart_5_0": 5, "heart_6_0": 6, "heart_7_0": 7, "heart_8_0": 8, "heart_9_0": 9, "heart_10_0": 10,
		"heart_J_0": 11, "heart_Q_0": 12, "heart_K_0": 13, "heart_A_0": 14, "spade_2_0": -1, "spade_3_0": 3, "spade_4_0": 4, "spade_5_0": 5,
		"spade_6_0": 6, "spade_7_0": 7, "spade_8_0": 8, "spade_9_0": 9, "spade_10_0": 10, "spade_J_0": 11, "spade_Q_0": 12, "spade_K_0": 13, "spade_A_0": 14,
		"club_2_1": -1, "club_3_1": 3, "club_4_1": 4, "club_5_1": 5,
		"club_6_1": 6, "club_7_1": 7, "club_8_1": 8, "club_9_1": 9, "club_10_1": 10, "club_J_1": 11, "club_Q_1": 12, "club_K_1": 13, "club_A_1": 14,
		"diamond_2_1": -1, "diamond_3_1": 3, "diamond_4_1": 4, "diamond_5_1": 5, "diamond_6_1": 6, "diamond_7_1": 7, "diamond_8_1": 8,
		"diamond_9_1": 9, "diamond_10_1": 10, "diamond_J_1": 11, "diamond_Q_1": 12, "diamond_K_1": 13, "diamond_A_1": 14, "heart_2_1": -1,
		"heart_3_1": 3, "heart_4_1": 4, "heart_5_1": 5, "heart_6_1": 6, "heart_7_1": 7, "heart_8_1": 8, "heart_9_1": 9, "heart_10_1": 10,
		"heart_J_1": 11, "heart_Q_1": 12, "heart_K_1": 13, "heart_A_1": 14, "spade_2_1": -1, "spade_3_1": 3, "spade_4_1": 4, "spade_5_1": 5,
		"spade_6_1": 6, "spade_7_1": 7, "spade_8_1": 8, "spade_9_1": 9, "spade_10_1": 10, "spade_J_1": 11, "spade_Q_1": 12, "spade_K_1": 13, "spade_A_1": 14}

	var temp string
	var temp1 string
	var card string

	for i := 0; i < len(godown); i++ {
		if cardsValue[godown[i]] != -1 {
			temp = godown[i]
			for j := i + 1; j < len(godown); j++ {
				if cardsValue[godown[j]] != -1 {
					temp1 = godown[j]
					break
				}
			}
			break
		}

	}
	if cardsValue[temp] < cardsValue[temp1] && cardsValue[checkedcard] != -1 {
		godown, card = godownMelds(godown, checkedcard)
		fmt.Println("List is: ", godown)
		fmt.Println("Card is: ", card)
		fmt.Println(" Melds(godown, checkedcard)")
	} else if cardsValue[temp] > cardsValue[temp1] && cardsValue[checkedcard] != -1 {
		godown, card = reverseMelds(godown, checkedcard)
		fmt.Println("List is: ", godown)
		fmt.Println("Card is: ", card)
		fmt.Println("reverseMelds(godown, checkedcard)")
	} else if cardsValue[temp] < cardsValue[temp1] && cardsValue[checkedcard] == -1 {
		fmt.Println("wild cardddddddddd22222211111111card")
		godown, card := wildCardEntry(godown, CardIndex, checkedcard)
		fmt.Println("List is: ", godown)
		fmt.Println("Card is: ", card)
	} else if cardsValue[temp] > cardsValue[temp1] && cardsValue[checkedcard] == -1 {
		fmt.Println("wild cardddddddddd22222222211111111card")
		godown, card := reversewildCardEntry(godown, CardIndex, checkedcard)
		fmt.Println("List is: ", godown)
		fmt.Println("Card is: ", card)

	} else {
		godown, card = threeAsoos(godown, checkedcard)
		fmt.Println("List is: ", godown)
		fmt.Println("Card is: ", card)
	}
	return godown, card

}

func godownMelds(godown []string, checkedcard string) ([]string, string) {
	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 3, "club_4_0": 4, "club_5_0": 5,
		"club_6_0": 6, "club_7_0": 7, "club_8_0": 8, "club_9_0": 9, "club_10_0": 10, "club_J_0": 11, "club_Q_0": 12, "club_K_0": 13, "club_A_0": 14,
		"diamond_2_0": -1, "diamond_3_0": 3, "diamond_4_0": 4, "diamond_5_0": 5, "diamond_6_0": 6, "diamond_7_0": 7, "diamond_8_0": 8,
		"diamond_9_0": 9, "diamond_10_0": 10, "diamond_J_0": 11, "diamond_Q_0": 12, "diamond_K_0": 13, "diamond_A_0": 14, "heart_2_0": -1,
		"heart_3_0": 3, "heart_4_0": 4, "heart_5_0": 5, "heart_6_0": 6, "heart_7_0": 7, "heart_8_0": 8, "heart_9_0": 9, "heart_10_0": 10,
		"heart_J_0": 11, "heart_Q_0": 12, "heart_K_0": 13, "heart_A_0": 14, "spade_2_0": -1, "spade_3_0": 3, "spade_4_0": 4, "spade_5_0": 5,
		"spade_6_0": 6, "spade_7_0": 7, "spade_8_0": 8, "spade_9_0": 9, "spade_10_0": 10, "spade_J_0": 11, "spade_Q_0": 12, "spade_K_0": 13, "spade_A_0": 14,
		"club_2_1": -1, "club_3_1": 3, "club_4_1": 4, "club_5_1": 5,
		"club_6_1": 6, "club_7_1": 7, "club_8_1": 8, "club_9_1": 9, "club_10_1": 10, "club_J_1": 11, "club_Q_1": 12, "club_K_1": 13, "club_A_1": 14,
		"diamond_2_1": -1, "diamond_3_1": 3, "diamond_4_1": 4, "diamond_5_1": 5, "diamond_6_1": 6, "diamond_7_1": 7, "diamond_8_1": 8,
		"diamond_9_1": 9, "diamond_10_1": 10, "diamond_J_1": 11, "diamond_Q_1": 12, "diamond_K_1": 13, "diamond_A_1": 14, "heart_2_1": -1,
		"heart_3_1": 3, "heart_4_1": 4, "heart_5_1": 5, "heart_6_1": 6, "heart_7_1": 7, "heart_8_1": 8, "heart_9_1": 9, "heart_10_1": 10,
		"heart_J_1": 11, "heart_Q_1": 12, "heart_K_1": 13, "heart_A_1": 14, "spade_2_1": -1, "spade_3_1": 3, "spade_4_1": 4, "spade_5_1": 5,
		"spade_6_1": 6, "spade_7_1": 7, "spade_8_1": 8, "spade_9_1": 9, "spade_10_1": 10, "spade_J_1": 11, "spade_Q_1": 12, "spade_K_1": 13, "spade_A_1": 14}
	wildCardSuit := ""
	var flag bool
	flag = false
	//var godown1 = []string{}
	card := ""
	for i := 0; i < len(godown); i++ {

		if !strings.Contains(godown[i], "joker") && !strings.Contains(godown[i], "_2") {

			if strings.Contains(godown[i], "heart") {
				wildCardSuit = "heart"
			} else if strings.Contains(godown[i], "spade") {
				wildCardSuit = "spade"
			} else if strings.Contains(godown[i], "diamond") {
				wildCardSuit = "diamond"
			} else if strings.Contains(godown[i], "club") {
				wildCardSuit = "club"
			}

		}
		// fmt.Println(wildCardSuit)
	}
	for i := 0; i < len(godown); i++ {
		if strings.Contains(checkedcard, wildCardSuit) {
			if strings.Contains(godown[i], "joker") {
				if i == 0 {
					if cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] {
						card = godown[i]
						// handList = append(handList, godown[i])
						godown[i] = checkedcard
						flag = true

					} else if cardsValue[godown[i+1]] == -1 {
						if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] {
							card = godown[i]
							godown[i] = checkedcard
							// handList = append(handList, godown[i])
							flag = true

						} else if cardsValue[godown[i+2]]-3 == cardsValue[checkedcard] {
							var temp string

							godown = append(godown, checkedcard)
							flag = true
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}

						} else if cardsValue[godown[i+2]]-1 == cardsValue[checkedcard] {
							fmt.Println("hfgfdhffhhjddfdddhdhfdfdfd")
							var temp = godown[i+1]
							godown[i+1] = checkedcard
							flag = true
							godown = append(godown, temp)
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}

						}

					} else if cardsValue[godown[i+1]]-2 == cardsValue[checkedcard] {
						flag = true
						var temp string
						godown = append(godown, checkedcard)
						for swp := 0; swp < len(godown)-1; swp++ {
							temp = godown[swp]
							godown[swp] = godown[len(godown)-1]
							godown[len(godown)-1] = temp

						}
					}

				} else if i <= len(godown)-1 {
					fmt.Println("card name", cardsValue[godown[i]])
					fmt.Println("card name", cardsValue[godown[i-1]]+1)
					fmt.Println(i)
					if cardsValue[godown[i-1]]+1 == cardsValue[checkedcard] {
						card = godown[i]
						// handList = append(handList, godown[i])
						godown[i] = checkedcard
						flag = true

					} else if i == len(godown)-1 && cardsValue[godown[i-1]] == -1 {
						fmt.Println("12odhdzzzzzzzzzzzzzzzzzzzzzzhsc678")
						if cardsValue[godown[i-2]]+3 == cardsValue[checkedcard] {
							godown = append(godown, checkedcard)
							flag = true

						} else if cardsValue[godown[i-1]]+2 == cardsValue[checkedcard] {

							godown = append(godown, checkedcard)
							flag = true

						} else if cardsValue[godown[i-2]]+2 == cardsValue[checkedcard] {
							card = godown[i]
							// handList = append(handList, godown[i])
							godown[i] = checkedcard
							flag = true

						}
					} else if cardsValue[godown[i-1]] == -1 && i == 1 {
						if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] {
							card = godown[i]
							// handList = append(handList, godown[i])
							godown[i] = checkedcard
							flag = true

						}
					} else if i == len(godown)-1 {
						if cardsValue[godown[i-1]]+2 == cardsValue[checkedcard] {

							godown = append(godown, checkedcard)
							flag = true

						}

					}

				}
			} else if strings.Contains(godown[i], "2_") {
				if i == 0 {

					if cardsValue[godown[i+1]]-1 == 3 {
						if cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] {
							if strings.Contains(checkedcard, wildCardSuit) {
								// handList = append(handList, godown[i])
								godown = append(godown, godown[i])
								godown[i] = checkedcard
								flag = true
							}

						}
					} else if cardsValue[godown[i+1]] == -1 {
						if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-2 != 3 {

							var temp string
							temp = godown[i]
							godown[i] = checkedcard
							godown = append(godown, temp)
							flag = true
							// var temp1 string
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}
						} else if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-2 == 3 {
							godown = append(godown, godown[i])
							godown[i] = checkedcard
							flag = true

						} else if cardsValue[godown[i+2]]-3 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-3 == 3 {
							godown = append(godown, checkedcard)
							temp := godown[i]
							godown[i] = checkedcard
							godown[len(godown)-1] = temp
							flag = true

						} else if cardsValue[godown[i+2]]-3 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-3 != 3 {

							var temp string
							temp = godown[i]
							godown[i] = checkedcard
							godown = append(godown, temp)
							flag = true
							// var temp1 string
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}
						}

					}
				} else if i <= len(godown)-1 {
					if i != len(godown)-1 && cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] || cardsValue[godown[i-1]]+1 == cardsValue[checkedcard] {

						if strings.Contains(godown[0], "joker_A_3_0") || cardsValue[godown[0]] == 3 {
							if cardsValue[godown[1]]-1 == 3 {
								godown = append(godown, godown[i])
								// godown = append(godown, checkedcard)
								godown[i] = checkedcard
								flag = true
							}

						} else {

							temp := godown[i]
							var temp1 string
							godown[i] = checkedcard
							godown = append(godown, temp)
							for swp := 0; swp < len(godown)-1; swp++ {
								temp1 = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp1
							}
							godown[0] = temp
							flag = true
						}
						// godown[i] = checkedcard
						// godown1 = godown
						// godown[0] = godown[i]
						// godown[i] = checkedcard
						// for j := 1; j == len(godown); j++ {
						// 	godown[j] = godown1[j-1]
						// }
						// godown[i] = checkedcard

					} else if i == len(godown)-1 {
						//fmt.Println("hgfhfhgghggjgfhjggfh")
						if cardsValue[godown[i-1]] == -1 {
							fmt.Println("hgfhfhgghggjgfhjggwxxxxxxxxxxxxxxxxxxxxfwhwkjkhfkhhkfh")
							if cardsValue[godown[i-2]]+3 == cardsValue[checkedcard] {
								fmt.Println("hgfhfhgghggjxxxxxxxxxxxxxxxxxxxxxxxxxgfhjggfh")
								godown = append(godown, checkedcard)
								flag = true
							} else if cardsValue[godown[i-2]]+2 == cardsValue[checkedcard] {
								//fmt.Println("hgfhfhgghggjgfhjggfhewfefwffwfweewffwfeww")
								temp2 := godown[i]
								godown[i] = checkedcard
								godown = append(godown, temp2)
								flag = true
								if cardsValue[godown[0]] != 3 {
									for swp := 0; swp < len(godown)-1; swp++ {
										temp2 = godown[swp]
										godown[swp] = godown[len(godown)-1]
										godown[len(godown)-1] = temp2
									}

								}

							}

						} else {
							fmt.Println("1234567890oiuhgvcxzsdfrtghjk")
							if cardsValue[godown[i-1]]+2 == cardsValue[checkedcard] {
								// fmt.Println("741258/5369")
								// temp3 := godown[i]
								// godown[i] = checkedcard
								godown = append(godown, checkedcard)
								flag = true
								// flag = true
								// if cardsValue[godown[0]] != 3 {
								// 	for swp := 0; swp < len(godown)-1; swp++ {
								// 		temp3 = godown[swp]
								// 		godown[swp] = godown[len(godown)-1]
								// 		godown[len(godown)-1] = temp3
								// 	}

								// }

							}
						}
					}

				}
			} else if cardsValue[godown[len(godown)-1]]+1 == cardsValue[checkedcard] && i == len(godown)-1 {
				flag = true
				godown = append(godown, checkedcard)

			} else if cardsValue[godown[0]]-1 == cardsValue[checkedcard] && i == 0 {
				flag = true
				godown = append(godown, checkedcard)
				var temp1 string
				for swp := 0; swp < len(godown)-1; swp++ {
					temp1 = godown[swp]
					godown[swp] = godown[len(godown)-1]
					godown[len(godown)-1] = temp1
				}

			}
			if flag == true {
				break
			}
		}
	}
	if flag == false {
		fmt.Println("SORRY YOU CANNOT ADJUST THIS CARD  IN SEQUENCE")
		godown = nil
	}
	// log.Println(godown)
	// log.Println(handList)
	return godown, card

}

func reverseMelds(godown []string, checkedcard string) ([]string, string) {
	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 14, "club_4_0": 13, "club_5_0": 12,
		"club_6_0": 11, "club_7_0": 10, "club_8_0": 9, "club_9_0": 8, "club_10_0": 7, "club_J_0": 6, "club_Q_0": 5, "club_K_0": 4, "club_A_0": 3,
		"diamond_2_0": -1, "diamond_3_0": 14, "diamond_4_0": 13, "diamond_5_0": 12, "diamond_6_0": 11, "diamond_7_0": 10, "diamond_8_0": 9,
		"diamond_9_0": 8, "diamond_10_0": 7, "diamond_J_0": 6, "diamond_Q_0": 5, "diamond_K_0": 4, "diamond_A_0": 3, "heart_2_0": -1,
		"heart_3_0": 14, "heart_4_0": 13, "heart_5_0": 12, "heart_6_0": 11, "heart_7_0": 10, "heart_8_0": 9, "heart_9_0": 8, "heart_10_0": 7,
		"heart_J_0": 6, "heart_Q_0": 5, "heart_K_0": 4, "heart_A_0": 3, "spade_2_0": -1, "spade_3_0": 14, "spade_4_0": 13, "spade_5_0": 12,
		"spade_6_0": 11, "spade_7_0": 10, "spade_8_0": 9, "spade_9_0": 8, "spade_10_0": 7, "spade_J_0": 6, "spade_Q_0": 5, "spade_K_0": 4, "spade_A_0": 3,
		"club_2_1": -1, "club_3_1": 14, "club_4_1": 13, "club_5_1": 12,
		"club_6_1": 11, "club_7_1": 10, "club_8_1": 9, "club_9_1": 8, "club_10_1": 7, "club_J_1": 6, "club_Q_1": 5, "club_K_1": 4, "club_A_1": 3,
		"diamond_2_1": -1, "diamond_3_1": 14, "diamond_4_1": 13, "diamond_5_1": 12, "diamond_6_1": 11, "diamond_7_1": 10, "diamond_8_1": 9,
		"diamond_9_1": 8, "diamond_10_1": 7, "diamond_J_1": 6, "diamond_Q_1": 5, "diamond_K_1": 4, "diamond_A_1": 3, "heart_2_1": -1,
		"heart_3_1": 14, "heart_4_1": 13, "heart_5_1": 12, "heart_6_1": 11, "heart_7_1": 10, "heart_8_1": 9, "heart_9_1": 8, "heart_10_1": 7,
		"heart_J_1": 6, "heart_Q_1": 5, "heart_K_1": 4, "heart_A_1": 3, "spade_2_1": -1, "spade_3_1": 14, "spade_4_1": 13, "spade_5_1": 12,
		"spade_6_1": 11, "spade_7_1": 10, "spade_8_1": 9, "spade_9_1": 8, "spade_10_1": 7, "spade_J_1": 6, "spade_Q_1": 5, "spade_K_1": 4, "spade_A_1": 3}

	wildCardSuit := ""
	var flag bool
	flag = false
	card := ""
	//var godown1 = []string{}
	for i := 0; i < len(godown); i++ {
		if !strings.Contains(godown[i], "joker") && !strings.Contains(godown[i], "_2") {
			if strings.Contains(godown[i], "heart") {
				wildCardSuit = "heart"
			} else if strings.Contains(godown[i], "spade") {
				wildCardSuit = "spade"
			} else if strings.Contains(godown[i], "diamond") {
				wildCardSuit = "diamond"
			} else if strings.Contains(godown[i], "club") {
				wildCardSuit = "club"
			}
		}
		fmt.Println(wildCardSuit)
	}
	for i := 0; i < len(godown); i++ {
		if strings.Contains(checkedcard, wildCardSuit) {
			if strings.Contains(godown[i], "joker") {
				if i == 0 {
					if cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] {
						card = godown[i]
						// handList = append(handList, godown[i])
						godown[i] = checkedcard
						flag = true
					} else if cardsValue[godown[i+1]] == -1 {
						if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] {
							godown[i] = checkedcard
							card = godown[i]
							// handList = append(handList, godown[i])
							flag = true
						} else if cardsValue[godown[i+2]]-3 == cardsValue[checkedcard] {
							fmt.Println("hfccccccccccccccccccccccccd")
							var temp string
							godown = append(godown, checkedcard)
							flag = true
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp
							}
						} else if cardsValue[godown[i+2]]-1 == cardsValue[checkedcard] {
							fmt.Println("hfgfdhffhhjddfdddhdhfdfdfd")
							var temp = godown[i+1]
							godown[i+1] = checkedcard
							flag = true
							godown = append(godown, temp)
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}

						}

					} else if cardsValue[godown[i+1]]-2 == cardsValue[checkedcard] {
						flag = true
						var temp string
						godown = append(godown, checkedcard)
						for swp := 0; swp < len(godown)-1; swp++ {
							temp = godown[swp]
							godown[swp] = godown[len(godown)-1]
							godown[len(godown)-1] = temp

						}
					}

				} else if i <= len(godown)-1 {
					fmt.Println("card name", cardsValue[godown[i]])
					fmt.Println("card name", cardsValue[godown[i-1]]+1)
					fmt.Println(i)
					if cardsValue[godown[i-1]]+1 == cardsValue[checkedcard] {
						card = godown[i]
						// handList = append(handList, godown[i])
						godown[i] = checkedcard
						flag = true

					} else if i == len(godown)-1 && cardsValue[godown[i-1]] == -1 {
						fmt.Println("12odhdhsc678")
						if cardsValue[godown[i-2]]+3 == cardsValue[checkedcard] {
							godown = append(godown, checkedcard)
							flag = true

						} else if cardsValue[godown[i-1]]+2 == cardsValue[checkedcard] {

							godown = append(godown, checkedcard)
							flag = true

						} else if cardsValue[godown[i-2]]+2 == cardsValue[checkedcard] {
							card = godown[i]
							// handList = append(handList, godown[i])
							godown[i] = checkedcard
							flag = true

						}
					} else if cardsValue[godown[i-1]] == -1 && i == 1 {
						if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] {
							card = godown[i]
							// handList = append(handList, godown[i])
							godown[i] = checkedcard
							flag = true

						}
					} else if i == len(godown)-1 {
						if cardsValue[godown[i-1]]+2 == cardsValue[checkedcard] {

							godown = append(godown, checkedcard)
							flag = true

						}

					}

				}
			} else if strings.Contains(godown[i], "2_") {
				if i == 0 {

					if cardsValue[godown[i+1]]-1 == 14 {
						if cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] {
							if strings.Contains(checkedcard, wildCardSuit) {
								// handList = append(handList, godown[i])
								godown = append(godown, godown[i])
								godown[i] = checkedcard
								flag = true
							}

						}
					} else if cardsValue[godown[i+1]] == -1 {
						if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-2 != 3 {

							var temp string
							temp = godown[i]
							godown[i] = checkedcard
							godown = append(godown, temp)
							flag = true
							// var temp1 string
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}
						} else if cardsValue[godown[i+2]]-2 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-2 == 3 {
							godown = append(godown, godown[i])
							godown[i] = checkedcard
							flag = true

						} else if cardsValue[godown[i+2]]-3 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-3 == 3 {
							godown = append(godown, checkedcard)
							temp := godown[i]
							godown[i] = checkedcard
							godown[len(godown)-1] = temp
							flag = true

						} else if cardsValue[godown[i+2]]-3 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-3 != 3 {

							var temp string
							temp = godown[i]
							godown[i] = checkedcard
							godown = append(godown, temp)
							flag = true
							// var temp1 string
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}
						}

					} else if cardsValue[godown[i+1]]-2 == cardsValue[checkedcard] && cardsValue[godown[i+2]]-2 != 3 {
						fmt.Println("hgfhefefefefefefefvfv khfkhhkfh")

						// var temp string
						// /temp = godown[i]
						// godown[i] = checkedcard
						godown = append(godown, checkedcard)
						flag = true
						var temp string
						// temp = godown[]
						for swp := 0; swp < len(godown)-1; swp++ {
							temp = godown[swp]
							godown[swp] = godown[len(godown)-1]
							godown[len(godown)-1] = temp

						}

					} else if cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] {
						var temp string
						temp = godown[i]
						if cardsValue[godown[i+1]]-1 != 3 {
							godown[i] = checkedcard
							godown = append(godown, temp)
							flag = true
							// var temp1 string
							for swp := 0; swp < len(godown)-1; swp++ {
								temp = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp

							}
						} else {
							flag = true
							temp = godown[i]
							godown[i] = checkedcard
							godown = append(godown, temp)

						}

					}
					// else if cardsValue[godown[i+2]]-3
				} else if i <= len(godown)-1 {

					if cardsValue[godown[i-1]]+1 == cardsValue[checkedcard] {
						var temp4 string
						temp4 = godown[i]
						godown[i] = checkedcard
						flag = true
						godown = append(godown, temp4)
						for swp := 0; swp < len(godown)-1; swp++ {
							temp4 = godown[swp]
							godown[swp] = godown[len(godown)-1]
							godown[len(godown)-1] = temp4
						}
					}
					if i == 0 /*len(godown)-1 */ && cardsValue[godown[i+1]]-1 == cardsValue[checkedcard] || cardsValue[godown[i-1]]+1 == cardsValue[checkedcard] {

						if strings.Contains(godown[0], "joker_A_3_0") || cardsValue[godown[0]] == 3 {
							if cardsValue[godown[1]]-1 == 3 {
								godown = append(godown, godown[i])
								// godown = append(godown, checkedcard)
								godown[i] = checkedcard
								flag = true
							}

						} else {

							temp := godown[i]
							var temp1 string
							godown[i] = checkedcard
							godown = append(godown, temp)
							for swp := 0; swp < len(godown)-1; swp++ {
								temp1 = godown[swp]
								godown[swp] = godown[len(godown)-1]
								godown[len(godown)-1] = temp1
							}
							godown[0] = temp
							flag = true
						}
						// godown[i] = checkedcard
						// godown1 = godown
						// godown[0] = godown[i]
						// godown[i] = checkedcard
						// for j := 1; j == len(godown); j++ {
						// 	godown[j] = godown1[j-1]
						// }
						// godown[i] = checkedcard

					} else if i == len(godown)-1 {
						//fmt.Println("hgfhfhgghggjgfhjggfh")
						if cardsValue[godown[i-1]] == -1 {
							fmt.Println("hgfhfhgghggjgfhjggwfwhwkjkhfkhhkfh")
							if cardsValue[godown[i-2]]+3 == cardsValue[checkedcard] {
								fmt.Println("hgfhfhgghggjgfhjggfh")
								godown = append(godown, checkedcard)
								flag = true
							} else if cardsValue[godown[i-2]]+2 == cardsValue[checkedcard] {
								//fmt.Println("hgfhfhgghggjgfhjggfhewfefwffwfweewffwfeww")
								temp2 := godown[i]
								godown[i] = checkedcard
								godown = append(godown, temp2)
								flag = true
								if cardsValue[godown[0]] != 3 {
									for swp := 0; swp < len(godown)-1; swp++ {
										temp2 = godown[swp]
										godown[swp] = godown[len(godown)-1]
										godown[len(godown)-1] = temp2
									}

								}

							}

						} else {
							fmt.Println("1234567890oiuhgvcxzsdfrtghjk")
							if cardsValue[godown[i-1]]+2 == cardsValue[checkedcard] {
								// fmt.Println("741258/5369")
								// temp3 := godown[i]
								// godown[i] = checkedcard
								godown = append(godown, checkedcard)
								flag = true
								// flag = true
								// if cardsValue[godown[0]] != 3 {
								// 	for swp := 0; swp < len(godown)-1; swp++ {
								// 		temp3 = godown[swp]
								// 		godown[swp] = godown[len(godown)-1]
								// 		godown[len(godown)-1] = temp3
								// 	}

								// }

							}
						}
					}

				}
			} else if cardsValue[godown[len(godown)-1]]+1 == cardsValue[checkedcard] && i == len(godown)-1 {
				flag = true
				godown = append(godown, checkedcard)

			} else if cardsValue[godown[0]]-1 == cardsValue[checkedcard] && i == 0 {
				flag = true
				godown = append(godown, checkedcard)
				var temp1 string
				for swp := 0; swp < len(godown)-1; swp++ {
					temp1 = godown[swp]
					godown[swp] = godown[len(godown)-1]
					godown[len(godown)-1] = temp1
				}

			}
			if flag == true {
				break
			}
		}
	}
	if flag == false {
		fmt.Println("SORRY YOU CANNOT ADJUST THIS CARD  IN SEQUENCE")
		godown = nil
	}
	log.Println(godown)
	return godown, card
}
func threeAsoos(godown []string, checkedcard string) ([]string, string) {
	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 3, "club_4_0": 4, "club_5_0": 5,
		"club_6_0": 6, "club_7_0": 7, "club_8_0": 8, "club_9_0": 9, "club_10_0": 10, "club_J_0": 11, "club_Q_0": 12, "club_K_0": 13, "club_A_0": 14,
		"diamond_2_0": -1, "diamond_3_0": 3, "diamond_4_0": 4, "diamond_5_0": 5, "diamond_6_0": 6, "diamond_7_0": 7, "diamond_8_0": 8,
		"diamond_9_0": 9, "diamond_10_0": 10, "diamond_J_0": 11, "diamond_Q_0": 12, "diamond_K_0": 13, "diamond_A_0": 14, "heart_2_0": -1,
		"heart_3_0": 3, "heart_4_0": 4, "heart_5_0": 5, "heart_6_0": 6, "heart_7_0": 7, "heart_8_0": 8, "heart_9_0": 9, "heart_10_0": 10,
		"heart_J_0": 11, "heart_Q_0": 12, "heart_K_0": 13, "heart_A_0": 14, "spade_2_0": -1, "spade_3_0": 3, "spade_4_0": 4, "spade_5_0": 5,
		"spade_6_0": 6, "spade_7_0": 7, "spade_8_0": 8, "spade_9_0": 9, "spade_10_0": 10, "spade_J_0": 11, "spade_Q_0": 12, "spade_K_0": 13, "spade_A_0": 14,
		"club_2_1": -1, "club_3_1": 3, "club_4_1": 4, "club_5_1": 5,
		"club_6_1": 6, "club_7_1": 7, "club_8_1": 8, "club_9_1": 9, "club_10_1": 10, "club_J_1": 11, "club_Q_1": 12, "club_K_1": 13, "club_A_1": 14,
		"diamond_2_1": -1, "diamond_3_1": 3, "diamond_4_1": 4, "diamond_5_1": 5, "diamond_6_1": 6, "diamond_7_1": 7, "diamond_8_1": 8,
		"diamond_9_1": 9, "diamond_10_1": 10, "diamond_J_1": 11, "diamond_Q_1": 12, "diamond_K_1": 13, "diamond_A_1": 14, "heart_2_1": -1,
		"heart_3_1": 3, "heart_4_1": 4, "heart_5_1": 5, "heart_6_1": 6, "heart_7_1": 7, "heart_8_1": 8, "heart_9_1": 9, "heart_10_1": 10,
		"heart_J_1": 11, "heart_Q_1": 12, "heart_K_1": 13, "heart_A_1": 14, "spade_2_1": -1, "spade_3_1": 3, "spade_4_1": 4, "spade_5_1": 5,
		"spade_6_1": 6, "spade_7_1": 7, "spade_8_1": 8, "spade_9_1": 9, "spade_10_1": 10, "spade_J_1": 11, "spade_Q_1": 12, "spade_K_1": 13, "spade_A_1": 14}

	threeAsoos := ""
	Suit := ""
	flag := false
	card := ""
	var temp int
	wildCounter := 0

	if strings.Contains(checkedcard, "heart") {
		Suit = "heart"
	} else if strings.Contains(checkedcard, "spade") {
		Suit = "spade"
	} else if strings.Contains(checkedcard, "diamond") {
		Suit = "diamond"
	} else if strings.Contains(checkedcard, "club") {
		Suit = "club"
	} else if strings.Contains(checkedcard, "joker") {
		Suit = "joker"
	}

	fmt.Println(" check", Suit)
	if cardsValue[godown[0]] == 3 || cardsValue[godown[1]] == 3 && cardsValue[checkedcard] == 3 {
		fmt.Println("firsxxxxxxxt check")
		threeAsoos = "3"
	} else if cardsValue[godown[0]] == 14 || cardsValue[godown[1]] == 14 && cardsValue[checkedcard] == 14 {
		threeAsoos = "A"
		fmt.Println("xxxxxxxxxxxxxxxxxxxxxxx222 check")
	}
	if threeAsoos != "" && cardsValue[checkedcard] == 3 || cardsValue[checkedcard] == 14 || cardsValue[checkedcard] == -1 {
		for i := 0; i < len(godown); i++ {

			if !strings.Contains(godown[i], Suit) && cardsValue[godown[i]] != -1 {
				fmt.Println("33 xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx3 check")
				flag = true
			} else if !strings.Contains(godown[i], Suit) && cardsValue[godown[i]] == -1 && wildCounter == 0 {
				fmt.Println("44 xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx4 check")
				card = godown[i]
				temp = i
				flag = true
				wildCounter++
			} else {
				flag = false
				fmt.Println("4////////////////////4 check")
				break
			}
			// else if wildcounter > 0 {
			// 	if !strings.Contains(godown[i], Suit) && cardsValue[godown[i]] != -1 {
			// 		flag = true

			// 	}

			// }

		}
	}

	if flag == true {
		if card == "" {
			fmt.Println("5555 check")
			if len(godown) <= 3 || cardsValue[checkedcard] != -1 {
				godown = append(godown, checkedcard)
			}
		} else {
			if len(godown) <= 3 {
				godown[temp] = checkedcard

				fmt.Println("666 check")
			}
		}
	} else {
		fmt.Println("Can not make asaoos ")
	}

	// fmt.Println(godown, card)
	log.Println(godown)
	return godown, card

}
func reversewildCardEntry(godown []string, loc int, checkedcard string) ([]string, string) {

	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 14, "club_4_0": 13, "club_5_0": 12,
		"club_6_0": 11, "club_7_0": 10, "club_8_0": 9, "club_9_0": 8, "club_10_0": 7, "club_J_0": 6, "club_Q_0": 5, "club_K_0": 4, "club_A_0": 3,
		"diamond_2_0": -1, "diamond_3_0": 14, "diamond_4_0": 13, "diamond_5_0": 12, "diamond_6_0": 11, "diamond_7_0": 10, "diamond_8_0": 9,
		"diamond_9_0": 8, "diamond_10_0": 7, "diamond_J_0": 6, "diamond_Q_0": 5, "diamond_K_0": 4, "diamond_A_0": 3, "heart_2_0": -1,
		"heart_3_0": 14, "heart_4_0": 13, "heart_5_0": 12, "heart_6_0": 11, "heart_7_0": 10, "heart_8_0": 9, "heart_9_0": 8, "heart_10_0": 7,
		"heart_J_0": 6, "heart_Q_0": 5, "heart_K_0": 4, "heart_A_0": 3, "spade_2_0": -1, "spade_3_0": 14, "spade_4_0": 13, "spade_5_0": 12,
		"spade_6_0": 11, "spade_7_0": 10, "spade_8_0": 9, "spade_9_0": 8, "spade_10_0": 7, "spade_J_0": 6, "spade_Q_0": 5, "spade_K_0": 4, "spade_A_0": 3,
		"club_2_1": -1, "club_3_1": 14, "club_4_1": 13, "club_5_1": 12,
		"club_6_1": 11, "club_7_1": 10, "club_8_1": 9, "club_9_1": 8, "club_10_1": 7, "club_J_1": 6, "club_Q_1": 5, "club_K_1": 4, "club_A_1": 3,
		"diamond_2_1": -1, "diamond_3_1": 14, "diamond_4_1": 13, "diamond_5_1": 12, "diamond_6_1": 11, "diamond_7_1": 10, "diamond_8_1": 9,
		"diamond_9_1": 8, "diamond_10_1": 7, "diamond_J_1": 6, "diamond_Q_1": 5, "diamond_K_1": 4, "diamond_A_1": 3, "heart_2_1": -1,
		"heart_3_1": 14, "heart_4_1": 13, "heart_5_1": 12, "heart_6_1": 11, "heart_7_1": 10, "heart_8_1": 9, "heart_9_1": 8, "heart_10_1": 7,
		"heart_J_1": 6, "heart_Q_1": 5, "heart_K_1": 4, "heart_A_1": 3, "spade_2_1": -1, "spade_3_1": 14, "spade_4_1": 13, "spade_5_1": 12,
		"spade_6_1": 11, "spade_7_1": 10, "spade_8_1": 9, "spade_9_1": 8, "spade_10_1": 7, "spade_J_1": 6, "spade_Q_1": 5, "spade_K_1": 4, "spade_A_1": 3}

	which_Wild := ""
	counter := 0
	card := ""
	var dublicate string
	flag := false
	var wildcrdloc int
	if strings.Contains(checkedcard, "_2") {
		which_Wild = "_2"
	} else if strings.Contains(checkedcard, "_joker") {
		which_Wild = "_joker"

	}
	for i := 0; i < len(godown); i++ {
		if cardsValue[godown[i]] == -1 {
			counter++
			dublicate = godown[i]
			fmt.Println(godown[i], cardsValue[godown[i]])
			fmt.Println(counter)
			wildcrdloc = i + 1

		}

	}
	fmt.Println("rev11111111111111111111111111111111", wildcrdloc, loc)
	fmt.Println("11111")
	fmt.Println(counter)
	if counter >= 2 || dublicate == checkedcard {
		fmt.Println("Already wild card exists")
	} else {
		if loc == wildcrdloc && strings.Contains(which_Wild, "_joker") {
			fmt.Println("111wwwwwwwwwwwwwwwww11")
			if cardsValue[godown[0]] != 3 {
				fmt.Println("1122222222222222222222222222222222222222")
				flag = true
				godown = append(godown, checkedcard)
				var temp string
				for swp := 0; swp < len(godown)-1; swp++ {
					temp = godown[swp]
					godown[swp] = godown[len(godown)-1]
					godown[len(godown)-1] = temp

				}

			} else {
				fmt.Println("555qqqqqqqqqq55555555555555555")
				flag = true
				godown = append(godown, checkedcard)
			}
		} else if loc == 0 && cardsValue[godown[0]] != 3 && cardsValue[godown[0]] != -1 {
			fmt.Println("333ssssssssssssss33333")
			flag = true
			godown = append(godown, checkedcard)
			var temp string
			for swp := 0; swp < len(godown)-1; swp++ {
				temp = godown[swp]
				godown[swp] = godown[len(godown)-1]
				godown[len(godown)-1] = temp

			}
		} else if loc == len(godown)+1 && cardsValue[godown[len(godown)-1]] != 14 && cardsValue[godown[len(godown)-1]] != -1 {
			fmt.Println("44444dddddddddddddddddddddd44")
			flag = true

			godown = append(godown, checkedcard)
		} else if loc == 0 && cardsValue[godown[0]] == -1 /* &&  cardsValue[godown[1]] != 4*/ {
			fmt.Println("5dddddddddddddddddddddd5555")
			if cardsValue[godown[1]] != 4 {
				fmt.Println("6666cccccccccccccccccccccc6666")
				flag = true
				godown = append(godown, checkedcard)
				var temp string
				for swp := 0; swp < len(godown)-1; swp++ {
					temp = godown[swp]
					godown[swp] = godown[len(godown)-1]
					godown[len(godown)-1] = temp

				}

			} else {
				fmt.Println("7777xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx77")
				fmt.Println("you cannot add this card in the")
			}
		} else if loc == len(godown)+1 && cardsValue[godown[len(godown)-1]] != -1 && cardsValue[godown[len(godown)-1]] != 14 {
			fmt.Println("rrrrrrrrrrrrrrxxxxxxxxrrrrr66")
			flag = true
			godown = append(godown, checkedcard)

		} else if loc > 0 && loc <= len(godown) {

			fmt.Println("you ca0000000qqqqqqqq0rd in the")
			if strings.Contains(godown[loc-1], "joker_") {
				flag = true
				card = godown[loc-1]
				godown[loc-1] = checkedcard
				// fmt.Println("you cannot add this card in the/////////")
			} else if strings.Contains(godown[loc-1], "_2_") {
				if cardsValue[godown[0]] != 3 {
					godown = append(godown, checkedcard)
					flag = true
				} else if cardsValue[godown[len(godown)-1]] == 14 {
					fmt.Println("987456321458628621458625")
					godown = append(godown, checkedcard)
					var temp string
					for swp := 0; swp < len(godown)-1; swp++ {
						temp = godown[swp]
						godown[swp] = godown[len(godown)-1]
						godown[len(godown)-1] = temp
						flag = true

					}

				}
			}

		} else if len(godown) < loc && cardsValue[godown[len(godown)-1]] != 14 && cardsValue[godown[len(godown)-1]] != -1 {
			flag = true
			godown = append(godown, checkedcard)

			fmt.Println("\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\ in the")

		} else if loc == len(godown)+1 && cardsValue[godown[len(godown)-2]] != 13 && cardsValue[godown[len(godown)-1]] == -1 {
			flag = true
			godown = append(godown, checkedcard)

			fmt.Println("ae")
		}
		//  else if len(godown) < loc && cardsValue[godown[len(godown)-1]] == -1 {

		// }

	}
	if flag == false {
		godown = nil
	}
	return godown, card

}
func wildCardEntry(godown []string, loc int, checkedcard string) ([]string, string) {
	which_Wild := checkedcard
	cardsValue := map[string]int{"joker_A_3_0": -1, "joker_A_3_1": -1, "club_2_0": -1, "club_3_0": 3, "club_4_0": 4, "club_5_0": 5,
		"club_6_0": 6, "club_7_0": 7, "club_8_0": 8, "club_9_0": 9, "club_10_0": 10, "club_J_0": 11, "club_Q_0": 12, "club_K_0": 13, "club_A_0": 14,
		"diamond_2_0": -1, "diamond_3_0": 3, "diamond_4_0": 4, "diamond_5_0": 5, "diamond_6_0": 6, "diamond_7_0": 7, "diamond_8_0": 8,
		"diamond_9_0": 9, "diamond_10_0": 10, "diamond_J_0": 11, "diamond_Q_0": 12, "diamond_K_0": 13, "diamond_A_0": 14, "heart_2_0": -1,
		"heart_3_0": 3, "heart_4_0": 4, "heart_5_0": 5, "heart_6_0": 6, "heart_7_0": 7, "heart_8_0": 8, "heart_9_0": 9, "heart_10_0": 10,
		"heart_J_0": 11, "heart_Q_0": 12, "heart_K_0": 13, "heart_A_0": 14, "spade_2_0": -1, "spade_3_0": 3, "spade_4_0": 4, "spade_5_0": 5,
		"spade_6_0": 6, "spade_7_0": 7, "spade_8_0": 8, "spade_9_0": 9, "spade_10_0": 10, "spade_J_0": 11, "spade_Q_0": 12, "spade_K_0": 13, "spade_A_0": 14,
		"club_2_1": -1, "club_3_1": 3, "club_4_1": 4, "club_5_1": 5,
		"club_6_1": 6, "club_7_1": 7, "club_8_1": 8, "club_9_1": 9, "club_10_1": 10, "club_J_1": 11, "club_Q_1": 12, "club_K_1": 13, "club_A_1": 14,
		"diamond_2_1": -1, "diamond_3_1": 3, "diamond_4_1": 4, "diamond_5_1": 5, "diamond_6_1": 6, "diamond_7_1": 7, "diamond_8_1": 8,
		"diamond_9_1": 9, "diamond_10_1": 10, "diamond_J_1": 11, "diamond_Q_1": 12, "diamond_K_1": 13, "diamond_A_1": 14, "heart_2_1": -1,
		"heart_3_1": 3, "heart_4_1": 4, "heart_5_1": 5, "heart_6_1": 6, "heart_7_1": 7, "heart_8_1": 8, "heart_9_1": 9, "heart_10_1": 10,
		"heart_J_1": 11, "heart_Q_1": 12, "heart_K_1": 13, "heart_A_1": 14, "spade_2_1": -1, "spade_3_1": 3, "spade_4_1": 4, "spade_5_1": 5,
		"spade_6_1": 6, "spade_7_1": 7, "spade_8_1": 8, "spade_9_1": 9, "spade_10_1": 10, "spade_J_1": 11, "spade_Q_1": 12, "spade_K_1": 13, "spade_A_1": 14}

	//which_Wild := ""
	counter := 0
	card := ""
	var dublicate string
	flag := false
	var wildcrdloc int

	// if strings.Contains(checkedcard, "_2") {
	// 	which_Wild = "_2"
	// } else if strings.Contains(checkedcard, "_joker") {
	// 	which_Wild = "_joker"

	// }
	fmt.Println("which..........._Wild", which_Wild)
	for i := 0; i < len(godown); i++ {
		fmt.Println("w-000000000000000000000000000000000000000000000000000000000000000_Wild")
		if cardsValue[godown[i]] == -1 {
			counter++
			dublicate = godown[i]
			fmt.Println(godown[i], cardsValue[godown[i]])
			fmt.Println(counter)
			wildcrdloc = i + 1

		}

	}
	fmt.Println("wildcrdloc", wildcrdloc, loc)

	fmt.Println(counter)
	if counter < 2 && dublicate == checkedcard {
		fmt.Println("Already wild card exists")
	} else {
		if loc == wildcrdloc && strings.Contains(which_Wild, "joker_") {
			fmt.Println("joker seen")
			if cardsValue[godown[0]] != 3 {
				fmt.Println("1st its ")
				flag = true
				godown = append(godown, checkedcard)
				var temp string
				for swp := 0; swp < len(godown)-1; swp++ {
					temp = godown[swp]
					godown[swp] = godown[len(godown)-1]
					godown[len(godown)-1] = temp

				}

			} else {
				fmt.Println("else of 1st first")
				flag = true
				godown = append(godown, checkedcard)
			}
		} else if loc == 0 && cardsValue[godown[0]] != 3 && cardsValue[godown[0]] != -1 {
			fmt.Println("33333333")
			flag = true
			godown = append(godown, checkedcard)
			var temp string
			for swp := 0; swp < len(godown)-1; swp++ {
				temp = godown[swp]
				godown[swp] = godown[len(godown)-1]
				godown[len(godown)-1] = temp

			}
		} else if loc == len(godown)+1 && cardsValue[godown[len(godown)-1]] != 14 && cardsValue[godown[len(godown)-1]] != -1 {
			fmt.Println("4444444")
			flag = true

			godown = append(godown, checkedcard)
		} else if loc == 0 && cardsValue[godown[0]] == -1 {
			fmt.Println("55555")
			if !strings.Contains(godown[1], "_4") {
				fmt.Println("66666666")
				flag = true
				godown = append(godown, checkedcard)
				var temp string
				for swp := 0; swp < len(godown)-1; swp++ {
					temp = godown[swp]
					godown[swp] = godown[len(godown)-1]
					godown[len(godown)-1] = temp

				}

			} else {
				fmt.Println("777777")
				fmt.Println("you cannot add this card in the")
			}
		} else if loc == len(godown)+1 && cardsValue[godown[len(godown)-1]] == -1 {

			if !strings.Contains(godown[len(godown)-1], "_A") {
				fmt.Println("88888")
				flag = true
				godown = append(godown, checkedcard)

			} else {
				fmt.Println("you cannot add this card in the")
			}
		} else if loc > 0 && loc <= len(godown) {

			fmt.Println("99999999999999999999999")
			if strings.Contains(godown[loc-1], "joker_") {
				flag = true
				card = godown[loc-1]
				godown[loc-1] = checkedcard
				// fmt.Println("you cannot add this card in the/////////")
			} else {

				flag = true
				godown = append(godown, checkedcard)
			}

		} else if len(godown) < loc && cardsValue[godown[len(godown)-1]] != 14 {
			flag = true
			godown = append(godown, checkedcard)

			fmt.Println("1010101010101010101010 in the")

		}
		//  else if len(godown) < loc && cardsValue[godown[len(godown)-1]] == -1 {

		// }

	}
	if flag == false {
		godown = nil
	}
	return godown, card

}
func submitLeadboardscore(stateMain interface{}, logger runtime.Logger, ctx context.Context, db *sql.DB, nk runtime.NakamaModule) {
	state := stateMain.(*MatchState)
	props := &Leaderboard.LeaderboardSubmitProps{}

	// id := state.players[state.OppHostteamPlayer1Name].userId
	// score := state.OppHostteamPlayer1
	// username := state.players[state.OppHostteamPlayer1Name].userName

	// props.Score = int64(state.OppHostteamPlayer1)
	// props.UserId = state.players[state.OppHostteamPlayer1Name].userId
	// props.UserName = state.players[state.OppHostteamPlayer1Name].userName
	// props.LeaderboardId = "banakilleader_board"
	// payload1, _ := json.Marshal(props)
	// logger.Info("payload1, _ := json.Marshal(props)", payload1)
	for _, plyr := range state.players {
		props.Score = plyr.totalScore
		props.UserId = plyr.userId
		props.UserName = plyr.userName
		props.LeaderboardId = "banakilleader_board"
		payload1, _ := json.Marshal(props)
		logger.Info("payload1, _ := json.Marshal(props)", payload1)
		rpcReturn, err := Leaderboard.LeaderboardBanakili(ctx, logger, db, nk, string(payload1))
		logger.Info("rpcReturn------------------formmatchloop-----", rpcReturn, err)
	}
}
func eventSessionStart(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	logger.Info("session start %v %v", ctx, evt)
}

func eventSessionEnd(ctx context.Context, logger runtime.Logger, evt *api.Event) {
	logger.Info("session end %v %v", ctx, evt)
}
