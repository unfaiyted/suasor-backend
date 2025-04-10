package jobs_test

import (
	"suasor/services/jobs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsContentTypeEnabled tests the IsContentTypeEnabled method
func TestIsContentTypeEnabled(t *testing.T) {
	job := &jobs.RecommendationJob{}

	// Test with empty content types (should return true for all)
	assert.True(t, job.IsContentTypeEnabled("", "movie"))
	assert.True(t, job.IsContentTypeEnabled("", "series"))
	assert.True(t, job.IsContentTypeEnabled("", "music"))

	// Test with specific content types
	assert.True(t, job.IsContentTypeEnabled("movie", "movie"))
	assert.False(t, job.IsContentTypeEnabled("movie", "series"))
	assert.False(t, job.IsContentTypeEnabled("movie", "music"))

	// Test with multiple content types
	assert.True(t, job.IsContentTypeEnabled("movie,series", "movie"))
	assert.True(t, job.IsContentTypeEnabled("movie,series", "series"))
	assert.False(t, job.IsContentTypeEnabled("movie,series", "music"))

	// Test with spaces
	assert.True(t, job.IsContentTypeEnabled("movie, series", "series"))
	assert.True(t, job.IsContentTypeEnabled(" movie , series ", "movie"))
}