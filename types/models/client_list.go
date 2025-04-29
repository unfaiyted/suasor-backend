// models/automation_client.go
package models

import (
	client "suasor/clients/types"
)

type ClientList struct {
	MediaClientList      `json:"media"`
	AutomationClientList `json:"automation"`
	AIClientList         `json:"ai"`

	Total int `json:"total"`
}

func NewClientList() *ClientList {
	return &ClientList{
		MediaClientList: MediaClientList{
			Emby:     map[uint64]*Client[*client.EmbyConfig]{},
			Jellyfin: map[uint64]*Client[*client.JellyfinConfig]{},
			Plex:     map[uint64]*Client[*client.PlexConfig]{},
			Subsonic: map[uint64]*Client[*client.SubsonicConfig]{},
		},
		AutomationClientList: AutomationClientList{
			Sonarr: map[uint64]*Client[*client.SonarrConfig]{},
			Radarr: map[uint64]*Client[*client.RadarrConfig]{},
			Lidarr: map[uint64]*Client[*client.LidarrConfig]{},
		},
		AIClientList: AIClientList{
			Claude: map[uint64]*Client[*client.ClaudeConfig]{},
			OpenAI: map[uint64]*Client[*client.OpenAIConfig]{},
			Ollama: map[uint64]*Client[*client.OllamaConfig]{},
		},
		Total: 0,
	}
}

func (c *ClientList) GetTotal() int {
	return c.Total
}

func (c *ClientList) GetSonarr() map[uint64]*Client[*client.SonarrConfig] {
	return c.Sonarr
}

func (c *ClientList) GetRadarr() map[uint64]*Client[*client.RadarrConfig] {
	return c.Radarr
}

func (c *ClientList) GetLidarr() map[uint64]*Client[*client.LidarrConfig] {
	return c.Lidarr
}

func (c *ClientList) GetClaude() map[uint64]*Client[*client.ClaudeConfig] {
	return c.Claude
}

func (c *ClientList) GetOpenAI() map[uint64]*Client[*client.OpenAIConfig] {
	return c.OpenAI
}

func (c *ClientList) GetOllama() map[uint64]*Client[*client.OllamaConfig] {
	return c.Ollama
}

type MediaClientList struct {
	Emby     map[uint64]*Client[*client.EmbyConfig]     `json:"emby"`
	Jellyfin map[uint64]*Client[*client.JellyfinConfig] `json:"jellyfin"`
	Plex     map[uint64]*Client[*client.PlexConfig]     `json:"plex"`
	Subsonic map[uint64]*Client[*client.SubsonicConfig] `json:"subsonic"`

	IDs map[uint64]client.ClientType `json:"ids"`

	Total int
}

func NewMediaClientList() *MediaClientList {
	return &MediaClientList{
		Emby:     map[uint64]*Client[*client.EmbyConfig]{},
		Jellyfin: map[uint64]*Client[*client.JellyfinConfig]{},
		Plex:     map[uint64]*Client[*client.PlexConfig]{},
		Subsonic: map[uint64]*Client[*client.SubsonicConfig]{},

		Total: 0,
	}
}

func (c *MediaClientList) GetTotal() int {
	return c.Total
}

func (c *MediaClientList) GetClientType(clientID uint64) (client.ClientType, bool) {
	if clientType, ok := c.IDs[clientID]; ok {
		return clientType, true
	}
	return "", false
}

func (c *MediaClientList) GetClientConfig(clientID uint64, clientType client.ClientType) client.ClientConfig {

	if clientType == client.ClientTypeEmby {
		return c.Emby[clientID].Config
	}
	if clientType == client.ClientTypeJellyfin {
		return c.Jellyfin[clientID].Config
	}
	if clientType == client.ClientTypePlex {
		return c.Plex[clientID].Config
	}
	if clientType == client.ClientTypeSubsonic {
		return c.Subsonic[clientID].Config
	}
	return nil
}

