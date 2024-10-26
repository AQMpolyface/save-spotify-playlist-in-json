package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
//	"time"
)

type Playlist struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type PlaylistInfo struct {
	Items []Playlist `json:"items"`
}

type PlaylistResponse struct {
	Items []TrackItem `json:"items"` // Corrected
}

type TrackItem struct {
	Track Track `json:"track"` // Corrected
}

type Track struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

const endpoint string = "https://api.spotify.com/v1/me/playlists"

var playlistFileJson string = "playlist.json"

const playlistUrl string = "https://api.spotify.com/v1/playlists/"

func main() {
	var help bool
	var token string
	flag.BoolVar(&help, "h", false, "input this flage to print this message of help")
	flag.StringVar(&token, "t", "", "token flag, (-t 'your_token')")
	flag.StringVar(&playlistFileJson, "f", "", "file to write to flag (-f 'filepath')")
	flag.Parse()

	if help {
		printHelp()
		os.Exit(0) // Exit the program after printing help
	}
	if token == "" {
		_, err := os.Stat("token")
		if os.IsNotExist(err) {
			fmt.Println("file doesnt exist", err)
			fmt.Println("please enter a token (with -t, or put it in the code), this program cant work wihout it")
			return
		}
		fileStats, err := os.Stat("token")
		if err != nil {
			fmt.Println("error checking file", err)
			return
		}

		if fileStats.Size() <= 0 {
			fmt.Println("please enter a token (with -t, or put it in the code), this program cant work wihout it")
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
			fmt.Println(token)
		}
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	/*if err != nil {for _, item := range tracks {
		fmt.Printf("Track Name: %s, ID: %s\n", item.Track.Name, item.Track.ID)
	}
		log.Fatal(err)
	}*/
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

	if resp.Header.Get("Retry-After") != "" {
		fmt.Println("rate limited by the spotify api, you ran the code too much, retry in %s:\n %s", resp.Header.Get("Retry-After"), string(body))
		return

	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("token probably need to be refreshed:", string(body))
		return
	}

  
	//os.WriteFile("data1.json", body, 0644)

	var data PlaylistInfo
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshaling data", err)
		return
	}
	// fmt.Println(data)

	client = &http.Client{}
  if _, err := os.Stat("playlist.json"); os.IsNotExist(err) {
        // Create the file if it does not exist
        _, err := os.Create("playlist.json")
        if err != nil {
            log.Fatalf("Failed to create file: %v", err)
        }
    }
  playlistFile, err := os.OpenFile("playlist.json",os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
if err != nil {
        log.Fatalf("Failed to open file: %v", err)
    }
    defer playlistFile.Close()
  startPlaylist := `{`
if _, err := playlistFile.WriteString(startPlaylist + "\n"); err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }
	for _, playlist := range data.Items {
		fmt.Printf("Playlist Name: %s, ID: %s\n", playlist.Name, playlist.Id)

playlistName := fmt.Sprintf(`"playlistname" : "%s",
    "playlistis" : "%s",
    "items" [`, playlist.Name, playlist.Id)
if _, err := playlistFile.WriteString(playlistName + "\n"); err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }
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

		var musicData PlaylistResponse
		err = json.Unmarshal(playlistBody, &musicData)
		if err != nil {
			fmt.Println("error unmarshaling playlist content:", err)
			return
		}

   for _, item := range musicData.Items {
//			time.Sleep(time.Second * 1)
      fmt.Println(item.Track.Name)
			fmt.Println(item.Track.ID)
songName := fmt.Sprintf(` {
  "song" : "%s",
  "id" : "%s"
        }`,item.Track.Name, item.Track.ID )
if _, err := playlistFile.WriteString(songName + "\n"); err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }
		}

endSong := `]`
    if _, err := playlistFile.WriteString(endSong + "\n"); err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }

		// Writing to the file
		/*fileWriter, err := os.OpenFile("playlist.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening playlist.json file:", err)
			fmt.Println("uwu")
		}
		defer fileWriter.Close()

		// Write the JSON data to the file
		jsonData, err := json.MarshalIndent(musicData, "", "  ")
		if err != nil {
			log.Fatal("error marshaling data to JSON:", err)
		}

		_, err = fileWriter.Write(jsonData)
		if err != nil {
			log.Fatal("error writing to playlist.json:", err)
		}
    endPlaylist := `
    }
    `
err = os.WriteFile(playlistFile, []byte(startPlaylist), 0644)
    if err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }*/
	break
  }

endFile := `
  }
  `
    if _, err := playlistFile.WriteString(endFile + "\n"); err != nil {
        log.Fatalf("Failed to write to file: %v", err)
    }



}

func printHelp() {
	fmt.Println("Usage: go run playlist.go [options]")
	fmt.Println("Options:")
	fmt.Println("  -h    Input this flag to print this message of help")
	fmt.Println("  -t    Token flag, (-t 'your_token')")
	fmt.Println("  -f    File to write to flag (-f 'filepath')")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  myprogram -t 'your_token' -f 'filepath'")

}
