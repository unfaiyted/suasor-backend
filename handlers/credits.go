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

// GetCreditsForMediaItem gets all credits for a media item
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

// GetCastForMediaItem gets cast credits for a media item
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

// GetCrewForMediaItem gets crew credits for a media item
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

// GetDirectorsForMediaItem gets director credits for a media item
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

// GetCreditsByPerson gets all credits for a person
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

// CreateCredit creates a new credit
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

// UpdateCredit updates an existing credit
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

// DeleteCredit deletes a credit
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

// CreateCreditsForMediaItem creates multiple credits for a media item
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

// GetCreditsByType gets credits for a media item by type (cast, crew, directors, etc.)
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