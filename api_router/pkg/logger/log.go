package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

var Log *logrus.Logger

func init() {

	logger := logrus.New()
	fileName := getFileDir()

	witter := rotateLog(fileName)

	logger.Out = witter

	Log = logger
}

// 获取日志输出路径
func getFileDir() string {
	now := time.Now()
	// 获取指定路径
	_, filePath, _, _ := runtime.Caller(0)
	logsPath := filepath.Join(filePath, "..", "..", "..", "logs")

	// 文件名称
	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logsPath, logFileName)

	// 查看文件是否存在，不存在则创建
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
		}
	}

	return fileName
}

// 日志本地文件分割
func rotateLog(fileName string) *rotatelogs.RotateLogs {
	witter, _ := rotatelogs.New(

		fileName+"%H%M",

		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithMaxAge(time.Duration(10)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(5)*time.Second),
	)

	return witter
}
