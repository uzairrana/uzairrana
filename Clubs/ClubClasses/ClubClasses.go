package ClubClasses

type ClubProps struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Desc             string `json:"desc"`
	Flag             string `json:"flag"`
	OwnerID          string `json:"ownerid"`
	ClubId           string `json:"clubid"`
	Level            string `json:"level"`
	TotalMembers     string `json:"totalmembers"`
	LevelProgressBar string `json:"levelprogressbar"`
	ClubPictureUrl   string `json:"url"`
}

type ClubJoinProps struct {
	ClubId string `json:"clubid"`
	UserId string `json:"userid"`
}
type ClubCreatePostClass struct {
	ClubId string `json:"clubid"`
	UserId string `json:"userid"`
	Desc   string `json:"desc"`
	Data   string `json:"data"`
}

type ClubPropsIds struct {
	ClubId     string `json:"clubid"`
	UserId     string `json:"userid"`
	Permission string `json:"permission"`
}
type ClubLeaveProps struct {
	ClubId   string `json:"clubid"`
	UserId   string `json:"userid"`
	Username string `json:"username"`
}
type ClubPosts struct {
	ClubId string `json:"clubid"`
	UserId string `json:"userid"`
	PostId string `json:"postid"`
}
type ClubMangersMetadata struct {
	Name     string   `json:"name"`
	ClubId   string   `json:"clubid"`
	Managers []string `json:"managers"`
	UserId   string   `json:"userId"`
}
type PostClass struct {
	PostsIDs []string `json:"postids"`
}
type ClubManagers struct {
	Managers []string `json:"managers"`
}
