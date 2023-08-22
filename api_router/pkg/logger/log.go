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
	// 如果实例存在则不用新建
	if Log != nil {
		fileName := getFileDir()
		witter := rotateLog(fileName)
		Log.Out = witter
		return
	}

	logger := logrus.New()
	fileName := getFileDir()
	witter := rotateLog(fileName)

	// 设置输出文件
	logger.Out = witter
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

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
		// 日志最长保留时间
		rotatelogs.WithMaxAge(time.Duration(12)*time.Hour),
		// 日志轮转的时间间隔
		rotatelogs.WithRotationTime(time.Duration(3)*time.Hour),
	)

	return witter
}
