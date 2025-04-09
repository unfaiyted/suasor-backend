package handlers

import (
	"net/http"
	"strconv"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"

	"github.com/gin-gonic/gin"
	"suasor/utils"
)

// JobHandler manages job-related requests
type JobHandler struct {
	jobService services.JobService
}

// NewJobHandler creates a new job handler
func NewJobHandler(jobService services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

// GetAllJobSchedules retrieves all job schedules
// @Summary Get all job schedules
// @Description Returns a list of all job schedules
// @Tags jobs
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse[[]models.JobSchedule]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/schedules [get]
func (h *JobHandler) GetAllJobSchedules(c *gin.Context) {
	schedules, err := h.jobService.GetAllJobSchedules(c.Request.Context())
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get job schedules")
		return
	}

	responses.RespondOK(c, schedules, "Job schedules retrieved successfully")
}

// GetJobScheduleByName retrieves a job schedule by name
// @Summary Get job schedule by name
// @Description Returns a specific job schedule by its name
// @Tags jobs
// @Accept json
// @Produce json
// @Param name path string true "Job name"
// @Success 200 {object} responses.APIResponse[models.JobSchedule]
// @Failure 404 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/schedules/{name} [get]
func (h *JobHandler) GetJobScheduleByName(c *gin.Context) {
	name := c.Param("name")
	schedule, err := h.jobService.GetJobScheduleByName(c.Request.Context(), name)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get job schedule")
		return
	}

	if schedule == nil {
		responses.RespondNotFound(c, nil, "Job schedule not found")
		return
	}

	responses.RespondOK(c, schedule, "Job schedule retrieved successfully")
}

// CreateJobSchedule creates a new job schedule
// @Summary Create a new job schedule
// @Description Creates a new job schedule
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body models.JobSchedule true "Job schedule to create"
// @Success 201 {object} responses.APIResponse[models.JobSchedule]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/schedules [post]
func (h *JobHandler) CreateJobSchedule(c *gin.Context) {
	var schedule models.JobSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		responses.RespondValidationError(c, err, "Invalid job schedule")
		return
	}

	// Check if a job schedule with this name already exists
	existingSchedule, err := h.jobService.GetJobScheduleByName(c.Request.Context(), schedule.JobName)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to check existing job schedule")
		return
	}

	if existingSchedule != nil {
		responses.RespondBadRequest(c, nil, "A job schedule with this name already exists")
		return
	}

	// Create the job schedule
	if err := h.jobService.CreateJobSchedule(c.Request.Context(), &schedule); err != nil {
		responses.RespondInternalError(c, err, "Failed to create job schedule")
		return
	}

	responses.RespondSuccess[models.JobSchedule](c, http.StatusCreated, schedule, "Job schedule created successfully")
}

// UpdateJobSchedule updates a job schedule
// @Summary Update job schedule
// @Description Updates an existing job schedule
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body requests.UpdateJobScheduleRequest true "Job schedule update"
// @Success 200 {object} responses.APIResponse[models.JobSchedule]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 404 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/schedules [put]
func (h *JobHandler) UpdateJobSchedule(c *gin.Context) {
	var req requests.UpdateJobScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err, "Invalid request")
		return
	}

	// Get the existing job schedule
	schedule, err := h.jobService.GetJobScheduleByName(c.Request.Context(), req.JobName)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get job schedule")
		return
	}

	if schedule == nil {
		responses.RespondNotFound(c, nil, "Job schedule not found")
		return
	}

	// Update the schedule fields
	schedule.Frequency = req.Frequency
	schedule.Enabled = req.Enabled

	// Save the updated schedule
	if err := h.jobService.UpdateJobSchedule(c.Request.Context(), schedule); err != nil {
		responses.RespondInternalError(c, err, "Failed to update job schedule")
		return
	}

	responses.RespondOK[struct{}](c, struct{}{}, "Job schedule updated successfully")
}

