package video

type VideoConfig struct {
	Track TrackInfo `json:"track"`
}

type TrackInfo struct {
	ContentId    string `json:"contentId"`
	PlaylistUrl  string `json:"playbackUrl"`
	PlaybackType string `json:"playbackType"`
}

func (track TrackInfo) Validate() bool {
	return track.ContentId != "" && track.PlaylistUrl != "" && track.PlaybackType != ""
}
