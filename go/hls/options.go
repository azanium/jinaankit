package hls

import (
	"fmt"
	"path/filepath"
	"strconv"
)

func getOptions(srcPath, targetPath, res string, segmentDuration string) ([]string, error) {
	config, err := getConfig(res)
	if err != nil {
		return nil, err
	}

	filenameTS := filepath.Join(targetPath, res+"_%03d.ts")
	filenameM3U8 := filepath.Join(targetPath, res+".m3u8")

	options := []string{
		"-hide_banner",
		"-y",
		"-i", srcPath,
		"-vf", config.Scale, //"scale=trunc(oh*a/2)*2:1080",
		"-c:a", "aac",
		"-ar", "48000",
		"-c:v", "h264",
		"-profile:v", "main",
		"-crf", "20",
		"-sc_threshold", "0",
		"-g", "48",
		"-keyint_min", "48",
		"-hls_time", segmentDuration,
		"-hls_playlist_type", "vod",
		"-b:v", config.VideoBitrate,
		"-maxrate", config.Maxrate,
		"-bufsize", config.BufSize,
		"-b:a", config.AudioBitrate,
		"-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2", // prevent not divisible by 2
		"-preset", "ultrafast",
		"-hls_segment_filename", filenameTS,
		filenameM3U8,
	}

	fmt.Printf("# Encoding params: %v\n", options)

	return options, nil
}

func getAudioOptions(srcPath, targetPath, bitRate string) ([]string, error) {
	outputFile := filepath.Join(targetPath, bitRate+"_%03d.m4a")
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
	return options, nil
}

func getThumbnailOptions(srcPath, targetPath string, thumbWidth int) ([]string, error) {
	width := strconv.Itoa(thumbWidth)
	options := []string{
		"-y",
		"-itsoffset", "-1",
		"-i", srcPath,
		"-vframes", "1",
		"-vf", "scale=" + width + ":-1",
		targetPath,
	}

	fmt.Printf("# Encoding params: %v\n", options)

	return options, nil
}
