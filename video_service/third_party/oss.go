package third_party

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

// 创建存储空间
func getBucket() (*oss.Bucket, error) {
	bucketName := viper.GetString("oss.bucketName")
	endPoint := viper.GetString("oss.endpoint")
	accessKeyId := viper.GetString("oss.accessKeyId")
	accessKeySecret := viper.GetString("oss.accessKeySecret")

	// 创建OSSClient实例
	client, err := oss.New(endPoint, accessKeyId, accessKeySecret)
	if err != nil {
		return nil, err
	}
	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// Upload 上传
func Upload(fileDir string, fileBytes []byte) error {
	bucket, err := getBucket()
	if err != nil {
		return err
	}

	// 上传文件
	file := bytes.NewReader(fileBytes)
	err = bucket.PutObject(fileDir, file)
	return err
}

// Delete 删除
func Delete(fileDir string) error {
	bucket, err := getBucket()
	if err != nil {
		return err
	}

	// 上传文件
	err = bucket.DeleteObject(fileDir)
	return err
}
