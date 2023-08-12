package third_party

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

func Upload(fileDir string, fileBytes []byte) error {

	bucketName := viper.GetString("oss.bucketName")
	endPoint := viper.GetString("oss.endpoint")
	accessKeyId := viper.GetString("oss.accessKeyId")
	accessKeySecret := viper.GetString("oss.accessKeySecret")

	// 创建OSSClient实例
	client, _ := oss.New(endPoint, accessKeyId, accessKeySecret)
	// 获取存储空间
	bucket, _ := client.Bucket(bucketName)
	// 上传文件
	file := bytes.NewReader(fileBytes)
	err := bucket.PutObject(fileDir, file)
	return err
}
