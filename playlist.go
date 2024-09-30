package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Playlist struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type PlaylistInfo struct {
	Items []Playlist `json:"items"`
}

type Track struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
}
type PlaylistResponse struct {
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
}

const token string =    "BQACkasNZ3hwIeXiLJNDcnPrULMVNz-wR3ALO6_HcNpgMcTMxqm-oFP3ONIBYxYCU-C821aqR6E9fOSLzdZYKXxF5BcxHOtKu1Mu05LbpMB63QKm1Uz2muEAr_pKOOM08axUhqUzHXzaQNmHTh05wuNA8hzsK5_3Aq6VpVGfO5dlTKisxhtHEHNni8S7LBEv-nLv7TPo68nxi7j65Qr1djJOpVoMuswcvCJwhkKqLAWnnfX5-KAT_oeqmJV9oPt71qGqQDGFUj2NhNul0UFZLTf6IrjCwh8Q"
const endpoint string = "https://api.spotify.com/v1/me/playlists"
const playlistFile string = "playlist.json"

func main() {

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error reading the body", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("token probably need to be refreshed:", string(body))
		return
	}

	//fmt.Println(string(body))

	var data PlaylistInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshaling data", err)
		return
	}
	// fmt.Println(data)
	if _, err := os.Stat("/path/to/whatever"); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
	}

	client = &http.Client{}
	for _, playlist := range data.Items {
		fmt.Printf("Playlist Name: %s, ID: %s\n", playlist.Name, playlist.Id)

		//playlistUrl := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlist.Id)
		fields := "items.track(id, name)"

		playlistReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=%s", playlist.Id, fields), nil)

		//   playlistReq, err := http.NewRequest("GET", fmt.Sprintf("%s%s?fields=%s", playlistUrl, playlist.Id, fields), nil)
		if err != nil {
			log.Fatal(err)
		}

		playlistReq.Header.Set("Authorization", "Bearer "+token)

		playlistResp, err := client.Do(playlistReq)
		if err != nil {
			log.Fatal(err)
		}
		defer playlistResp.Body.Close()
		playlistBody, err := io.ReadAll(playlistResp.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			return
		}
		fmt.Println(string(playlistBody))
fileWriter, err := os.OpenFile(playlistFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening automatonplayer.md file:", err)
		}

    defer fileWriter.Close()
    
    /*dataToWrite, err  := json.Marshal(playlistBody)
          if err != nil {
        fmt.Println("error marshalling body", err)
        return
    }*/

	_, err = fileWriter.Write(playlistBody)
		if err != nil {
			log.Fatal("error writing to automatonplayer.md:", err)
		}
		//var data PlaylistResponse
	}
}