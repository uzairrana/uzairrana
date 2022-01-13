package Authentication

type UnlockAvatar struct {
	ARRAY []BundleDetails `json:"Bundles"`
}

type UnlockBundle struct {
	ARRAY []BundleDetails `json:"Bundles"`
}

type BundleDetails struct {
	BundleName string
}

type EmailAuthentication struct {
	EMAIL      string `string:"EMAIL"`
	PASSWORD   string `string:"PASSWORD"`
	DOB        string `string:"DOB"`
	FIRST_NAME string `string:"FIRST_NAME"`
	LAST_NAME  string `string:"LAST_NAME"`
}

type FacebookAuthentication struct {
	USERNAME string `string:"USERNAME"`
}

// type Match struct {
// 	MatchId          string
// 	Inprogress       string
// 	MinBuyIn         string
// 	MaxBuyIn         string
// 	SmallBlindAmount string
// 	BigBlindAmount   string
// 	NumberOfPlayers  string
// 	MatchCreatedTime string
// 	HostName         string
// }

type IPM struct {
	InProgressMatches []string `json:"InProgressMatches"`
}
type EM struct {
	EndedMatches []string `json:"EndedMatches"`
}

type CM struct {
	ClubMatches []string `json:"ClubMatches"`
}

//In Progress or not, Min/Max Buy in, Small/Big Blind Total Players
