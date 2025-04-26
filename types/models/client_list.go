// models/automation_client.go
package models

import (
	client "suasor/clients/types"
)

type ClientList struct {
	Emby     []*Client[*client.EmbyConfig]     `json:"emby"`
	Jellyfin []*Client[*client.JellyfinConfig] `json:"jellyfin"`
	Plex     []*Client[*client.PlexConfig]     `json:"plex"`
	Subsonic []*Client[*client.SubsonicConfig] `json:"subsonic"`
	Sonarr   []*Client[*client.SonarrConfig]   `json:"sonarr"`
	Radarr   []*Client[*client.RadarrConfig]   `json:"radarr"`
	Lidarr   []*Client[*client.LidarrConfig]   `json:"lidarr"`
	Claude   []*Client[*client.ClaudeConfig]   `json:"claude"`
	OpenAI   []*Client[*client.OpenAIConfig]   `json:"openai"`
	Ollama   []*Client[*client.OllamaConfig]   `json:"ollama"`

	Total int `json:"total"`
}

func (c *ClientList) AddEmby(client *Client[*client.EmbyConfig]) {
	c.Emby = append(c.Emby, client)
	c.Total++
}

func (c *ClientList) AddEmbyArray(clients []*Client[*client.EmbyConfig]) {
	c.Emby = append(c.Emby, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddJellyfin(client *Client[*client.JellyfinConfig]) {
	c.Jellyfin = append(c.Jellyfin, client)
	c.Total++
}

func (c *ClientList) AddJellyfinArray(clients []*Client[*client.JellyfinConfig]) {
	c.Jellyfin = append(c.Jellyfin, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddPlex(client *Client[*client.PlexConfig]) {
	c.Plex = append(c.Plex, client)
	c.Total++
}

func (c *ClientList) AddPlexArray(clients []*Client[*client.PlexConfig]) {
	c.Plex = append(c.Plex, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddSubsonic(client *Client[*client.SubsonicConfig]) {
	c.Subsonic = append(c.Subsonic, client)
	c.Total++
}

func (c *ClientList) AddSubsonicArray(clients []*Client[*client.SubsonicConfig]) {
	c.Subsonic = append(c.Subsonic, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddSonarr(client *Client[*client.SonarrConfig]) {
	c.Sonarr = append(c.Sonarr, client)
	c.Total++
}

func (c *ClientList) AddSonarrArray(clients []*Client[*client.SonarrConfig]) {
	c.Sonarr = append(c.Sonarr, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddRadarr(client *Client[*client.RadarrConfig]) {
	c.Radarr = append(c.Radarr, client)
	c.Total++
}

func (c *ClientList) AddRadarrArray(clients []*Client[*client.RadarrConfig]) {
	c.Radarr = append(c.Radarr, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddLidarr(client *Client[*client.LidarrConfig]) {
	c.Lidarr = append(c.Lidarr, client)
	c.Total++
}

func (c *ClientList) AddLidarrArray(clients []*Client[*client.LidarrConfig]) {
	c.Lidarr = append(c.Lidarr, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddClaude(client *Client[*client.ClaudeConfig]) {
	c.Claude = append(c.Claude, client)
	c.Total++
}

func (c *ClientList) AddClaudeArray(clients []*Client[*client.ClaudeConfig]) {
	c.Claude = append(c.Claude, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddOpenAI(client *Client[*client.OpenAIConfig]) {
	c.OpenAI = append(c.OpenAI, client)
	c.Total++
}

func (c *ClientList) AddOpenAIArray(clients []*Client[*client.OpenAIConfig]) {
	c.OpenAI = append(c.OpenAI, clients...)
	c.Total += len(clients)
}

func (c *ClientList) AddOllama(client *Client[*client.OllamaConfig]) {
	c.Ollama = append(c.Ollama, client)
	c.Total++
}

func (c *ClientList) AddOllamaArray(clients []*Client[*client.OllamaConfig]) {
	c.Ollama = append(c.Ollama, clients...)
	c.Total += len(clients)
}

func (c *ClientList) GetTotal() int {
	return c.Total
}

func (c *ClientList) GetEmby() []*Client[*client.EmbyConfig] {
	return c.Emby
}

func (c *ClientList) GetJellyfin() []*Client[*client.JellyfinConfig] {
	return c.Jellyfin
}

func (c *ClientList) GetPlex() []*Client[*client.PlexConfig] {
	return c.Plex
}

func (c *ClientList) GetSubsonic() []*Client[*client.SubsonicConfig] {
	return c.Subsonic
}

func (c *ClientList) GetSonarr() []*Client[*client.SonarrConfig] {
	return c.Sonarr
}

func (c *ClientList) GetRadarr() []*Client[*client.RadarrConfig] {
	return c.Radarr
}

func (c *ClientList) GetLidarr() []*Client[*client.LidarrConfig] {
	return c.Lidarr
}

func (c *ClientList) GetClaude() []*Client[*client.ClaudeConfig] {
	return c.Claude
}

func (c *ClientList) GetOpenAI() []*Client[*client.OpenAIConfig] {
	return c.OpenAI
}

func (c *ClientList) GetOllama() []*Client[*client.OllamaConfig] {
	return c.Ollama
}
