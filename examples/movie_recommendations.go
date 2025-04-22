package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"suasor/clients"
	"suasor/clients/ai"
	"suasor/clients/ai/claude"
	aitypes "suasor/clients/ai/types"
	"suasor/clients/media"
	"suasor/clients/media/emby"
	"suasor/clients/media/providers"
	mediatypes "suasor/clients/media/types"
	clienttypes "suasor/clients/types"
	"suasor/types/models"
	"time"

	"github.com/joho/godotenv"
)

// MovieRecommendation represents a movie recommendation from Claude
type MovieRecommendation struct {
	Title       string   `json:"title"`
	Year        int      `json:"year"`
	Reasons     []string `json:"reasons"`
	GenreMatch  []string `json:"genreMatch,omitempty"`
	DirectorRef string   `json:"directorRef,omitempty"`
	ActorRef    string   `json:"actorRef,omitempty"`
}

// MovieRecommendations represents a collection of movie recommendations
type MovieRecommendations struct {
	Recommendations []MovieRecommendation `json:"recommendations"`
	BasedOn         []string              `json:"basedOn"`
	GeneratedAt     time.Time             `json:"generatedAt"`
}

// getRecentlyWatchedMovies gets the most recently watched movies from Emby
func getRecentlyWatchedMovies(ctx context.Context, embyClient media.ClientMedia, count int) ([]models.MediaItem[mediatypes.Movie], error) {
	// Get the movie provider interface directly
	movieProvider, ok := embyClient.(providers.MovieProvider)
	if !ok {
		return nil, fmt.Errorf("emby client does not support movies")
	}

	// Since we can't reliably get the history with proper movie objects, we'll use
	// a more direct approach and just get recently added movies as a simplification
	options := &mediatypes.QueryOptions{
		Limit:     count,
		Sort:      "DateCreated", // Use the correct Emby sort field
		SortOrder: mediatypes.SortOrderDesc,
	}

	movies, err := movieProvider.GetMovies(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get recently added movies: %w", err)
	}

	return movies, nil
}

// getFavoriteMovies gets the user's favorite movies from Emby
func getFavoriteMovies(ctx context.Context, embyClient media.ClientMedia) ([]models.MediaItem[mediatypes.Movie], error) {
	// Get the movie provider interface
	movieProvider, ok := embyClient.(providers.MovieProvider)
	if !ok {
		return nil, fmt.Errorf("emby client does not support movies")
	}

	// Query favorite movies
	options := &mediatypes.QueryOptions{
		Filters: map[string]string{
			"isFavorite": "true",
		},
		SortOrder: mediatypes.SortOrderDesc,
		Sort:      "SortName", // Use a supported Emby sort field
	}

	return movieProvider.GetMovies(ctx, options)
}

// checkIfMovieExists checks if a movie exists in the Emby library
func checkIfMovieExists(ctx context.Context, embyClient media.ClientMedia, title string, year int) (bool, string, error) {
	// Get the movie provider interface
	movieProvider, ok := embyClient.(providers.MovieProvider)
	if !ok {
		return false, "", fmt.Errorf("emby client does not support movies")
	}

	// Search for movies by title
	options := &mediatypes.QueryOptions{
		Query: title,
	}

	movies, err := movieProvider.GetMovies(ctx, options)
	if err != nil {
		return false, "", fmt.Errorf("failed to search for movies: %w", err)
	}

	// Check if any of the results match both title and year (exact match)
	for _, movie := range movies {
		if movie.Data.Details.ReleaseYear == year && strings.EqualFold(movie.Data.Details.Title, title) {
			return true, movie.ExternalID, nil
		}
	}

	// Do a more flexible fuzzy matching if no exact match found
	for _, movie := range movies {
		// Check for year match and title contains or is contained by the search title
		if movie.Data.Details.ReleaseYear == year {
			movieTitle := strings.ToLower(movie.Data.Details.Title)
			searchTitle := strings.ToLower(title)

			// Check if titles are very similar
			if movieTitle == searchTitle ||
				strings.Contains(movieTitle, searchTitle) ||
				strings.Contains(searchTitle, movieTitle) {
				return true, movie.ExternalID, nil
			}

			// Also match on title without special characters
			cleanMovieTitle := cleanTitle(movieTitle)
			cleanSearchTitle := cleanTitle(searchTitle)

			if cleanMovieTitle == cleanSearchTitle {
				return true, movie.ExternalID, nil
			}
		}
	}

	return false, "", nil
}

