// router/item.go
package router

import (
	"github.com/gin-gonic/gin"
	"suasor/di/container"
	"suasor/handlers"
)

// RegisterMediaItemRoutes configures routes for media items, people, and credits
func RegisterPeopleBasedRoutes(r *gin.RouterGroup, c *container.Container) {

	// People handlers
	peopleHandler := container.MustGet[*handlers.PeopleHandler](c)

	// People routes
	people := r.Group("/people")
	{
		people.GET("", peopleHandler.SearchPeople)
		people.GET("/popular", peopleHandler.GetPopularPeople)
		people.GET("/roles/:role", peopleHandler.GetPeopleByRole)
		people.GET("/:personID", peopleHandler.GetPersonByID)
		people.GET("/:personID/credits", peopleHandler.GetPersonWithCredits)
		people.GET("/:personID/credits/grouped", peopleHandler.GetPersonCreditsGrouped)
		people.POST("", peopleHandler.CreatePerson)
		people.PUT("/:personID", peopleHandler.UpdatePerson)
		people.DELETE("/:personID", peopleHandler.DeletePerson)
		people.POST("/import", peopleHandler.ImportPerson)
		people.POST("/:personID/external-ids", peopleHandler.AddExternalIDToPerson)
	}

	// Credit handlers
	creditHandler := container.MustGet[*handlers.CreditHandler](c)

	// Credit routes
	credits := r.Group("/credits")
	{
		credits.GET("/media/:itemID", creditHandler.GetCreditsForMediaItem)
		credits.GET("/media/:itemID/cast", creditHandler.GetCastForMediaItem)
		credits.GET("/media/:itemID/crew", creditHandler.GetCrewForMediaItem)
		credits.GET("/media/:itemID/directors", creditHandler.GetDirectorsForMediaItem)
		credits.GET("/media/:itemID/:type", creditHandler.GetCreditsByType)
		credits.GET("/person/:personID", creditHandler.GetCreditsByPerson)
		credits.POST("", creditHandler.CreateCredit)
		credits.POST("/media/:itemID", creditHandler.CreateCreditsForMediaItem)
		credits.PUT("/:creditID", creditHandler.UpdateCredit)
		credits.DELETE("/:creditID", creditHandler.DeleteCredit)
	}
}
