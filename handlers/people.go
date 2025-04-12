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

// PeopleHandler handles people-related requests
type PeopleHandler struct {
	personService *services.PersonService
}

// NewPeopleHandler creates a new people handler
func NewPeopleHandler(personService *services.PersonService) *PeopleHandler {
	return &PeopleHandler{
		personService: personService,
	}
}

// GetPersonByID gets a person by ID
func (h *PeopleHandler) GetPersonByID(c *gin.Context) {
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
	
	// Get person
	person, err := h.personService.GetPersonByID(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get person")
		responses.RespondInternalError(c, err, "Failed to get person")
		return
	}
	
	if person == nil {
		responses.RespondNotFound(c, errors.New("person not found"), "Person not found")
		return
	}
	
	// Return person
	c.JSON(http.StatusOK, person)
}

// GetPersonWithCredits gets a person with their credits
func (h *PeopleHandler) GetPersonWithCredits(c *gin.Context) {
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
	
	// Get person with credits
	person, credits, err := h.personService.GetPersonWithCredits(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get person with credits")
		responses.RespondInternalError(c, err, "Failed to get person with credits")
		return
	}
	
	if person == nil {
		responses.RespondNotFound(c, errors.New("person not found"), "Person not found")
		return
	}
	
	// Return person with credits
	c.JSON(http.StatusOK, gin.H{
		"person":  person,
		"credits": credits,
	})
}

// SearchPeople searches for people by name
func (h *PeopleHandler) SearchPeople(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get query from request
	query := c.Query("q")
	if query == "" {
		responses.RespondBadRequest(c, errors.New("missing search query"), "Search query is required")
		return
	}
	
	// Get limit from request
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit")
		responses.RespondBadRequest(c, err, "Invalid limit")
		return
	}
	
	// Search people
	people, err := h.personService.SearchPeople(ctx, query, limit)
	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("Failed to search people")
		responses.RespondInternalError(c, err, "Failed to search people")
		return
	}
	
	// Return people
	c.JSON(http.StatusOK, people)
}

// GetPopularPeople gets popular people
func (h *PeopleHandler) GetPopularPeople(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get limit from request
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error().Err(err).Str("limit", limitStr).Msg("Invalid limit")
		responses.RespondBadRequest(c, err, "Invalid limit")
		return
	}
	
	// Get popular people
	people, err := h.personService.GetPopularPeople(ctx, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get popular people")
		responses.RespondInternalError(c, err, "Failed to get popular people")
		return
	}
	
	// Return people
	c.JSON(http.StatusOK, people)
}

// GetPeopleByRole gets people by role
func (h *PeopleHandler) GetPeopleByRole(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Get role from path
	role := c.Param("role")
	if role == "" {
		responses.RespondBadRequest(c, errors.New("missing role"), "Role is required")
		return
	}
	
	// Get people by role
	people, err := h.personService.GetPeopleByRole(ctx, role)
	if err != nil {
		log.Error().Err(err).Str("role", role).Msg("Failed to get people by role")
		responses.RespondInternalError(c, err, "Failed to get people by role")
		return
	}
	
	// Return people
	c.JSON(http.StatusOK, people)
}

// CreatePerson creates a new person
func (h *PeopleHandler) CreatePerson(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Parse request body
	var req requests.CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Create person object
	person := &models.Person{
		Name:        req.Name,
		Photo:       req.Photo,
		DateOfBirth: req.DateOfBirth,
		DateOfDeath: req.DateOfDeath,
		Gender:      req.Gender,
		Biography:   req.Biography,
		Birthplace:  req.Birthplace,
		KnownFor:    req.KnownFor,
	}
	
	// Add external IDs
	for _, extID := range req.ExternalIDs {
		person.AddExternalID(extID.Source, extID.ID)
	}
	
	// Create person
	createdPerson, err := h.personService.CreatePerson(ctx, person)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create person")
		responses.RespondInternalError(c, err, "Failed to create person")
		return
	}
	
	// Return created person
	c.JSON(http.StatusCreated, createdPerson)
}

