package video

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/mmktomato/go-twmedia/twparser/domutil"
	"github.com/mmktomato/go-twmedia/util"
)

var authRegex = regexp.MustCompile(`authorization:\s*['"]([^'"]+)['"]`)

type TrackInfo struct {
	ContentId   string `json:"contentId"`
	PlaylistUrl string `json:"playbackUrl"`
}
type videoConfig struct {
	Track TrackInfo `json:"track"`
}

// ParseVideo returns *TrackInfo. It contains playlist url and content id.
// Typically playlist url has `.m3u8` extension.
func ParseVideo(tweetId string, r io.Reader) (ret *TrackInfo, err error) {
	// TODO: write test code. needs mock for `Fetch` because video.GetAuthToken uses it.

	err = domutil.Tokenize(r, func(token html.Token) (bool, error) {
		switch token.Type {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			if token.Data == "script" {
				authToken, err := getAuthToken(token.Attr)
				if err != nil {
					return false, err
				}
				if authToken != "" {
					ret, err = fetchTrackInfo(tweetId, authToken)

					return false, err
				}
			}
		}
		return true, nil
	})
	return ret, err
}

func SavePlaylist(track *TrackInfo) error {
	u, err := url.Parse(track.PlaylistUrl)
	if err != nil {
		return err
	}

	baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	return util.Fetch(track.PlaylistUrl, func(r io.Reader) error {
		playlist, listType, err := m3u8.DecodeFrom(r, true)
		if err != nil {
			return err
		}

		switch listType {
		case m3u8.MEDIA:
			mediapl := playlist.(*m3u8.MediaPlaylist)
			// TODO: download and save
			fmt.Printf("%+v\n", mediapl)
		case m3u8.MASTER:
			masterpl := playlist.(*m3u8.MasterPlaylist)
			variant := findBiggestVideo(masterpl)
			nextTrack := &TrackInfo{track.ContentId, baseUrl + variant.URI}
			return SavePlaylist(nextTrack)
		}
		return nil
	})
}

func getAuthToken(attrs []html.Attribute) (string, error) {
	jsurl, err := parseScriptAttr(attrs)
	if err != nil {
		return "", err
	}

	token, err := extractAuthToken(jsurl)
	if err != nil {
		return "", err
	}
	return token, nil
}

func fetchTrackInfo(tweetId, authToken string) (*TrackInfo, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/videos/tweet/config/%s.json", tweetId)
	var ret *TrackInfo = nil
	err := util.FetchWithHeader(url, map[string]string{"authorization": authToken}, func(r io.Reader) error {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		var vconf videoConfig
		if err = json.Unmarshal(buf, &vconf); err != nil {
			return err
		}

		ret = &vconf.Track
		return nil
	})
	return ret, err
}

func parseScriptAttr(attrs []html.Attribute) (ret string, err error) {
	err = domutil.FindAttr(attrs, "src", func(srcAttr html.Attribute) error {
		if strings.Contains(srcAttr.Val, "TwitterVideoPlayerIframe") {
			ret = srcAttr.Val
		}
		return nil
	})
	return ret, err
}

func extractAuthToken(jsurl string) (ret string, err error) {
	err = util.Fetch(jsurl, func(r io.Reader) error {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		found := authRegex.FindSubmatch(buf)
		if 1 < len(found) {
			ret = string(found[1])
		}
		return nil
	})
	return ret, err
}

func findBiggestVideo(masterpl *m3u8.MasterPlaylist) *m3u8.Variant {
	s := masterpl.Variants[:]
	sort.Slice(s, func(i, j int) bool {
		return s[i].Bandwidth < s[j].Bandwidth
	})
	return s[len(s)-1]
}
