package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"suasor/clients/metadata/tmdb"
	"suasor/clients/types"

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
	// Get API key from environment or .env file
	apiKey := os.Getenv("TMDB_API_KEY")

	// If not set, try to load from .env file
	if apiKey == "" {
		// Read from .env file in the project root
		envFilePath := "/home/faiyt/codebase/suasor/.env"
		envContent, err := os.ReadFile(envFilePath)
		if err != nil {
			fmt.Printf("Failed to read .env file: %v\n", err)
			fmt.Println("Please set TMDB_API_KEY environment variable")
			return
		}

		// Parse the content for TMDB_API_KEY
		lines := strings.Split(string(envContent), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "TMDB_API_KEY=") {
				apiKey = strings.TrimPrefix(line, "TMDB_API_KEY=")
				apiKey = strings.Trim(apiKey, "\"' ")
				break
			}
		}

		if apiKey == "" {
			fmt.Println("TMDB_API_KEY not found in .env file")
			return
		}
	}

	// Create TMDB config
	config := types.NewTMDBConfig()
	config.Name = "TMDB Client"
	config.Enabled = true
	config.ApiKey = apiKey

	// Create TMDB client
	client, err := tmdb.NewClient(config)
	if err != nil {
		fmt.Printf("Failed to create TMDB client: %v\n", err)
		return
	}

	ctx := context.Background()

	// Search for movies
	movies, err := client.SearchMovies(ctx, "Inception")
	if err != nil {
		fmt.Printf("Failed to search movies: %v\n", err)
		return
	}

	fmt.Printf("Found %d movies matching 'Inception'\n", len(movies))

	if len(movies) > 0 {
		movie := movies[0]
		fmt.Printf("\nMovie Details:\n")
		fmt.Printf("ID: %s\n", movie.ID)
		fmt.Printf("Title: %s\n", movie.Title)
		fmt.Printf("Release Date: %s\n", movie.ReleaseDate)
		fmt.Printf("Overview: %s\n", movie.Overview)

		// Get movie details
		movieDetails, err := client.GetMovie(ctx, movie.ID)
		if err != nil {
			fmt.Printf("Failed to get movie details: %v\n", err)
			return
		}

		fmt.Printf("\nAdditional Details:\n")
		fmt.Printf("Runtime: %d minutes\n", movieDetails.Runtime)
		fmt.Printf("Vote Average: %.1f\n", movieDetails.VoteAverage)
		fmt.Printf("Vote Count: %d\n", movieDetails.VoteCount)

		// Get recommendations
		recommendations, err := client.GetMovieRecommendations(ctx, movie.ID)
		if err != nil {
			fmt.Printf("Failed to get movie recommendations: %v\n", err)
			return
		}

		fmt.Printf("\nRecommendations:\n")
		for i, rec := range recommendations {
			if i >= 5 {
				break
			}
			fmt.Printf("- %s (%s)\n", rec.Title, rec.ReleaseDate)
		}
	}

	// Search for TV shows
	tvShows, err := client.SearchTVShows(ctx, "Stranger Things")
	if err != nil {
		fmt.Printf("Failed to search TV shows: %v\n", err)
		return
	}

	fmt.Printf("\nFound %d TV shows matching 'Stranger Things'\n", len(tvShows))

	if len(tvShows) > 0 {
		show := tvShows[0]
		fmt.Printf("\nTV Show Details:\n")
		fmt.Printf("ID: %s\n", show.ID)
		fmt.Printf("Name: %s\n", show.Name)
		fmt.Printf("First Air Date: %s\n", show.FirstAirDate)
		fmt.Printf("Overview: %s\n", show.Overview)
	}

	// Get trending movies
	trending, err := client.GetTrendingMovies(ctx)
	if err != nil {
		fmt.Printf("Failed to get trending movies: %v\n", err)
		return
	}

	fmt.Printf("\nTrending Movies:\n")
	for i, movie := range trending {
		if i >= 5 {
			break
		}
		fmt.Printf("- %s (%s)\n", movie.Title, movie.ReleaseDate)
	}
}
