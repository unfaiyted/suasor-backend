// models/automation_client.go
package models

import (
	client "suasor/clients/types"
)

type ClientList struct {
	emby     []*Client[*client.EmbyConfig]
	jellyfin []*Client[*client.JellyfinConfig]
	plex     []*Client[*client.PlexConfig]
	subsonic []*Client[*client.SubsonicConfig]
	sonarr   []*Client[*client.SonarrConfig]
	radarr   []*Client[*client.RadarrConfig]
	lidarr   []*Client[*client.LidarrConfig]
	claude   []*Client[*client.ClaudeConfig]
	openAI   []*Client[*client.OpenAIConfig]
	ollama   []*Client[*client.OllamaConfig]

	total int
}

func (c *ClientList) AddEmby(client *Client[*client.EmbyConfig]) {
	c.emby = append(c.emby, client)
	c.total++
}

func (c *ClientList) AddEmbyArray(clients []*Client[*client.EmbyConfig]) {
	c.emby = append(c.emby, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddJellyfin(client *Client[*client.JellyfinConfig]) {
	c.jellyfin = append(c.jellyfin, client)
	c.total++
}

func (c *ClientList) AddJellyfinArray(clients []*Client[*client.JellyfinConfig]) {
	c.jellyfin = append(c.jellyfin, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddPlex(client *Client[*client.PlexConfig]) {
	c.plex = append(c.plex, client)
	c.total++
}

func (c *ClientList) AddPlexArray(clients []*Client[*client.PlexConfig]) {
	c.plex = append(c.plex, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddSubsonic(client *Client[*client.SubsonicConfig]) {
	c.subsonic = append(c.subsonic, client)
	c.total++
}

func (c *ClientList) AddSubsonicArray(clients []*Client[*client.SubsonicConfig]) {
	c.subsonic = append(c.subsonic, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddSonarr(client *Client[*client.SonarrConfig]) {
	c.sonarr = append(c.sonarr, client)
	c.total++
}

func (c *ClientList) AddSonarrArray(clients []*Client[*client.SonarrConfig]) {
	c.sonarr = append(c.sonarr, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddRadarr(client *Client[*client.RadarrConfig]) {
	c.radarr = append(c.radarr, client)
	c.total++
}

func (c *ClientList) AddRadarrArray(clients []*Client[*client.RadarrConfig]) {
	c.radarr = append(c.radarr, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddLidarr(client *Client[*client.LidarrConfig]) {
	c.lidarr = append(c.lidarr, client)
	c.total++
}

func (c *ClientList) AddLidarrArray(clients []*Client[*client.LidarrConfig]) {
	c.lidarr = append(c.lidarr, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddClaude(client *Client[*client.ClaudeConfig]) {
	c.claude = append(c.claude, client)
	c.total++
}

func (c *ClientList) AddClaudeArray(clients []*Client[*client.ClaudeConfig]) {
	c.claude = append(c.claude, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddOpenAI(client *Client[*client.OpenAIConfig]) {
	c.openAI = append(c.openAI, client)
	c.total++
}

func (c *ClientList) AddOpenAIArray(clients []*Client[*client.OpenAIConfig]) {
	c.openAI = append(c.openAI, clients...)
	c.total += len(clients)
}

func (c *ClientList) AddOllama(client *Client[*client.OllamaConfig]) {
	c.ollama = append(c.ollama, client)
	c.total++
}

func (c *ClientList) AddOllamaArray(clients []*Client[*client.OllamaConfig]) {
	c.ollama = append(c.ollama, clients...)
	c.total += len(clients)
}

func (c *ClientList) GetTotal() int {
	return c.total
}

func (c *ClientList) GetEmby() []*Client[*client.EmbyConfig] {
	return c.emby
}

func (c *ClientList) GetJellyfin() []*Client[*client.JellyfinConfig] {
	return c.jellyfin
}

func (c *ClientList) GetPlex() []*Client[*client.PlexConfig] {
	return c.plex
}

func (c *ClientList) GetSubsonic() []*Client[*client.SubsonicConfig] {
	return c.subsonic
}

func (c *ClientList) GetSonarr() []*Client[*client.SonarrConfig] {
	return c.sonarr
}

func (c *ClientList) GetRadarr() []*Client[*client.RadarrConfig] {
	return c.radarr
}

func (c *ClientList) GetLidarr() []*Client[*client.LidarrConfig] {
	return c.lidarr
}

func (c *ClientList) GetClaude() []*Client[*client.ClaudeConfig] {
	return c.claude
}

func (c *ClientList) GetOpenAI() []*Client[*client.OpenAIConfig] {
	return c.openAI
}

func (c *ClientList) GetOllama() []*Client[*client.OllamaConfig] {
	return c.ollama
}