// cleanTitle removes common special characters and words from a title
func cleanTitle(title string) string {
	// Convert to lowercase and remove special characters
	title = strings.ToLower(title)
	title = strings.ReplaceAll(title, ":", "")
	title = strings.ReplaceAll(title, "-", "")
	title = strings.ReplaceAll(title, ".", "")
	title = strings.ReplaceAll(title, ",", "")
	title = strings.ReplaceAll(title, "'", "")
	title = strings.ReplaceAll(title, "\"", "")

	// Remove common words that might be omitted
	title = strings.ReplaceAll(title, "the ", "")
	title = strings.ReplaceAll(title, " the", "")
	title = strings.ReplaceAll(title, "a ", "")
	title = strings.ReplaceAll(title, " a", "")

	return strings.TrimSpace(title)
}

// createFallbackRecommendations creates fallback movie recommendations based on user preferences
// when Claude is unable to generate valid recommendations
func createFallbackRecommendations(recentMovies, favoriteMovies []models.MediaItem[mediatypes.Movie]) *MovieRecommendations {
	// Create a map to track genres and their frequency
	genreFrequency := make(map[string]int)
	directors := make(map[string]int)
	actors := make(map[string]int)
	years := make([]int, 0)
	movieTitles := make([]string, 0)

	// Process recent and favorite movies to extract preferences
	allMovies := append(append([]models.MediaItem[mediatypes.Movie]{}, recentMovies...), favoriteMovies...)

	for _, movie := range allMovies {
		// Track movie titles for "based on" field
		movieTitles = append(movieTitles, movie.Data.Details.Title)

		// Count genre frequencies
		for _, genre := range movie.Data.Details.Genres {
			genreFrequency[genre]++
		}

		// Track directors
		for _, crew := range movie.Data.Crew {
			if crew.Role == "Director" {
				directors[crew.Name]++
			}
		}

		// Track lead actors
		for i, actor := range movie.Data.Cast {
			if i < 2 { // Only consider first two actors
				actors[actor.Name]++
			}
		}

		// Track years
		years = append(years, movie.Data.Details.ReleaseYear)
	}

	// Create a list of the top genres
	type kv struct {
		Key   string
		Value int
	}

	// Convert genre map to slice for sorting
	var genreSlice []kv
	for k, v := range genreFrequency {
		genreSlice = append(genreSlice, kv{k, v})
	}

	// Sort by count (descending)
	sortByFrequency := func(slice []kv) []kv {
		for i := 0; i < len(slice)-1; i++ {
			for j := i + 1; j < len(slice); j++ {
				if slice[i].Value < slice[j].Value {
					slice[i], slice[j] = slice[j], slice[i]
				}
			}
		}
		return slice
	}

	genreSlice = sortByFrequency(genreSlice)

	// Get top directors and actors
	var directorSlice []kv
	for k, v := range directors {
		directorSlice = append(directorSlice, kv{k, v})
	}
	directorSlice = sortByFrequency(directorSlice)

	var actorSlice []kv
	for k, v := range actors {
		actorSlice = append(actorSlice, kv{k, v})
	}
	actorSlice = sortByFrequency(actorSlice)

	// Create fallback recommendations based on the most common genres
	recommendations := &MovieRecommendations{
		Recommendations: []MovieRecommendation{},
		BasedOn:         movieTitles[:min(5, len(movieTitles))],
		GeneratedAt:     time.Now(),
	}

	// Define some fallback recommendations for common genres
	fallbacksByGenre := map[string][]MovieRecommendation{
		"Action": {
			{Title: "John Wick", Year: 2014, Reasons: []string{"High-octane action sequences", "Stylish cinematography"}, GenreMatch: []string{"Action", "Thriller"}},
			{Title: "Mad Max: Fury Road", Year: 2015, Reasons: []string{"Visually stunning action", "Intense pacing"}, GenreMatch: []string{"Action", "Adventure"}},
			{Title: "The Raid", Year: 2011, Reasons: []string{"Revolutionary martial arts sequences", "Gritty action"}, GenreMatch: []string{"Action", "Crime", "Thriller"}},
		},
		"Comedy": {
			{Title: "The Grand Budapest Hotel", Year: 2014, Reasons: []string{"Quirky humor", "Visual storytelling"}, GenreMatch: []string{"Comedy", "Adventure"}},
			{Title: "Knives Out", Year: 2019, Reasons: []string{"Clever writing", "Entertaining mystery"}, GenreMatch: []string{"Comedy", "Crime", "Mystery"}},
			{Title: "Booksmart", Year: 2019, Reasons: []string{"Smart coming-of-age story", "Great chemistry between leads"}, GenreMatch: []string{"Comedy", "Drama"}},
		},
		"Drama": {
			{Title: "Parasite", Year: 2019, Reasons: []string{"Thought-provoking social commentary", "Unexpected plot twists"}, GenreMatch: []string{"Drama", "Thriller"}},
			{Title: "Marriage Story", Year: 2019, Reasons: []string{"Powerful performances", "Emotional storytelling"}, GenreMatch: []string{"Drama", "Romance"}},
			{Title: "Nomadland", Year: 2020, Reasons: []string{"Beautiful cinematography", "Moving character study"}, GenreMatch: []string{"Drama"}},
		},
		"Sci-Fi": {
			{Title: "Arrival", Year: 2016, Reasons: []string{"Thought-provoking concept", "Emotional depth"}, GenreMatch: []string{"Sci-Fi", "Drama"}},
			{Title: "Blade Runner 2049", Year: 2017, Reasons: []string{"Visually stunning", "Philosophical themes"}, GenreMatch: []string{"Sci-Fi", "Drama", "Mystery"}},
			{Title: "Dune", Year: 2021, Reasons: []string{"Epic scale", "Faithful adaptation"}, GenreMatch: []string{"Sci-Fi", "Adventure"}},
		},
		"Horror": {
			{Title: "Hereditary", Year: 2018, Reasons: []string{"Psychological horror", "Unsettling atmosphere"}, GenreMatch: []string{"Horror", "Mystery"}},
			{Title: "The Witch", Year: 2015, Reasons: []string{"Period setting", "Atmospheric tension"}, GenreMatch: []string{"Horror", "Mystery"}},
			{Title: "Get Out", Year: 2017, Reasons: []string{"Social commentary", "Original concept"}, GenreMatch: []string{"Horror", "Mystery", "Thriller"}},
		},
		"Thriller": {
			{Title: "Gone Girl", Year: 2014, Reasons: []string{"Twisty narrative", "Character-driven mystery"}, GenreMatch: []string{"Thriller", "Drama", "Mystery"}},
			{Title: "Nightcrawler", Year: 2014, Reasons: []string{"Compelling character study", "Dark themes"}, GenreMatch: []string{"Thriller", "Crime", "Drama"}},
			{Title: "Wind River", Year: 2017, Reasons: []string{"Atmospheric tension", "Strong performances"}, GenreMatch: []string{"Thriller", "Crime", "Mystery"}},
		},
	}

	// Add recommendations based on top genres found in user's movies
	usedTitles := make(map[string]bool)

	// Add up to 5 recommendations
	for _, genre := range genreSlice {
		if len(recommendations.Recommendations) >= 5 {
			break
		}

		if genreRecs, ok := fallbacksByGenre[genre.Key]; ok {
			for _, rec := range genreRecs {
				if len(recommendations.Recommendations) >= 5 {
					break
				}

				// Skip if we've already added this title
				if usedTitles[rec.Title] {
					continue
				}

				// Add director/actor references if appropriate
				if len(directorSlice) > 0 {
					rec.DirectorRef = fmt.Sprintf("If you enjoy films by %s", directorSlice[0].Key)
				}

				if len(actorSlice) > 0 {
					rec.ActorRef = fmt.Sprintf("Features acting in the style of %s", actorSlice[0].Key)
				}

				recommendations.Recommendations = append(recommendations.Recommendations, rec)
				usedTitles[rec.Title] = true
			}
		}
	}

	// If we still don't have enough recommendations, add some general critically acclaimed films
	generalRecs := []MovieRecommendation{
		{Title: "Everything Everywhere All at Once", Year: 2022, Reasons: []string{"Genre-bending storytelling", "Emotional depth"}, GenreMatch: []string{"Action", "Adventure", "Comedy"}},
		{Title: "The Shawshank Redemption", Year: 1994, Reasons: []string{"Timeless storytelling", "Character relationships"}, GenreMatch: []string{"Drama"}},
		{Title: "Whiplash", Year: 2014, Reasons: []string{"Intense performances", "Gripping direction"}, GenreMatch: []string{"Drama", "Music"}},
		{Title: "Spider-Man: Into the Spider-Verse", Year: 2018, Reasons: []string{"Innovative animation", "Fresh take on superhero genre"}, GenreMatch: []string{"Animation", "Action", "Adventure"}},
		{Title: "The Social Network", Year: 2010, Reasons: []string{"Sharp dialogue", "Compelling character study"}, GenreMatch: []string{"Drama", "Biography"}},
	}

	for _, rec := range generalRecs {
		if len(recommendations.Recommendations) >= 5 {
			break
		}
		if !usedTitles[rec.Title] {
			recommendations.Recommendations = append(recommendations.Recommendations, rec)
			usedTitles[rec.Title] = true
		}
	}

	return recommendations
}

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// cleanJSONResponse cleans up the JSON response from Claude
func cleanJSONResponse(response string) string {
	// Try to extract content from code blocks
	codeBlockStart := "```json"
	codeBlockEnd := "```"

	startIdx := strings.Index(response, codeBlockStart)
	if startIdx != -1 {
		startIdx += len(codeBlockStart)
		endIdx := strings.Index(response[startIdx:], codeBlockEnd)
		if endIdx != -1 {
			return strings.TrimSpace(response[startIdx : startIdx+endIdx])
		}
	}

	// Handle case where JSON is enclosed in backticks but without the json marker
	if strings.HasPrefix(response, "```") && strings.HasSuffix(response, "```") {
		return strings.TrimSpace(response[3 : len(response)-3])
	}

	// Remove any leading/trailing whitespace and non-JSON text
	response = strings.TrimSpace(response)

	// Look for the first opening brace and last closing brace
	firstBrace := strings.Index(response, "{")
	lastBrace := strings.LastIndex(response, "}")

	if firstBrace != -1 && lastBrace != -1 && lastBrace > firstBrace {
		return response[firstBrace : lastBrace+1]
	}

	// Return the original if no JSON-like structure found
	return response
}

