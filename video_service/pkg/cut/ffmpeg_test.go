package cut

import (
	"testing"
	"video/config"
	"video/third_party"
)

func TestCover(t *testing.T) {
	config.InitConfig()
	videoURL := "http://tiny-tiktok.oss-cn-chengdu.aliyuncs.com/video1.mp4"

	imageBytes, _ := Cover(videoURL, "00:00:05")

	third_party.Upload("output.jpg", imageBytes)
}
