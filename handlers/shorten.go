package handlers

import (
	"net/http"
	"suasor/models"
	"suasor/services"
	"suasor/utils"

	"github.com/gin-gonic/gin"
)

type ShortenHandler struct {
	service services.ShortenService
}

func NewShortenHandler(service services.ShortenService) *ShortenHandler {
	return &ShortenHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create a shortened URL
// @Description Creates a new shortened URL from a long URL, with optional custom code and expiration. If no custom code is provided, one will be generated.
// @Tags shorten
// @Accept json
// @Produce json
// @Param request body models.ShortenRequest true "URL to shorten"
// @Example request
//
//	{
//	  "originalUrl": "https://example.com/some/very/long/path/that/needs/shortening",
//	  "customCode": "mycode",
//	  "expiresAfter": 7
//	}
//
// @Success 201 {object} models.APIResponse[models.ShortenData] "Successfully created shortened URL"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "shorten": {
//	      "originalUrl": "https://example.com/some/very/long/path/that/needs/shortening",
//	      "shortCode": "mycode",
//	      "expiresAt": "2023-06-15T10:30:45Z"
//	    }
//	  },
//	  "message": "URL shortened successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format or short code already exists"
// @Example response
//
//	{
//	  "error": "bad_request",
//	  "message": "The specified short code already exists",
//	  "details": {
//	    "error": "short code already exists"
//	  },
//	  "timestamp": "2023-06-08T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /shorten [post]
func (h *ShortenHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for URL shortening")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Str("originalUrl", req.OriginalURL).
		Str("customCode", req.CustomCode).
		Int("expiresAfter", req.ExpiresAfter).
		Msg("Creating shortened URL")

	result, err := h.service.Create(ctx, req)
	if err != nil {
		// Check if this is a code conflict error
		if err.Error() == "short code already exists" {
			log.Warn().Err(err).Str("customCode", req.CustomCode).Msg("Short code already exists")
			utils.RespondBadRequest(c, err, "The specified short code already exists")
			return
		}

		log.Error().Err(err).Str("originalUrl", req.OriginalURL).Msg("Failed to create shortened URL")
		utils.RespondInternalError(c, err, "Failed to create shortened URL")
		return
	}

	log.Info().Str("shortCode", result.Shorten.ShortCode).Str("shortUrl", result.ShortURL).Msg("Successfully created shortened URL")

	utils.RespondCreated(c, result, "URL shortened successfully")
}

// Update godoc
// @Summary Update a shortened URL
// @Description Updates an existing shortened URL by its short code
// @Tags shorten
// @Accept json
// @Produce json
// @Param code path string true "Short code identifier" example:"abc123"
// @Param request body models.ShortenRequest true "Updated URL data"
// @Example request
//
//	{
//	  "originalUrl": "https://example.com/updated/path",
//	  "expiresAfter": 14
//	}
//
// @Success 200 {object} models.APIResponse[models.ShortenData] "Successfully updated shortened URL"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "shorten": {
//	      "originalUrl": "https://example.com/updated/path",
//	      "shortCode": "abc123",
//	      "expiresAt": "2023-06-22T10:30:45Z"
//	    }
//	  },
//	  "message": "URL updated successfully"
//	}
//
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format"
// @Failure 404 {object} models.ErrorResponse[error] "Short URL not found"
// @Example response
//
//	{
//	  "error": "not_found",
//	  "message": "The specified short URL was not found",
//	  "details": {
//	    "error": "short URL not found"
//	  },
//	  "timestamp": "2023-06-08T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /shorten/{code} [put]
func (h *ShortenHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	code := c.Param("code")
	if code == "" {
		log.Warn().Msg("Missing short code in update request")
		utils.RespondBadRequest(c, nil, "Short code is required")
		return
	}

	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Str("code", code).Msg("Invalid request format for URL update")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().Str("code", code).Str("originalUrl", req.OriginalURL).
		Int("expiresAfter", req.ExpiresAfter).Msg("Updating shortened URL")

	result, err := h.service.Update(ctx, code, req)
	if err != nil {
		if err.Error() == "short URL not found" {
			log.Warn().Str("code", code).Msg("Short URL not found for update")
			utils.RespondNotFound(c, err, "The specified short URL was not found")
		} else {
			log.Error().Err(err).Str("code", code).Msg("Failed to update shortened URL")
			utils.RespondInternalError(c, err, "Failed to update shortened URL")
		}
		return
	}

	log.Info().Str("code", code).Str("shortUrl", result.ShortURL).Msg("Successfully updated shortened URL")

	utils.RespondOK(c, result, "URL updated successfully")
}

// Delete godoc
// @Summary Delete a shortened URL
// @Description Deletes an existing shortened URL by its short code
// @Tags shorten
// @Param code path string true "Short code identifier" example:"abc123"
// @Success 204 "No Content - URL successfully deleted"
// @Failure 400 {object} models.ErrorResponse[error] "Bad request - missing code parameter"
// @Example response
//
//	{
//	  "error": "bad_request",
//	  "message": "Short code is required",
//	  "details": {},
//	  "timestamp": "2023-06-08T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 404 {object} models.ErrorResponse[error] "Short URL not found"
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /shorten/{code} [delete]
func (h *ShortenHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	code := c.Param("code")
	if code == "" {
		log.Warn().Msg("Missing short code in delete request")
		utils.RespondBadRequest(c, nil, "Short code is required")
		return
	}

	log.Info().Str("code", code).Msg("Deleting shortened URL")

	err := h.service.Delete(ctx, code)
	if err != nil {
		if err.Error() == "short URL not found" {
			log.Warn().Str("code", code).Msg("Short URL not found for deletion")
			utils.RespondNotFound(c, err, "The specified short URL was not found")
		} else {
			log.Error().Err(err).Str("code", code).Msg("Failed to delete shortened URL")
			utils.RespondInternalError(c, err, "Failed to delete shortened URL")
		}
		return
	}

	log.Info().Str("code", code).Msg("Successfully deleted shortened URL")
	c.Status(http.StatusNoContent)
}

// Redirect godoc
// @Summary Redirect to original URL
// @Description Redirects to the original URL from a short code
// @Tags shorten
// @Param code path string true "Short code identifier" example:"abc123"
// @Success 302 "Found - Redirects to the original URL"
// @Header 302 {string} Location "The URL to redirect to"
// @Failure 400 {object} models.ErrorResponse[error] "Bad request - missing code parameter"
// @Failure 404 {object} models.ErrorResponse[error] "Short URL not found or has expired"
// @Example response
//
//	{
//	  "error": "not_found",
//	  "message": "The specified short URL was not found or has expired",
//	  "details": {
//	    "error": "shortened URL has expired"
//	  },
//	  "timestamp": "2023-06-08T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Router /shorten/{code} [get]
func (h *ShortenHandler) Redirect(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	code := c.Param("code")
	if code == "" {
		log.Warn().Msg("Missing short code in redirect request")
		utils.RespondBadRequest(c, nil, "Short code is required")
		return
	}

	log.Info().Str("code", code).Msg("Redirecting to original URL")

	url, err := h.service.GetOriginalURL(ctx, code)
	if err != nil {
		log.Warn().Err(err).Str("code", code).Msg("Failed to retrieve original URL for redirect")
		utils.RespondNotFound(c, err, "The specified short URL was not found or has expired")
		return
	}

	log.Info().Str("code", code).Str("originalUrl", url).Msg("Successfully redirecting to original URL")
	c.Redirect(http.StatusFound, url)
}

// GetByOriginalURL godoc
// @Summary Check if a URL is already shortened
// @Description Checks if an original URL already has a short code and optionally creates one if it doesn't exist
// @Tags shorten
// @Accept json
// @Produce json
// @Param request body models.GetByOriginalURLRequest true "Original URL to check"
// @Example request
//
//	{
//	  "originalUrl": "https://example.com/some/very/long/path",
//	  "createIfNotExists": true,
//	  "customCode": "mycode",
//	  "expiresAfter": 7
//	}
//
// @Success 200 {object} models.APIResponse[models.ShortenData] "Successfully retrieved shortened URL information"
// @Example response
//
//	{
//	  "success": true,
//	  "data": {
//	    "shorten": {
//	      "originalUrl": "https://example.com/some/very/long/path",
//	      "shortCode": "abc123",
//	      "expiresAt": "2023-06-15T10:30:45Z"
//	    }
//	  },
//	  "message": "URL information retrieved successfully"
//	}
//
// @Success 201 {object} models.APIResponse[models.ShortenData] "Successfully created new shortened URL"
// @Failure 400 {object} models.ErrorResponse[error] "Invalid request format"
// @Failure 404 {object} models.ErrorResponse[error] "Original URL not found and createIfNotExists is false"
// @Example response
//
//	{
//	  "error": "not_found",
//	  "message": "No shortened URL exists for the provided original URL",
//	  "details": {},
//	  "timestamp": "2023-06-08T10:30:45Z",
//	  "requestId": "c7f3305d-8c9a-4b9b-b701-3b9a1e36c1f0"
//	}
//
// @Failure 500 {object} models.ErrorResponse[error] "Server error"
// @Router /shorten/lookup [post]
func (h *ShortenHandler) GetByOriginalURL(c *gin.Context) {
	ctx := c.Request.Context()
	log := utils.LoggerFromContext(ctx)

	var req models.GetByOriginalURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request format for URL lookup")
		utils.RespondValidationError(c, err)
		return
	}

	log.Info().
		Str("originalUrl", req.OriginalURL).
		Bool("createIfNotExists", req.CreateIfNotExists).
		Msg("Looking up original URL")

	result, found, err := h.service.GetByOriginalUrl(ctx, req.OriginalURL)

	if err != nil {
		log.Error().Err(err).Str("originalUrl", req.OriginalURL).Msg("Failed to lookup URL")
		utils.RespondInternalError(c, err, "Failed to lookup URL")
		return
	}

	if !found {
		if req.CreateIfNotExists {
			// Create new shortened URL if requested
			shortenReq := models.ShortenRequest{
				OriginalURL:  req.OriginalURL,
				ExpiresAfter: req.ExpiresAfter,
				CustomCode:   req.CustomCode,
			}

			result, err = h.service.Create(ctx, shortenReq)
			if err != nil {
				log.Error().Err(err).Str("originalUrl", req.OriginalURL).Msg("Failed to create shortened URL")
				utils.RespondInternalError(c, err, "Failed to create shortened URL")
				return
			}

			log.Info().
				Str("originalUrl", req.OriginalURL).
				Str("shortCode", result.Shorten.ShortCode).
				Msg("Created new shortened URL during lookup")

			utils.RespondCreated(c, result, "New shortened URL created")
			return
		}

		// URL not found and no creation requested
		log.Info().Str("originalUrl", req.OriginalURL).Msg("No shortened URL found for original URL")
		utils.RespondNotFound(c, nil, "No shortened URL exists for the provided original URL")
		return
	}

	log.Info().
		Str("originalUrl", req.OriginalURL).
		Str("shortCode", result.Shorten.ShortCode).
		Msg("Successfully found shortened URL for original URL")

	utils.RespondOK(c, result, "URL details retrieved successfully")
}
