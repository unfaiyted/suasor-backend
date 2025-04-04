
/*
 * Emby Server REST API
 *
 * Explore the Emby Server API
 *
 */
package embyclient

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"github.com/antihax/optional"
)

// Linger please
var (
	_ context.Context
)

type AudioServiceApiService service
/*
AudioServiceApiService Gets an audio stream
Requires authentication as user
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param streamFileName
 * @param id Item Id
 * @param container Container
 * @param optional nil or *AudioServiceApiGetAudioByIdByStreamfilenameOpts - Optional Parameters:
     * @param "DeviceProfileId" (optional.String) -  Optional. The dlna device profile id to utilize.
     * @param "DeviceId" (optional.String) -  The device id of the client requesting. Used to stop encoding processes when needed.
     * @param "AudioCodec" (optional.String) -  Optional. Specify a audio codec to encode to, e.g. mp3. If omitted the server will auto-select using the url&#x27;s extension. Options: aac, mp3, vorbis, wma.
     * @param "EnableAutoStreamCopy" (optional.Bool) -  Whether or not to allow automatic stream copy if requested values match the original source. Defaults to true.
     * @param "AudioSampleRate" (optional.Int32) -  Optional. Specify a specific audio sample rate, e.g. 44100
     * @param "AudioBitRate" (optional.Int32) -  Optional. Specify an audio bitrate to encode to, e.g. 128000. If omitted this will be left to encoder defaults.
     * @param "AudioChannels" (optional.Int32) -  Optional. Specify a specific number of audio channels to encode to, e.g. 2
     * @param "MaxAudioChannels" (optional.Int32) -  Optional. Specify a maximum number of audio channels to encode to, e.g. 2
     * @param "Static" (optional.Bool) -  Optional. If true, the original file will be streamed statically without any encoding. Use either no url extension or the original file extension. true/false
     * @param "CopyTimestamps" (optional.Bool) -  Whether or not to copy timestamps when transcoding with an offset. Defaults to false.
     * @param "StartTimeTicks" (optional.Int64) -  Optional. Specify a starting offset, in ticks. 1ms &#x3D; 10000 ticks.
     * @param "Width" (optional.Int32) -  Optional. The fixed horizontal resolution of the encoded video.
     * @param "Height" (optional.Int32) -  Optional. The fixed vertical resolution of the encoded video.
     * @param "MaxWidth" (optional.Int32) -  Optional. The maximum horizontal resolution of the encoded video.
     * @param "MaxHeight" (optional.Int32) -  Optional. The maximum vertical resolution of the encoded video.
     * @param "VideoBitRate" (optional.Int32) -  Optional. Specify a video bitrate to encode to, e.g. 500000. If omitted this will be left to encoder defaults.
     * @param "SubtitleStreamIndex" (optional.Int32) -  Optional. The index of the subtitle stream to use. If omitted no subtitles will be used.
     * @param "SubtitleMethod" (optional.Interface of SubtitleDeliveryMethod) -  Optional. Specify the subtitle delivery method.
     * @param "MaxVideoBitDepth" (optional.Int32) -  Optional.
     * @param "VideoCodec" (optional.String) -  Optional. Specify a video codec to encode to, e.g. h264. If omitted the server will auto-select using the url&#x27;s extension. Options: h264, mpeg4, theora, vpx, wmv.
     * @param "AudioStreamIndex" (optional.Int32) -  Optional. The index of the audio stream to use. If omitted the first audio stream will be used.
     * @param "VideoStreamIndex" (optional.Int32) -  Optional. The index of the video stream to use. If omitted the first video stream will be used.

*/

type AudioServiceApiGetAudioByIdByStreamfilenameOpts struct {
    DeviceProfileId optional.String
    DeviceId optional.String
    AudioCodec optional.String
    EnableAutoStreamCopy optional.Bool
    AudioSampleRate optional.Int32
    AudioBitRate optional.Int32
    AudioChannels optional.Int32
    MaxAudioChannels optional.Int32
    Static optional.Bool
    CopyTimestamps optional.Bool
    StartTimeTicks optional.Int64
    Width optional.Int32
    Height optional.Int32
    MaxWidth optional.Int32
    MaxHeight optional.Int32
    VideoBitRate optional.Int32
    SubtitleStreamIndex optional.Int32
    SubtitleMethod optional.Interface
    MaxVideoBitDepth optional.Int32
    VideoCodec optional.String
    AudioStreamIndex optional.Int32
    VideoStreamIndex optional.Int32
}

