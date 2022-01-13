package Leaderboard

type LeaderboardSubmitProps struct {
	LeaderboardId string `json:"leaderboardId"`
	UserId        string `json:"userId"`
	// Id  string    `json:"id"`
	Score    float32  `json:"score"`
	UserName string `json:"username"`
}
