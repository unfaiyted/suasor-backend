package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mediatypes "suasor/clients/media/types"
	"suasor/repository"
	"suasor/services/scheduler"
	"suasor/types/models"
)

// UserActivityAnalysisJob analyzes user activity patterns to generate insights
type UserActivityAnalysisJob struct {
	jobRepo             repository.JobRepository
	userRepo            repository.UserRepository
	configRepo          repository.UserConfigRepository
	userMovieDataRepo   repository.UserMediaItemDataRepository[*mediatypes.Movie]
	userSeriesDataRepo  repository.UserMediaItemDataRepository[*mediatypes.Series]
	userEpisodeDataRepo repository.UserMediaItemDataRepository[*mediatypes.Episode]
	userMusicDataRepo   repository.UserMediaItemDataRepository[*mediatypes.Track]
	movieRepo           repository.CoreMediaItemRepository[*mediatypes.Movie]
	seriesRepo          repository.CoreMediaItemRepository[*mediatypes.Series]
	musicRepo           repository.CoreMediaItemRepository[*mediatypes.Track]
}

// NewUserActivityAnalysisJob creates a new user activity analysis job
func NewUserActivityAnalysisJob(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	configRepo repository.UserConfigRepository,
	userMovieDataRepo repository.UserMediaItemDataRepository[*mediatypes.Movie],
	userSeriesDataRepo repository.UserMediaItemDataRepository[*mediatypes.Series],
	userEpisodeDataRepo repository.UserMediaItemDataRepository[*mediatypes.Episode],
	userMusicDataRepo repository.UserMediaItemDataRepository[*mediatypes.Track],
	movieRepo repository.CoreMediaItemRepository[*mediatypes.Movie],
	seriesRepo repository.CoreMediaItemRepository[*mediatypes.Series],
	musicRepo repository.CoreMediaItemRepository[*mediatypes.Track],
) *UserActivityAnalysisJob {
	return &UserActivityAnalysisJob{
		jobRepo:             jobRepo,
		userRepo:            userRepo,
		configRepo:          configRepo,
		userMovieDataRepo:   userMovieDataRepo,
		userSeriesDataRepo:  userSeriesDataRepo,
		userEpisodeDataRepo: userEpisodeDataRepo,
		userMusicDataRepo:   userMusicDataRepo,
		movieRepo:           movieRepo,
		seriesRepo:          seriesRepo,
		musicRepo:           musicRepo,
	}
}

// Name returns the unique name of the job
func (j *UserActivityAnalysisJob) Name() string {
	return "system.user.activity.analysis"
}

// Schedule returns when the job should next run
func (j *UserActivityAnalysisJob) Schedule() time.Duration {
	// Default to weekly
	return 7 * 24 * time.Hour
}

// Execute runs the activity analysis job
func (j *UserActivityAnalysisJob) Execute(ctx context.Context) error {
	log.Println("Starting user activity analysis job")

	// Get all users
	users, err := j.userRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}

	// Process each user
	for _, user := range users {
		if err := j.analyzeUserActivity(ctx, user); err != nil {
			log.Printf("Error analyzing activity for user %s: %v", user.Username, err)
			// Continue with other users even if one fails
			continue
		}
	}

	// Run global analytics across all users
	if err := j.analyzeGlobalTrends(ctx); err != nil {
		log.Printf("Error analyzing global trends: %v", err)
	}

	log.Println("User activity analysis job completed")
	return nil
}

