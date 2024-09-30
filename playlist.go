package main

import (
	"encoding/json"
//	"errors"
	"flag"
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

const token string =    "Your_token_here"
const endpoint string = "https://api.spotify.com/v1/me/playlists"
const playlistFile string = "playlist.json"

func main() {
    var token string
    flag.StringVar(&token, "", "t" , "token flag")

  flag.Parse()

    if token == "Your_token_here" || token == "" {
    fmt.Println("please enter a token (with -t, or put it in the program), this program cant work wihout it")
    return

  }
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
