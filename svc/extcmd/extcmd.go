package extcmd

import (
	"fmt"
	"os/exec"
)

type ExternalCmdService interface {
	RunFfmpeg(string, string) error
}

type ExternalCmdServiceImpl struct{}

func NewExternalCmdServiceImpl() *ExternalCmdServiceImpl {
	return &ExternalCmdServiceImpl{}
}

func (svc *ExternalCmdServiceImpl) RunFfmpeg(playlistUrl, outFilename string) error {
	// TODO: unit test

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

	// TODO: hide ffmpeg's output to verbose log
	fmt.Println(string(out))

	return nil
}
