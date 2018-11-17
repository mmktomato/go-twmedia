package video

import (
	"io"
	"io/ioutil"
	"regexp"

	"github.com/mmktomato/go-twmedia/util"
)

var authRegex = regexp.MustCompile(`authorization:\s*['"]([^'"]+)['"]`)

func GetAuthToken(jsurl string) (ret string, err error) {
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