func (a *AudioServiceApiService) GetAudioByIdByStreamfilename(ctx context.Context, streamFileName string, id string, container string, localVarOptionals *AudioServiceApiGetAudioByIdByStreamfilenameOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/Audio/{Id}/{StreamFileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"StreamFileName"+"}", fmt.Sprintf("%v", streamFileName), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.DeviceProfileId.IsSet() {
		localVarQueryParams.Add("DeviceProfileId", parameterToString(localVarOptionals.DeviceProfileId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.DeviceId.IsSet() {
		localVarQueryParams.Add("DeviceId", parameterToString(localVarOptionals.DeviceId.Value(), ""))
	}
	localVarQueryParams.Add("Container", parameterToString(container, ""))
	if localVarOptionals != nil && localVarOptionals.AudioCodec.IsSet() {
		localVarQueryParams.Add("AudioCodec", parameterToString(localVarOptionals.AudioCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.EnableAutoStreamCopy.IsSet() {
		localVarQueryParams.Add("EnableAutoStreamCopy", parameterToString(localVarOptionals.EnableAutoStreamCopy.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioSampleRate.IsSet() {
		localVarQueryParams.Add("AudioSampleRate", parameterToString(localVarOptionals.AudioSampleRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioBitRate.IsSet() {
		localVarQueryParams.Add("AudioBitRate", parameterToString(localVarOptionals.AudioBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioChannels.IsSet() {
		localVarQueryParams.Add("AudioChannels", parameterToString(localVarOptionals.AudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxAudioChannels.IsSet() {
		localVarQueryParams.Add("MaxAudioChannels", parameterToString(localVarOptionals.MaxAudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Static.IsSet() {
		localVarQueryParams.Add("Static", parameterToString(localVarOptionals.Static.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CopyTimestamps.IsSet() {
		localVarQueryParams.Add("CopyTimestamps", parameterToString(localVarOptionals.CopyTimestamps.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StartTimeTicks.IsSet() {
		localVarQueryParams.Add("StartTimeTicks", parameterToString(localVarOptionals.StartTimeTicks.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Width.IsSet() {
		localVarQueryParams.Add("Width", parameterToString(localVarOptionals.Width.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Height.IsSet() {
		localVarQueryParams.Add("Height", parameterToString(localVarOptionals.Height.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxWidth.IsSet() {
		localVarQueryParams.Add("MaxWidth", parameterToString(localVarOptionals.MaxWidth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxHeight.IsSet() {
		localVarQueryParams.Add("MaxHeight", parameterToString(localVarOptionals.MaxHeight.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoBitRate.IsSet() {
		localVarQueryParams.Add("VideoBitRate", parameterToString(localVarOptionals.VideoBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleStreamIndex.IsSet() {
		localVarQueryParams.Add("SubtitleStreamIndex", parameterToString(localVarOptionals.SubtitleStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleMethod.IsSet() {
		localVarQueryParams.Add("SubtitleMethod", parameterToString(localVarOptionals.SubtitleMethod.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxVideoBitDepth.IsSet() {
		localVarQueryParams.Add("MaxVideoBitDepth", parameterToString(localVarOptionals.MaxVideoBitDepth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoCodec.IsSet() {
		localVarQueryParams.Add("VideoCodec", parameterToString(localVarOptionals.VideoCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioStreamIndex.IsSet() {
		localVarQueryParams.Add("AudioStreamIndex", parameterToString(localVarOptionals.AudioStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoStreamIndex.IsSet() {
		localVarQueryParams.Add("VideoStreamIndex", parameterToString(localVarOptionals.VideoStreamIndex.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			
			localVarQueryParams.Add("api_key", key)
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}
/*
AudioServiceApiService Gets an audio stream
Requires authentication as user
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Item Id
 * @param container Container
 * @param optional nil or *AudioServiceApiGetAudioByIdStreamOpts - Optional Parameters:
     * @param "DeviceProfileId" (optional.String) -  Optional. The dlna device profile id to utilize.
     * @param "DeviceId" (optional.String) -  The device id of the client requesting. Used to stop encoding processes when needed.
     * @param "AudioCodec" (optional.String) -  Optional. Specify a audio codec to encode to, e.g. mp3. If omitted the server will auto-select using the url&#x27;s extension. Options: aac, mp3, vorbis, wma.
     * @param "EnableAutoStreamCopy" (optional.Bool) -  Whether or not to allow automatic stream copy if requested values match the original source. Defaults to true.
     * @param "AudioSampleRate" (optional.Int32) -  Optional. Specify a specific audio sample rate, e.g. 44100
     * @param "AudioBitRate" (optional.Int32) -  Optional. Specify an audio bitrate to encode to, e.g. 128000. If omitted this will be left to encoder defaults.
     * @param "AudioChannels" (optional.Int32) -  Optional. Specify a specific number of audio channels to encode to, e.g. 2
     * @param "MaxAudioChannels" (optional.Int32) -  Optional. Specify a maximum number of audio channels to encode to, e.g. 2
     * @param "Static" (optional.Bool) -  Optional. If true, the original file will be streamed statically without any encoding. Use either no url extension or the original file extension. true/false
     * @param "CopyTimestamps" (optional.Bool) -  Whether or not to copy timestamps when transcoding with an offset. Defaults to false.
     * @param "StartTimeTicks" (optional.Int64) -  Optional. Specify a starting offset, in ticks. 1ms &#x3D; 10000 ticks.
     * @param "Width" (optional.Int32) -  Optional. The fixed horizontal resolution of the encoded video.
     * @param "Height" (optional.Int32) -  Optional. The fixed vertical resolution of the encoded video.
     * @param "MaxWidth" (optional.Int32) -  Optional. The maximum horizontal resolution of the encoded video.
     * @param "MaxHeight" (optional.Int32) -  Optional. The maximum vertical resolution of the encoded video.
     * @param "VideoBitRate" (optional.Int32) -  Optional. Specify a video bitrate to encode to, e.g. 500000. If omitted this will be left to encoder defaults.
     * @param "SubtitleStreamIndex" (optional.Int32) -  Optional. The index of the subtitle stream to use. If omitted no subtitles will be used.
     * @param "SubtitleMethod" (optional.Interface of SubtitleDeliveryMethod) -  Optional. Specify the subtitle delivery method.
     * @param "MaxVideoBitDepth" (optional.Int32) -  Optional.
     * @param "VideoCodec" (optional.String) -  Optional. Specify a video codec to encode to, e.g. h264. If omitted the server will auto-select using the url&#x27;s extension. Options: h264, mpeg4, theora, vpx, wmv.
     * @param "AudioStreamIndex" (optional.Int32) -  Optional. The index of the audio stream to use. If omitted the first audio stream will be used.
     * @param "VideoStreamIndex" (optional.Int32) -  Optional. The index of the video stream to use. If omitted the first video stream will be used.

*/

type AudioServiceApiGetAudioByIdStreamOpts struct {
    DeviceProfileId optional.String
    DeviceId optional.String
    AudioCodec optional.String
    EnableAutoStreamCopy optional.Bool
    AudioSampleRate optional.Int32
    AudioBitRate optional.Int32
    AudioChannels optional.Int32
    MaxAudioChannels optional.Int32
    Static optional.Bool
    CopyTimestamps optional.Bool
    StartTimeTicks optional.Int64
    Width optional.Int32
    Height optional.Int32
    MaxWidth optional.Int32
    MaxHeight optional.Int32
    VideoBitRate optional.Int32
    SubtitleStreamIndex optional.Int32
    SubtitleMethod optional.Interface
    MaxVideoBitDepth optional.Int32
    VideoCodec optional.String
    AudioStreamIndex optional.Int32
    VideoStreamIndex optional.Int32
}

func (a *AudioServiceApiService) GetAudioByIdStream(ctx context.Context, id string, container string, localVarOptionals *AudioServiceApiGetAudioByIdStreamOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/Audio/{Id}/stream"
	localVarPath = strings.Replace(localVarPath, "{"+"Id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.DeviceProfileId.IsSet() {
		localVarQueryParams.Add("DeviceProfileId", parameterToString(localVarOptionals.DeviceProfileId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.DeviceId.IsSet() {
		localVarQueryParams.Add("DeviceId", parameterToString(localVarOptionals.DeviceId.Value(), ""))
	}
	localVarQueryParams.Add("Container", parameterToString(container, ""))
	if localVarOptionals != nil && localVarOptionals.AudioCodec.IsSet() {
		localVarQueryParams.Add("AudioCodec", parameterToString(localVarOptionals.AudioCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.EnableAutoStreamCopy.IsSet() {
		localVarQueryParams.Add("EnableAutoStreamCopy", parameterToString(localVarOptionals.EnableAutoStreamCopy.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioSampleRate.IsSet() {
		localVarQueryParams.Add("AudioSampleRate", parameterToString(localVarOptionals.AudioSampleRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioBitRate.IsSet() {
		localVarQueryParams.Add("AudioBitRate", parameterToString(localVarOptionals.AudioBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioChannels.IsSet() {
		localVarQueryParams.Add("AudioChannels", parameterToString(localVarOptionals.AudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxAudioChannels.IsSet() {
		localVarQueryParams.Add("MaxAudioChannels", parameterToString(localVarOptionals.MaxAudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Static.IsSet() {
		localVarQueryParams.Add("Static", parameterToString(localVarOptionals.Static.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CopyTimestamps.IsSet() {
		localVarQueryParams.Add("CopyTimestamps", parameterToString(localVarOptionals.CopyTimestamps.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StartTimeTicks.IsSet() {
		localVarQueryParams.Add("StartTimeTicks", parameterToString(localVarOptionals.StartTimeTicks.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Width.IsSet() {
		localVarQueryParams.Add("Width", parameterToString(localVarOptionals.Width.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Height.IsSet() {
		localVarQueryParams.Add("Height", parameterToString(localVarOptionals.Height.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxWidth.IsSet() {
		localVarQueryParams.Add("MaxWidth", parameterToString(localVarOptionals.MaxWidth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxHeight.IsSet() {
		localVarQueryParams.Add("MaxHeight", parameterToString(localVarOptionals.MaxHeight.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoBitRate.IsSet() {
		localVarQueryParams.Add("VideoBitRate", parameterToString(localVarOptionals.VideoBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleStreamIndex.IsSet() {
		localVarQueryParams.Add("SubtitleStreamIndex", parameterToString(localVarOptionals.SubtitleStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleMethod.IsSet() {
		localVarQueryParams.Add("SubtitleMethod", parameterToString(localVarOptionals.SubtitleMethod.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxVideoBitDepth.IsSet() {
		localVarQueryParams.Add("MaxVideoBitDepth", parameterToString(localVarOptionals.MaxVideoBitDepth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoCodec.IsSet() {
		localVarQueryParams.Add("VideoCodec", parameterToString(localVarOptionals.VideoCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioStreamIndex.IsSet() {
		localVarQueryParams.Add("AudioStreamIndex", parameterToString(localVarOptionals.AudioStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoStreamIndex.IsSet() {
		localVarQueryParams.Add("VideoStreamIndex", parameterToString(localVarOptionals.VideoStreamIndex.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			
			localVarQueryParams.Add("api_key", key)
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}
/*
AudioServiceApiService Gets an audio stream
Requires authentication as user
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Item Id
 * @param container Container
 * @param optional nil or *AudioServiceApiGetAudioByIdStreamByContainerOpts - Optional Parameters:
     * @param "DeviceProfileId" (optional.String) -  Optional. The dlna device profile id to utilize.
     * @param "DeviceId" (optional.String) -  The device id of the client requesting. Used to stop encoding processes when needed.
     * @param "AudioCodec" (optional.String) -  Optional. Specify a audio codec to encode to, e.g. mp3. If omitted the server will auto-select using the url&#x27;s extension. Options: aac, mp3, vorbis, wma.
     * @param "EnableAutoStreamCopy" (optional.Bool) -  Whether or not to allow automatic stream copy if requested values match the original source. Defaults to true.
     * @param "AudioSampleRate" (optional.Int32) -  Optional. Specify a specific audio sample rate, e.g. 44100
     * @param "AudioBitRate" (optional.Int32) -  Optional. Specify an audio bitrate to encode to, e.g. 128000. If omitted this will be left to encoder defaults.
     * @param "AudioChannels" (optional.Int32) -  Optional. Specify a specific number of audio channels to encode to, e.g. 2
     * @param "MaxAudioChannels" (optional.Int32) -  Optional. Specify a maximum number of audio channels to encode to, e.g. 2
     * @param "Static" (optional.Bool) -  Optional. If true, the original file will be streamed statically without any encoding. Use either no url extension or the original file extension. true/false
     * @param "CopyTimestamps" (optional.Bool) -  Whether or not to copy timestamps when transcoding with an offset. Defaults to false.
     * @param "StartTimeTicks" (optional.Int64) -  Optional. Specify a starting offset, in ticks. 1ms &#x3D; 10000 ticks.
     * @param "Width" (optional.Int32) -  Optional. The fixed horizontal resolution of the encoded video.
     * @param "Height" (optional.Int32) -  Optional. The fixed vertical resolution of the encoded video.
     * @param "MaxWidth" (optional.Int32) -  Optional. The maximum horizontal resolution of the encoded video.
     * @param "MaxHeight" (optional.Int32) -  Optional. The maximum vertical resolution of the encoded video.
     * @param "VideoBitRate" (optional.Int32) -  Optional. Specify a video bitrate to encode to, e.g. 500000. If omitted this will be left to encoder defaults.
     * @param "SubtitleStreamIndex" (optional.Int32) -  Optional. The index of the subtitle stream to use. If omitted no subtitles will be used.
     * @param "SubtitleMethod" (optional.Interface of SubtitleDeliveryMethod) -  Optional. Specify the subtitle delivery method.
     * @param "MaxVideoBitDepth" (optional.Int32) -  Optional.
     * @param "VideoCodec" (optional.String) -  Optional. Specify a video codec to encode to, e.g. h264. If omitted the server will auto-select using the url&#x27;s extension. Options: h264, mpeg4, theora, vpx, wmv.
     * @param "AudioStreamIndex" (optional.Int32) -  Optional. The index of the audio stream to use. If omitted the first audio stream will be used.
     * @param "VideoStreamIndex" (optional.Int32) -  Optional. The index of the video stream to use. If omitted the first video stream will be used.

*/

type AudioServiceApiGetAudioByIdStreamByContainerOpts struct {
    DeviceProfileId optional.String
    DeviceId optional.String
    AudioCodec optional.String
    EnableAutoStreamCopy optional.Bool
    AudioSampleRate optional.Int32
    AudioBitRate optional.Int32
    AudioChannels optional.Int32
    MaxAudioChannels optional.Int32
    Static optional.Bool
    CopyTimestamps optional.Bool
    StartTimeTicks optional.Int64
    Width optional.Int32
    Height optional.Int32
    MaxWidth optional.Int32
    MaxHeight optional.Int32
    VideoBitRate optional.Int32
    SubtitleStreamIndex optional.Int32
    SubtitleMethod optional.Interface
    MaxVideoBitDepth optional.Int32
    VideoCodec optional.String
    AudioStreamIndex optional.Int32
    VideoStreamIndex optional.Int32
}

func (a *AudioServiceApiService) GetAudioByIdStreamByContainer(ctx context.Context, id string, container string, localVarOptionals *AudioServiceApiGetAudioByIdStreamByContainerOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/Audio/{Id}/stream.{Container}"
	localVarPath = strings.Replace(localVarPath, "{"+"Id"+"}", fmt.Sprintf("%v", id), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Container"+"}", fmt.Sprintf("%v", container), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.DeviceProfileId.IsSet() {
		localVarQueryParams.Add("DeviceProfileId", parameterToString(localVarOptionals.DeviceProfileId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.DeviceId.IsSet() {
		localVarQueryParams.Add("DeviceId", parameterToString(localVarOptionals.DeviceId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioCodec.IsSet() {
		localVarQueryParams.Add("AudioCodec", parameterToString(localVarOptionals.AudioCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.EnableAutoStreamCopy.IsSet() {
		localVarQueryParams.Add("EnableAutoStreamCopy", parameterToString(localVarOptionals.EnableAutoStreamCopy.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioSampleRate.IsSet() {
		localVarQueryParams.Add("AudioSampleRate", parameterToString(localVarOptionals.AudioSampleRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioBitRate.IsSet() {
		localVarQueryParams.Add("AudioBitRate", parameterToString(localVarOptionals.AudioBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioChannels.IsSet() {
		localVarQueryParams.Add("AudioChannels", parameterToString(localVarOptionals.AudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxAudioChannels.IsSet() {
		localVarQueryParams.Add("MaxAudioChannels", parameterToString(localVarOptionals.MaxAudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Static.IsSet() {
		localVarQueryParams.Add("Static", parameterToString(localVarOptionals.Static.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CopyTimestamps.IsSet() {
		localVarQueryParams.Add("CopyTimestamps", parameterToString(localVarOptionals.CopyTimestamps.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StartTimeTicks.IsSet() {
		localVarQueryParams.Add("StartTimeTicks", parameterToString(localVarOptionals.StartTimeTicks.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Width.IsSet() {
		localVarQueryParams.Add("Width", parameterToString(localVarOptionals.Width.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Height.IsSet() {
		localVarQueryParams.Add("Height", parameterToString(localVarOptionals.Height.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxWidth.IsSet() {
		localVarQueryParams.Add("MaxWidth", parameterToString(localVarOptionals.MaxWidth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxHeight.IsSet() {
		localVarQueryParams.Add("MaxHeight", parameterToString(localVarOptionals.MaxHeight.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoBitRate.IsSet() {
		localVarQueryParams.Add("VideoBitRate", parameterToString(localVarOptionals.VideoBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleStreamIndex.IsSet() {
		localVarQueryParams.Add("SubtitleStreamIndex", parameterToString(localVarOptionals.SubtitleStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleMethod.IsSet() {
		localVarQueryParams.Add("SubtitleMethod", parameterToString(localVarOptionals.SubtitleMethod.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxVideoBitDepth.IsSet() {
		localVarQueryParams.Add("MaxVideoBitDepth", parameterToString(localVarOptionals.MaxVideoBitDepth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoCodec.IsSet() {
		localVarQueryParams.Add("VideoCodec", parameterToString(localVarOptionals.VideoCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioStreamIndex.IsSet() {
		localVarQueryParams.Add("AudioStreamIndex", parameterToString(localVarOptionals.AudioStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoStreamIndex.IsSet() {
		localVarQueryParams.Add("VideoStreamIndex", parameterToString(localVarOptionals.VideoStreamIndex.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			
			localVarQueryParams.Add("api_key", key)
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}
/*
AudioServiceApiService Gets an audio stream
Requires authentication as user
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param streamFileName
 * @param id Item Id
 * @param container Container
 * @param optional nil or *AudioServiceApiHeadAudioByIdByStreamfilenameOpts - Optional Parameters:
     * @param "DeviceProfileId" (optional.String) -  Optional. The dlna device profile id to utilize.
     * @param "DeviceId" (optional.String) -  The device id of the client requesting. Used to stop encoding processes when needed.
     * @param "AudioCodec" (optional.String) -  Optional. Specify a audio codec to encode to, e.g. mp3. If omitted the server will auto-select using the url&#x27;s extension. Options: aac, mp3, vorbis, wma.
     * @param "EnableAutoStreamCopy" (optional.Bool) -  Whether or not to allow automatic stream copy if requested values match the original source. Defaults to true.
     * @param "AudioSampleRate" (optional.Int32) -  Optional. Specify a specific audio sample rate, e.g. 44100
     * @param "AudioBitRate" (optional.Int32) -  Optional. Specify an audio bitrate to encode to, e.g. 128000. If omitted this will be left to encoder defaults.
     * @param "AudioChannels" (optional.Int32) -  Optional. Specify a specific number of audio channels to encode to, e.g. 2
     * @param "MaxAudioChannels" (optional.Int32) -  Optional. Specify a maximum number of audio channels to encode to, e.g. 2
     * @param "Static" (optional.Bool) -  Optional. If true, the original file will be streamed statically without any encoding. Use either no url extension or the original file extension. true/false
     * @param "CopyTimestamps" (optional.Bool) -  Whether or not to copy timestamps when transcoding with an offset. Defaults to false.
     * @param "StartTimeTicks" (optional.Int64) -  Optional. Specify a starting offset, in ticks. 1ms &#x3D; 10000 ticks.
     * @param "Width" (optional.Int32) -  Optional. The fixed horizontal resolution of the encoded video.
     * @param "Height" (optional.Int32) -  Optional. The fixed vertical resolution of the encoded video.
     * @param "MaxWidth" (optional.Int32) -  Optional. The maximum horizontal resolution of the encoded video.
     * @param "MaxHeight" (optional.Int32) -  Optional. The maximum vertical resolution of the encoded video.
     * @param "VideoBitRate" (optional.Int32) -  Optional. Specify a video bitrate to encode to, e.g. 500000. If omitted this will be left to encoder defaults.
     * @param "SubtitleStreamIndex" (optional.Int32) -  Optional. The index of the subtitle stream to use. If omitted no subtitles will be used.
     * @param "SubtitleMethod" (optional.Interface of SubtitleDeliveryMethod) -  Optional. Specify the subtitle delivery method.
     * @param "MaxVideoBitDepth" (optional.Int32) -  Optional.
     * @param "VideoCodec" (optional.String) -  Optional. Specify a video codec to encode to, e.g. h264. If omitted the server will auto-select using the url&#x27;s extension. Options: h264, mpeg4, theora, vpx, wmv.
     * @param "AudioStreamIndex" (optional.Int32) -  Optional. The index of the audio stream to use. If omitted the first audio stream will be used.
     * @param "VideoStreamIndex" (optional.Int32) -  Optional. The index of the video stream to use. If omitted the first video stream will be used.

*/

type AudioServiceApiHeadAudioByIdByStreamfilenameOpts struct {
    DeviceProfileId optional.String
    DeviceId optional.String
    AudioCodec optional.String
    EnableAutoStreamCopy optional.Bool
    AudioSampleRate optional.Int32
    AudioBitRate optional.Int32
    AudioChannels optional.Int32
    MaxAudioChannels optional.Int32
    Static optional.Bool
    CopyTimestamps optional.Bool
    StartTimeTicks optional.Int64
    Width optional.Int32
    Height optional.Int32
    MaxWidth optional.Int32
    MaxHeight optional.Int32
    VideoBitRate optional.Int32
    SubtitleStreamIndex optional.Int32
    SubtitleMethod optional.Interface
    MaxVideoBitDepth optional.Int32
    VideoCodec optional.String
    AudioStreamIndex optional.Int32
    VideoStreamIndex optional.Int32
}

func (a *AudioServiceApiService) HeadAudioByIdByStreamfilename(ctx context.Context, streamFileName string, id string, container string, localVarOptionals *AudioServiceApiHeadAudioByIdByStreamfilenameOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Head")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/Audio/{Id}/{StreamFileName}"
	localVarPath = strings.Replace(localVarPath, "{"+"StreamFileName"+"}", fmt.Sprintf("%v", streamFileName), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.DeviceProfileId.IsSet() {
		localVarQueryParams.Add("DeviceProfileId", parameterToString(localVarOptionals.DeviceProfileId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.DeviceId.IsSet() {
		localVarQueryParams.Add("DeviceId", parameterToString(localVarOptionals.DeviceId.Value(), ""))
	}
	localVarQueryParams.Add("Container", parameterToString(container, ""))
	if localVarOptionals != nil && localVarOptionals.AudioCodec.IsSet() {
		localVarQueryParams.Add("AudioCodec", parameterToString(localVarOptionals.AudioCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.EnableAutoStreamCopy.IsSet() {
		localVarQueryParams.Add("EnableAutoStreamCopy", parameterToString(localVarOptionals.EnableAutoStreamCopy.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioSampleRate.IsSet() {
		localVarQueryParams.Add("AudioSampleRate", parameterToString(localVarOptionals.AudioSampleRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioBitRate.IsSet() {
		localVarQueryParams.Add("AudioBitRate", parameterToString(localVarOptionals.AudioBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioChannels.IsSet() {
		localVarQueryParams.Add("AudioChannels", parameterToString(localVarOptionals.AudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxAudioChannels.IsSet() {
		localVarQueryParams.Add("MaxAudioChannels", parameterToString(localVarOptionals.MaxAudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Static.IsSet() {
		localVarQueryParams.Add("Static", parameterToString(localVarOptionals.Static.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CopyTimestamps.IsSet() {
		localVarQueryParams.Add("CopyTimestamps", parameterToString(localVarOptionals.CopyTimestamps.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StartTimeTicks.IsSet() {
		localVarQueryParams.Add("StartTimeTicks", parameterToString(localVarOptionals.StartTimeTicks.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Width.IsSet() {
		localVarQueryParams.Add("Width", parameterToString(localVarOptionals.Width.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Height.IsSet() {
		localVarQueryParams.Add("Height", parameterToString(localVarOptionals.Height.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxWidth.IsSet() {
		localVarQueryParams.Add("MaxWidth", parameterToString(localVarOptionals.MaxWidth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxHeight.IsSet() {
		localVarQueryParams.Add("MaxHeight", parameterToString(localVarOptionals.MaxHeight.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoBitRate.IsSet() {
		localVarQueryParams.Add("VideoBitRate", parameterToString(localVarOptionals.VideoBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleStreamIndex.IsSet() {
		localVarQueryParams.Add("SubtitleStreamIndex", parameterToString(localVarOptionals.SubtitleStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleMethod.IsSet() {
		localVarQueryParams.Add("SubtitleMethod", parameterToString(localVarOptionals.SubtitleMethod.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxVideoBitDepth.IsSet() {
		localVarQueryParams.Add("MaxVideoBitDepth", parameterToString(localVarOptionals.MaxVideoBitDepth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoCodec.IsSet() {
		localVarQueryParams.Add("VideoCodec", parameterToString(localVarOptionals.VideoCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioStreamIndex.IsSet() {
		localVarQueryParams.Add("AudioStreamIndex", parameterToString(localVarOptionals.AudioStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoStreamIndex.IsSet() {
		localVarQueryParams.Add("VideoStreamIndex", parameterToString(localVarOptionals.VideoStreamIndex.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			
			localVarQueryParams.Add("api_key", key)
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}
/*
AudioServiceApiService Gets an audio stream
Requires authentication as user
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Item Id
 * @param container Container
 * @param optional nil or *AudioServiceApiHeadAudioByIdStreamOpts - Optional Parameters:
     * @param "DeviceProfileId" (optional.String) -  Optional. The dlna device profile id to utilize.
     * @param "DeviceId" (optional.String) -  The device id of the client requesting. Used to stop encoding processes when needed.
     * @param "AudioCodec" (optional.String) -  Optional. Specify a audio codec to encode to, e.g. mp3. If omitted the server will auto-select using the url&#x27;s extension. Options: aac, mp3, vorbis, wma.
     * @param "EnableAutoStreamCopy" (optional.Bool) -  Whether or not to allow automatic stream copy if requested values match the original source. Defaults to true.
     * @param "AudioSampleRate" (optional.Int32) -  Optional. Specify a specific audio sample rate, e.g. 44100
     * @param "AudioBitRate" (optional.Int32) -  Optional. Specify an audio bitrate to encode to, e.g. 128000. If omitted this will be left to encoder defaults.
     * @param "AudioChannels" (optional.Int32) -  Optional. Specify a specific number of audio channels to encode to, e.g. 2
     * @param "MaxAudioChannels" (optional.Int32) -  Optional. Specify a maximum number of audio channels to encode to, e.g. 2
     * @param "Static" (optional.Bool) -  Optional. If true, the original file will be streamed statically without any encoding. Use either no url extension or the original file extension. true/false
     * @param "CopyTimestamps" (optional.Bool) -  Whether or not to copy timestamps when transcoding with an offset. Defaults to false.
     * @param "StartTimeTicks" (optional.Int64) -  Optional. Specify a starting offset, in ticks. 1ms &#x3D; 10000 ticks.
     * @param "Width" (optional.Int32) -  Optional. The fixed horizontal resolution of the encoded video.
     * @param "Height" (optional.Int32) -  Optional. The fixed vertical resolution of the encoded video.
     * @param "MaxWidth" (optional.Int32) -  Optional. The maximum horizontal resolution of the encoded video.
     * @param "MaxHeight" (optional.Int32) -  Optional. The maximum vertical resolution of the encoded video.
     * @param "VideoBitRate" (optional.Int32) -  Optional. Specify a video bitrate to encode to, e.g. 500000. If omitted this will be left to encoder defaults.
     * @param "SubtitleStreamIndex" (optional.Int32) -  Optional. The index of the subtitle stream to use. If omitted no subtitles will be used.
     * @param "SubtitleMethod" (optional.Interface of SubtitleDeliveryMethod) -  Optional. Specify the subtitle delivery method.
     * @param "MaxVideoBitDepth" (optional.Int32) -  Optional.
     * @param "VideoCodec" (optional.String) -  Optional. Specify a video codec to encode to, e.g. h264. If omitted the server will auto-select using the url&#x27;s extension. Options: h264, mpeg4, theora, vpx, wmv.
     * @param "AudioStreamIndex" (optional.Int32) -  Optional. The index of the audio stream to use. If omitted the first audio stream will be used.
     * @param "VideoStreamIndex" (optional.Int32) -  Optional. The index of the video stream to use. If omitted the first video stream will be used.

*/

type AudioServiceApiHeadAudioByIdStreamOpts struct {
    DeviceProfileId optional.String
    DeviceId optional.String
    AudioCodec optional.String
    EnableAutoStreamCopy optional.Bool
    AudioSampleRate optional.Int32
    AudioBitRate optional.Int32
    AudioChannels optional.Int32
    MaxAudioChannels optional.Int32
    Static optional.Bool
    CopyTimestamps optional.Bool
    StartTimeTicks optional.Int64
    Width optional.Int32
    Height optional.Int32
    MaxWidth optional.Int32
    MaxHeight optional.Int32
    VideoBitRate optional.Int32
    SubtitleStreamIndex optional.Int32
    SubtitleMethod optional.Interface
    MaxVideoBitDepth optional.Int32
    VideoCodec optional.String
    AudioStreamIndex optional.Int32
    VideoStreamIndex optional.Int32
}

func (a *AudioServiceApiService) HeadAudioByIdStream(ctx context.Context, id string, container string, localVarOptionals *AudioServiceApiHeadAudioByIdStreamOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Head")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/Audio/{Id}/stream"
	localVarPath = strings.Replace(localVarPath, "{"+"Id"+"}", fmt.Sprintf("%v", id), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.DeviceProfileId.IsSet() {
		localVarQueryParams.Add("DeviceProfileId", parameterToString(localVarOptionals.DeviceProfileId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.DeviceId.IsSet() {
		localVarQueryParams.Add("DeviceId", parameterToString(localVarOptionals.DeviceId.Value(), ""))
	}
	localVarQueryParams.Add("Container", parameterToString(container, ""))
	if localVarOptionals != nil && localVarOptionals.AudioCodec.IsSet() {
		localVarQueryParams.Add("AudioCodec", parameterToString(localVarOptionals.AudioCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.EnableAutoStreamCopy.IsSet() {
		localVarQueryParams.Add("EnableAutoStreamCopy", parameterToString(localVarOptionals.EnableAutoStreamCopy.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioSampleRate.IsSet() {
		localVarQueryParams.Add("AudioSampleRate", parameterToString(localVarOptionals.AudioSampleRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioBitRate.IsSet() {
		localVarQueryParams.Add("AudioBitRate", parameterToString(localVarOptionals.AudioBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioChannels.IsSet() {
		localVarQueryParams.Add("AudioChannels", parameterToString(localVarOptionals.AudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxAudioChannels.IsSet() {
		localVarQueryParams.Add("MaxAudioChannels", parameterToString(localVarOptionals.MaxAudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Static.IsSet() {
		localVarQueryParams.Add("Static", parameterToString(localVarOptionals.Static.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CopyTimestamps.IsSet() {
		localVarQueryParams.Add("CopyTimestamps", parameterToString(localVarOptionals.CopyTimestamps.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StartTimeTicks.IsSet() {
		localVarQueryParams.Add("StartTimeTicks", parameterToString(localVarOptionals.StartTimeTicks.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Width.IsSet() {
		localVarQueryParams.Add("Width", parameterToString(localVarOptionals.Width.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Height.IsSet() {
		localVarQueryParams.Add("Height", parameterToString(localVarOptionals.Height.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxWidth.IsSet() {
		localVarQueryParams.Add("MaxWidth", parameterToString(localVarOptionals.MaxWidth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxHeight.IsSet() {
		localVarQueryParams.Add("MaxHeight", parameterToString(localVarOptionals.MaxHeight.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoBitRate.IsSet() {
		localVarQueryParams.Add("VideoBitRate", parameterToString(localVarOptionals.VideoBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleStreamIndex.IsSet() {
		localVarQueryParams.Add("SubtitleStreamIndex", parameterToString(localVarOptionals.SubtitleStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleMethod.IsSet() {
		localVarQueryParams.Add("SubtitleMethod", parameterToString(localVarOptionals.SubtitleMethod.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxVideoBitDepth.IsSet() {
		localVarQueryParams.Add("MaxVideoBitDepth", parameterToString(localVarOptionals.MaxVideoBitDepth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoCodec.IsSet() {
		localVarQueryParams.Add("VideoCodec", parameterToString(localVarOptionals.VideoCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioStreamIndex.IsSet() {
		localVarQueryParams.Add("AudioStreamIndex", parameterToString(localVarOptionals.AudioStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoStreamIndex.IsSet() {
		localVarQueryParams.Add("VideoStreamIndex", parameterToString(localVarOptionals.VideoStreamIndex.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			
			localVarQueryParams.Add("api_key", key)
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}
/*
AudioServiceApiService Gets an audio stream
Requires authentication as user
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param id Item Id
 * @param container Container
 * @param optional nil or *AudioServiceApiHeadAudioByIdStreamByContainerOpts - Optional Parameters:
     * @param "DeviceProfileId" (optional.String) -  Optional. The dlna device profile id to utilize.
     * @param "DeviceId" (optional.String) -  The device id of the client requesting. Used to stop encoding processes when needed.
     * @param "AudioCodec" (optional.String) -  Optional. Specify a audio codec to encode to, e.g. mp3. If omitted the server will auto-select using the url&#x27;s extension. Options: aac, mp3, vorbis, wma.
     * @param "EnableAutoStreamCopy" (optional.Bool) -  Whether or not to allow automatic stream copy if requested values match the original source. Defaults to true.
     * @param "AudioSampleRate" (optional.Int32) -  Optional. Specify a specific audio sample rate, e.g. 44100
     * @param "AudioBitRate" (optional.Int32) -  Optional. Specify an audio bitrate to encode to, e.g. 128000. If omitted this will be left to encoder defaults.
     * @param "AudioChannels" (optional.Int32) -  Optional. Specify a specific number of audio channels to encode to, e.g. 2
     * @param "MaxAudioChannels" (optional.Int32) -  Optional. Specify a maximum number of audio channels to encode to, e.g. 2
     * @param "Static" (optional.Bool) -  Optional. If true, the original file will be streamed statically without any encoding. Use either no url extension or the original file extension. true/false
     * @param "CopyTimestamps" (optional.Bool) -  Whether or not to copy timestamps when transcoding with an offset. Defaults to false.
     * @param "StartTimeTicks" (optional.Int64) -  Optional. Specify a starting offset, in ticks. 1ms &#x3D; 10000 ticks.
     * @param "Width" (optional.Int32) -  Optional. The fixed horizontal resolution of the encoded video.
     * @param "Height" (optional.Int32) -  Optional. The fixed vertical resolution of the encoded video.
     * @param "MaxWidth" (optional.Int32) -  Optional. The maximum horizontal resolution of the encoded video.
     * @param "MaxHeight" (optional.Int32) -  Optional. The maximum vertical resolution of the encoded video.
     * @param "VideoBitRate" (optional.Int32) -  Optional. Specify a video bitrate to encode to, e.g. 500000. If omitted this will be left to encoder defaults.
     * @param "SubtitleStreamIndex" (optional.Int32) -  Optional. The index of the subtitle stream to use. If omitted no subtitles will be used.
     * @param "SubtitleMethod" (optional.Interface of SubtitleDeliveryMethod) -  Optional. Specify the subtitle delivery method.
     * @param "MaxVideoBitDepth" (optional.Int32) -  Optional.
     * @param "VideoCodec" (optional.String) -  Optional. Specify a video codec to encode to, e.g. h264. If omitted the server will auto-select using the url&#x27;s extension. Options: h264, mpeg4, theora, vpx, wmv.
     * @param "AudioStreamIndex" (optional.Int32) -  Optional. The index of the audio stream to use. If omitted the first audio stream will be used.
     * @param "VideoStreamIndex" (optional.Int32) -  Optional. The index of the video stream to use. If omitted the first video stream will be used.

*/

type AudioServiceApiHeadAudioByIdStreamByContainerOpts struct {
    DeviceProfileId optional.String
    DeviceId optional.String
    AudioCodec optional.String
    EnableAutoStreamCopy optional.Bool
    AudioSampleRate optional.Int32
    AudioBitRate optional.Int32
    AudioChannels optional.Int32
    MaxAudioChannels optional.Int32
    Static optional.Bool
    CopyTimestamps optional.Bool
    StartTimeTicks optional.Int64
    Width optional.Int32
    Height optional.Int32
    MaxWidth optional.Int32
    MaxHeight optional.Int32
    VideoBitRate optional.Int32
    SubtitleStreamIndex optional.Int32
    SubtitleMethod optional.Interface
    MaxVideoBitDepth optional.Int32
    VideoCodec optional.String
    AudioStreamIndex optional.Int32
    VideoStreamIndex optional.Int32
}

func (a *AudioServiceApiService) HeadAudioByIdStreamByContainer(ctx context.Context, id string, container string, localVarOptionals *AudioServiceApiHeadAudioByIdStreamByContainerOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Head")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/Audio/{Id}/stream.{Container}"
	localVarPath = strings.Replace(localVarPath, "{"+"Id"+"}", fmt.Sprintf("%v", id), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"Container"+"}", fmt.Sprintf("%v", container), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.DeviceProfileId.IsSet() {
		localVarQueryParams.Add("DeviceProfileId", parameterToString(localVarOptionals.DeviceProfileId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.DeviceId.IsSet() {
		localVarQueryParams.Add("DeviceId", parameterToString(localVarOptionals.DeviceId.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioCodec.IsSet() {
		localVarQueryParams.Add("AudioCodec", parameterToString(localVarOptionals.AudioCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.EnableAutoStreamCopy.IsSet() {
		localVarQueryParams.Add("EnableAutoStreamCopy", parameterToString(localVarOptionals.EnableAutoStreamCopy.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioSampleRate.IsSet() {
		localVarQueryParams.Add("AudioSampleRate", parameterToString(localVarOptionals.AudioSampleRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioBitRate.IsSet() {
		localVarQueryParams.Add("AudioBitRate", parameterToString(localVarOptionals.AudioBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioChannels.IsSet() {
		localVarQueryParams.Add("AudioChannels", parameterToString(localVarOptionals.AudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxAudioChannels.IsSet() {
		localVarQueryParams.Add("MaxAudioChannels", parameterToString(localVarOptionals.MaxAudioChannels.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Static.IsSet() {
		localVarQueryParams.Add("Static", parameterToString(localVarOptionals.Static.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.CopyTimestamps.IsSet() {
		localVarQueryParams.Add("CopyTimestamps", parameterToString(localVarOptionals.CopyTimestamps.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.StartTimeTicks.IsSet() {
		localVarQueryParams.Add("StartTimeTicks", parameterToString(localVarOptionals.StartTimeTicks.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Width.IsSet() {
		localVarQueryParams.Add("Width", parameterToString(localVarOptionals.Width.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Height.IsSet() {
		localVarQueryParams.Add("Height", parameterToString(localVarOptionals.Height.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxWidth.IsSet() {
		localVarQueryParams.Add("MaxWidth", parameterToString(localVarOptionals.MaxWidth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxHeight.IsSet() {
		localVarQueryParams.Add("MaxHeight", parameterToString(localVarOptionals.MaxHeight.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoBitRate.IsSet() {
		localVarQueryParams.Add("VideoBitRate", parameterToString(localVarOptionals.VideoBitRate.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleStreamIndex.IsSet() {
		localVarQueryParams.Add("SubtitleStreamIndex", parameterToString(localVarOptionals.SubtitleStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.SubtitleMethod.IsSet() {
		localVarQueryParams.Add("SubtitleMethod", parameterToString(localVarOptionals.SubtitleMethod.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.MaxVideoBitDepth.IsSet() {
		localVarQueryParams.Add("MaxVideoBitDepth", parameterToString(localVarOptionals.MaxVideoBitDepth.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoCodec.IsSet() {
		localVarQueryParams.Add("VideoCodec", parameterToString(localVarOptionals.VideoCodec.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.AudioStreamIndex.IsSet() {
		localVarQueryParams.Add("AudioStreamIndex", parameterToString(localVarOptionals.AudioStreamIndex.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.VideoStreamIndex.IsSet() {
		localVarQueryParams.Add("VideoStreamIndex", parameterToString(localVarOptionals.VideoStreamIndex.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if ctx != nil {
		// API Key Authentication
		if auth, ok := ctx.Value(ContextAPIKey).(APIKey); ok {
			var key string
			if auth.Prefix != "" {
				key = auth.Prefix + " " + auth.Key
			} else {
				key = auth.Key
			}
			
			localVarQueryParams.Add("api_key", key)
		}
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}
