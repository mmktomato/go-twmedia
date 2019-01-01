package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/grafov/m3u8"
	"github.com/mmktomato/go-twmedia/svc/extcmd"
	"github.com/mmktomato/go-twmedia/util"
	"github.com/mmktomato/go-twmedia/util/domutil"
)

var authRegex = regexp.MustCompile(`authorization:\s*['"]([^'"]+)['"]`)

type TrackInfo struct {
	ContentId   string `json:"contentId"`
	PlaylistUrl string `json:"playbackUrl"`
}

type videoConfig struct {
	Track TrackInfo `json:"track"`
}

type VideoService interface {
	ParseVideo(string, io.Reader) (*TrackInfo, error)
	SavePlaylist(*TrackInfo) error
}

type VideoServiceImpl struct {
	extcmdService extcmd.ExternalCmdService
	httpClient    *util.HttpClient
}

func NewVideoServiceImpl(extcmdService extcmd.ExternalCmdService, httpClient *util.HttpClient) *VideoServiceImpl {
	return &VideoServiceImpl{extcmdService, httpClient}
}

// ParseVideo returns *TrackInfo. It contains playlist url and content id.
// Typically playlist url has `.m3u8` extension.
func (svc *VideoServiceImpl) ParseVideo(tweetId string, r io.Reader) (ret *TrackInfo, err error) {
	// TODO: write test code. needs mock for `Fetch` because video.GetAuthToken uses it.

	err = domutil.Tokenize(r, func(token html.Token) (bool, error) {
		switch token.Type {
		case html.StartTagToken:
			fallthrough
		case html.SelfClosingTagToken:
			if token.Data == "script" {
				authToken, err := svc.getAuthToken(token.Attr)
				if err != nil {
					return false, err
				}
				if authToken != "" {
					ret, err = svc.fetchTrackInfo(tweetId, authToken)

					return false, err
				}
			}
		}
		return true, nil
	})
	return ret, err
}

func (svc *VideoServiceImpl) SavePlaylist(track *TrackInfo) error {
	// TODO: validate `track`. track.PlaylistUrl and track.ContentId.
	// move the validation to `fetchTrackInfo` func.
	u, err := url.Parse(track.PlaylistUrl)
	if err != nil {
		return err
	}

	baseUrl := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	return svc.httpClient.Fetch(track.PlaylistUrl, func(r io.Reader) error {
		playlist, listType, err := m3u8.DecodeFrom(r, true)
		if err != nil {
			return err
		}

		switch listType {
		case m3u8.MEDIA:
			outFilename := track.ContentId + ".mp4"
			err = svc.extcmdService.RunFfmpeg(track.PlaylistUrl, outFilename)
			if err != nil {
				return err
			}
			fmt.Println(outFilename)
		case m3u8.MASTER:
			masterpl := playlist.(*m3u8.MasterPlaylist)
			if len(masterpl.Variants) < 1 {
				return errors.New("No variants found")
			}
			variant := svc.findBiggestVideo(masterpl.Variants)
			nextTrack := &TrackInfo{track.ContentId, baseUrl + variant.URI}
			return svc.SavePlaylist(nextTrack)
		}
		return nil
	})
}

func (svc *VideoServiceImpl) getAuthToken(attrs []html.Attribute) (string, error) {
	jsurl, err := svc.parseScriptAttr(attrs)
	if err != nil {
		return "", err
	}

	token, err := svc.extractAuthToken(jsurl)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (svc *VideoServiceImpl) fetchTrackInfo(tweetId, authToken string) (*TrackInfo, error) {
	url := fmt.Sprintf("https://api.twitter.com/1.1/videos/tweet/config/%s.json", tweetId)
	var ret *TrackInfo = nil
	headers := map[string]string{"authorization": authToken}

	if cookie := svc.findCookie(); cookie != "" {
		if csrfToken := svc.findCsrfTokenFromCookie(cookie); csrfToken != "" {
			headers["x-csrf-token"] = csrfToken
		}
	}

	err := svc.httpClient.FetchWithHeader(url, headers, func(r io.Reader) error {
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

func (svc *VideoServiceImpl) parseScriptAttr(attrs []html.Attribute) (ret string, err error) {
	err = domutil.FindAttr(attrs, "src", func(srcAttr html.Attribute) error {
		if strings.Contains(srcAttr.Val, "TwitterVideoPlayerIframe") {
			ret = srcAttr.Val
		}
		return nil
	})
	return ret, err
}

func (svc *VideoServiceImpl) extractAuthToken(jsurl string) (ret string, err error) {
	err = svc.httpClient.Fetch(jsurl, func(r io.Reader) error {
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

func (svc *VideoServiceImpl) findBiggestVideo(variants []*m3u8.Variant) *m3u8.Variant {
	s := variants
	sort.Slice(s, func(i, j int) bool {
		return s[i].Bandwidth < s[j].Bandwidth
	})
	return s[len(s)-1]
}

func (svc *VideoServiceImpl) findCsrfTokenFromCookie(cookie string) string {
	pairs := strings.Split(cookie, ";")
	for _, pair := range pairs {
		s := strings.TrimSpace(pair)
		if strings.HasPrefix(s, "ct0=") {
			return strings.Split(s, "=")[1]
		}
	}
	return ""
}

func (svc *VideoServiceImpl) findCookie() string {
	arr := []string{"cookie", "Cookie", "COOKIE"}
	for _, el := range arr {
		if cookie, ok := svc.httpClient.FindHeader(el); ok {
			return cookie
		}
	}
	return ""
}