// RunJobManually triggers a job to run immediately
// @Summary Run job manually
// @Description Triggers a job to run immediately
// @Tags jobs
// @Accept json
// @Produce json
// @Param name path string true "Job name"
// @Success 202 {object} responses.APIResponse[any]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 404 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/{name}/run [post]
func (h *JobHandler) RunJobManually(c *gin.Context) {
	name := c.Param("name")
	log := utils.LoggerFromContext(c)
	log.Info().Str("job", name).Msg("Running job manually")

	// Validate job exists by trying to get its schedule
	schedule, err := h.jobService.GetJobScheduleByName(c.Request.Context(), name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get job schedule")
		responses.RespondInternalError(c, err, "Failed to get job schedule")
		return
	}

	if schedule == nil {
		log.Error().Msg("Job not found")
		responses.RespondNotFound(c, nil, "Job not found")
		return
	}

	// Run the job
	if err := h.jobService.RunJobManually(c.Request.Context(), name); err != nil {
		log.Error().Err(err).Msg("Failed to run job")
		responses.RespondInternalError(c, err, "Failed to run job")
		return
	}

	log.Info().Msg("Job started successfully")
	responses.RespondSuccess[struct{}](c, http.StatusAccepted, struct{}{}, "Job started successfully")
}

// GetRecentJobRuns retrieves recent job runs
// @Summary Get recent job runs
// @Description Returns a list of recent job runs
// @Tags jobs
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results (default 50)"
// @Success 200 {object} responses.APIResponse[[]models.JobRun]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/runs [get]
func (h *JobHandler) GetRecentJobRuns(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	runs, err := h.jobService.GetRecentJobRuns(c.Request.Context(), limit)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get job runs")
		return
	}

	responses.RespondOK(c, runs, "Job runs retrieved successfully")
}

// GetMediaSyncJobs retrieves job runs for the current user
// @Summary Get job runs for current user
// @Description Returns a list of job runs for the current user
// @Tags jobs
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results (default 50)"
// @Success 200 {object} responses.APIResponse[[]models.MediaSyncJob]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/media-sync [get]
func (h *JobHandler) GetMediaSyncJobs(c *gin.Context) {
	// Get the user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondInternalError(c, nil, "User ID not found in context")
		return
	}

	jobs, err := h.jobService.GetMediaSyncJobs(c.Request.Context(), userID.(uint64))
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get media sync jobs")
		return
	}

	responses.RespondOK(c, jobs, "Media sync jobs retrieved successfully")
}

// SetupMediaSyncJob creates or updates a media sync job
// @Summary Setup media sync job
// @Description Creates or updates a media sync job for the current user
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body requests.SetupMediaSyncJobRequest true "Media sync job setup"
// @Success 200 {object} responses.APIResponse[any]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/media-sync [post]
func (h *JobHandler) SetupMediaSyncJob(c *gin.Context) {
	var req requests.SetupMediaSyncJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err, "Invalid request")
		return
	}

	// Get the user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondInternalError(c, nil, "User ID not found in context")
		return
	}

	if err := h.jobService.SetupMediaSyncJob(c.Request.Context(), userID.(uint64), req.ClientID, req.ClientType, req.MediaType, req.Frequency); err != nil {
		responses.RespondInternalError(c, err, "Failed to setup media sync job")
		return
	}

	responses.RespondOK[struct{}](c, struct{}{}, "Media sync job setup successfully")
}

// RunMediaSyncJob runs a media sync job manually
// @Summary Run media sync job manually
// @Description Runs a media sync job manually for the current user
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body requests.RunMediaSyncJobRequest true "Media sync job run"
// @Success 202 {object} responses.APIResponse[any]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/media-sync/run [post]
func (h *JobHandler) RunMediaSyncJob(c *gin.Context) {
	var req requests.RunMediaSyncJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err, "Invalid request")
		return
	}

	// Get the user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondInternalError(c, nil, "User ID not found in context")
		return
	}

	if err := h.jobService.RunMediaSyncJob(c.Request.Context(), userID.(uint64), req.ClientID, req.MediaType); err != nil {
		responses.RespondInternalError(c, err, "Failed to run media sync job")
		return
	}

	responses.RespondSuccess[struct{}](c, http.StatusAccepted, struct{}{}, "Media sync job started successfully")
}

