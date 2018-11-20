package video

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

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

func GetAuthToken(attrs []html.Attribute) (string, error) {
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

func GetTrackInfo(tweetId, authToken string) (*TrackInfo, error) {
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
