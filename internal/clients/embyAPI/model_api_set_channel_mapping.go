/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

type ApiSetChannelMapping struct {
	TunerChannelId string `json:"TunerChannelId,omitempty"`
	ProviderChannelId string `json:"ProviderChannelId,omitempty"`
}
