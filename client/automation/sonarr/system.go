package sonarr

import (
	"context"
	"fmt"

	sonarr "github.com/devopsarr/sonarr-go/sonarr"
	"suasor/utils"

	"suasor/client/automation/types"
)

func (s *SonarrClient) GetSystemStatus(ctx context.Context) (types.SystemStatus, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("baseURL", s.config.BaseURL).
		Msg("Retrieving system status from Sonarr server")

	// Call the Sonarr API
	log.Debug().Msg("Making API request to Sonarr server for system status")

	statusResult, resp, err := s.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", s.config.BaseURL).
			Str("apiEndpoint", "/system/status").
			Int("statusCode", 0).
			Msg("Failed to fetch system status from Sonarr")
		return types.SystemStatus{}, fmt.Errorf("failed to fetch system status: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("version", statusResult.GetVersion()).
		Msg("Successfully retrieved system status from Sonarr")

	// Convert to our internal type
	status := types.SystemStatus{
		Version:     statusResult.GetVersion(),
		StartupPath: statusResult.GetStartupPath(),
		AppData:     statusResult.GetAppData(),
		OsName:      statusResult.GetOsName(),
		Branch:      statusResult.GetBranch(),
	}

	return status, nil
}

func (s *SonarrClient) ExecuteCommand(ctx context.Context, command types.Command) (types.CommandResult, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", s.ClientID).
		Str("clientType", string(s.ClientType)).
		Str("commandName", command.Name).
		Msg("Executing command in Sonarr")

	// Create command
	newCommand := sonarr.NewCommandResource()
	newCommand.SetName(command.Name)

	// Execute command
	cmdResult, resp, err := s.client.CommandAPI.CreateCommand(ctx).CommandResource(*newCommand).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("commandName", command.Name).
			Msg("Failed to execute command in Sonarr")
		return types.CommandResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("commandId", cmdResult.GetId()).
		Str("commandName", cmdResult.GetName()).
		Str("status", string(cmdResult.GetStatus())).
		Msg("Successfully initiated command in Sonarr")

	return types.CommandResult{
		ID:        int64(cmdResult.GetId()),
		Name:      cmdResult.GetName(),
		Status:    string(cmdResult.GetStatus()),
		StartedAt: cmdResult.GetStarted(),
	}, nil
}
