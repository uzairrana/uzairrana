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

package tarneeb

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"strconv"
	"strings"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type DuringMatchForDummy struct{}

type Team struct {
	teamScore   int
	member1     string
	member2     string
	teamBid     int
	roundScore  int
	kaboot      string
	acceptedBid bool
}

type Player struct {
	//userId    string
	displayName    string
	userName       string
	placeHolderUrl string
	bid            int
	teamMate       string
	teamNum        string
	points         int
	seatPosition   int
	cards          []string
	presence       runtime.Presence
	score          int
	cardWithValue  map[string]int //uzair
	bot            bool
	response       bool
	opcode         int64
}
type MatchState struct {
	debug         bool
	joinedPlayers int
	userNames     []string
	players       map[string]*Player
	teams         map[string]*Team
	//team              [2]Team
	deck               [52]string
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
	bidPassCounter     int
	bidRoundCounter    int
	firstTurnsCounter  int
	firstBiddingPerson string
	turnRoundCounter   int
	gameRoundCounter   int
	cardWithValue      map[string]int

	highestBid       int
	highestBidPerson string
	minBid           int

	trumpSuit         string
	highestCard       string
	highestCardPlayer string
	highestCardValue  int
	botOpcode         int64
	currentBot        string
	nextBot           string

	//roundScoreSeq    [4]int
	currentTurn           string
	nextTurn              string
	matchExitCounter      int
	matchExitBool         bool
	currentPlayerTurnCard []string

	sign string
}
type GetTeamMate struct {
	X string `json:"x"`
	Y bool   `json:"y"`
}
type Before_GamePlay struct {
	FirstPlayer     bool            `json:"FirstPlayer"`
	PlayersJoined   int             `json:"PlayersJoined"`
	HostUserName    string          `json:"HostUserName"`
	HostDisplayName string          `json:"HostDisplayName"`
	JoinedPlayers   []Name_UserName `json:"JoinedPlayers"`

	started          bool        `json:"Started"`
	SetTeamMate      string      `json:"SetTeamMate"`
	TeamMateResponse GetTeamMate `json:"TeamMateResponse"`

	StartingBidValue int `json:"StartingBidValue"`

	SittingArangement []string `json:"SittingArangement"`

	Dealer string `json:"Dealer"`
	Turn   string `josn:"Turn"`

	BidValue         int    `json:"BidValueInt"`
	BidValueString   string `json:"BidValue"`
	BidPass          bool   `json:"BidPass"`
	HighestBid       int    `json:"HighestBid"`
	HighestBidPerson string `json:"HighestBidPerson"`

	Kaboot         string `json:"kaboot"`
	TrumpSuit      string `json:"TrumpSuit"`
	PlayerUserName string `json:"PlayerUserName"`

	Cards                  []string  `json:"Cards"`
	ThrownCard             string    `json:"ThrownCard"`
	Chat                   string    `json:"Chat"`
	RoundScore             [4]string `json:"RoundScore"`
	Team1                  int       `json:"Team1"`
	Team2                  int       `json:"Team2"`
	Exit                   string    `json:"Exit"`
	Score_Team1            string    `json:"Score_Team1"`
	Score_Team2            string    `json:"Score_Team2"`
	Winner1                string    `json:"Winner1"`
	Winner2                string    `json:"Winner2"`
	RemoveCards            string    `json:"RemoveCards"`
	HighestCardPerson      string    `json:"HighestCardPerson"`
	HighestCardPersonScore int       `json:"HighestCardPersonScore"`
	HostTeam               int       `json:"HostTeam"`
	HostOpponentTeam       int       `json:"HostOpponentTeam"`
}

type Name_UserName struct {
	Name      string
	UserName  string
	AvatarUrl string
}

type Tarneeb struct{}

func (t *Tarneeb) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	//var debug bool
	debug := false
	// if d, ok := params["debug"]; ok {
	// 	if dv, ok := d.(bool); ok {
	// 		debug = dv
	// 	}
	// }
	//logger.Info("___________match init, starting")
	state := &MatchState{
		debug:              debug,
		hostSignal:         true,
		joinedPlayers:      0,
		players:            make(map[string]*Player),
		teams:              make(map[string]*Team),
		deck:               [52]string{"club_A", "club_2", "club_3", "club_4", "club_5", "club_6", "club_7", "club_8", "club_9", "club_10", "club_J", "club_Q", "club_K", "diamond_A", "diamond_2", "diamond_3", "diamond_4", "diamond_5", "diamond_6", "diamond_7", "diamond_8", "diamond_9", "diamond_10", "diamond_J", "diamond_Q", "diamond_K", "heart_A", "heart_2", "heart_3", "heart_4", "heart_5", "heart_6", "heart_7", "heart_8", "heart_9", "heart_10", "heart_J", "heart_Q", "heart_K", "spade_A", "spade_2", "spade_3", "spade_4", "spade_5", "spade_6", "spade_7", "spade_8", "spade_9", "spade_10", "spade_J", "spade_Q", "spade_K"},
		playersLimit:       4,
		playersLimitSignal: true,
		bidFlag:            false,
		highestBid:         0,
		cardWithValue:      make(map[string]int), //uzair
		minBid:             7,
		sittingArangement:  [4]string{"BAtXtNDFgZ", "BAtXtNDFgZ", "BAtXtNDFgZ", "BAtXtNDFgZ"},
		firstTurnsCounter:  0,
	}

	if state.debug {
		logger.Info("match init, starting with debug: %v", state.debug)
	}
	tickRate := 1

	shuffle(state.deck, state)
	label := "TarneebGame"
	return state, tickRate, label
}

func (t *Tarneeb) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
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

func (t *Tarneeb) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, presences []runtime.Presence) interface{} {
	state := stateMain.(*MatchState)
	if stateMain.(*MatchState).debug {
		for _, presence := range presences {
			logger.Info("match join username %v user_id %v session_id %v node %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId())

		}
	}

	logger.Info("___________Match Join, starting", state.joinedPlayers)

	// if state.joinedPlayers == 1 {
	// 	state.hostUserName = presence.GetUsername()
	// }

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USERNAME).(string)
	if !ok {

		logger.Info("presence isssss: %s", ok, userID)
		// return "", errors.New("Invalid context")
	}

	//logger.Info("*********************************-", state.joinedPlayers)

	return state
}

func (t *Tarneeb) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, presences []runtime.Presence) interface{} {
	state := stateMain.(*MatchState)

	for _, presence := range presences {
		logger.Info("match leave username %v user_id %v session_id %v node %v", presence.GetUsername(), presence.GetUserId(), presence.GetSessionId(), presence.GetNodeId())
		state.players[presence.GetUsername()].bot = true

		if state.players[presence.GetUsername()].response {
			state.botOpcode = state.players[presence.GetUsername()].opcode
		}
	}

	state.matchExitCounter--
	state.matchExitBool = true

	logger.Info("___________Match Leave, starting ", state.matchExitCounter)
	//state.players[state.hostUserName].bot = true

	return state
}

