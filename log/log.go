package log

import (
	"sync"

	"github.com/seachenjy/go-comment/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//Logger log
var (
	logger *logrus.Logger
	m      *sync.RWMutex = &sync.RWMutex{}
	inited bool
)

// Init 日志
func Init() {
	m.Lock()
	defer m.Unlock()
	if inited {
		return
	}
	inited = true
	writer, err := rotatelogs.New(
		config.Cfg.LOG.Path+"/log.%Y-%m-%d.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		// rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(config.Cfg.LOG.RotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(config.Cfg.LOG.MaxAge),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)
	errWriter, err := rotatelogs.New(
		config.Cfg.LOG.Path+"/log.errors.%Y-%m-%d.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		// rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(config.Cfg.LOG.RotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(config.Cfg.LOG.MaxAge),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)
	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	logrus.SetLevel(logrus.WarnLevel)

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: errWriter,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: errWriter,
	}, &logrus.TextFormatter{DisableColors: true})

	logger = logrus.New()
	logger.AddHook(lfsHook)
}

//GetLogger return logger
func GetLogger() *logrus.Logger {
	if !inited {
		Init()
	}
	return logger
}
