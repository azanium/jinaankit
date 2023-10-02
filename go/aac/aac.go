package aac

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Option struct {
	FfmpegPath string
	SrcPath    string
	TargetPath string
	BitRate    string
	Hash       string
}

// GenerateHLS will generate HLS file based on resolution presets.
// The available resolutions are: 360p, 480p, 720p and 1080p.
func GenerateAAC(opt Option) (string, error) {
	options, filename, err := getOptions(opt.SrcPath, opt.TargetPath, opt.Hash, opt.BitRate)
	if err != nil {
		return "", err
	}

	return filename, GenerateFfmpegAAC(opt.FfmpegPath, options)
}

// GenerateHLSCustom will generate HLS using the flexible options params.s
// options is array of string that accepted by ffmpeg command
func GenerateFfmpegAAC(ffmpegPath string, options []string) error {
	cmd := exec.Command(ffmpegPath, options...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return err
}

func getOptions(srcPath, targetPath, hash, bitRate string) ([]string, string, error) {
	filename := hash + "_" + bitRate + ".m4a"
	outputFile := filepath.Join(targetPath, filename)
	options := []string{
		"-hide_banner",
		"-y",
		"-i", srcPath,
		"-c:a", "aac",
		"-c:v", "copy",
		"-vcodec", "copy",
		"-b:a", bitRate,
		outputFile,
	}
	fmt.Printf("encoding params: %v\n", options)
	return options, filename, nil
}
