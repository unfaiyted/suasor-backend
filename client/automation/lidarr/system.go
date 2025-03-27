package lidarr

import (
	"context"
	"fmt"

	lidarr "github.com/devopsarr/lidarr-go/lidarr"
	"suasor/client/automation/types"
	"suasor/utils"
)

func (l *LidarrClient) GetSystemStatus(ctx context.Context) (types.SystemStatus, error) {
	// Get logger from context
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("baseURL", l.config.BaseURL).
		Msg("Retrieving system status from Lidarr server")

	// Call the Lidarr API
	log.Debug().Msg("Making API request to Lidarr server for system status")

	statusResult, resp, err := l.client.SystemAPI.GetSystemStatus(ctx).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", l.config.BaseURL).
			Str("apiEndpoint", "/system/status").
			Int("statusCode", 0).
			Msg("Failed to fetch system status from Lidarr")
		return types.SystemStatus{}, fmt.Errorf("failed to fetch system status: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Str("version", statusResult.GetVersion()).
		Msg("Successfully retrieved system status from Lidarr")

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

func (l *LidarrClient) ExecuteCommand(ctx context.Context, command types.Command) (types.CommandResult, error) {
	log := utils.LoggerFromContext(ctx)

	log.Info().
		Uint64("clientID", l.ClientID).
		Str("clientType", string(l.ClientType)).
		Str("commandName", command.Name).
		Msg("Executing command in Lidarr")

	// Create command
	newCommand := lidarr.NewCommandResource()
	newCommand.SetName(command.Name)

	// Execute command
	cmdResult, resp, err := l.client.CommandAPI.CreateCommand(ctx).CommandResource(*newCommand).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("commandName", command.Name).
			Msg("Failed to execute command in Lidarr")
		return types.CommandResult{}, fmt.Errorf("failed to execute command: %w", err)
	}

	log.Info().
		Int("statusCode", resp.StatusCode).
		Int32("commandId", cmdResult.GetId()).
		Str("commandName", cmdResult.GetName()).
		Str("status", string(cmdResult.GetStatus())).
		Msg("Successfully initiated command in Lidarr")

	return types.CommandResult{
		ID:        int64(cmdResult.GetId()),
		Name:      cmdResult.GetName(),
		Status:    string(cmdResult.GetStatus()),
		StartedAt: cmdResult.GetStarted(),
	}, nil
}
