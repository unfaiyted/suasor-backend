package automation

import (
	// "context"
	"errors"
	// "fmt"
	base "suasor/client"
	// media "suasor/client/media/types"
	types "suasor/client/types"
	// "suasor/types/models"
)

var ErrFeatureNotSupported = errors.New("feature not supported by this automation client")

type AutomationClient interface {
	SupportsMovies() bool
	SupportsTVShows() bool
	SupportsMusic() bool
}

type BaseAutomationClient struct {
	base.BaseClient
	ClientType types.AutomationClientType
}

// Default caity implementations (all false by default)
func (m *BaseAutomationClient) SupportsMovies() bool  { return false }
func (m *BaseAutomationClient) SupportsTVShows() bool { return false }
func (m *BaseAutomationClient) SupportsMusic() bool   { return false }
