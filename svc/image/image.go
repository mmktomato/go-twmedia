package image

import (
	neturl "net/url"
	"strings"

	"github.com/mmktomato/go-twmedia/util"
)

type ImageService interface {
	IsTargetImage(string) bool
	GetImageFilename(string) (string, error)
}

type ImageServiceImpl struct {
	logger *util.TinyLogger
}

func NewImageServiceImpl(logger *util.TinyLogger) *ImageServiceImpl {
	return &ImageServiceImpl{logger}
}

func (svc *ImageServiceImpl) IsTargetImage(url string) bool {
	return strings.HasSuffix(url, ":large")
}

func (svc *ImageServiceImpl) GetImageFilename(url string) (string, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return "", err
	}

	tokens := strings.Split(u.Path, "/")
	ret := tokens[len(tokens)-1]
	return strings.Split(ret, ":")[0], nil // foo.jpg:large
}
