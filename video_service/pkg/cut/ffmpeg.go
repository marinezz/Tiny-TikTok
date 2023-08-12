package cut

import (
	"bytes"
	"os"
	"os/exec"
)

// Cover 截取视频为图片
func Cover(videoURL string, timeOffset string) ([]byte, error) {
	cmd := exec.Command(
		"ffmpeg", "-i", videoURL, "-ss", timeOffset, "-vframes", "1", "-q:v", "2", "-f", "image2", "pipe:1",
	)
	cmd.Stderr = os.Stderr

	var outputBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return outputBuffer.Bytes(), nil
}