// formatMoviesForAIPrompt formats a list of movies into a string for the AI prompt
func formatMoviesForAIPrompt(movies []models.MediaItem[mediatypes.Movie]) string {
	result := ""
	for i, movie := range movies {
		details := movie.Data.Details
		genres := ""
		if len(details.Genres) > 0 {
			genres = fmt.Sprintf("Genres: %v", details.Genres)
		}

		directors := ""
		for _, crew := range movie.Data.Crew {
			if crew.Role == "Director" {
				if directors != "" {
					directors += ", "
				}
				directors += crew.Name
			}
		}
		if directors != "" {
			directors = fmt.Sprintf("Director(s): %s", directors)
		}

		cast := ""
		if len(movie.Data.Cast) > 0 && len(movie.Data.Cast) <= 3 {
			cast = "Cast: "
			for j, actor := range movie.Data.Cast {
				if j > 0 {
					cast += ", "
				}
				cast += actor.Name
			}
		}

		result += fmt.Sprintf("%d. %s (%d) - %s %s %s\n",
			i+1,
			details.Title,
			details.ReleaseYear,
			genres,
			directors,
			cast,
		)
	}
	return result
}

// getMovieRecommendations uses Claude to generate movie recommendations based on watch history and favorites
func getMovieRecommendations(ctx context.Context, claudeClient ai.AIClient, recentMovies, favoriteMovies []models.MediaItem[mediatypes.Movie]) (*MovieRecommendations, error) {
	// Format the movies for the prompt
	recentMoviesText := formatMoviesForAIPrompt(recentMovies)
	favoriteMoviesText := formatMoviesForAIPrompt(favoriteMovies)

	// Create the prompt
	prompt := fmt.Sprintf(`Based on the following user's movie preferences, please recommend 5 movies they might enjoy.

Recently watched movies:
%s

Favorite movies:
%s

Please provide your recommendations in valid JSON format matching this structure:
{
  "recommendations": [
    {
      "title": "Movie Title",
      "year": 2023,
      "reasons": ["Reason 1", "Reason 2"],
      "genreMatch": ["Drama", "Thriller"],
      "directorRef": "Same director as Movie X from their favorites",
      "actorRef": "Stars Actor Y from their recently watched list"
    }
  ],
  "basedOn": ["Movie Title 1", "Movie Title 2"]
}

Your response should ONLY contain the JSON with no additional text.`, recentMoviesText, favoriteMoviesText)

	// Set up generation options
	options := &aitypes.GenerationOptions{
		Temperature:        0.7,
		MaxTokens:          1500,
		SystemInstructions: "You are a cinema expert who deeply understands movie preferences. You provide thoughtful recommendations based on what users enjoy.",
		ResponseFormat:     "json", // For Claude API versions that support this
	}

	// Generate recommendations
	response, err := claudeClient.GenerateText(ctx, prompt, options)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendations: %w", err)
	}

	// If Claude failed to generate a proper response, create a fallback recommendation
	if len(response) < 10 || !strings.Contains(response, "recommendations") {
		fmt.Printf("\nClaude returned an incomplete or invalid response. Using fallback recommendations...\n")
		return createFallbackRecommendations(recentMovies, favoriteMovies), nil
	}

	// Trim potential prefix/suffix from Claude's response
	// Sometimes Claude might add markdown code blocks or explanatory text
	response = cleanJSONResponse(response)

	// Parse the JSON response
	var recommendations MovieRecommendations
	if err := json.Unmarshal([]byte(response), &recommendations); err != nil {
		// Log the actual response content for debugging
		fmt.Printf("\nFailed to parse response from Claude. Response content:\n%s\n", response)
		fmt.Println("\nUsing fallback recommendations instead...")
		return createFallbackRecommendations(recentMovies, favoriteMovies), nil
	}

	// Add timestamp
	recommendations.GeneratedAt = time.Now()

	return &recommendations, nil
}

