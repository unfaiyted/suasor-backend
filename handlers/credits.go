package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

// CreditHandler handles credit-related requests
type CreditHandler struct {
	creditService *services.CreditService
}

// NewCreditHandler creates a new credit handler
func NewCreditHandler(creditService *services.CreditService) *CreditHandler {
	return &CreditHandler{
		creditService: creditService,
	}
}

// GetCreditsForMediaItem godoc
// @Summary Get all credits for a media item
// @Description Retrieves all credits (cast and crew) associated with a specific media item
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaItemID path int true "Media Item ID"
// @Success 200 {array} models.Credit "Credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid media item ID"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/media/{mediaItemID} [get]
func (h *CreditHandler) GetCreditsForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get media item ID from path
	mediaItemIDStr := c.Param("mediaItemID")
	mediaItemID, err := strconv.ParseUint(mediaItemIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("mediaItemID", mediaItemIDStr).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}
	
	// Get credits
	credits, err := h.creditService.GetCreditsForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get credits")
		responses.RespondInternalError(c, err, "Failed to get credits")
		return
	}
	
	// Return credits
	c.JSON(http.StatusOK, credits)
}

// GetCastForMediaItem godoc
// @Summary Get cast for a media item
// @Description Retrieves all cast credits associated with a specific media item
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaItemID path int true "Media Item ID"
// @Success 200 {array} models.Credit "Cast credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid media item ID"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/media/{mediaItemID}/cast [get]
func (h *CreditHandler) GetCastForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get media item ID from path
	mediaItemIDStr := c.Param("mediaItemID")
	mediaItemID, err := strconv.ParseUint(mediaItemIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("mediaItemID", mediaItemIDStr).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}
	
	// Get cast
	cast, err := h.creditService.GetCastForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get cast")
		responses.RespondInternalError(c, err, "Failed to get cast")
		return
	}
	
	// Return cast
	c.JSON(http.StatusOK, cast)
}

// GetCrewForMediaItem godoc
// @Summary Get crew for a media item
// @Description Retrieves all crew credits associated with a specific media item, optionally filtered by department
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaItemID path int true "Media Item ID"
// @Param department query string false "Filter by department (e.g., 'Directing', 'Writing')"
// @Success 200 {array} models.Credit "Crew credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid media item ID"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/media/{mediaItemID}/crew [get]
func (h *CreditHandler) GetCrewForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get media item ID from path
	mediaItemIDStr := c.Param("mediaItemID")
	mediaItemID, err := strconv.ParseUint(mediaItemIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("mediaItemID", mediaItemIDStr).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}
	
	// Get department from query (optional)
	department := c.Query("department")
	
	var crew []models.Credit
	
	// Get crew, filtered by department if provided
	if department != "" {
		crew, err = h.creditService.GetCrewByDepartment(ctx, mediaItemID, department)
	} else {
		crew, err = h.creditService.GetCrewForMediaItem(ctx, mediaItemID)
	}
	
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get crew")
		responses.RespondInternalError(c, err, "Failed to get crew")
		return
	}
	
	// Return crew
	c.JSON(http.StatusOK, crew)
}

// GetDirectorsForMediaItem godoc
// @Summary Get directors for a media item
// @Description Retrieves all director credits associated with a specific media item
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaItemID path int true "Media Item ID"
// @Success 200 {array} models.Credit "Director credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid media item ID"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/media/{mediaItemID}/directors [get]
func (h *CreditHandler) GetDirectorsForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get media item ID from path
	mediaItemIDStr := c.Param("mediaItemID")
	mediaItemID, err := strconv.ParseUint(mediaItemIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("mediaItemID", mediaItemIDStr).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}
	
	// Get directors
	directors, err := h.creditService.GetDirectorsForMediaItem(ctx, mediaItemID)
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Msg("Failed to get directors")
		responses.RespondInternalError(c, err, "Failed to get directors")
		return
	}
	
	// Return directors
	c.JSON(http.StatusOK, directors)
}

