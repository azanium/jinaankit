package hls

// import (
// 	"os"
// 	"os/exec"
// 	"path"
// 	"testing"
// )

// func TestCmdExecuteFfmpeg(t *testing.T) {
// 	base, _ := os.Getwd()

// 	targetPath := path.Join(base, "static")
// 	srcPath := path.Join(base, "static", "Love.mov")
// 	ffmpeg, _ := exec.LookPath("ffmpeg")

// 	err := GenerateHLS(ffmpeg, srcPath, targetPath, "480p")
// 	if err != nil {
// 		panic(err)
// 	}
// }
