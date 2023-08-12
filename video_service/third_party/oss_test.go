package third_party

import (
	"io/ioutil"
	"testing"
	"video/config"
)

func TestUpload(t *testing.T) {

	config.InitConfig()
	filePath := "D:\\Project\\video\\output_image.jpg"
	fileByte, _ := ioutil.ReadFile(filePath)
	Upload("hello", fileByte)
}
