package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// UpdateBroadcastTitle updates the title of the specified broadcast.
func UpdateBroadcastTitle(client *http.Client, newTitle string, broadcastId string) {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Error creating YouTube service: %v", err)
	}

	// Create the live broadcast resource to update.
	broadcast := &youtube.LiveBroadcast{
		Id: broadcastId,
		Snippet: &youtube.LiveBroadcastSnippet{
			Title: newTitle,
		},
	}

	// Call the Update method to modify the broadcast title.
	call := youtubeService.LiveBroadcasts.Update([]string{"snippet"}, broadcast)
	if _, err := call.Do(); err != nil {
		log.Fatalf("Error updating broadcast title: %v", err)
	}

	log.Printf("Broadcast title updated to: %s", newTitle)
}

// GetClient retrieves an HTTP client for the YouTube API.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens.
	tok, err := tokenFromFile("token.json")
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken("token.json", tok)
	}
	return config.Client(context.Background(), tok)
}

// GetActiveBroadcastId retrieves the ID of the active live broadcast.
func GetActiveBroadcastId(client *http.Client) (string, error) {
	println("%v", client)
	panic("")
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("Error creating YouTube service: %v", err)
	}

	// Retrieve the list of live broadcasts.
	call := youtubeService.LiveBroadcasts.List([]string{"id"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		return "", fmt.Errorf("Error retrieving broadcasts: %v", err)
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("No active broadcasts found")
	}

	// Return the ID of the first active broadcast.
	return response.Items[0].Id, nil
}

// tokenFromFile retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Visit the URL for the auth dialog: %v", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// saveToken saves a token to a file.
func saveToken(file string, token *oauth2.Token) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	title := flag.String("title", "", "New title for the broadcast")
	flag.Parse()

	if *title == "" {
		log.Fatal("You must provide a title")
	}

	// Read the credentials from the JSON file
	configData, err := os.ReadFile("token.json") // Update the path accordingly
	if err != nil {
		log.Fatalf("Error reading credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(configData, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	client := getClient(config)

	// Get the active broadcast ID
	broadcastId, err := GetActiveBroadcastId(client)
	if err != nil {
		log.Fatalf("Error getting active broadcast ID: %v", err)
	}

	// Update the broadcast title
	UpdateBroadcastTitle(client, *title, broadcastId)
}
