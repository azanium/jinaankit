package hls

import (
	"os"
	"os/exec"
	"strconv"
)

type HLSOption struct {
	FfmpegPath      string
	SrcPath         string
	TargetPath      string
	Resolution      string
	SegmentDuration int
}

type ThumbnailOption struct {
	FfmpegPath     string
	SrcPath        string
	TargetPath     string
	ThumbnailWidth int
}

// GenerateHLS will generate HLS file based on resolution presets.
// The available resolutions are: 360p, 480p, 720p and 1080p.
func GenerateHLS(option HLSOption) error {
	segmentDuration := strconv.Itoa(option.SegmentDuration)
	options, err := getOptions(option.SrcPath, option.TargetPath, option.Resolution, segmentDuration)
	if err != nil {
		return err
	}

	return GenerateHLSCustom(option.FfmpegPath, options)
}

func GenerateThumbnail(option ThumbnailOption) error {
	options, err := getThumbnailOptions(option.SrcPath, option.TargetPath, option.ThumbnailWidth)
	if err != nil {
		return err
	}

	return GenerateHLSCustom(option.FfmpegPath, options)
}

// GenerateHLSCustom will generate HLS using the flexible options params.s
// options is array of string that accepted by ffmpeg command
func GenerateHLSCustom(ffmpegPath string, options []string) error {
	cmd := exec.Command(ffmpegPath, options...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}