func (t *Tarneeb) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, stateMain interface{}, messages []runtime.MatchData) interface{} {
	state := stateMain.(*MatchState)
	if state.debug {
		logger.Info("match loop match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match loop match_id %v message count %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), len(messages))
	}

	//logger.Info("msg username", state.hostUserName)

	if state.matchExitCounter <= 0 && state.matchExitBool {
		logger.Info("_____going to kil________-", state.matchExitCounter)
		logger.Info("______going to kil_______-", state.matchExitBool)

		return nil
	}
	HostPlyr := &Before_GamePlay{}
	dataFromBot := &Before_GamePlay{}
	var dataFromBot_byte []byte

	if state.botOpcode != 0 {
		//dataFromBot_byte = setBotJsonData(state, logger, dataFromBot)
		var err error
		if dataFromBot_byte, err = json.Marshal(setBotJsonData(state, logger, dataFromBot)); err != nil {
			logger.Info("Error is: ", err)

		} else {
			logger.Info("dataFromBot successfully marshal: ", dataFromBot_byte)

		}
	}
	//logger.Info("_____________-", state.matchExitCounter)
	//logger.Info("_____________-", state.matchExitBool)

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
				if state.joinedPlayers == state.playersLimit && state.players[state.hostUserName].bot {
					state.botOpcode = 401
					logger.Info("state.botOpcode is: ", state.botOpcode)

				} else {
					state.players[state.hostUserName].response = true
					state.players[state.hostUserName].opcode = 401
				}
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

		decoded := &DuringMatchForDummy{}

		if err := json.Unmarshal(msg.GetData(), &decoded); err != nil {
			logger.Info("Error is: ", err)

		}

		data = msg.GetData()
		opcd = msg.GetOpCode()

		logger.Info("=opcd After decode=", opcd)
	}
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
			state.players[state.hostUserName].response = false
			state.players[state.hostUserName].opcode = 0

			logger.Info("Data is________________________: ", state.hostUserName)

			dispatcher.BroadcastMessage(102, Json_started, nil, nil, true)
			cardWithValue(state, logger)

			//Bot func
			if state.players[state.hostUserName].bot {
				state.botOpcode = 402
				logger.Info("state.botOpcode is: ", state.botOpcode)

			} else {
				state.players[state.hostUserName].response = true
				state.players[state.hostUserName].opcode = 402
			}
			//Bot func

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

				state.players[state.hostUserName].response = false
				state.players[state.hostUserName].opcode = 0
				logger.Info("402 is state.players[TeamMate.PlayerUserName].presence_______: ", state.players[TeamMate.SetTeamMate].presence)

				presence := []runtime.Presence{state.players[TeamMate.SetTeamMate].presence}

				if Json_TeamMateReq, err := json.Marshal(TeamMate.SetTeamMate); err != nil {
					logger.Info("Error is: ", err)
				} else {
					logger.Info("402 sseendd 103 ________________________: ", state.hostUserName)

					dispatcher.BroadcastMessage(103, Json_TeamMateReq, presence, nil, true)

					//IF Requested Person IS NOT AVAILABLE

					//Bot func
					if state.players[TeamMate.SetTeamMate].bot {
						state.botOpcode = 403
						state.nextBot = TeamMate.SetTeamMate
						logger.Info("state.botOpcode is: ", state.botOpcode)

					} else {
						state.players[TeamMate.SetTeamMate].response = true
						state.players[TeamMate.SetTeamMate].opcode = 403
						state.nextBot = TeamMate.SetTeamMate

					}
					//Bot func
				}
			}

		}

		// opCOde = 0
		// opcd = 0
	case 403:

		// After team mate slction
		logger.Info("403    TeamMate.TeamMateResponse.y___________________: ")

		TeamMate := &Before_GamePlay{}
		if err := json.Unmarshal(data, &TeamMate); err != nil {
			logger.Info("Error is: ", err)
		} else {

			TeamMate.SetTeamMate = TeamMate.TeamMateResponse.X

			if TeamMate.SetTeamMate != "" {

				state.players[TeamMate.SetTeamMate].response = false
				state.players[TeamMate.SetTeamMate].opcode = 0
				state.nextBot = ""

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
						sittingArrangement.JoinedPlayers = append(sittingArrangement.JoinedPlayers, Nam_UserNam)
					}

					if Json_sittingArrangement, err := json.Marshal(sittingArrangement); err != nil {
						logger.Info("Error is: ", err)
					} else {
						dispatcher.BroadcastMessage(104, Json_sittingArrangement, nil, nil, true)

						// Score := &Before_GamePlay{}

						// Score.Team1 = 1 //state.teams[state.players[state.hostUserName].teamNum].roundScore
						// Score.Team2 = 2 //state.teams[state.players[state.hostOpponent].teamNum].roundScore

						// if Json_Score, err := json.Marshal(Score); err != nil {
						// 	logger.Info("Error is: ", err)
						// } else {
						// 	dispatcher.BroadcastMessage(113, Json_Score, nil, nil, true)

						// }
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

				} else {

					logger.Info("rejected Resending the 102 ___: ", state.hostUserName)

					Started := &Before_GamePlay{}
					Started.started = true
					if Json_started, err := json.Marshal(Started); err != nil {
						logger.Info("Error is: ", err)

					} else {
						logger.Info("to show again team list_________: ", state.hostUserName)

						dispatcher.BroadcastMessage(102, Json_started, nil, nil, true)
						if state.players[state.hostUserName].bot {
							state.botOpcode = 402
							logger.Info("state.botOpcode is: ", state.botOpcode)

						} else {
							state.players[state.hostUserName].response = true
							state.players[state.hostUserName].opcode = 402
						}

					}
				}
			}

		}

		// opCOde = 0
		// opcd = 0
	case 404:

		logger.Info("404    bbiiidddddddddddddddd___________________: ", data)

		logger.Info("BidFlag  is: ", state.bidFlag)
		bid := &Before_GamePlay{}
		//bidding

		if err := json.Unmarshal(data, &bid); err != nil {
			logger.Info("Error is: ", err)
		} else if bid.PlayerUserName != "" {

			state.players[bid.PlayerUserName].response = false
			state.players[bid.PlayerUserName].opcode = 0
			state.currentBot = ""

			logger.Info("minBidminBid is: ", state.minBid)
			bid.BidValue, _ = strconv.Atoi(bid.BidValueString)
			if (bid.BidValue >= state.minBid && bid.BidValue < 13 && bid.PlayerUserName != "") || (bid.BidValue == 13 && state.minBid >= 10 && bid.PlayerUserName != "") {
				state.players[bid.PlayerUserName].bid = bid.BidValue
				logger.Info("PlayerUserNamePlayerUserName is: ", bid.PlayerUserName)

				state.minBid = bid.BidValue
				if bid.BidValue > state.highestBid {
					logger.Info("highestBidhighestBid is: ", bid.PlayerUserName)

					state.highestBid = bid.BidValue
					state.highestBidPerson = bid.PlayerUserName

					state.bidRoundCounter++
					if bid.BidValue == 13 || state.bidRoundCounter == 4 {
						state.bidFlag = false
						state.bidRoundCounter = 4
					}

					bid.BidPass = true

					//state.bidRoundCounter++
					bid.Turn = bidTurnSequence(state.players[bid.PlayerUserName].seatPosition, state, logger)

					if state.bidFlag {
						if Json_Turn, err := json.Marshal(bid); err != nil {
							logger.Info("Error is: ", err)
						} else {
							logger.Info("11111ooooooOOOOOOoooooo7777___44OOO44 ", err)

							dispatcher.BroadcastMessage(107, Json_Turn, nil, nil, true)
							if state.players[bid.Turn].bot {
								state.currentBot = bid.Turn
								state.botOpcode = 404
							} else {
								state.players[bid.Turn].response = true
								state.players[bid.Turn].opcode = 404
								state.currentBot = bid.Turn
							}
						}
					}
				}

			} else if bid.BidValue == 0 && bid.PlayerUserName != "" {

				if state.bidPassCounter == 3 {

					bid.BidPass = false
					state.bidRoundCounter++

					if state.bidRoundCounter == 4 {
						state.bidFlag = false
					}

					logger.Info("state.bidFlag   is: ", state.bidFlag)

					if state.bidFlag { //qa quick chng
						bid.Turn = bidTurnSequence(state.players[bid.PlayerUserName].seatPosition, state, logger)
						if Json_Turn, err := json.Marshal(bid); err != nil {
							logger.Info("Error is: ", err)
						} else {
							logger.Info("111110000oooooo00000077777 BIDPAAASS_404: ", bid.PlayerUserName)

							dispatcher.BroadcastMessage(107, Json_Turn, nil, nil, true)
							if state.players[bid.Turn].bot {
								state.currentBot = bid.Turn
								state.botOpcode = 404
							} else {
								state.players[bid.Turn].response = true
								state.players[bid.Turn].opcode = 404
								state.currentBot = bid.Turn
							}
						}
					}
				} else if state.bidPassCounter < 3 {

					state.players[bid.PlayerUserName].bid = bid.BidValue
					state.bidPassCounter++
					logger.Info("bidPassCounter  is: ", state.bidPassCounter)

					if state.bidPassCounter == 3 {
						bid.BidPass = false

					} else {
						bid.BidPass = true
					}
					state.bidRoundCounter++

					if state.bidRoundCounter == 4 {
						state.bidFlag = false
					}

					if state.bidFlag {
						bid.Turn = bidTurnSequence(state.players[bid.PlayerUserName].seatPosition, state, logger)
						if Json_Turn, err := json.Marshal(bid); err != nil {
							logger.Info("Error is: ", err)
						} else {
							logger.Info("111110000oooooo00000077777 is_____404 BID PASCOUNTER < 3: ", bid.PlayerUserName)

							dispatcher.BroadcastMessage(107, Json_Turn, nil, nil, true)
							if state.players[bid.Turn].bot {
								state.currentBot = bid.Turn
								state.botOpcode = 404
							} else {
								state.players[bid.Turn].response = true
								state.players[bid.Turn].opcode = 404
								state.currentBot = bid.Turn
							}
						}
					}
				}

			}

		}

		logger.Info("players[bid.PlayerUserName].bid___ ()()()()___ ): ", state.players[bid.PlayerUserName].bid)

		logger.Info("state.bidRoundCounter++____________ ): ", state.bidRoundCounter)
		logger.Info("state.!state.bidFlag++____________ ): ", state.bidFlag)

		if !state.bidFlag && state.bidRoundCounter >= 4 && state.highestBidPerson != "" {
			logger.Info("405   highest bid___________________: ", data)

			state.firstBiddingPerson = state.highestBidPerson
			bid.HighestBid = state.highestBid
			bid.HighestBidPerson = state.highestBidPerson
			if Json_Turn, err := json.Marshal(bid); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(108, Json_Turn, nil, nil, true)
				if state.players[bid.HighestBidPerson].bot {
					state.botOpcode = 405
					state.currentBot = bid.HighestBidPerson
				} else {
					state.players[bid.HighestBidPerson].response = true
					state.players[bid.HighestBidPerson].opcode = 405
				}
			}
		}

		// opCOde = 0
		// opcd = 0
	case 405:
		//Trump suit
		logger.Info("405    Trruuummp suit___________________: ", data)

		TrumSuit := &Before_GamePlay{}
		state.trumpSuit = ""
		if err := json.Unmarshal(data, &TrumSuit); err != nil {
			logger.Info("Error is: ", err)
		} else {
			logger.Info("Data is: ", TrumSuit)

			state.players[state.highestBidPerson].response = false
			state.players[state.highestBidPerson].opcode = 0

		}

		if TrumSuit.TrumpSuit != "" {
			state.trumpSuit = TrumSuit.TrumpSuit
		}

		updateCardValues(state, logger, 40, state.trumpSuit)

		if Json_TrumSuit, err := json.Marshal(TrumSuit); err != nil {
			logger.Info("Error is: ", err)
		} else {
			dispatcher.BroadcastMessage(109, Json_TrumSuit, nil, nil, true)

			for _, plyr := range state.players {
				if plyr.bid > state.players[plyr.teamMate].bid {
					if plyr.teamNum == "team1" {
						state.teams["team1"].teamBid = plyr.bid
					} else {
						state.teams["team2"].teamBid = plyr.bid

					}
				} else {
					if plyr.teamNum == "team1" {
						state.teams["team1"].teamBid = state.players[plyr.teamMate].bid
					} else {
						state.teams["team2"].teamBid = state.players[plyr.teamMate].bid

					}
				}
			}

			logger.Info("Team1_____: ", state.teams["team1"].teamBid)
			logger.Info("Team2_____: ", state.teams["team2"].teamBid)

			// First implemented the bid mistakenly separate of each team (highest bid of the person, inside team bid) unawarenss of game rule, after clearing the rule set the same value of highest bid in both teams.
			// state.teams["team1"].teamBid = state.highestBid
			// state.teams["team2"].teamBid = state.highestBid

			if state.teams["team1"].teamBid > state.teams["team2"].teamBid {

				state.teams["team1"].acceptedBid = true
				state.teams["team2"].acceptedBid = false

				//	state.players[state.players[state.hostUserName].teamMate].acceptedBid = true

			} else {

				state.teams["team2"].acceptedBid = true
				state.teams["team1"].acceptedBid = false

				//state.players[state.hostOpponent].acceptedBid = true
				//state.players[state.players[state.hostOpponent].teamMate].acceptedBid = true
			}

			TrumSuit.Turn = state.highestBidPerson //bidTurnSequence(state.players[state.highestBidPerson].seatPosition, state, logger)
			if Json_Turn, err := json.Marshal(TrumSuit); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(110, Json_Turn, nil, nil, true)

				state.currentTurn = TrumSuit.Turn

				logger.Info("sending 1111100000: ", state.currentTurn)

				presence := []runtime.Presence{state.players[TrumSuit.Turn].presence}

				//first turn so all card are allowed.
				if Json_Turn_Specifically, err := json.Marshal(TrumSuit); err != nil {
					logger.Info("Error is: ", err)
				} else {
					logger.Info("sending 1111111111: ", state.currentTurn)

					dispatcher.BroadcastMessage(111, Json_Turn_Specifically, presence, nil, true)

					if state.players[TrumSuit.Turn].bot {
						state.botOpcode = 406
						state.currentBot = TrumSuit.Turn
					} else {

						state.players[TrumSuit.Turn].response = true
						state.players[TrumSuit.Turn].opcode = 406
						state.currentBot = TrumSuit.Turn

					}

				}
				state.turnRoundCounter++
				state.gameRoundCounter++
				logger.Info("turnRoundCounter++  1111: ", state.turnRoundCounter)

			}

		}

		// opCOde = 0
		// opcd = 0
	case 406:

		logger.Info("406   card___________________: ", state.currentTurn)

		ThrownCrd := &Before_GamePlay{}

		if err := json.Unmarshal(data, &ThrownCrd); err != nil {
			logger.Info("Error is: ", err)
		} else {
			logger.Info("Data is: ", ThrownCrd)
			logger.Info("card send is: ", ThrownCrd.ThrownCard)

			state.players[ThrownCrd.PlayerUserName].response = false
			state.players[ThrownCrd.PlayerUserName].opcode = 0
			state.currentBot = ""

		}
		logger.Info("state.currentTurn is: ", state.currentTurn)
		logger.Info("ThrownCrd.PlayerUserName: ", ThrownCrd.PlayerUserName)

		if state.currentTurn == ThrownCrd.PlayerUserName {

			logger.Info("_____ state.gameRoundCounter_____________: ", state.gameRoundCounter)

			cardThrow(state, logger, ThrownCrd, dispatcher)

			logger.Info("_____ state.turnRoundCounter_____: ", state.turnRoundCounter)

		}

		// opCOde = 0
		// opcd = 0
	case 407:
		logger.Info("chaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaat ")

		Cht := &Before_GamePlay{}

		if err := json.Unmarshal(data, &Cht); err != nil {
			logger.Info("Error is: ", err)
		} else {
			logger.Info("Data is: ", Cht)
		}

		if Json_chat, err := json.Marshal(Cht); err != nil {
			logger.Info("Error is: ", err)
		} else {
			dispatcher.BroadcastMessage(116, Json_chat, nil, nil, true)
		}
		// opCOde = 0
		// opcd = 0
	case 600:
		logger.Info("Exxxxit: ")

		Ex := &Before_GamePlay{}
		Ex.Exit = "Exit"

		if Json_Exit, err := json.Marshal(Ex); err != nil {
			logger.Info("Error is: ", err)
		} else {
			dispatcher.BroadcastMessage(600, Json_Exit, nil, nil, true)

		}
		return nil

		// opCOde = 0
		// opcd = 0
	case 700:
		logger.Info(" ++++++++++++++++++++++++++++++++%++++++++ ")

		Ex := &Before_GamePlay{}
		Ex.Exit = "Exit"

		if Json_Exit, err := json.Marshal(Ex); err != nil {
			logger.Info("Error is: ", err)
		} else {
			dispatcher.BroadcastMessage(700, Json_Exit, nil, nil, true)

		}
	}
	opCOde = 0
	opcd = 0

	return state
}

