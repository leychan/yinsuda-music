package main

import (
	"fmt"
	"os"

	"github.com/leychan/yinsuda-music/pkg/client"
)

func main() {
	appId := os.Getenv("yinsuda_appid")
	appSecret := os.Getenv("yinsuda_appsecret")
	host := os.Getenv("yinsuda_host")

	if appId == "" || appSecret == "" || host == "" {
		fmt.Println("Usage: APP_ID=... APP_SECRET=... API_HOST=... go run cmd/demo/main.go")
		fmt.Println("Note: API_HOST should be the base URL, verifying with your API provider (e.g., https://test-api.yinsuda.com)")
		os.Exit(1)
	}

	fmt.Printf("Initializing client with AppID: %s, Host: %s\n", appId, host)
	c := client.NewClient(appId, appSecret, host)

	// 1. Test GetSongList
	fmt.Println("\n>>> Testing GetSongList...")
	listReq := &client.GetSongListRequest{
		Limit: 5,
	}
	listResp, err := c.GetSongList(listReq)
	if err != nil {
		fmt.Printf("[-] GetSongList failed: %v\n", err)
		// If auth fails here, it will exit
		return
	}
	fmt.Printf("[+] GetSongList success. Retrieved %d songs.\n", len(listResp.SongList))
	if len(listResp.SongList) > 0 {
		s := listResp.SongList[0]
		fmt.Printf("    First Song: %s (ID: %s)\n", s.SongName, s.SongId)
	}

	// // 2. Test SearchSong (if enabled)
	// fmt.Println("\n>>> Testing SearchSong (Query: 'a')...")
	searchReq := &client.SearchSongRequest{
		SearchText: "住在心里",
		SearchType: 1,
		Offset:     0,
		Limit:      5,
		Status:     1,
	}
	searchResp, err := c.SearchSong(searchReq)
	if err == nil {
		fmt.Printf("[+] SearchSong success. Total: %d\n", searchResp.Total)
	} else {
		fmt.Printf("[-] SearchSong failed (might not be enabled): %v\n", err)
	}
}
