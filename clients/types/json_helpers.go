package types

import (
	"encoding/json"
	"fmt"
)

// UnmarshalConfigJSON is a helper function to unmarshal JSON into config structs
// with embedded interfaces like ClientConfig or ClientMediaConfig
func UnmarshalConfigJSON(data []byte, config any) error {
	// Extract the details field as JSON raw message
	var temp struct {
		Details json.RawMessage `json:"details"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("error unmarshaling details: %w", err)
	}

	// Process based on config type
	switch c := config.(type) {
	// Media clients
	case *EmbyConfig:
		return unmarshalMediaClientConfig(data, temp.Details, c)
	case *JellyfinConfig:
		return unmarshalMediaClientConfig(data, temp.Details, c)
	case *PlexConfig:
		return unmarshalMediaClientConfig(data, temp.Details, c)
	case *SubsonicConfig:
		return unmarshalMediaClientConfig(data, temp.Details, c)

	// Automation clients
	case *RadarrConfig:
		return unmarshalAutomationConfig(data, temp.Details, c)
	case *SonarrConfig:
		return unmarshalAutomationConfig(data, temp.Details, c)
	case *LidarrConfig:
		return unmarshalAutomationConfig(data, temp.Details, c)

	// AI clients
	case *ClaudeConfig:
		return unmarshalAIConfig(data, temp.Details, c)
	case *OpenAIConfig:
		return unmarshalAIConfig(data, temp.Details, c)
	case *OllamaConfig:
		return unmarshalAIConfig(data, temp.Details, c)

	// Metadata clients
	case *TMDBConfig:
		return unmarshalMetadataConfig(data, temp.Details, c)

	default:
		return fmt.Errorf("unsupported config type: %T", config)
	}
}

// Helper for unmarshal media client configs
func unmarshalMediaClientConfig[T interface {
	*EmbyConfig | *JellyfinConfig | *PlexConfig | *SubsonicConfig
}](data []byte, details json.RawMessage, config T) error {
	// First structure for root-level fields
	baseStruct := struct {
		UserID   string `json:"userID,omitempty"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Token    string `json:"token,omitempty"`
	}{}

	if err := json.Unmarshal(data, &baseStruct); err != nil {
		return fmt.Errorf("error unmarshaling base fields: %w", err)
	}

	// Second structure for API payload format where fields are in the config object
	configStruct := struct {
		Config struct {
			Username string `json:"username,omitempty"`
			UserID   string `json:"userID,omitempty"`
			Password string `json:"password,omitempty"`
			Token    string `json:"token,omitempty"`
		} `json:"config,omitempty"`
	}{}

	// Try to unmarshal for API format where fields are inside config
	_ = json.Unmarshal(data, &configStruct)

	// Use config fields if root fields are empty
	if baseStruct.Username == "" && configStruct.Config.Username != "" {
		baseStruct.Username = configStruct.Config.Username
	}
	if baseStruct.UserID == "" && configStruct.Config.UserID != "" {
		baseStruct.UserID = configStruct.Config.UserID
	}
	if baseStruct.Password == "" && configStruct.Config.Password != "" {
		baseStruct.Password = configStruct.Config.Password
	}
	if baseStruct.Token == "" && configStruct.Config.Token != "" {
		baseStruct.Token = configStruct.Config.Token
	}

	// Now handle the details which contains the ClientMediaConfig
	if len(details) > 0 {
		// Create a proper clientMediaConfig
		mediaConfig := &clientMediaConfig{}

		if err := json.Unmarshal(details, mediaConfig); err != nil {
			return fmt.Errorf("error unmarshaling media config: %w", err)
		}

		// Now set the fields based on the config type
		switch c := any(config).(type) {
		case *EmbyConfig:
			c.ClientMediaConfig = mediaConfig
			c.UserID = baseStruct.UserID
			c.Username = baseStruct.Username
		case *JellyfinConfig:
			c.ClientMediaConfig = mediaConfig
			c.UserID = baseStruct.UserID
			c.Username = baseStruct.Username
		case *PlexConfig:
			c.ClientMediaConfig = mediaConfig
			c.Username = baseStruct.Username
			c.Token = baseStruct.Token
		case *SubsonicConfig:
			c.ClientMediaConfig = mediaConfig
			c.Username = baseStruct.Username
			c.Password = baseStruct.Password
		}
	}

	return nil
}

// Helper for unmarshal automation client configs
func unmarshalAutomationConfig[T interface {
	*RadarrConfig | *SonarrConfig | *LidarrConfig
}](data []byte, details json.RawMessage, config T) error {
	// Create a temporary struct to handle basic fields
	baseStruct := struct {
		UserID   string `json:"userID,omitempty"`
		Username string `json:"username"`
	}{}

	if err := json.Unmarshal(data, &baseStruct); err != nil {
		return fmt.Errorf("error unmarshaling base fields: %w", err)
	}

	// Now handle the details which contains the ClientAutomationConfig
	if len(details) > 0 {
		// Create a proper clientAutomationConfig
		automationConfig := &clientAutomationConfig{}

		if err := json.Unmarshal(details, automationConfig); err != nil {
			return fmt.Errorf("error unmarshaling automation config: %w", err)
		}

		// Now set the fields based on the config type
		switch c := any(config).(type) {
		case *RadarrConfig:
			c.ClientAutomationConfig = automationConfig
		case *SonarrConfig:
			c.ClientAutomationConfig = automationConfig
		case *LidarrConfig:
			c.ClientAutomationConfig = automationConfig
		}
	}

	return nil
}

// Helper for unmarshal AI client configs
func unmarshalAIConfig[T interface {
	*ClaudeConfig | *OpenAIConfig | *OllamaConfig
}](data []byte, details json.RawMessage, config T) error {
	// Create a temporary struct to handle basic fields
	baseStruct := struct{}{}

	if err := json.Unmarshal(data, &baseStruct); err != nil {
		return fmt.Errorf("error unmarshaling base fields: %w", err)
	}

	// Now handle the details which contains the ClientAIConfig
	if len(details) > 0 {
		// Create a proper clientAIConfig
		aiConfig := &clientAIConfig{}

		if err := json.Unmarshal(details, aiConfig); err != nil {
			return fmt.Errorf("error unmarshaling AI config: %w", err)
		}

		// Now set the fields based on the config type
		switch c := any(config).(type) {
		case *ClaudeConfig:
			c.AIClientConfig = aiConfig
		case *OpenAIConfig:
			c.AIClientConfig = aiConfig
		case *OllamaConfig:
			c.AIClientConfig = aiConfig
		}
	}

	return nil
}

// Helper for unmarshal metadata client configs
func unmarshalMetadataConfig[T interface {
	*TMDBConfig
}](data []byte, details json.RawMessage, config T) error {
	// Create a temporary struct to handle basic fields
	baseStruct := struct{}{}

	if err := json.Unmarshal(data, &baseStruct); err != nil {
		return fmt.Errorf("error unmarshaling base fields: %w", err)
	}

	// Now handle the details which contains the ClientMetadataConfig
	if len(details) > 0 {
		// Create a proper clientMetadataConfig
		metadataConfig := &clientMetadataConfig{}

		if err := json.Unmarshal(details, metadataConfig); err != nil {
			return fmt.Errorf("error unmarshaling metadata config: %w", err)
		}

		// Now set the fields based on the config type
		switch c := any(config).(type) {
		case *TMDBConfig:
			c.ClientMetadataConfig = metadataConfig
		}
	}

	return nil
}
