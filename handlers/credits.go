package handlers

import (
	"errors"
	"suasor/services"
	"suasor/types/models"
	"suasor/types/requests"
	"suasor/types/responses"
	"suasor/utils/logger"

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
//
//	@Summary		Get all credits for a media item
//	@Description	Retrieves all credits (cast and crew) associated with a specific media item
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemID	path		int												true	"Media Item ID"
//	@Success		200		{array}		responses.APIResponse[[]models.Credit]									"Credits retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid media item ID"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{itemID} [get]
func (h *CreditHandler) GetCreditsForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get media item ID from path
	itemID, err := checkItemID(c, "itemID")
	if err != nil {
		return
	}

	// Get credits
	credits, err := h.creditService.GetCreditsForMediaItem(ctx, itemID)
	if handleServiceError(c, err, "Failed to get credits", "", "Failed to get credits") {
		return
	}

	// Return credits
	responses.RespondOK(c, credits, "Credits retrieved successfully")
}

// GetCastForMediaItem godoc
//
//	@Summary		Get cast for a media item
//	@Description	Retrieves all cast credits associated with a specific media item
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemID	path		int												true	"Media Item ID"
//	@Success		200		{array}		responses.APIResponse[[]models.Credit]									"Cast credits retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid media item ID"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{itemID}/cast [get]
func (h *CreditHandler) GetCastForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get media item ID from path
	itemID, err := checkItemID(c, "itemID")
	if err != nil {
		return
	}

	// Get cast
	cast, err := h.creditService.GetCastForMediaItem(ctx, itemID)
	if handleServiceError(c, err, "Failed to get cast", "", "Failed to get cast") {
		return
	}

	// Return cast
	responses.RespondOK(c, cast, "Cast credits retrieved successfully")
}

// GetCrewForMediaItem godoc
//
//	@Summary		Get crew for a media item
//	@Description	Retrieves all crew credits associated with a specific media item, optionally filtered by department
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemID		path		int												true	"Media Item ID"
//	@Param			department	query		string											false	"Filter by department (e.g., 'Directing', 'Writing')"
//	@Success		200			{array}		responses.APIResponse[[]models.Credit]									"Crew credits retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid media item ID"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{itemID}/crew [get]
func (h *CreditHandler) GetCrewForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get media item ID from path
	itemID, err := checkItemID(c, "itemID")
	if err != nil {
		return
	}

	// Get department from query (optional)
	departmentStr := c.Query("department")
	department := models.MediaDepartment(departmentStr)

	var crew []*models.Credit

	// Get crew, filtered by department if provided
	if department != "" {
		crew, err = h.creditService.GetCrewByDepartment(ctx, itemID, department)
	} else {
		crew, err = h.creditService.GetCrewForMediaItem(ctx, itemID)
	}

	if handleServiceError(c, err, "Failed to get crew", "", "Failed to get crew") {
		return
	}

	// Return crew
	responses.RespondOK(c, crew, "Crew credits retrieved successfully")
}

// GetDirectorsForMediaItem godoc
//
//	@Summary		Get directors for a media item
//	@Description	Retrieves all director credits associated with a specific media item
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemID	path		int												true	"Media Item ID"
//	@Success		200		{array}		responses.APIResponse[[]models.Credit]									"Director credits retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid media item ID"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{itemID}/directors [get]
func (h *CreditHandler) GetDirectorsForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get media item ID from path
	itemID, err := checkItemID(c, "itemID")
	if err != nil {
		return
	}

	// Get directors
	directors, err := h.creditService.GetDirectorsForMediaItem(ctx, itemID)
	if handleServiceError(c, err, "Failed to get directors", "", "Failed to get directors") {
		return
	}

	// Return directors
	responses.RespondOK(c, directors, "Director credits retrieved successfully")
}

// GetCreditsByPerson godoc
//
//	@Summary		Get all credits for a person
//	@Description	Retrieves all credits associated with a specific person
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			personID	path		int												true	"Person ID"
//	@Success		200			{array}		responses.APIResponse[[]models.Credit]									"Credits retrieved successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid person ID"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/person/{personID} [get]
func (h *CreditHandler) GetCreditsByPerson(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Get credits
	credits, err := h.creditService.GetCreditsByPerson(ctx, personID)
	if handleServiceError(c, err, "Failed to get credits", "", "Failed to get credits") {
		return
	}

	// Return credits
	responses.RespondOK(c, credits, "Credits retrieved successfully")
}