// UpdatePerson updates an existing person
func (h *PeopleHandler) UpdatePerson(c *gin.Context) {
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
	
	// Parse request body
	var req requests.UpdatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Get existing person
	existingPerson, err := h.personService.GetPersonByID(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get person")
		responses.RespondInternalError(c, err, "Failed to get person")
		return
	}
	
	if existingPerson == nil {
		responses.RespondNotFound(c, errors.New("person not found"), "Person not found")
		return
	}
	
	// Update fields if provided
	if req.Name != "" {
		existingPerson.Name = req.Name
	}
	if req.Photo != "" {
		existingPerson.Photo = req.Photo
	}
	if req.DateOfBirth != nil {
		existingPerson.DateOfBirth = req.DateOfBirth
	}
	if req.DateOfDeath != nil {
		existingPerson.DateOfDeath = req.DateOfDeath
	}
	if req.Gender != "" {
		existingPerson.Gender = req.Gender
	}
	if req.Biography != "" {
		existingPerson.Biography = req.Biography
	}
	if req.Birthplace != "" {
		existingPerson.Birthplace = req.Birthplace
	}
	if req.KnownFor != "" {
		existingPerson.KnownFor = req.KnownFor
	}
	
	// Add external IDs
	for _, extID := range req.ExternalIDs {
		existingPerson.AddExternalID(extID.Source, extID.ID)
	}
	
	// Update person
	updatedPerson, err := h.personService.UpdatePerson(ctx, existingPerson)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update person")
		
		// Check for specific errors
		if errors.Is(err, errors.New("person not found")) {
			responses.RespondNotFound(c, err, "Person not found")
			return
		}
		
		responses.RespondInternalError(c, err, "Failed to update person")
		return
	}
	
	// Return updated person
	c.JSON(http.StatusOK, updatedPerson)
}

// DeletePerson deletes a person
func (h *PeopleHandler) DeletePerson(c *gin.Context) {
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
	
	// Delete person
	if err := h.personService.DeletePerson(ctx, personID); err != nil {
		log.Error().Err(err).Msg("Failed to delete person")
		
		// Check for specific errors
		if errors.Is(err, errors.New("person not found")) {
			responses.RespondNotFound(c, err, "Person not found")
			return
		}
		
		responses.RespondInternalError(c, err, "Failed to delete person")
		return
	}
	
	// Return success
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPersonCreditsGrouped gets a person's credits grouped by type
func (h *PeopleHandler) GetPersonCreditsGrouped(c *gin.Context) {
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
	
	// Get person credits grouped
	creditsGrouped, err := h.personService.GetPersonCreditsGrouped(ctx, personID)
	if err != nil {
		log.Error().Err(err).Uint64("personID", personID).Msg("Failed to get person credits")
		responses.RespondInternalError(c, err, "Failed to get person credits")
		return
	}
	
	// Return credits grouped
	c.JSON(http.StatusOK, creditsGrouped)
}

// ImportPerson imports a person from an external source
func (h *PeopleHandler) ImportPerson(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)
	
	// Parse request body
	var req requests.ImportPersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Create person object from request
	person := &models.Person{
		Name:        req.PersonData.Name,
		Photo:       req.PersonData.Photo,
		DateOfBirth: req.PersonData.DateOfBirth,
		DateOfDeath: req.PersonData.DateOfDeath,
		Gender:      req.PersonData.Gender,
		Biography:   req.PersonData.Biography,
		Birthplace:  req.PersonData.Birthplace,
		KnownFor:    req.PersonData.KnownFor,
	}
	
	// Add external IDs
	for _, extID := range req.PersonData.ExternalIDs {
		person.AddExternalID(extID.Source, extID.ID)
	}
	
	// Import person
	importedPerson, err := h.personService.ImportPerson(ctx, req.Source, req.ExternalID, person)
	if err != nil {
		log.Error().Err(err).Msg("Failed to import person")
		responses.RespondInternalError(c, err, "Failed to import person")
		return
	}
	
	// Return imported person
	c.JSON(http.StatusOK, importedPerson)
}

// AddExternalIDToPerson adds an external ID to a person
func (h *PeopleHandler) AddExternalIDToPerson(c *gin.Context) {
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
	
	// Parse request body
	var req requests.ExternalIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		responses.RespondBadRequest(c, err, "Invalid request body")
		return
	}
	
	// Add external ID
	if err := h.personService.AddExternalIDToPerson(ctx, personID, req.Source, req.ID); err != nil {
		log.Error().Err(err).Msg("Failed to add external ID")
		responses.RespondInternalError(c, err, "Failed to add external ID")
		return
	}
	
	// Return success
	c.JSON(http.StatusOK, gin.H{"success": true})
}