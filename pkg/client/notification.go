package client

// NotificationType constants
const (
	NotifyTypeSong        = "SONG"
	NotifyTypeSongList    = "SONG_LIST"
	NotifyTypeRankingList = "RANKING_LIST"
)

// Notification represents the common structure of the notification payload.
// Note: Depending on NotifyType, either Songs or Codes will be populated.
// Users can unmarshal into this struct or specific ones.
type Notification struct {
	NotifyId   string       `json:"notifyId"`
	NotifyTime string       `json:"notifyTime"`
	NotifyType string       `json:"notifyType"`
	AppId      interface{}  `json:"appId"` // Doc says Long, example shows string/long? Safest to use json.Number or interface{} then cast, or string if quoted. Example says "1009232".
	Songs      []SongChange `json:"songs,omitempty"` // For SONG type
	Codes      []string     `json:"codes,omitempty"` // For SONG_LIST/RANKING_LIST type
}

type SongChange struct {
	SongId     string `json:"songId"`
	ChangeDate string `json:"changeDate"`
}

// NotificationResponse is the response the client must send back.
type NotificationResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