// CreateCredit godoc
//
//	@Summary		Create a new credit
//	@Description	Creates a new credit associating a person with a media item
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		requests.CreateCreditRequest					true	"Credit information"
//	@Success		201		{object}	responses.APIResponse[models.Credit]									"Credit created successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid request format"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits [post]
func (h *CreditHandler) CreateCredit(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var req requests.CreateCreditRequest
	if !checkJSONBinding(c, &req) {
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
	if handleServiceError(c, err, "Failed to create credit", "", "Failed to create credit") {
		return
	}

	// Return created credit
	responses.RespondCreated(c, createdCredit, "Credit created successfully")
}

// UpdateCredit godoc
//
//	@Summary		Update an existing credit
//	@Description	Updates a credit record with the provided information
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			creditID	path		int												true	"Credit ID"
//	@Param			request		body		requests.UpdateCreditRequest					true	"Updated credit information"
//	@Success		200			{object}	models.Credit									"Credit updated successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid credit ID or request format"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Credit not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{creditID} [put]
func (h *CreditHandler) UpdateCredit(c *gin.Context) {
	ctx := c.Request.Context()

	// Get credit ID from path
	creditID, err := checkItemID(c, "creditID")
	if err != nil {
		return
	}

	// Parse request body
	var req requests.UpdateCreditRequest
	if !checkJSONBinding(c, &req) {
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
	if handleServiceError(c, err, "Failed to update credit", "credit not found", "Failed to update credit") {
		return
	}

	// Return updated credit
	responses.RespondOK(c, updatedCredit, "Credit updated successfully")
}

// DeleteCredit godoc
//
//	@Summary		Delete a credit
//	@Description	Deletes a credit record by ID
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			creditID path		int true	"Credit ID"
//	@Success		200			{object}	responses.SuccessResponse "Credit deleted successfully"
//	@Failure		400			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid credit ID"
//	@Failure		404			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Credit not found"
//	@Failure		500			{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{creditID} [delete]
func (h *CreditHandler) DeleteCredit(c *gin.Context) {
	ctx := c.Request.Context()

	// Get credit ID from path
	creditID, err := checkItemID(c, "creditID")
	if err != nil {
		return
	}

	// Delete credit
	err = h.creditService.DeleteCredit(ctx, creditID)
	if handleServiceError(c, err, "Failed to delete credit", "credit not found", "Failed to delete credit") {
		return
	}

	// Return success
	responses.RespondOK(c, gin.H{"success": true}, "Credit deleted successfully")
}

// CreateCreditsForMediaItem godoc
//
//	@Summary		Create multiple credits for a media item
//	@Description	Creates multiple credits for a specific media item in a single operation
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemID	path		int												true	"Media Item ID"
//	@Param			request	body		requests.CreateCreditRequest					true	"Multiple credits information"
//	@Success		201		{array}		responses.APIResponse[models.Credit]									"Credits created successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid media item ID or request format"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{itemID} [post]
func (h *CreditHandler) CreateCreditsForMediaItem(c *gin.Context) {
	ctx := c.Request.Context()

	// Get media item ID from path
	itemID, err := checkItemID(c, "itemID")
	if err != nil {
		return
	}

	// Parse request body
	var req requests.CreateCreditsRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	// Create credits
	var credits []*models.Credit
	for _, creditReq := range req.Credits {
		credit := models.Credit{
			PersonID:    creditReq.PersonID,
			MediaItemID: itemID, // Use the path parameter
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
		credits = append(credits, &credit)
	}

	createdCredits, err := h.creditService.CreateCreditsForMediaItem(ctx, itemID, credits)
	if handleServiceError(c, err, "Failed to create credits", "", "Failed to create credits") {
		return
	}

	// Return created credits
	responses.RespondCreated(c, createdCredits, "Credits created successfully")
}

// GetCreditsByType godoc
//
//	@Summary		Get credits by type for a media item
//	@Description	Retrieves credits for a media item filtered by type (cast, crew, directors)
//	@Tags			credits
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			itemID	path		int												true	"Media Item ID"
//	@Param			type	path		string											true	"Credit type (cast, crew, directors)"
//	@Success		200		{array}		responses.APIResponse[models.Credit]									"Credits retrieved successfully"
//	@Failure		400		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Invalid media item ID or credit type"
//	@Failure		500		{object}	responses.ErrorResponse[responses.ErrorDetails]	"Server error"
//	@Router			/media/credits/{itemID}/{type} [get]
func (h *CreditHandler) GetCreditsByType(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.LoggerFromContext(ctx)

	// Get media item ID from path
	itemID, err := checkItemID(c, "itemID")
	if err != nil {
		return
	}

	// Get type from path
	creditType := c.Param("type")

	var credits []*models.Credit

	// Get credits based on type
	switch creditType {
	case "cast":
		credits, err = h.creditService.GetCastForMediaItem(ctx, itemID)
	case "crew":
		credits, err = h.creditService.GetCrewForMediaItem(ctx, itemID)
	case "directors":
		credits, err = h.creditService.GetDirectorsForMediaItem(ctx, itemID)
	default:
		log.Error().Str("type", creditType).Msg("Invalid credit type")
		responses.RespondBadRequest(c, errors.New("invalid credit type"), "Invalid credit type")
		return
	}

	if handleServiceError(c, err, "Failed to get credits", "", "Failed to get credits") {
		return
	}

	// Return credits
	responses.RespondOK(c, credits, "Credits retrieved successfully")
}
