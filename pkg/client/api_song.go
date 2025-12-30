package client

import (
	"strings"
)

// --- Models ---

type ImagePathMap struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Album struct {
	AlbumId          string         `json:"albumId"`
	AlbumName        string         `json:"albumName"`
	ImagePathMapList []ImagePathMap `json:"imagePathMapList"`
}

type Artist struct {
	ArtistId   string `json:"artistId"`
	ArtistName string `json:"artistName"`
}

type Lrc struct {
	Type string `json:"type"` // txt, lrc, qrc
	Url  string `json:"url"`
}

type Copyright struct {
	SceneId        string   `json:"sceneId"`
	TerminalIdList []string `json:"terminalIdList"`
}

type Song struct {
	SongId         string      `json:"songId"`
	SongName       string      `json:"songName"`
	CompanyName    string      `json:"companyName"`
	PublicTime     string      `json:"publicTime"`
	Version        string      `json:"version"`
	Duration       int         `json:"duration"` // Seconds
	Status         int         `json:"status"`   // 1-Available, 0-Unavailable
	GrantStatus    int         `json:"grantStatus"`
	TakeDownReason string      `json:"takeDownReason"`
	GrantStartTime string      `json:"grantStartTime"`
	Sequence       int         `json:"sequence"`
	Language       string      `json:"language"`
	PitchUrl       string      `json:"pitchUrl"`
	ChorusStartMS  int         `json:"chorusStartMS"`
	ChorusEndMS    int         `json:"chorusEndMS"`
	CopyrightList  []Copyright `json:"copyrightList"`
	Album          Album       `json:"album"`
	ArtistList     []Artist    `json:"artistList"`
	LrcList        []Lrc       `json:"lrcList"`
}

type MediaInfo struct {
	FileType      string `json:"fileType"`
	Complete      int    `json:"complete"`
	Url           string `json:"url"`
	Expire        string `json:"expire"`
	StartSecond   string `json:"startSecond"`
	EndSecond     string `json:"endSecond"`
}

// --- Requests & Responses ---

// 1. GetSongList

type GetSongListRequest struct {
	QueryInfo string `json:"queryInfo"` // Empty for start, "END" for finished
	Limit     int    `json:"limit,omitempty"`
	SearchText string `json:"searchText,omitempty"`
	Offset int `json:"offset,omitempty"`
}

type GetSongListResponse struct {
	NextQueryInfo string `json:"nextQueryInfo"`
	SongList      []Song `json:"songList"`
}

// 2. GetSongInfo

type GetSongInfoRequest struct {
	SongIdListStr string `json:"songIdListStr"`
}

type GetSongInfoResponse struct {
	SongList []Song `json:"songList"`
}

// 3. GetSongUrl

type GetSongUrlRequest struct {
	SongId     string `json:"songId"`
	IdentityId string `json:"identityId,omitempty"`
}

type GetSongUrlResponse struct {
	MediaList []MediaInfo `json:"mediaList"`
}

// 4. SearchSong

type SearchSongRequest struct {
	SearchText string `json:"searchText"`
	SearchType int    `json:"searchType"` // 1: Full, 2: SongName, 3: ArtistName
	Status     int    `json:"status,omitempty"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

type SearchSongResponse struct {
	Total    int    `json:"total"`
	SongList []Song `json:"songList"`
}

// --- Methods ---

func (c *Client) GetSongList(req *GetSongListRequest) (*GetSongListResponse, error) {
	if req.Limit == 0 {
		req.Limit = 100
	}
	var result struct {
		Data GetSongListResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/getSongList", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) GetSongInfo(songIds []string) (*GetSongInfoResponse, error) {
	req := &GetSongInfoRequest{
		SongIdListStr: strings.Join(songIds, ","),
	}
	var result struct {
		Data GetSongInfoResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/getSongInfo", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) GetSongUrl(req *GetSongUrlRequest) (*GetSongUrlResponse, error) {
	var result struct {
		Data GetSongUrlResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/getSongUrl", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (c *Client) SearchSong(req *SearchSongRequest) (*SearchSongResponse, error) {
	var result struct {
		Data SearchSongResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/searchSong", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