func setBotJsonData(stateMain interface{}, logger runtime.Logger, jsonData *Before_GamePlay) *Before_GamePlay {
	state := stateMain.(*MatchState)

	if state.botOpcode == 402 {
		// rand.Seed(time.Now().UnixNano())
		// min := 0
		// max := (len(state.userNames) - 1)
		// teamMate := rand.Intn(max-min+1) + min
		jsonData.SetTeamMate = state.userNames[2]
	} else if state.botOpcode == 403 {
		jsonData.TeamMateResponse.X = state.nextBot
		jsonData.TeamMateResponse.Y = true
	} else if state.botOpcode == 404 {

		jsonData.PlayerUserName = state.currentBot
		if state.minBid < 13 {
			jsonData.BidValueString = strconv.Itoa(state.minBid + 1)
		}

	} else if state.botOpcode == 405 {

		jsonData.TrumpSuit = "club"
	} else if state.botOpcode == 406 {

		jsonData.PlayerUserName = state.currentBot
		jsonData.ThrownCard = ""

		logger.Info("state.currentPlayerTurnCard iss +++++++++ooo: ", state.currentPlayerTurnCard)
		logger.Info("state.players[jsonData.PlayerUserName].cards[i] iss +++++++++ooo: ", state.players[jsonData.PlayerUserName].cards)

		for i := 0; i < len(state.players[jsonData.PlayerUserName].cards) && jsonData.ThrownCard == ""; i++ {

			if len(state.currentPlayerTurnCard) > 0 && i < len(state.currentPlayerTurnCard) {
				jsonData.ThrownCard = state.currentPlayerTurnCard[i]
			} else {
				jsonData.ThrownCard = state.players[jsonData.PlayerUserName].cards[i]
			}

			logger.Info("boooOOOoooott Caarrdd iss +++++++++: ", jsonData.ThrownCard)
		}
		//jsonData.ThrownCard = state
	}

	// var Json_BotData []byte
	// if Json_BotData, err := json.Marshal(jsonData); err != nil {
	// 	logger.Info("Error is: ", err)
	// } else {
	// 	logger.Info("Bot Data in byte is: ", Json_BotData)
	// }
	return jsonData
}

