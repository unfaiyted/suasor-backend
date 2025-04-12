package handlers

import (
	"github.com/gin-gonic/gin"
)

// CalendarHandler handles all calendar operations
type CalendarHandler struct {
	// Probably going to need to create a calendar service that handles most of the logic.
}

// NewCalendarHandler creates a new calendar handler
func NewCalendarHandler() *CalendarHandler {
	return &CalendarHandler{}
}

// We should be able to get all calendar items for a user or all users.
// We should be able to pull the calendar items from our configured calendar sources.
// Specifically we have the calendar info in our automation integrations.
// We also should be able to pull up items that are upcoming from the tmdb metadata service integration.
// Maybe look into calendar items that are based on items we recommened or think the user would be interesting that are coming up. Use the recommendations to look for items that are coming out and map those to the calendar.

// GetCalendarItems handles retrieving calendar items
func (h *CalendarHandler) GetCalendarItems(c *gin.Context) {
	// We should be able to get calendar items by a user ID and a start and end date.
	// The calendar items should be returned in a list of calendar items.

	// We will need to create a calendar service that handles most of the logic.
}
