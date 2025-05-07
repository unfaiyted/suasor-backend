package main

import (
	"context"
	"fmt"
	"suasor/clients/media/emby"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	"suasor/clients/types"
)

func main() {
	ctx := context.Background()
	
	// Create a client - in real code this would come from the factory
	client := &emby.EmbyClient{}
	
	// Check if the client implements PlaylistProvider interface
	_, isPlaylistProvider := interface{}(client).(providers.PlaylistProvider)
	fmt.Printf("Client implements PlaylistProvider: %v\n", isPlaylistProvider)
	
	// Create an adapter
	adapter := providers.NewPlaylistListAdapter(client)
	
	// The adapter should implement ListProvider[*mediatypes.Playlist]
	_, isListProvider := interface{}(adapter).(providers.ListProvider[*mediatypes.Playlist])
	fmt.Printf("Adapter implements ListProvider[*mediatypes.Playlist]: %v\n", isListProvider)
	
	// Check direct cast which is what list_sync_service.go does (this should fail)
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("As expected, panic on direct cast: %v\n", r)
			}
		}()
		
		var listProvider providers.ListProvider[mediatypes.ListData]
		listProvider = adapter.(providers.ListProvider[mediatypes.ListData])
		fmt.Printf("Unexpectedly succeeded: %v\n", listProvider != nil)
	}()
	
	fmt.Println("Test completed")
}