func (t *Tarneeb) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	if state.(*MatchState).debug {
		logger.Info("match terminate match_id %v tick %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), tick)
		logger.Info("match terminate match_id %v grace seconds %v", ctx.Value(runtime.RUNTIME_CTX_MATCH_ID), graceSeconds)
	}
	logger.Info("___________Match MatchTerminate, starting")

	return state
}

func shuffle(deck [52]string, stateMain interface{}) [52]string {
	state := stateMain.(*MatchState)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	for i, card := range deck {
		state.deck[i] = card
	}
	return deck
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
	var emptyMap map[string]int
	emptyMap = make(map[string]int)

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
				player := &Player{uzr.GetDisplayName(), uzr.GetUsername(), uzr.GetAvatarUrl(), -1, "", "", 0, 0, arr, presence, 0, emptyMap, false, false, 0}
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

func randomValue(mix int, max int) int {
	rand.Seed(time.Now().UnixNano())
	min, max := 10, 30
	//fmt.Println(rand.Intn(max - min + 1) + min)

	randomValue := rand.Intn(max-min+1) + min
	return randomValue

}

func arrangePositions(stateMain interface{}, logger runtime.Logger) {
	state := stateMain.(*MatchState)

	i := 0

	state.sittingArangement[0] = state.hostUserName
	state.players[state.hostUserName].seatPosition = 0
	state.players[state.hostUserName].teamNum = "team1"

	state.sittingArangement[2] = state.players[state.hostUserName].teamMate
	state.players[state.players[state.hostUserName].teamMate].seatPosition = 2
	state.players[state.players[state.hostUserName].teamMate].teamNum = "team1"

	state.teams["team1"] = &Team{0, state.sittingArangement[0], state.sittingArangement[2], 0, 0, "", false}

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

			state.teams["team2"] = &Team{0, state.sittingArangement[1], state.sittingArangement[3], 0, 0, "", false}

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
	state.firstBiddingPerson = state.dealer
	logger.Info("Dealer is: ", state.dealer)
}

func bidTurnSequence(currentPosition int, stateMain interface{}, logger runtime.Logger) string {
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

func cardDealing(stateMain interface{}, logger runtime.Logger) {

	state := stateMain.(*MatchState)
	i := 0
	for _, card := range state.deck {

		if i < len(state.sittingArangement[i]) {
			state.players[state.sittingArangement[i]].cards = append(state.players[state.sittingArangement[i]].cards, card)
			state.players[state.sittingArangement[i]].cardWithValue[card] = state.cardWithValue[card]

		}
		if i == 3 {
			i = 0
		} else {
			i++
		}

	}

	logger.Info("cardDealing_______________________")

	// for _, plyr := range state.players {
	// 	logger.Info(" state.players is__: ", len(state.players[plyr.userName].cardWithValue))

	// 	//state.players[plyr.userName].cardWithValue = nil
	// 	// for _, card := range state.players[plyr.userName].cards {

	// 	// 	//logger.Info(" cards are_____: ", card)

	// 	// }

	// }

}

func resetingCardValue(stateMain interface{}, logger runtime.Logger) {
	state := stateMain.(*MatchState)

	for _, plyr := range state.players {
		logger.Info("1ssstt looOOOOooopp state.players is__: ", len(state.players[plyr.userName].cardWithValue))

		for i, _ := range state.players[plyr.userName].cardWithValue {
			state.players[plyr.userName].cardWithValue[i] = state.cardWithValue[i]
		}
	}

}

func cardWithValue(stateMain interface{}, logger runtime.Logger) {

	state := stateMain.(*MatchState)
	for _, plyr := range state.players {

		state.players[plyr.userName].cardWithValue = make(map[string]int)

	}

	state.cardWithValue = nil
	state.cardWithValue = make(map[string]int)

	state.cardWithValue = map[string]int{"club_A": 14, "club_2": 2, "club_3": 3, "club_4": 4, "club_5": 5, "club_6": 6, "club_7": 7,
		"club_8": 8, "club_9": 9, "club_10": 10, "club_J": 11, "club_Q": 12, "club_K": 13, "diamond_A": 14,
		"diamond_2": 2, "diamond_3": 3, "diamond_4": 4, "diamond_5": 5, "diamond_6": 6, "diamond_7": 7,
		"diamond_8": 8, "diamond_9": 9, "diamond_10": 10, "diamond_J": 11, "diamond_Q": 12, "diamond_K": 13,
		"heart_A": 14, "heart_2": 2, "heart_3": 3, "heart_4": 4, "heart_5": 5, "heart_6": 6, "heart_7": 7,
		"heart_8": 8, "heart_9": 9, "heart_10": 10, "heart_J": 11, "heart_Q": 12, "heart_K": 13, "spade_A": 14,
		"spade_2": 2, "spade_3": 3, "spade_4": 4, "spade_5": 5, "spade_6": 6, "spade_7": 7, "spade_8": 8,
		"spade_9": 9, "spade_10": 10, "spade_J": 11, "spade_Q": 12, "spade_K": 13}
}

//Uzair's code
func updateCardValues(stateMain interface{}, logger runtime.Logger, incrementalValue int, prioritySign string) {
	state := stateMain.(*MatchState)
	logger.Info("updateCardValues is__: ", state.sign)

	for _, plyr := range state.players {
		logger.Info("1ssstt looOOOOooopp state.players is__: ", len(state.players[plyr.userName].cardWithValue))

		for i, cardValue := range state.players[plyr.userName].cardWithValue {
			if strings.Contains(i, prioritySign) {

				state.players[plyr.userName].cardWithValue[i] = cardValue + incrementalValue

			}
			////logger.Info("card value____is: ", state.players[plyr.userName].cardWithValue[i])

		}
	}

}

func nillingPlayerCards(stateMain interface{}, logger runtime.Logger) {
	state := stateMain.(*MatchState)

	for _, plyr := range state.players {
		logger.Info("1ssstt looOOOOooopp state.players is__: ", len(state.players[plyr.userName].cardWithValue))

		state.players[plyr.userName].cardWithValue = nil
		state.players[plyr.userName].cards = nil

	}

}
func cardThrow(stateMain interface{}, logger runtime.Logger, obj *Before_GamePlay, dispatcher runtime.MatchDispatcher) {
	state := stateMain.(*MatchState)

	obj.Cards = nil
	//var sign string
	logger.Info("state.sign is: ", state.sign)
	logger.Info("_____ state.turnRoundCounter_____: ", state.turnRoundCounter)

	if state.turnRoundCounter <= 1 {

		if strings.Contains(obj.ThrownCard, "club") {
			state.sign = "club"
		} else if strings.Contains(obj.ThrownCard, "diamond") {
			state.sign = "diamond"
		} else if strings.Contains(obj.ThrownCard, "heart") {
			state.sign = "heart"
		} else if strings.Contains(obj.ThrownCard, "spade") {
			state.sign = "spade"
		}

		updateCardValues(state, logger, 20, state.sign)

	} else {
	}
	cardComparis(state, logger, obj)

	logger.Info("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@state.firstTurnsCounter After card compar: ", state.turnRoundCounter)

	// if state.turnRoundCounter <= 4 {
	// 	obj.Turn = bidTurnSequence(state.players[obj.PlayerUserName].seatPosition, state, logger)

	// }

	//hitting card removing 112 before sending last card information (110) ----> now shifting after last card information (110)

	// if state.turnRoundCounter == 4 {

	// 	//	logger.Info("_if state.turnRoundCounter == 4___: ", state.turnRoundCounter)

	// 	state.turnRoundCounter = 0
	// 	state.firstTurnsCounter = 0
	// 	state.players[state.highestCardPlayer].score++

	// 	resetingCardValue(state, logger)                     //cardWithValue(state, logger)                         // reseting the value of all cards, actually for reseting the first selected card of new turn.
	// 	updateCardValues(state, logger, 40, state.trumpSuit) // For keeping back the value of trumpSuit cards after reseting the all cards.

	// 	HighestCrdPrsn := &Before_GamePlay{}
	// 	HighestCrdPrsn.HighestCardPerson = state.highestCardPlayer
	// 	HighestCrdPrsn.HighestCardPersonScore = state.players[state.highestCardPlayer].score
	// 	state.sign = ""

	// 	if Json_HighestCrdPrsn, err := json.Marshal(HighestCrdPrsn); err != nil {
	// 		logger.Info("Error is: ", err)
	// 	} else {

	// 		logger.Info(" ________sending 112: ", HighestCrdPrsn)

	// 		time.Sleep(10 * time.Second)

	// 		dispatcher.BroadcastMessage(112, Json_HighestCrdPrsn, nil, nil, true) //change 113 to 112 later after score testing
	// 		logger.Info("1111222 is:_________________delayyyy___________________ ")
	// 		time.Sleep(10 * time.Second)

	// 		//state.firstTurnsCounter++
	// 	}

	// }     //hitting card removing 112 before sending last card information (110) ----> now shifting after last card information (110)

	obj.Cards = nil //to hide card broadcasting on other apps.
	state.currentPlayerTurnCard = nil
	obj.Turn = state.nextTurn

	if Json_ThrownCrd, err := json.Marshal(obj); err != nil {
		logger.Info("Error is: ", err)
	} else {

		logger.Info("110 data is: ", state.players[obj.Turn].cards)

		if state.players[obj.PlayerUserName].bot {
			time.Sleep(2 * time.Second)
		}
		dispatcher.BroadcastMessage(110, Json_ThrownCrd, nil, nil, true)

		/////////////////////////////////////////

		if state.turnRoundCounter == 4 {

			//	logger.Info("_if state.turnRoundCounter == 4___: ", state.turnRoundCounter)

			state.turnRoundCounter = 0
			state.firstTurnsCounter = 0
			state.players[state.highestCardPlayer].score++

			resetingCardValue(state, logger)                     //cardWithValue(state, logger)                         // reseting the value of all cards, actually for reseting the first selected card of new turn.
			updateCardValues(state, logger, 40, state.trumpSuit) // For keeping back the value of trumpSuit cards after reseting the all cards.

			HighestCrdPrsn := &Before_GamePlay{}
			HighestCrdPrsn.HighestCardPerson = state.highestCardPlayer
			HighestCrdPrsn.HighestCardPersonScore = state.players[state.highestCardPlayer].score
			state.sign = ""

			if Json_HighestCrdPrsn, err := json.Marshal(HighestCrdPrsn); err != nil {
				logger.Info("Error is: ", err)
			} else {

				logger.Info(" ________sending 112: ", HighestCrdPrsn)

				//time.Sleep(2 * time.Second)

				dispatcher.BroadcastMessage(112, Json_HighestCrdPrsn, nil, nil, true) //change 113 to 112 later after score testing
				logger.Info("1111222 is:_________________delayyyy___________________ ")
				time.Sleep(2 * time.Second)

				//state.firstTurnsCounter++
			}

		} ////////////////////////////////////////////////////////

		for _, allowCarddata := range state.players[obj.Turn].cards {
			////logger.Info("data cards are: ", data)

			if strings.Contains(allowCarddata, state.sign) {
				obj.Cards = append(obj.Cards, allowCarddata)
				state.currentPlayerTurnCard = append(state.currentPlayerTurnCard, allowCarddata)
			}
		}
	}

	if state.gameRoundCounter < 52 {

		//state.currentTurn = obj.Turn

		logger.Info("PLayer cards are: ", obj.Cards)

		presence := []runtime.Presence{state.players[obj.Turn].presence}

		if Json_ThrownCrd, err := json.Marshal(obj); err != nil {
			logger.Info("Error is: ", err)
		} else {
			logger.Info("111 data is: ", obj)

			dispatcher.BroadcastMessage(111, Json_ThrownCrd, presence, nil, true)
			if state.players[obj.Turn].bot {
				state.botOpcode = 406
				state.currentBot = obj.Turn
			} else {

				state.players[obj.Turn].response = true
				state.players[obj.Turn].opcode = 406
				state.currentBot = obj.Turn

			}
			logger.Info("____turnRoundCounter: ", state.turnRoundCounter)
			state.currentTurn = obj.Turn

			state.turnRoundCounter++
			state.gameRoundCounter++
		}
		//}
	} else {

		logger.Info("$$$$$$$state.highest bid____________: ", state.highestBid)
		logger.Info("$$$$$$$state.highest bidr____________: ", state.highestBidPerson)

		logger.Info("_____state.turnRoundCounter____________: ", state.turnRoundCounter)
		logger.Info("_____state.RoundScore__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].roundScore)
		logger.Info("_____state.biiiiiiiiiiiiiiiiiiiiiiiiiiiddddddddddddddddddddddddd__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].teamBid)
		logger.Info("_____state.bTeaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaamscorrrrrrrrrrrrrrrrre__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].teamScore)
		logger.Info("_____state.teamScore__state.teams['team1'].acceptedBid ", state.teams["team1"].acceptedBid)

		logger.Info("_________________++++++++++++++++++++++++++++++++++++___________________________")

		logger.Info("_____state.RoundScore__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
		logger.Info("_____state.RoundScore__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
		logger.Info("_____state.biiiiiiiiiiiiiiiiiiiiiiiiiiiddddddddddddddddddddddddd__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].teamBid)
		logger.Info("_____state.bTeaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaamscorrrrrrrrrrrrrrrrre_______HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].teamScore)
		logger.Info("_____state.teamScore__state.teams['team2'].acceptedBid ", state.teams["team2"].acceptedBid)

		state.teams[state.players[state.hostUserName].teamNum].teamScore = state.players[state.hostUserName].score + state.players[state.players[state.hostUserName].teamMate].score
		state.teams[state.players[state.hostOpponent].teamNum].teamScore = state.players[state.hostOpponent].score + state.players[state.players[state.hostOpponent].teamMate].score

		if state.teams[state.players[state.hostUserName].teamNum].teamScore < state.highestBid && state.highestBid != 13 && state.teams["team1"].acceptedBid {
			state.teams[state.players[state.hostUserName].teamNum].roundScore = state.teams[state.players[state.hostUserName].teamNum].roundScore + (state.teams[state.players[state.hostUserName].teamNum].teamScore - state.highestBid)
			// first_score := state.teams[state.players[state.hostUserName].teamNum].roundScore
			logger.Info("players[state.hostUserName].teamNum].teamScore < state.teams[state.players[sta: ")

		} else if state.highestBid == 13 && state.teams[state.players[state.hostUserName].teamNum].teamScore == 13 && state.teams["team1"].acceptedBid {
			state.teams[state.players[state.hostUserName].teamNum].roundScore = 26
			logger.Info("kaboot is: ")

			kab := &Before_GamePlay{}
			kab.Kaboot = "kaboot"
			if Json_kaboot, err := json.Marshal(kab); err != nil {
				logger.Info("Error is: ", err)
			} else {

				dispatcher.BroadcastMessage(114, Json_kaboot, nil, nil, true)
				//state.firstTurnsCounter++
			}
		} else if state.highestBid == 13 && state.teams[state.players[state.hostUserName].teamNum].teamScore != 13 && state.teams["team1"].acceptedBid {
			state.teams[state.players[state.hostUserName].teamNum].roundScore = state.teams[state.players[state.hostUserName].teamNum].roundScore - 16

		} else {

			logger.Info(" elseeee 1111 ", state.teams[state.players[state.hostUserName].teamNum].roundScore)

			state.teams[state.players[state.hostUserName].teamNum].roundScore = state.teams[state.players[state.hostUserName].teamNum].roundScore + state.teams[state.players[state.hostUserName].teamNum].teamScore
		}

		//ooooooOOOOOOOooooopppoooOOOOooonent
		if state.teams[state.players[state.hostOpponent].teamNum].teamScore < state.highestBid && state.highestBid != 13 && state.teams["team2"].acceptedBid {
			state.teams[state.players[state.hostOpponent].teamNum].roundScore = state.teams[state.players[state.hostOpponent].teamNum].roundScore + (state.teams[state.players[state.hostOpponent].teamNum].teamScore - state.highestBid)
			logger.Info("state.teams[state.players[state.hostOpponent].teamNum].roundScore ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)

		} else if state.teams[state.players[state.hostOpponent].teamNum].teamScore == 13 && state.highestBid == 13 && state.teams["team2"].acceptedBid {
			state.teams[state.players[state.hostOpponent].teamNum].roundScore = 26
			logger.Info("state.teams[state.players[state.hostOpponent].teamNum].roundScore = 26 ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
			//kaboot
			logger.Info("kaboot is: ")
			kab := &Before_GamePlay{}
			kab.Kaboot = "kaboot"
			if Json_kaboot, err := json.Marshal(kab); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(114, Json_kaboot, nil, nil, true)
			}

		} else if state.highestBid == 13 && state.teams[state.players[state.hostOpponent].teamNum].teamScore != 13 && state.teams["team2"].acceptedBid {
			state.teams[state.players[state.hostOpponent].teamNum].roundScore = state.teams[state.players[state.hostOpponent].teamNum].roundScore - 16
			logger.Info(" else if_________________________--state.teams[state.players[state.hostOpponent].teamNum].teamBid == 13", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
		} else {
			state.teams[state.players[state.hostOpponent].teamNum].roundScore = state.teams[state.players[state.hostOpponent].teamNum].roundScore + state.teams[state.players[state.hostOpponent].teamNum].teamScore
			logger.Info("else_______________________-----------------------------------s state.teams[state.players[state.hostOpponent].teamNum].roundScore ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
		}
		logger.Info("_____state.teamScore__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].teamScore)
		logger.Info("_____state.RoundScore__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].roundScore)

		logger.Info("______________________________________________________________________")

		logger.Info("_____state.teamScore__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].teamScore)
		logger.Info("_____state.RoundScore__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)

		Score := &Before_GamePlay{}

		Score.Team1 = state.teams[state.players[state.hostUserName].teamNum].roundScore
		Score.Team2 = state.teams[state.players[state.hostOpponent].teamNum].roundScore

		if Json_Score, err := json.Marshal(Score); err != nil {
			logger.Info("Error is: ", err)
		} else {

			dispatcher.BroadcastMessage(113, Json_Score, nil, nil, true)
			time.Sleep(4 * time.Second)

		}
		state.players[state.hostUserName].score = 0
		state.players[state.hostOpponent].score = 0
		state.players[state.players[state.hostUserName].teamMate].score = 0
		state.players[state.players[state.hostOpponent].teamMate].score = 0

		if Score.Team1 >= 31 && Score.Team1 > Score.Team2 {
			//Scor := &Before_GamePlay{}
			Score.Score_Team1 = "Score_Team1"
			Score.Winner1 = state.teams[state.players[state.hostUserName].teamNum].member1
			Score.Winner2 = state.teams[state.players[state.hostUserName].teamNum].member2

			if Json_Score_Team1, err := json.Marshal(Score); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(115, Json_Score_Team1, nil, nil, true)
			}
		} else if Score.Team2 >= 31 && Score.Team2 > Score.Team1 {
			Scor := &Before_GamePlay{}
			Scor.Score_Team2 = "Score_Team2"
			Scor.Winner1 = state.teams[state.players[state.hostOpponent].teamNum].member1
			Scor.Winner2 = state.teams[state.players[state.hostOpponent].teamNum].member2
			if Json_Score_Team2, err := json.Marshal(Scor); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(115, Json_Score_Team2, nil, nil, true)
				//state.firstTurnsCounter++
			}
		} else {

			logger.Info("$$$$$$$state.highest bid____________: ", state.highestBid)
			logger.Info("$$$$$$$state.highest bidr____________: ", state.highestBidPerson)

			logger.Info("$$$$$$$state.turnRoundCounter____________: ", state.turnRoundCounter)
			logger.Info("$$$$$$$$$$state.RoundScore__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].roundScore)
			logger.Info("$$$$$$$$$state.biiiiiiiiiiiiiiiiiiiiiiiiiiiddddddddddddddddddddddddd__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].teamBid)
			logger.Info("$$$$$$$$$state.bTeaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaamscorrrrrrrrrrrrrrrrre__HostUsername ", state.teams[state.players[state.hostUserName].teamNum].teamScore)
			logger.Info("_$$$$$$$state.teamScore__state.teams['team1'].acceptedBid ", state.teams["team1"].acceptedBid)

			logger.Info("___________$$$$$$$++++++++++++++++++++++++++++++++++++___________________________")

			//logger.Info("_____state.RoundScore__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
			logger.Info("$$$$$$$$state.RoundScore__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].roundScore)
			logger.Info("$$$$$$$state.biiiiiiiiiiiiiiiiiiiiiiiiiiiddddddddddddddddddddddddd__HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].teamBid)
			logger.Info("$$$$$$state.bTeaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaamscorrrrrrrrrrrrrrrrre_______HostOpponent ", state.teams[state.players[state.hostOpponent].teamNum].teamScore)
			logger.Info("$$$$$$$$state.teamScore__state.teams['team2'].acceptedBid ", state.teams["team2"].acceptedBid)

			Score.Team1 = 0
			Score.Team2 = 0
			GameNextRound := &Before_GamePlay{}

			state.bidRoundCounter = 0
			state.gameRoundCounter = 0
			//state.firstTurnsCounter = 0
			state.turnRoundCounter = 0
			state.highestCardValue = 0
			state.highestCard = ""
			state.highestCardPlayer = ""
			state.minBid = 7
			state.highestBid = 0
			state.bidPassCounter = 0

			nillingPlayerCards(state, logger)
			//state.deck = nil
			state.deck = [52]string{"club_A", "club_2", "club_3", "club_4", "club_5", "club_6", "club_7", "club_8", "club_9", "club_10", "club_J", "club_Q", "club_K", "diamond_A", "diamond_2", "diamond_3", "diamond_4", "diamond_5", "diamond_6", "diamond_7", "diamond_8", "diamond_9", "diamond_10", "diamond_J", "diamond_Q", "diamond_K", "heart_A", "heart_2", "heart_3", "heart_4", "heart_5", "heart_6", "heart_7", "heart_8", "heart_9", "heart_10", "heart_J", "heart_Q", "heart_K", "spade_A", "spade_2", "spade_3", "spade_4", "spade_5", "spade_6", "spade_7", "spade_8", "spade_9", "spade_10", "spade_J", "spade_Q", "spade_K"}

			shuffle(state.deck, state)
			cardWithValue(state, logger)
			B4GamePlay := &Before_GamePlay{}
			B4GamePlay.Dealer = state.dealer

			if Json_Dealer, err := json.Marshal(B4GamePlay); err != nil {
				logger.Info("Error is: ", err)
			} else {
				dispatcher.BroadcastMessage(105, Json_Dealer, nil, nil, true)
			}
			//cardDealing(state, logger)
			afterSitting(state, logger, GameNextRound, dispatcher)

		}
	}
}

func cardComparis(stateMain interface{}, logger runtime.Logger, obj *Before_GamePlay) {
	state := stateMain.(*MatchState)
	logger.Info("cardComparis__ is: ", obj)

	for i, card := range state.players[obj.PlayerUserName].cards {
		//logger.Info("card__ is: ", card)

		if card == obj.ThrownCard {
			logger.Info("state.highestCardValue: ", state.highestCardValue)

			logger.Info("state.players[obj.PlayerUserName].cardWithValue[c is: ", state.players[obj.PlayerUserName].cardWithValue[card])

			if state.players[obj.PlayerUserName].cardWithValue[card] > state.highestCardValue {
				logger.Info("state.players[obj.PlayerUserName].cardWithValue[c is: ", obj)

				state.highestCardValue = state.players[obj.PlayerUserName].cardWithValue[card]
				state.highestCardPlayer = obj.PlayerUserName
				state.highestCard = obj.ThrownCard
				state.nextTurn = bidTurnSequence(state.players[obj.PlayerUserName].seatPosition, state, logger)

			} else {
				state.nextTurn = bidTurnSequence(state.players[obj.PlayerUserName].seatPosition, state, logger)
			}
			remove(state.players[obj.PlayerUserName].cards, i)

			//logger.Info("NOaw Card is______: ", state.players[obj.PlayerUserName].cards[i])
			break

		}
	}

	if state.turnRoundCounter == 4 {
		logger.Info("cleeearrinng_____________________________: ")
		state.nextTurn = state.highestCardPlayer
		state.highestCardValue = 0
	}

	logger.Info("state.nextTurn is: ", state.nextTurn)

}

func afterSitting(stateMain interface{}, logger runtime.Logger, B4GamePlay *Before_GamePlay, dispatcher runtime.MatchDispatcher) {
	state := stateMain.(*MatchState)

	cardDealing(state, logger)
	for _, plyr := range state.players {

		for _, Card := range plyr.cards {

			B4GamePlay.Cards = append(B4GamePlay.Cards, Card)
			//logger.Info("card Player is: ", plyr)

		}
		//logger.Info("card Player is: ", plyr)
		//logger.Info("card Player cards are: ", B4GamePlay.Cards)

		B4GamePlay.PlayerUserName = plyr.userName
		if Json_PlayerCards, err := json.Marshal(B4GamePlay); err != nil {
			logger.Info("Error is: ", err)
		} else {
			presence := []runtime.Presence{plyr.presence}
			dispatcher.BroadcastMessage(106, Json_PlayerCards, presence, nil, true)
		}
		B4GamePlay.Cards = nil
	}

	state.bidFlag = true
	B4GamePlay.BidValue = -1
	B4GamePlay.BidPass = true
	B4GamePlay.Turn = bidTurnSequence(state.players[state.firstBiddingPerson].seatPosition, state, logger)
	//state.firstBiddingPerson = B4GamePlay.Turn
	state.dealer = state.firstBiddingPerson
	if Json_Turn, err := json.Marshal(B4GamePlay); err != nil {
		logger.Info("Error is: ", err)
	} else {
		logger.Info("11111ooooooOOOOOOoooooo7777___aFTERsTTING ", err)

		dispatcher.BroadcastMessage(107, Json_Turn, nil, nil, true)
		//state.bidRoundCounter++

		if state.players[B4GamePlay.Turn].bot {
			state.currentBot = B4GamePlay.Turn
			state.botOpcode = 404
		} else {
			state.players[B4GamePlay.Turn].response = true
			state.players[B4GamePlay.Turn].opcode = 404
			state.currentBot = B4GamePlay.Turn
		}

	}
}

func remove1(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func remove(s []string, i int) {
	s[i] = ""
}
