package hls

import (
	"fmt"
)

type config struct {
	Name         string
	Scale        string
	VideoBitrate string
	Maxrate      string
	BufSize      string
	AudioBitrate string
	Resolution   string
	Bandwidth    string
}

var preset = map[string]*config{
	"240p": {
		Name:         "240p",
		Scale:        "scale=w=426:h=240:force_original_aspect_ratio=decrease",
		VideoBitrate: "350k",
		Maxrate:      "375k",
		BufSize:      "1200k",
		AudioBitrate: "64k",
		Resolution:   "426x240",
		Bandwidth:    "200000", // 300K
	},
	"360p": {
		Name:         "360p",
		Scale:        "scale=w=640:h=360:force_original_aspect_ratio=decrease",
		VideoBitrate: "800k",
		Maxrate:      "856k",
		BufSize:      "1200k",
		AudioBitrate: "96k",
		Resolution:   "640x360",
		Bandwidth:    "800000", // 500K, 800K
	},
	"480p": {
		Name:         "480p",
		Scale:        "scale=w=842:h=480:force_original_aspect_ratio=decrease",
		VideoBitrate: "1400k",
		Maxrate:      "1498k",
		BufSize:      "2100k",
		AudioBitrate: "128k",
		Resolution:   "842x480",
		Bandwidth:    "1400000", // 1MBps
	},
	"720p": {
		Name:         "720p",
		Scale:        "scale=w=1280:h=720:force_original_aspect_ratio=decrease",
		VideoBitrate: "2800k",
		Maxrate:      "2996k",
		BufSize:      "4200k",
		AudioBitrate: "128k",
		Resolution:   "1280x720",
		Bandwidth:    "2800000", // 2.5MBps
	},
	"1080p": {
		Name:         "1080p",
		Scale:        "scale=w=1920:h=1080:force_original_aspect_ratio=decrease",
		VideoBitrate: "5000k",
		Maxrate:      "5350k",
		BufSize:      "7500k",
		AudioBitrate: "192k",
		Resolution:   "1920x1080",
		Bandwidth:    "5000000", // 5 Mbps
	},
}

// getConfig return config from the available preset
func getConfig(res string) (*config, error) {
	cfg, ok := preset[res]
	if !ok {
		return nil, fmt.Errorf("preset '%s' not found", res)
	}

	return cfg, nil
}
