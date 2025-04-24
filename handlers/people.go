package handlers

import (
	"errors"
	"net/http"
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

// GetPersonByID godoc
// @Summary Get person by ID
// @Description Retrieves a specific person by their ID
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Success 200 {object} models.Person "Person retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Person not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/{personID} [get]
func (h *PeopleHandler) GetPersonByID(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Get person
	person, err := h.personService.GetPersonByID(ctx, personID)
	if err != nil {
		handleServiceError(c, err, "Retrieving person", "Person not found", "Failed to get person")
		return
	}

	if person == nil {
		responses.RespondNotFound(c, errors.New("person not found"), "Person not found")
		return
	}

	// Return person
	c.JSON(http.StatusOK, person)
}

// GetPersonWithCredits godoc
// @Summary Get person with their credits
// @Description Retrieves a specific person along with all their credits
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Success 200 {object} map[string]interface{} "Person and their credits retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Person not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/{personID}/credits [get]
func (h *PeopleHandler) GetPersonWithCredits(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Get person with credits
	person, credits, err := h.personService.GetPersonWithCredits(ctx, personID)
	if handleServiceError(c, err, "Failed to get person with credits", "", "Failed to get person with credits") {
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

// SearchPeople godoc
// @Summary Search for people by name
// @Description Searches for people whose names match the provided query
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query"
// @Param limit query int false "Maximum number of results to return" default(20)
// @Success 200 {array} models.Person "People retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Missing search query or invalid limit"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people [get]
func (h *PeopleHandler) SearchPeople(c *gin.Context) {
	ctx := c.Request.Context()

	// Get query from request
	query := c.Query("q")
	if query == "" {
		responses.RespondBadRequest(c, errors.New("missing search query"), "Search query is required")
		return
	}

	// Get limit from request
	limit := utils.GetLimit(c, 20, 100, true)

	// Search people
	people, err := h.personService.SearchPeople(ctx, query, limit)
	if err != nil {
		handleServiceError(c, err, "Searching people", "", "Failed to search people")
		return
	}

	// Return people
	c.JSON(http.StatusOK, people)
}

// GetPopularPeople godoc
// @Summary Get popular people
// @Description Retrieves a list of popular people, sorted by popularity
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of results to return" default(20)
// @Success 200 {array} models.Person "Popular people retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid limit"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/popular [get]
func (h *PeopleHandler) GetPopularPeople(c *gin.Context) {
	ctx := c.Request.Context()

	// Get limit from request
	limit := utils.GetLimit(c, 20, 100, true)

	// Get popular people
	people, err := h.personService.GetPopularPeople(ctx, limit)
	if handleServiceError(c, err, "Failed to get popular people", "", "Failed to get popular people") {
		return
	}

	// Return people
	c.JSON(http.StatusOK, people)
}

// GetPeopleByRole godoc
// @Summary Get people by role
// @Description Retrieves people filtered by their professional role (Actor, Director, etc.)
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role path string true "Role to filter by (e.g., 'Actor', 'Director')"
// @Success 200 {array} models.Person "People retrieved successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Missing role parameter"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/roles/{role} [get]
func (h *PeopleHandler) GetPeopleByRole(c *gin.Context) {
	ctx := c.Request.Context()

	// Get role from path
	role := c.Param("role")
	if role == "" {
		responses.RespondBadRequest(c, errors.New("missing role"), "Role is required")
		return
	}

	// Get people by role
	people, err := h.personService.GetPeopleByRole(ctx, role)
	if handleServiceError(c, err, "Failed to get people by role", "", "Failed to get people by role") {
		return
	}

	// Return people
	c.JSON(http.StatusOK, people)
}

// CreatePerson godoc
// @Summary Create a new person
// @Description Creates a new person record with the provided information
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.CreatePersonRequest true "Person information"
// @Success 201 {object} models.Person "Person created successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people [post]
func (h *PeopleHandler) CreatePerson(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var req requests.CreatePersonRequest
	if !checkJSONBinding(c, &req) {
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
	if handleServiceError(c, err, "Failed to create person", "", "Failed to create person") {
		return
	}

	// Return created person
	c.JSON(http.StatusCreated, createdPerson)
}

// UpdatePerson godoc
// @Summary Update an existing person
// @Description Updates a person record with the provided information
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Param request body requests.UpdatePersonRequest true "Updated person information"
// @Success 200 {object} models.Person "Person updated successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID or request format"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Person not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/{personID} [put]
func (h *PeopleHandler) UpdatePerson(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Parse request body
	var req requests.UpdatePersonRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	// Get existing person
	existingPerson, err := h.personService.GetPersonByID(ctx, personID)
	if handleServiceError(c, err, "Failed to get person", "", "Failed to get person") {
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
	if handleServiceError(c, err, "Failed to update person", "person not found", "Failed to update person") {
		return
	}

	// Return updated person
	c.JSON(http.StatusOK, updatedPerson)
}

// DeletePerson godoc
// @Summary Delete a person
// @Description Deletes a person record by ID
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Success 200 {object} map[string]bool "Person deleted successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID"
// @Failure 404 {object} responses.ErrorResponse[responses.ErrorDetails] "Person not found"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/{personID} [delete]
func (h *PeopleHandler) DeletePerson(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Delete person
	err = h.personService.DeletePerson(ctx, personID)
	if handleServiceError(c, err, "Failed to delete person", "person not found", "Failed to delete person") {
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPersonCreditsGrouped godoc
// @Summary Get a person's credits grouped by type
// @Description Retrieves a person's credits organized by department and role
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Success 200 {object} map[string]map[string][]models.Credit "Credits grouped by department and role"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/{personID}/credits/grouped [get]
func (h *PeopleHandler) GetPersonCreditsGrouped(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Get person credits grouped
	creditsGrouped, err := h.personService.GetPersonCreditsGrouped(ctx, personID)
	if handleServiceError(c, err, "Failed to get person credits", "", "Failed to get person credits") {
		return
	}

	// Return credits grouped
	c.JSON(http.StatusOK, creditsGrouped)
}

// ImportPerson godoc
// @Summary Import a person from an external source
// @Description Imports a person from an external source with the provided details
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body requests.ImportPersonRequest true "Person import information"
// @Success 200 {object} models.Person "Person imported successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid request format"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/import [post]
func (h *PeopleHandler) ImportPerson(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var req requests.ImportPersonRequest
	if !checkJSONBinding(c, &req) {
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
	if handleServiceError(c, err, "Failed to import person", "", "Failed to import person") {
		return
	}

	// Return imported person
	c.JSON(http.StatusOK, importedPerson)
}

// AddExternalIDToPerson godoc
// @Summary Add external ID to person
// @Description Adds or updates an external ID reference for a person
// @Tags people
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param personID path int true "Person ID"
// @Param request body requests.ExternalIDRequest true "External ID information"
// @Success 200 {object} map[string]bool "External ID added successfully"
// @Failure 400 {object} responses.ErrorResponse[responses.ErrorDetails] "Invalid person ID or request format"
// @Failure 500 {object} responses.ErrorResponse[responses.ErrorDetails] "Server error"
// @Router /api/v1/people/{personID}/external-ids [post]
func (h *PeopleHandler) AddExternalIDToPerson(c *gin.Context) {
	ctx := c.Request.Context()

	// Get person ID from path
	personID, err := checkItemID(c, "personID")
	if err != nil {
		return
	}

	// Parse request body
	var req requests.ExternalIDRequest
	if !checkJSONBinding(c, &req) {
		return
	}

	// Add external ID
	err = h.personService.AddExternalIDToPerson(ctx, personID, req.Source, req.ID)
	if handleServiceError(c, err, "Failed to add external ID", "", "Failed to add external ID") {
		return
	}

	// Return success
	c.JSON(http.StatusOK, gin.H{"success": true})
}