// GetCreditsByPerson godoc
// @Summary Get all credits for a person
// @Description Retrieves all credits associated with a specific person
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Success 200 {array} models.Credit "Credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/person/{personID} [get]
func (h *CreditHandler) GetCreditsByPerson(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get person ID from path
	personIDStr := c.Param("personID")
	personID, err := strconv.ParseUint(personIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("personID", personIDStr).Msg("Invalid person ID")
		responses.RespondBadRequest(c, err, "Invalid person ID")
		return
	}
	
	// Get credits
	credits, err := h.creditService.GetCreditsByPerson(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get credits")
		responses.RespondInternalError(c, err, "Failed to get credits")
		return
	}
	
	// Return credits
	c.JSON(http.StatusOK, credits)
}

// CreateCredit godoc
// @Summary Create a new credit
// @Description Creates a new credit associating a person with a media item
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CreateCreditRequest true "Credit information"
// @Success 201 {object} models.Credit "Credit created successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits [post]
func (h *CreditHandler) CreateCredit(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Parse request body
	var req requests.CreateCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Create credit
	credit := models.Credit{
		PersonID:    req.PersonID,
		MediaItemID: req.MediaItemID,
		Name:        req.Name,
		Role:        req.Role,
		Character:   req.Character,
		Department:  req.Department,
		Job:         req.Job,
		Order:       req.Order,
		IsCast:      req.IsCast,
		IsCrew:      req.IsCrew,
		IsGuest:     req.IsGuest,
		IsCreator:   req.IsCreator,
		IsArtist:    req.IsArtist,
	}
	
	createdCredit, err := h.creditService.CreateCredit(ctx, &credit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create credit")
		responses.RespondInternalError(c, err, "Failed to create credit")
		return
	}
	
	// Return created credit
	c.JSON(http.StatusCreated, createdCredit)
}

// UpdateCredit godoc
// @Summary Update an existing credit
// @Description Updates a credit record with the provided information
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param creditID path int true "Credit ID"
// @Param request body requests.UpdateCreditRequest true "Updated credit information"
// @Success 200 {object} models.Credit "Credit updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid credit ID or request format"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Credit not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/{creditID} [put]
func (h *CreditHandler) UpdateCredit(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get credit ID from path
	creditIDStr := c.Param("creditID")
	creditID, err := strconv.ParseUint(creditIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("creditID", creditIDStr).Msg("Invalid credit ID")
		responses.RespondBadRequest(c, err, "Invalid credit ID")
		return
	}
	
	// Parse request body
	var req requests.UpdateCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Create credit object
	credit := models.Credit{
		BaseModel: models.BaseModel{
			ID: creditID,
		},
		PersonID:    req.PersonID,
		MediaItemID: req.MediaItemID,
		Name:        req.Name,
		Role:        req.Role,
		Character:   req.Character,
		Department:  req.Department,
		Job:         req.Job,
		Order:       req.Order,
		IsCast:      req.IsCast,
		IsCrew:      req.IsCrew,
		IsGuest:     req.IsGuest,
		IsCreator:   req.IsCreator,
		IsArtist:    req.IsArtist,
	}
	
	// Update credit
	updatedCredit, err := h.creditService.UpdateCredit(ctx, &credit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update credit")
		
		// Check for specific errors
		if errors.Is(err, errors.New("credit not found")) {
			responses.RespondNotFound(c, err, "Credit not found")
			return
		}
		
		responses.RespondInternalError(c, err, "Failed to update credit")
		return
	}
	
	// Return updated credit
	c.JSON(http.StatusOK, updatedCredit)
}