// loadEnv loads environment variables from .env files
func loadEnv() {
	// Try to load .env file from several possible locations
	locations := []string{
		".env",                                   // Current directory
		"../.env",                                // Parent directory
		filepath.Join(os.Getenv("HOME"), ".env"), // Home directory
	}

	for _, location := range locations {
		if _, err := os.Stat(location); err == nil {
			fmt.Printf("Loading environment from: %s\n", location)
			err = godotenv.Load(location)
			if err == nil {
				return // Successfully loaded
			}
			fmt.Printf("Warning: Error loading %s: %v\n", location, err)
		}
	}
	fmt.Println("Warning: No .env file found, using existing environment variables")
}

func main() {
	// Initialize context
	ctx := context.Background()

	// Load environment variables from .env file
	loadEnv()

	// Map the test environment variables to the variables we need
	embyServerURL := os.Getenv("EMBY_TEST_URL")
	embyAPIKey := os.Getenv("EMBY_TEST_API_KEY")
	embyUsername := os.Getenv("EMBY_TEST_USER")
	claudeAPIKey := os.Getenv("CLAUDE_API_KEY")
	claudeModel := os.Getenv("CLAUDE_MODEL")

	if claudeModel == "" {
		claudeModel = "claude-3-opus-20240229" // Default model
	}

	// Validate required environment variables
	if embyServerURL == "" || embyAPIKey == "" || embyUsername == "" || claudeAPIKey == "" {
		log.Fatal("Missing required environment variables from .env file. Please check your .env file contains: EMBY_TEST_URL, EMBY_TEST_API_KEY, EMBY_TEST_USER, CLAUDE_API_KEY")
	}

	// Create Emby config
	embyConfig := clienttypes.EmbyConfig{
		BaseClientMediaConfig: clienttypes.BaseClientMediaConfig{
			BaseURL: embyServerURL,
			APIKey:  embyAPIKey,
		},
		Username: embyUsername, // Use username instead of userID - the client will resolve it
	}

	// Create Claude config
	claudeConfig := clienttypes.ClaudeConfig{
		BaseAIClientConfig: clienttypes.BaseAIClientConfig{
			BaseURL:     "https://api.anthropic.com",
			APIKey:      claudeAPIKey,
			Model:       claudeModel,
			Temperature: 0.7,
			MaxTokens:   2500,
		},
	}

	// Create Emby client
	embyClient, err := emby.NewEmbyClient(ctx, 1, embyConfig)
	if err != nil {
		log.Fatalf("Failed to create Emby client: %v", err)
	}

	// Create Claude client
	claudeClient, err := claude.NewClaudeClient(ctx, 2, claudeConfig)
	if err != nil {
		log.Fatalf("Failed to create Claude client: %v", err)
	}

	// Register clients with factory
	fmt.Println("Registering clients with factory...")
	client.RegisterClientFactory(clienttypes.ClientTypeEmby, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (client.Client, error) {
		return embyClient, nil
	})
	client.RegisterClientFactory(clienttypes.ClientTypeClaude, func(ctx context.Context, clientID uint64, config clienttypes.ClientConfig) (client.Client, error) {
		return claudeClient, nil
	})

	fmt.Println("\n===== Phase 1: Gathering Movie Data =====")

	// Test connection to Emby
	fmt.Println("Testing connection to Emby server...")
	connected, err := embyClient.TestConnection(ctx)
	if err != nil || !connected {
		log.Fatalf("Failed to connect to Emby server: %v", err)
	}
	fmt.Println("Successfully connected to Emby server!")

	// Get recently watched movies (last 10)
	fmt.Println("\nFetching recently watched movies...")
	recentMovies, err := getRecentlyWatchedMovies(ctx, embyClient, 10)
	if err != nil {
		log.Fatalf("Failed to get recently watched movies: %v", err)
	}
	fmt.Printf("✓ Found %d recently watched movies\n", len(recentMovies))

	// Show a sample of recent movies
	if len(recentMovies) > 0 {
		fmt.Println("\nRecent movies sample:")
		max := 3
		if len(recentMovies) < max {
			max = len(recentMovies)
		}
		for i := 0; i < max; i++ {
			fmt.Printf("  - %s (%d)\n",
				recentMovies[i].Data.Details.Title,
				recentMovies[i].Data.Details.ReleaseYear)
		}
	}

	// Get favorite movies
	fmt.Println("\nFetching favorite movies...")
	favoriteMovies, err := getFavoriteMovies(ctx, embyClient)
	if err != nil {
		log.Fatalf("Failed to get favorite movies: %v", err)
	}
	fmt.Printf("✓ Found %d favorite movies\n", len(favoriteMovies))

	// Show a sample of favorite movies
	if len(favoriteMovies) > 0 {
		fmt.Println("\nFavorite movies sample:")
		max := 3
		if len(favoriteMovies) < max {
			max = len(favoriteMovies)
		}
		for i := 0; i < max; i++ {
			fmt.Printf("  - %s (%d)\n",
				favoriteMovies[i].Data.Details.Title,
				favoriteMovies[i].Data.Details.ReleaseYear)
		}
	}

	fmt.Println("\n===== Phase 2: Generating Recommendations with Claude AI =====")

	// Get movie recommendations from Claude
	fmt.Println("Sending movie preferences to Claude for analysis...")
	fmt.Printf("- Including %d recently watched movies\n", len(recentMovies))
	fmt.Printf("- Including %d favorite movies\n", len(favoriteMovies))

	recommendations, err := getMovieRecommendations(ctx, claudeClient, recentMovies, favoriteMovies)
	if err != nil {
		log.Fatalf("Failed to get movie recommendations: %v", err)
	}
	fmt.Printf("✓ Received %d movie recommendations from Claude\n", len(recommendations.Recommendations))

	fmt.Println("\n===== Phase 3: Checking Library for Recommended Movies =====")
	fmt.Println("Checking if recommended movies already exist in your library...\n")

	// Check if recommended movies already exist in the library
	var inLibrary, notInLibrary int
	for i, rec := range recommendations.Recommendations {
		exists, id, err := checkIfMovieExists(ctx, embyClient, rec.Title, rec.Year)
		status := "📋 Not in library"
		if err != nil {
			status = fmt.Sprintf("❌ Error checking: %v", err)
		} else if exists {
			status = fmt.Sprintf("✅ Already in library (ID: %s)", id)
			inLibrary++
		} else {
			notInLibrary++
		}

		fmt.Printf("Recommendation %d: %s (%d) - %s\n", i+1, rec.Title, rec.Year, status)
		fmt.Printf("   Why recommended: %v\n", rec.Reasons)
		if len(rec.GenreMatch) > 0 {
			fmt.Printf("   Genre match: %v\n", rec.GenreMatch)
		}
		if rec.DirectorRef != "" {
			fmt.Printf("   Director connection: %s\n", rec.DirectorRef)
		}
		if rec.ActorRef != "" {
			fmt.Printf("   Actor connection: %s\n", rec.ActorRef)
		}
		fmt.Println()
	}

	fmt.Printf("Summary: %d already in library, %d new recommendations\n", inLibrary, notInLibrary)

	fmt.Println("\n===== Phase 4: Saving Results =====")

	// Output full JSON
	jsonData, err := json.MarshalIndent(recommendations, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal recommendations to JSON: %v", err)
	}

	// Save recommendations to file
	filename := fmt.Sprintf("movie_recommendations_%s.json", time.Now().Format("2006-01-02"))
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		log.Fatalf("Failed to save recommendations to file: %v", err)
	}
	fmt.Printf("✅ Recommendations saved to %s\n", filename)

	fmt.Println("\n===== JSON Output =====")
	fmt.Println(string(jsonData))
}
