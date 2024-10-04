// package uploadyoutube

// import (
// 	"context"
// 	"log"
// 	"net/http"

// 	"google.golang.org/api/option"
// 	"google.golang.org/api/youtube/v3"
// )

// func UpdateBroadcastTitle(client *http.Client, newTitle string, broadcastId string) {
// 	ctx := context.Background()
// 	youtubeService, err := youtube.NewService(ctx, option.WithHTTPClient(client))
// 	if err != nil {
// 		log.Fatalf("Error creating YouTube service: %v", err)
// 	}

// 	broadcastSnippet := &youtube.LiveBroadcastSnippet{
// 		Title: newTitle,
// 	}

// 	call := youtubeService.LiveBroadcasts.Update("snippet", broadcastSnippet)
// 	if _, err := call.Do(); err != nil {
// 		log.Fatalf("Error updating broadcast title: %v", err)
// 	}

// 	log.Printf("Broadcast title updated to: %s", newTitle)
// }