// DeleteCredit godoc
// @Summary Delete a credit
// @Description Deletes a credit record by ID
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param creditID path int true "Credit ID"
// @Success 200 {object} map[string]bool "Credit deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid credit ID"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Credit not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/{creditID} [delete]
func (h *CreditHandler) DeleteCredit(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get credit ID from path
	creditIDStr := c.Param("creditID")
	creditID, err := strconv.ParseUint(creditIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("creditID", creditIDStr).Msg("Invalid credit ID")
		responses.RespondBadRequest(c, err, "Invalid credit ID")
		return
	}
	
	// Delete credit
	if err := h.creditService.DeleteCredit(ctx, creditID); err != nil {
		log.Error().Err(err).Msg("Failed to delete credit")
		
		// Check for specific errors
		if errors.Is(err, errors.New("credit not found")) {
			responses.RespondNotFound(c, err, "Credit not found")
			return
		}
		
		responses.RespondInternalError(c, err, "Failed to delete credit")
		return
	}
	
	// Return success
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// CreateCreditsForMediaItem godoc
// @Summary Create multiple credits for a media item
// @Description Creates multiple credits for a specific media item in a single operation
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaItemID path int true "Media Item ID"
// @Param request body requests.CreateCreditsRequest true "Multiple credits information"
// @Success 201 {array} models.Credit "Credits created successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid media item ID or request format"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/media/{mediaItemID} [post]
func (h *CreditHandler) CreateCreditsForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get media item ID from path
	mediaItemIDStr := c.Param("mediaItemID")
	mediaItemID, err := strconv.ParseUint(mediaItemIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("mediaItemID", mediaItemIDStr).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}
	
	// Parse request body
	var req requests.CreateCreditsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Create credits
	var credits []models.Credit
	for _, creditReq := range req.Credits {
		credit := models.Credit{
			PersonID:    creditReq.PersonID,
			MediaItemID: mediaItemID, // Use the path parameter
			Name:        creditReq.Name,
			Role:        creditReq.Role,
			Character:   creditReq.Character,
			Department:  creditReq.Department,
			Job:         creditReq.Job,
			Order:       creditReq.Order,
			IsCast:      creditReq.IsCast,
			IsCrew:      creditReq.IsCrew,
			IsGuest:     creditReq.IsGuest,
			IsCreator:   creditReq.IsCreator,
			IsArtist:    creditReq.IsArtist,
		}
		credits = append(credits, credit)
	}
	
	createdCredits, err := h.creditService.CreateCreditsForMediaItem(ctx, mediaItemID, credits)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create credits")
		responses.RespondInternalError(c, err, "Failed to create credits")
		return
	}
	
	// Return created credits
	c.JSON(http.StatusCreated, createdCredits)
}

// GetCreditsByType godoc
// @Summary Get credits by type for a media item
// @Description Retrieves credits for a media item filtered by type (cast, crew, directors)
// @Tags credits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mediaItemID path int true "Media Item ID"
// @Param type path string true "Credit type (cast, crew, directors)"
// @Success 200 {array} models.Credit "Credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid media item ID or credit type"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /credits/media/{mediaItemID}/{type} [get]
func (h *CreditHandler) GetCreditsByType(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get media item ID from path
	mediaItemIDStr := c.Param("mediaItemID")
	mediaItemID, err := strconv.ParseUint(mediaItemIDStr, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("mediaItemID", mediaItemIDStr).Msg("Invalid media item ID")
		responses.RespondBadRequest(c, err, "Invalid media item ID")
		return
	}
	
	// Get type from path
	creditType := c.Param("type")
	
	var credits []models.Credit
	
	// Get credits based on type
	switch creditType {
	case "cast":
		credits, err = h.creditService.GetCastForMediaItem(ctx, mediaItemID)
	case "crew":
		credits, err = h.creditService.GetCrewForMediaItem(ctx, mediaItemID)
	case "directors":
		credits, err = h.creditService.GetDirectorsForMediaItem(ctx, mediaItemID)
	default:
		log.Error().Str("type", creditType).Msg("Invalid credit type")
		responses.RespondBadRequest(c, errors.New("invalid credit type"), "Invalid credit type")
		return
	}
	
	if err != nil {
		log.Error().Err(err).Uint64("mediaItemID", mediaItemID).Str("type", creditType).Msg("Failed to get credits")
		responses.RespondInternalError(c, err, "Failed to get credits")
		return
	}
	
	// Return credits
	c.JSON(http.StatusOK, credits)
}