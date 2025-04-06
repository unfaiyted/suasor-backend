package main

import (
	"fmt"
	"os"

	tmdb "github.com/cyruzin/golang-tmdb"

	"github.com/joho/godotenv"
	"path/filepath"
)

func init() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",    // Current directory
		"../.env", // Project root
		filepath.Join(os.Getenv("HOME"), "claude_test.env"), // Home directory
	}

	for _, location := range locations {
		err := godotenv.Load(location)
		if err == nil {
			fmt.Printf("Loaded environment from: %s\n", location)
			break
		}
	}
}

func main() {
	// Use environment variable for the API key
	apiKey := os.Getenv("TMDB_API_KEY")
	if apiKey == "" {
		fmt.Println("TMDB_API_KEY environment variable is not set")
		return
	}

	// Initialize the TMDB client
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		fmt.Printf("Failed to initialize TMDB client: %v\n", err)
		return
	}

	// Search for a movie
	options := map[string]string{
		"language": "en-US",
		"page":     "1",
	}

	searchResults, err := tmdbClient.GetSearchMovies("Inception", options)
	if err != nil {
		fmt.Printf("Failed to search for movies: %v\n", err)
		return
	}

	fmt.Printf("Found %d results for 'Inception'\n", len(searchResults.Results))

	// Display the first movie result
	if len(searchResults.Results) > 0 {
		movie := searchResults.Results[0]
		fmt.Printf("Movie ID: %d\n", movie.ID)
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Printf("Overview: %s\n", movie.Overview)

		// Get movie details
		movieDetails, err := tmdbClient.GetMovieDetails(int(movie.ID), options)
		if err != nil {
			fmt.Printf("Failed to get movie details: %v\n", err)
			return
		}

		fmt.Printf("\nMovie Details:\n")
		fmt.Printf("Budget: %d\n", movieDetails.Budget)
		fmt.Printf("Revenue: %d\n", movieDetails.Revenue)
		fmt.Printf("Runtime: %d minutes\n", movieDetails.Runtime)
		fmt.Printf("Status: %s\n", movieDetails.Status)

		if len(movieDetails.Genres) > 0 {
			fmt.Println("\nGenres:")
			for _, genre := range movieDetails.Genres {
				fmt.Printf("- %s\n", genre.Name)
			}
		}
	}
}

