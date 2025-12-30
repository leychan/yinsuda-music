package client

// --- Models ---

type PlayListInfo struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"` // ImgUrl in example, doc says url
	ImgUrl      string `json:"imgUrl"` // Doc uses both names in different places
	Status      int    `json:"status"` // 0-Unavailable, 1-Available
}

type PlayListDetail struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	ImgUrl      string `json:"imgUrl"`
	SongList    []struct {
		SongId string `json:"songId"`
	} `json:"songList"`
}

// --- Requests & Responses ---

type PageRequest struct {
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type QuerySongListResponse struct {
	Total int            `json:"total"`
	List  []PlayListInfo `json:"list"`
}

type QuerySongListDetailRequest struct {
	Code string `json:"code"`
}

type QuerySongListDetailResponse PlayListDetail

// --- Methods ---

// 5. QuerySongListPage

func (c *Client) QuerySongListPage(req *PageRequest) (*QuerySongListResponse, error) {
	var result struct {
		Data QuerySongListResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/querySongListPage", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// 6. QuerySongListDetail

func (c *Client) QuerySongListDetail(code string) (*QuerySongListDetailResponse, error) {
	req := &QuerySongListDetailRequest{Code: code}
	var result struct {
		Data QuerySongListDetailResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/querySongListDetail", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// 7. QueryRankingListPage
// Logic is identical to QuerySongListPage

func (c *Client) QueryRankingListPage(req *PageRequest) (*QuerySongListResponse, error) {
	var result struct {
		Data QuerySongListResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/queryRankingListPage", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// 8. QueryRankingListDetail
// Logic is identical to QuerySongListDetail

func (c *Client) QueryRankingListDetail(code string) (*QuerySongListDetailResponse, error) {
	req := &QuerySongListDetailRequest{Code: code}
	var result struct {
		Data QuerySongListDetailResponse `json:"data"`
	}
	err := c.Do("POST", "/mcrc-sas/yinsuda/queryRankingListDetail", nil, req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
