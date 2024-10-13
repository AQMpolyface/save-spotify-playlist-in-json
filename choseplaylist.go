package main

import (
  //"net/http"
  "encoding/json"
  "fmt"
  "os"
)
type PlaylistItem struct {
  Item []struct {
    Tracks struct {
      Name string `json:"name`
    Id string `json:id`
  } `json:"tracks"`
} `json:"items"`
}
var playlistFile string = "playlist.json"
  func main() {

  jsonData, err := os.ReadFile("playlist.json")
      if err != nil {
          fmt.Println("error reading file", err)
    return
  }

var playlistItem PlaylistItem


  err = json.Unmarshal(jsonData, &playlistItem)
    if err != nil {
    fmt.Println("error unmarshalling data:", err)
    return
  }
fmt.Println(playlistItem.Item[2].Tracks.Name)
}
