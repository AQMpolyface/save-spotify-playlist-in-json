
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Track struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TrackItem struct {
	Track Track `json:"track"`
}

type TracksResponse struct {
	Items []TrackItem `json:"items"`
}

type Playlist struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Tracks []TrackItem `json:"tracks"` // This will be populated later
}

type PlaylistInfo struct {
	Items []Playlist `json:"items"`
}

type PlaylistsData struct {
	Playlists []Playlist `json:"playlists"`
}
const endpoint string = "https://api.spotify.com/v1/me/playlists"

func main() {
	var playlistFile string
	var help bool
	var token string
	flag.BoolVar(&help, "h", false, "input this flag to print this message of help")
	flag.StringVar(&token, "t", "", "token flag, (-t 'your_token')")
	flag.StringVar(&playlistFile, "f", "playlist.json", "file to write to flag (-f 'filepath')")
	flag.Parse()

	if help {
		printHelp()
		os.Exit(0) // Exit the program after printing help
	}
	if token == "" {
		_, err := os.Stat("token")
		if os.IsNotExist(err) {
			fmt.Println("file doesn't exist", err)
			fmt.Println("please enter a token (with -t, or put it in the code), this program can't work without it")
			return
		}
		fileStats, err := os.Stat("token")
		if err != nil {
			fmt.Println("error checking file", err)
			return
		}

		if fileStats.Size() <= 0 {
			fmt.Println("please enter a token (with -t, or put it in the code), this program can't work without it")
			return
		} else {
			fmt.Println("File exists, using the token saved in the token file")
			bytetoken, err := os.ReadFile("token")
			if err != nil {
				fmt.Println("error reading token file", err)
				return
			}
			token = string(bytetoken)
			fmt.Println(token)
		}
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

	if resp.Header.Get("Retry-After") != "" {
		fmt.Printf("rate limited by the Spotify API, you ran the code too much, retry in %s:\n %s", resp.Header.Get("Retry-After"), string(body))
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("token probably needs to be refreshed:", string(body))
		return
	}

	var data PlaylistInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshaling data", err)
		return
	}

	// Create a PlaylistsData structure to hold all playlists and their tracks
	playlistsData := PlaylistsData{}

	for _, playlist := range data.Items {
		// Fetch tracks for each playlist
		fields := "items(track(name,id))" // Corrected fields to get track details
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=%s", playlist.ID, fields)
		playlistReq, err := http.NewRequest("GET", url, nil)
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

		// Use the new TrackResponse struct to unmarshal the track data
		var trackData TracksResponse
		err = json.Unmarshal(playlistBody, &trackData)
		if err != nil {
			fmt.Println("error unmarshaling playlist content:", err)
			return
		}

		// Assign the tracks to the playlist
		playlist.Tracks = trackData.Items
		playlistsData.Playlists = append(playlistsData.Playlists,playlist)
	}

	// Write the playlists data to the specified JSON file
	file, err := os.Create(playlistFile)
	if err != nil {
		log.Fatal("error creating file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Set indentation for pretty printing
	if err := encoder.Encode(playlistsData); err != nil {
		log.Fatal("error encoding JSON to file:", err)
	}

	fmt.Printf("Playlists and tracks have been successfully written to %s\n", playlistFile)
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  -h        Show help")
	fmt.Println("  -t       Spotify API token (required)")
	fmt.Println("  -f       Output file for playlists (default: playlist.json)")
}
