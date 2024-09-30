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
	Name    string `json:"name"`
	Id  	string `json:"id"`
	Artists []struct {
		Name string `json:"name"`
	} `json:"artists"`
}


type Tracks struct {
	Items []struct {
		Track Track `json:"track"`
	} `json:"items"`
}

type PlaylistResponse struct {
	Tracks Tracks `json:"tracks"`
}

const endpoint string = "https://api.spotify.com/v1/me/playlists"
const playlistFile string = "playlist.json"
const playlistUrl string = "https://api.spotify.com/v1/playlists/"

var token string

func main() {

	flag.StringVar(&token, "t", "", "token flag")
	flag.Parse()

	if token == "" {
		_, err := os.Stat("token")
		if os.IsNotExist(err) {
			fmt.Println("please enter a token (with -t, or put it in the program), this program cant work wihout it")
			fmt.Scan(&token)
			return
		} else if err != nil {
			fmt.Println("Error checking file:", err)
		} else {
			fmt.Println("File exists, using the token saved in the token file")
			bytetoken, err := os.ReadFile("token")
			if err != nil {
				fmt.Println("error reading token file", err)
				return
			}
			token = string(bytetoken)
		}
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
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
		//fields := "items.track(id, name)"  ?fields=%s*

		//playlistReq, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlist.Id), nil)

		//   playlistReq, err := http.NewRequest("GET", fmt.Sprintf("%s%s?fields=%s", playlistUrl, playlist.Id, fields), nil)
		//if err != nil {
		//	log.Fatal(err)
		//}
		fields := "items.track(name,id)"
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=%s", playlist.Id, fields)
		playlistReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		playlistReq.Header.Set("Authorization", "Bearer "+token)
		//playlistReq.Header.Set("User-Agent", "curl/7.64.1")
		client := &http.Client{}
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
			log.Fatal("error writing to playlist.json:", err)
		}
		var playlistResponse PlaylistResponse
		err = json.Unmarshal(body, &playlistResponse)
		if err != nil {
			log.Fatal("error unmarshaling data", err)
		}

		// Loop through the items and print the track ID and name
		for _, item := range playlistResponse.Tracks.Items {
			fmt.Println("startinf to parse songname and id")
			songName := item.Track.Name
			idname := item.Track.Name
			fmt.Printf("Track Name: %s, Track ID: %s\n", songName, idname)
		}
	}

	}