// analyzeUserActivity analyzes activity for a single user
func (j *UserActivityAnalysisJob) analyzeUserActivity(ctx context.Context, user models.User) error {
	// Skip inactive users
	if !user.Active {
		log.Printf("Skipping inactive user: %s", user.Username)
		return nil
	}

	// Get user configuration
	config, err := j.configRepo.GetUserConfig(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("error getting user config: %w", err)
	}

	// Check if activity analysis is enabled for the user
	if !config.ActivityAnalysisEnabled {
		log.Printf("Activity analysis not enabled for user: %s", user.Username)
		return nil
	}

	// Create a job run record for this user
	now := time.Now()
	jobRun := &models.JobRun{
		JobName:   j.Name(),
		JobType:   models.JobTypeAnalysis,
		Status:    models.JobStatusRunning,
		StartTime: &now,
		UserID:    &user.ID,
		Metadata:  fmt.Sprintf(`{"userId":%d,"username":"%s","type":"activityAnalysis"}`, user.ID, user.Username),
	}

	if err := j.jobRepo.CreateJobRun(ctx, jobRun); err != nil {
		log.Printf("Error creating job run record: %v", err)
		return err
	}

	// Process each type of analysis
	var jobError error
	analysisResults := map[string]interface{}{}

	// Analyze watch/listen times
	timeAnalysis, err := j.analyzeActivityTimes(ctx, user.ID)
	if err != nil {
		log.Printf("Error analyzing activity times: %v", err)
		jobError = err
	} else {
		analysisResults["activityTimes"] = timeAnalysis
	}

	// Analyze genre preferences
	genreAnalysis, err := j.analyzeGenrePreferences(ctx, user.ID)
	if err != nil {
		log.Printf("Error analyzing genre preferences: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		analysisResults["genrePreferences"] = genreAnalysis
	}

	// Analyze binge watching behavior
	bingeAnalysis, err := j.analyzeBingeWatching(ctx, user.ID)
	if err != nil {
		log.Printf("Error analyzing binge watching: %v", err)
		if jobError == nil {
			jobError = err
		}
	} else {
		analysisResults["bingeWatching"] = bingeAnalysis
	}

	// Complete the job run
	status := models.JobStatusCompleted
	errorMessage := ""
	if jobError != nil {
		status = models.JobStatusFailed
		errorMessage = jobError.Error()
	}

	// Store analysis results
	resultsJSON, _ := json.Marshal(analysisResults)
	j.storeUserAnalysisResults(ctx, user.ID, string(resultsJSON))

	j.completeJobRun(ctx, jobRun.ID, status, errorMessage)
	return nil
}

// completeJobRun finalizes a job run with status and error info
func (j *UserActivityAnalysisJob) completeJobRun(ctx context.Context, jobRunID uint64, status models.JobStatus, errorMsg string) {
	if err := j.jobRepo.CompleteJobRun(ctx, jobRunID, status, errorMsg); err != nil {
		log.Printf("Error completing job run: %v", err)
	}
}

// analyzeActivityTimes analyzes when a user typically consumes media
func (j *UserActivityAnalysisJob) analyzeActivityTimes(ctx context.Context, userID uint64) (map[string]interface{}, error) {
	log.Printf("Analyzing activity times for user %d", userID)

	// In a real implementation, we would:
	// 1. Query the user's media play history
	// 2. Analyze patterns in days of week and times of day
	// 3. Identify peak viewing/listening times
	// 4. Analyze how long sessions typically last

	// Mock implementation
	return map[string]interface{}{
		"peakDays":              []string{"Saturday", "Sunday"},
		"peakHours":             []int{20, 21, 22}, // 8pm-11pm
		"averageSessionMinutes": 85,
		"weekdayVsWeekend": map[string]float64{
			"weekday": 0.35, // 35% of activity on weekdays
			"weekend": 0.65, // 65% of activity on weekends
		},
	}, nil
}

// analyzeGenrePreferences analyzes a user's genre preferences
func (j *UserActivityAnalysisJob) analyzeGenrePreferences(ctx context.Context, userID uint64) (map[string]interface{}, error) {
	log.Printf("Analyzing genre preferences for user %d", userID)

	// In a real implementation, we would:
	// 1. Query the user's media play history
	// 2. Extract genre information from the media items
	// 3. Analyze which genres are most frequently watched/listened to
	// 4. Look for patterns in genre combinations

	// Mock implementation
	return map[string]interface{}{
		"movies": map[string]float64{
			"Sci-Fi":    0.32,
			"Action":    0.28,
			"Comedy":    0.18,
			"Drama":     0.12,
			"Animation": 0.10,
		},
		"series": map[string]float64{
			"Comedy":      0.35,
			"Drama":       0.25,
			"Sci-Fi":      0.15,
			"Documentary": 0.15,
			"Crime":       0.10,
		},
		"music": map[string]float64{
			"Rock":       0.40,
			"Electronic": 0.30,
			"Jazz":       0.20,
			"Classical":  0.10,
		},
	}, nil
}

// analyzeBingeWatching analyzes a user's binge watching behavior
func (j *UserActivityAnalysisJob) analyzeBingeWatching(ctx context.Context, userID uint64) (map[string]interface{}, error) {
	log.Printf("Analyzing binge watching behavior for user %d", userID)

	// In a real implementation, we would:
	// 1. Query the user's media play history
	// 2. Look for patterns of watching multiple episodes in a row
	// 3. Identify series that are most often binge-watched
	// 4. Calculate statistics on binge watching frequency

	// Mock implementation
	return map[string]interface{}{
		"bingeFrequency":          "weekly",
		"averageEpisodesPerBinge": 4.2,
		"mostBingedSeries": []string{
			"Stranger Things",
			"Breaking Bad",
			"The Office",
		},
		"bingesByDayOfWeek": map[string]float64{
			"Friday":   0.30,
			"Saturday": 0.45,
			"Sunday":   0.15,
			"Other":    0.10,
		},
	}, nil
}

// analyzeGlobalTrends analyzes trends across all users
func (j *UserActivityAnalysisJob) analyzeGlobalTrends(ctx context.Context) error {
	log.Println("Analyzing global trends across all users")

	// In a real implementation, we would:
	// 1. Aggregate data across all users
	// 2. Look for popular content, genres, viewing times
	// 3. Identify trends in user behavior
	// 4. Store the results for use in recommendations and reports

	// Mock implementation
	globalTrends := map[string]interface{}{
		"popularMovies": []string{
			"The Avengers",
			"Inception",
			"The Dark Knight",
		},
		"popularSeries": []string{
			"Game of Thrones",
			"Breaking Bad",
			"Stranger Things",
		},
		"popularGenres": map[string]float64{
			"Action": 0.25,
			"Drama":  0.20,
			"Sci-Fi": 0.18,
			"Comedy": 0.15,
			"Other":  0.22,
		},
		"peakActivityDays": map[string]float64{
			"Saturday": 0.25,
			"Sunday":   0.20,
			"Friday":   0.18,
			"Other":    0.37,
		},
	}

	// Store the global trends
	trendsJSON, _ := json.Marshal(globalTrends)
	j.storeGlobalTrends(ctx, string(trendsJSON))

	return nil
}

// storeUserAnalysisResults stores the analysis results for a user
func (j *UserActivityAnalysisJob) storeUserAnalysisResults(ctx context.Context, userID uint64, resultsJSON string) {
	// In a real implementation, we would:
	// 1. Store the results in a database table
	// 2. Track historical analysis for trend detection
	// 3. Make the results available to the user and other services

	log.Printf("Stored activity analysis results for user %d", userID)
}

// storeGlobalTrends stores the global trend analysis
func (j *UserActivityAnalysisJob) storeGlobalTrends(ctx context.Context, trendsJSON string) {
	// In a real implementation, we would:
	// 1. Store the results in a database table
	// 2. Make the results available to administrators and recommendation services

	log.Println("Stored global trend analysis results")
}

// SetupUserActivityAnalysisSchedule creates or updates an activity analysis schedule
func (j *UserActivityAnalysisJob) SetupUserActivityAnalysisSchedule(ctx context.Context, frequency string) error {
	// Check if job already exists
	existing, err := j.jobRepo.GetJobSchedule(ctx, j.Name())
	if err != nil {
		return fmt.Errorf("error checking for existing job: %w", err)
	}

	// If job exists, update it
	if existing != nil {
		existing.Frequency = frequency
		existing.Enabled = frequency != string(scheduler.FrequencyManual)
		return j.jobRepo.UpdateJobSchedule(ctx, existing)
	}

	// Create a new job schedule
	schedule := &models.JobSchedule{
		JobName:     j.Name(),
		JobType:     models.JobTypeAnalysis,
		Frequency:   frequency,
		Enabled:     frequency != string(scheduler.FrequencyManual),
		LastRunTime: nil, // Never run yet
	}

	return j.jobRepo.CreateJobSchedule(ctx, schedule)
}

// RunManualAnalysis runs the activity analysis job manually
func (j *UserActivityAnalysisJob) RunManualAnalysis(ctx context.Context) error {
	return j.Execute(ctx)
}

// GetUserActivityReport generates a report of a user's activity
func (j *UserActivityAnalysisJob) GetUserActivityReport(ctx context.Context, userID uint64) (map[string]interface{}, error) {
	// In a real implementation, we would:
	// 1. Retrieve the stored analysis results for the user
	// 2. Format them into a readable report
	// 3. Include historical trends and comparisons

	// Mock implementation
	return map[string]interface{}{
		"lastUpdated": time.Now().Format(time.RFC3339),
		"totalWatchTime": map[string]interface{}{
			"hours": 127.5,
			"trend": "+12% from last month",
		},
		"favoriteGenres":   []string{"Sci-Fi", "Action", "Comedy"},
		"mostWatchedShow":  "Stranger Things",
		"mostWatchedMovie": "Inception",
		"activityByDayOfWeek": map[string]float64{
			"Monday":    0.05,
			"Tuesday":   0.08,
			"Wednesday": 0.07,
			"Thursday":  0.10,
			"Friday":    0.15,
			"Saturday":  0.30,
			"Sunday":    0.25,
		},
		"recommendations": []string{
			"Based on your viewing habits, you might enjoy watching on Wednesday evenings",
			"You tend to enjoy sci-fi content the most",
			"You typically watch 2.5 hours at a time on weekends",
		},
	}, nil
}

// GetGlobalActivityReport generates a report of global activity
func (j *UserActivityAnalysisJob) GetGlobalActivityReport(ctx context.Context) (map[string]interface{}, error) {
	// In a real implementation, we would:
	// 1. Retrieve the stored global analysis results
	// 2. Format them into a readable report
	// 3. Include historical trends

	// Mock implementation
	return map[string]interface{}{
		"lastUpdated":                 time.Now().Format(time.RFC3339),
		"totalActiveUsers":            125,
		"averageUserWatchTimePerWeek": 8.5,
		"mostPopularContent": map[string][]string{
			"movies": {"The Avengers", "Inception", "The Dark Knight"},
			"series": {"Game of Thrones", "Breaking Bad", "Stranger Things"},
			"music":  {"Pink Floyd", "Taylor Swift", "The Beatles"},
		},
		"peakUsageTimes": map[string]string{
			"dailyPeak":   "8pm - 11pm",
			"weeklyPeak":  "Saturday evening",
			"monthlyPeak": "Last weekend of the month",
		},
	}, nil
}
