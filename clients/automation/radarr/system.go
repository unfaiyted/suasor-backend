package radarr

import (
	"context"
	"fmt"

	radarr "github.com/devopsarr/radarr-go/radarr"
	"suasor/clients/automation/types"
	"suasor/utils/logger"
)

// GetSystemStatus retrieves system information from Radarr
func (r *RadarrClient) GetSystemStatus(ctx context.Context) (types.SystemStatus, error) {
	// Get logger from context
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("baseURL", r.config.BaseURL).
		Msg("Retrieving system status from Radarr server")

	// Call the Radarr API
	log.Debug().Msg("Making API request to Radarr server for system status")

	statusResult, resp, err := r.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", r.config.BaseURL).
			Str("apiEndpoint", "/system/status").
			Int("statusCode", 0).
			Msg("Failed to fetch system status from Radarr")
		return types.SystemStatus{}, fmt.Errorf("failed to fetch system status: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("version", statusResult.GetVersion()).
		Msg("Successfully retrieved system status from Radarr")

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

func (r *RadarrClient) ExecuteCommand(ctx context.Context, command types.Command) (types.CommandResult, error) {
	log := logger.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", r.ClientID).
		Str("clientType", string(r.ClientType)).
		Str("commandName", command.Name).
		Msg("Executing command in Radarr")

	// Create command
	newCommand := radarr.NewCommandResource()
	newCommand.SetName(command.Name)

	// Add command-specific parameters
	// if command.Parameters != nil {
	// 	// Convert map to appropriate format if needed
	// 	switch command.Name {
	// 	case "MoviesSearch":
	// 		if movieIds, ok := command.Parameters["movieIds"].([]int64); ok {
	// 			int32Ids := make([]int32, len(movieIds))
	// 			for i, id := range movieIds {
	// 				int32Ids[i] = int32(id)
	// 			}
	// 			newCommand.SetMovieIds(int32Ids)
	// 		}
	// 		// Add other command-specific parameter handling as needed
	// 	}
	// }

	// Execute command
	cmdResult, resp, err := r.client.CommandAPI.CreateCommand(ctx).CommandResource(*newCommand).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("commandName", command.Name).
			Msg("Failed to execute command in Radarr")
		return types.CommandResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("commandId", cmdResult.GetId()).
		Str("commandName", cmdResult.GetName()).
		Str("status", string(cmdResult.GetStatus())).
		Msg("Successfully initiated command in Radarr")

	return types.CommandResult{
		ID:        int64(cmdResult.GetId()),
		Name:      cmdResult.GetName(),
		Status:    string(cmdResult.GetStatus()),
		StartedAt: cmdResult.GetStarted(),
	}, nil
}
