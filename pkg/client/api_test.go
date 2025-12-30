package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_SongAPIs(t *testing.T) {
	mux := http.NewServeMux()
	
	// Mock Auth (Reuse logic from client_test.go if we were refactoring, but for now duplicate concise version)
	mux.HandleFunc("/oauth2/token", func(w http.ResponseWriter, r *http.Request) {
		resp := TokenResponse{Code: "0", Success: true, Data: TokenData{AccessToken: "token", Expire: 3600}}
		json.NewEncoder(w).Encode(resp)
	})

	// Mock GetSongList
	mux.HandleFunc("/mcrc-sas/yinsuda/getSongList", func(w http.ResponseWriter, r *http.Request) {
		var req GetSongListRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.Limit != 100 {
			t.Errorf("Expected default limit 100, got %d", req.Limit)
		}

		resp := BaseResponse{
			Code: 0, Success: true,
			Data: GetSongListResponse{
				NextQueryInfo: "END",
				SongList: []Song{{SongId: "S1", SongName: "Name1"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Mock GetSongInfo
	mux.HandleFunc("/mcrc-sas/yinsuda/getSongInfo", func(w http.ResponseWriter, r *http.Request) {
		var req GetSongInfoRequest
		json.NewDecoder(r.Body).Decode(&req)
		
		if req.SongIdListStr != "S1,S2" {
			t.Errorf("Expected S1,S2, got %s", req.SongIdListStr)
		}

		resp := BaseResponse{
			Code: 0, Success: true,
			Data: GetSongInfoResponse{
				SongList: []Song{{SongId: "S1"}, {SongId: "S2"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Mock Playlist
	mux.HandleFunc("/mcrc-sas/yinsuda/querySongListPage", func(w http.ResponseWriter, r *http.Request) {
		resp := BaseResponse{
			Code: 0, Success: true,
			Data: QuerySongListResponse{
				Total: 10,
				List: []PlayListInfo{{Code: "P1", Title: "Playlist1"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("appId", "secret", server.URL)

	// Test GetSongList
	listResp, err := client.GetSongList(&GetSongListRequest{Limit: 100})
	if err != nil {
		t.Fatalf("GetSongList failed: %v", err)
	}
	if len(listResp.SongList) != 1 || listResp.SongList[0].SongId != "S1" {
		t.Errorf("GetSongList unexpected result")
	}

	// Test GetSongInfo
	infoResp, err := client.GetSongInfo([]string{"S1", "S2"})
	if err != nil {
		t.Fatalf("GetSongInfo failed: %v", err)
	}
	if len(infoResp.SongList) != 2 {
		t.Errorf("GetSongInfo unexpected count")
	}

	// Test Playlist
	plResp, err := client.QuerySongListPage(&PageRequest{Offset: 0, Length: 10})
	if err != nil {
		t.Fatalf("QuerySongListPage failed: %v", err)
	}
	if plResp.Total != 10 {
		t.Errorf("Playlist total mismatch")
	}
}