// GetUserRecommendations retrieves recommendations for the current user
// @Summary Get recommendations for current user
// @Description Returns a list of recommendations for the current user
// @Tags jobs
// @Accept json
// @Produce json
// @Param active query bool false "Only return active recommendations (default true)"
// @Param limit query int false "Limit number of results (default 50)"
// @Success 200 {object} responses.APIResponse[[]models.Recommendation]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/recommendations [get]
func (h *JobHandler) GetUserRecommendations(c *gin.Context) {
	// Get the user ID from the context
	userID, exists := c.Get("userID")
	if !exists {
		responses.RespondInternalError(c, nil, "User ID not found in context")
		return
	}

	// Parse active parameter
	activeStr := c.DefaultQuery("active", "true")
	active := activeStr == "true"

	// Parse limit
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	recommendations, err := h.jobService.GetUserRecommendations(c.Request.Context(), userID.(uint64), active, limit)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get recommendations")
		return
	}

	responses.RespondOK(c, recommendations, "Recommendations retrieved successfully")
}

// GetJobRunProgress retrieves progress information for a specific job run
// @Summary Get job run progress
// @Description Returns progress information for a specific job run
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Job Run ID"
// @Success 200 {object} responses.APIResponse[models.JobRun]
// @Failure 404 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/runs/{id}/progress [get]
func (h *JobHandler) GetJobRunProgress(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid job run ID")
		return
	}

	jobRun, err := h.jobService.GetJobRunByID(c.Request.Context(), id)
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get job run")
		return
	}

	if jobRun == nil {
		responses.RespondNotFound(c, nil, "Job run not found")
		return
	}

	responses.RespondOK(c, jobRun, "Job run progress retrieved successfully")
}

// GetActiveJobRuns retrieves all currently running jobs
// @Summary Get all active job runs
// @Description Returns a list of all currently running jobs
// @Tags jobs
// @Accept json
// @Produce json
// @Success 200 {object} responses.APIResponse[[]models.JobRun]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/active [get]
func (h *JobHandler) GetActiveJobRuns(c *gin.Context) {
	runs, err := h.jobService.GetActiveJobRuns(c.Request.Context())
	if err != nil {
		responses.RespondInternalError(c, err, "Failed to get active job runs")
		return
	}

	responses.RespondOK(c, runs, "Active job runs retrieved successfully")
}

// DismissRecommendation marks a recommendation as dismissed
// @Summary Dismiss recommendation
// @Description Marks a recommendation as dismissed
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Recommendation ID"
// @Success 200 {object} responses.APIResponse[any]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/recommendations/{id}/dismiss [post]
func (h *JobHandler) DismissRecommendation(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid recommendation ID")
		return
	}

	if err := h.jobService.DismissRecommendation(c.Request.Context(), id); err != nil {
		responses.RespondInternalError(c, err, "Failed to dismiss recommendation")
		return
	}

	responses.RespondOK[struct{}](c, struct{}{}, "Recommendation dismissed successfully")
}

// UpdateRecommendationViewedStatus updates whether a recommendation has been viewed
// @Summary Update recommendation viewed status
// @Description Updates whether a recommendation has been viewed
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path int true "Recommendation ID"
// @Param request body requests.UpdateRecommendationViewedRequest true "Viewed status update"
// @Success 200 {object} responses.APIResponse[any]
// @Failure 400 {object} responses.ErrorResponse[error]
// @Failure 500 {object} responses.ErrorResponse[error]
// @Router /jobs/recommendations/{id}/viewed [put]
func (h *JobHandler) UpdateRecommendationViewedStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		responses.RespondBadRequest(c, err, "Invalid recommendation ID")
		return
	}

	var req requests.UpdateRecommendationViewedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondValidationError(c, err, "Invalid request")
		return
	}

	if err := h.jobService.UpdateRecommendationViewedStatus(c.Request.Context(), id, req.Viewed); err != nil {
		responses.RespondInternalError(c, err, "Failed to update recommendation viewed status")
		return
	}

	responses.RespondOK[struct{}](c, struct{}{}, "Recommendation viewed status updated successfully")
}

