package video

import (
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		track TrackInfo
		valid bool
	}{
		{track: TrackInfo{"contentId", "playlistUrl", "playbackType"}, valid: true},
		{track: TrackInfo{"", "playlistUrl", "playbackType"}, valid: false},
		{track: TrackInfo{"contentId", "", "playbackType"}, valid: false},
		{track: TrackInfo{"contentId", "playlistUrl", ""}, valid: false},
	}

	for i, tt := range tests {
		if tt.track.Validate() != tt.valid {
			t.Errorf("%d: Validate doesn't return expected value.", i)
		}
	}
}
