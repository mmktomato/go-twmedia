package extcmd

import (
	"os/exec"

	"github.com/mmktomato/go-twmedia/util"
)

type ExternalCmdService interface {
	RunFfmpeg(string, string) error
}

type ExternalCmdServiceImpl struct {
	logger *util.TinyLogger
}

func NewExternalCmdServiceImpl(logger *util.TinyLogger) *ExternalCmdServiceImpl {
	return &ExternalCmdServiceImpl{logger}
}

func (svc *ExternalCmdServiceImpl) RunFfmpeg(playlistUrl, outFilename string) error {
	// ffmpeg -i <playlistUrl> -movflags faststart -c copy -f mpegts <outFilename>
	// ffmpeg -i <playlistUrl> -movflags faststart -c copy -acodec aac -r 60 -bsf:a aac_adtstoasc -f mpegts <outFilename>

	out, err := exec.Command(
		"ffmpeg", "-i", playlistUrl,
		"-movflags", "faststart",
		"-c", "copy",
		"-f", "mpegts",
		outFilename,
	).CombinedOutput()
	if err != nil {
		return err
	}

	svc.logger.Verboseln(string(out))

	return nil
}