func (c *MediaClientList) AddEmby(client *Client[*client.EmbyConfig]) {
	c.Emby[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *MediaClientList) AddEmbyArray(clients []*Client[*client.EmbyConfig]) {
	for _, client := range clients {
		c.AddEmby(client)
	}
	c.Total += len(clients)
}

func (c *MediaClientList) AddJellyfin(client *Client[*client.JellyfinConfig]) {
	c.Jellyfin[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *MediaClientList) AddJellyfinArray(clients []*Client[*client.JellyfinConfig]) {
	for _, client := range clients {
		c.AddJellyfin(client)
	}
	c.Total += len(clients)
}

func (c *MediaClientList) AddPlex(client *Client[*client.PlexConfig]) {
	c.Plex[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *MediaClientList) AddPlexArray(clients []*Client[*client.PlexConfig]) {
	for _, client := range clients {
		c.AddPlex(client)
	}
	c.Total += len(clients)
}

func (c *MediaClientList) AddSubsonic(client *Client[*client.SubsonicConfig]) {
	c.Subsonic[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *MediaClientList) AddSubsonicArray(clients []*Client[*client.SubsonicConfig]) {
	for _, client := range clients {
		c.AddSubsonic(client)
	}
	c.Total += len(clients)
}

func (c *MediaClientList) GetEmby() map[uint64]*Client[*client.EmbyConfig] {
	return c.Emby
}

func (c *MediaClientList) GetEmbyArray() []*Client[*client.EmbyConfig] {
	var clients []*Client[*client.EmbyConfig]
	for _, client := range c.Emby {
		clients = append(clients, client)
	}
	return clients
}

func (c *MediaClientList) GetJellyfin() map[uint64]*Client[*client.JellyfinConfig] {
	return c.Jellyfin
}

func (c *MediaClientList) GetJellyfinArray() []*Client[*client.JellyfinConfig] {
	var clients []*Client[*client.JellyfinConfig]
	for _, client := range c.Jellyfin {
		clients = append(clients, client)
	}
	return clients
}

func (c *MediaClientList) GetPlex() map[uint64]*Client[*client.PlexConfig] {
	return c.Plex
}

func (c *MediaClientList) GetPlexArray() []*Client[*client.PlexConfig] {
	var clients []*Client[*client.PlexConfig]
	for _, client := range c.Plex {
		clients = append(clients, client)
	}
	return clients
}

func (c *MediaClientList) GetSubsonic() map[uint64]*Client[*client.SubsonicConfig] {
	return c.Subsonic
}

func (c *MediaClientList) GetSubsonicArray() []*Client[*client.SubsonicConfig] {
	var clients []*Client[*client.SubsonicConfig]
	for _, client := range c.Subsonic {
		clients = append(clients, client)
	}
	return clients
}

type AIClientList struct {
	Claude map[uint64]*Client[*client.ClaudeConfig] `json:"claude"`
	OpenAI map[uint64]*Client[*client.OpenAIConfig] `json:"openai"`
	Ollama map[uint64]*Client[*client.OllamaConfig] `json:"ollama"`

	IDs   map[uint64]client.ClientType `json:"ids"`
	Total int                          `json:"total"`
}

func NewAIClientList() *AIClientList {
	return &AIClientList{
		Claude: map[uint64]*Client[*client.ClaudeConfig]{},
		OpenAI: map[uint64]*Client[*client.OpenAIConfig]{},
		Ollama: map[uint64]*Client[*client.OllamaConfig]{},
		Total:  0,
	}
}

func (c *AIClientList) AddClaude(client *Client[*client.ClaudeConfig]) {
	c.Claude[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *AIClientList) AddClaudeArray(clients []*Client[*client.ClaudeConfig]) {
	for _, client := range clients {
		c.AddClaude(client)
	}
	c.Total += len(clients)
}

func (c *AIClientList) AddOpenAI(client *Client[*client.OpenAIConfig]) {
	c.OpenAI[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *AIClientList) AddOpenAIArray(clients []*Client[*client.OpenAIConfig]) {
	for _, client := range clients {
		c.AddOpenAI(client)
	}
	c.Total += len(clients)
}

func (c *AIClientList) AddOllama(client *Client[*client.OllamaConfig]) {
	c.Ollama[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *AIClientList) AddOllamaArray(clients []*Client[*client.OllamaConfig]) {
	for _, client := range clients {
		c.AddOllama(client)
	}
	c.Total += len(clients)
}

type AutomationClientList struct {
	Sonarr map[uint64]*Client[*client.SonarrConfig] `json:"sonarr"`
	Radarr map[uint64]*Client[*client.RadarrConfig] `json:"radarr"`
	Lidarr map[uint64]*Client[*client.LidarrConfig] `json:"lidarr"`

	IDs map[uint64]client.ClientType `json:"ids"`

	Total int `json:"total"`
}

func NewAutomationClientList() *AutomationClientList {
	return &AutomationClientList{
		Sonarr: map[uint64]*Client[*client.SonarrConfig]{},
		Radarr: map[uint64]*Client[*client.RadarrConfig]{},
		Lidarr: map[uint64]*Client[*client.LidarrConfig]{},
		Total:  0,
	}
}

func (c *AutomationClientList) AddSonarr(client *Client[*client.SonarrConfig]) {
	c.Sonarr[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *AutomationClientList) AddSonarrArray(clients []*Client[*client.SonarrConfig]) {
	for _, client := range clients {
		c.AddSonarr(client)
	}
	c.Total += len(clients)
}

func (c *AutomationClientList) AddRadarr(client *Client[*client.RadarrConfig]) {
	c.Radarr[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *AutomationClientList) AddRadarrArray(clients []*Client[*client.RadarrConfig]) {
	for _, client := range clients {
		c.AddRadarr(client)
	}
	c.Total += len(clients)
}

func (c *AutomationClientList) AddLidarr(client *Client[*client.LidarrConfig]) {
	c.Lidarr[client.ID] = client
	c.IDs[client.ID] = client.GetType()
	c.Total++
}

func (c *AutomationClientList) AddLidarrArray(clients []*Client[*client.LidarrConfig]) {
	for _, client := range clients {
		c.AddLidarr(client)
	}
	c.Total += len(clients)
}

type MetadataClientList struct {
	// Tmdb map[uint64]Client[*client.TmdbConfig]
}
