package log

import (
	"IM-Service/src/configs/err"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

var logger *Logger
var conSoleRule *string

const (
	CONSOLE_FILE = "CONSOLE_FILE"
	CONSOLE      = "CONSOLE"
	FILE         = "FILE"
	CLOSE        = "CLOSE"
)

type Logger struct {
	*logrus.Logger
	Pid int
}

func InitLog(logFile string, con string) {
	conSoleRule = &con
	logger = loggerInit(logFile, 6)
}

func defaultLog() {
	InitLog("../logs", CONSOLE_FILE)
}

type UTCFormatter struct {
	logrus.Formatter
}

func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.In(time.UTC)
	return u.Formatter.Format(e)
}

func loggerInit(dir string, logLevel uint32) *Logger {
	var logger = logrus.New()
	logger.SetLevel(logrus.Level(logLevel))

	var formatter logrus.Formatter = &nested.Formatter{
		TimestampFormat: "01-02 15:04:05",
		HideKeys:        true,
		FieldsOrder:     []string{"FilePath"},
	}

	formatter = UTCFormatter{formatter}
	logger.SetFormatter(formatter)
	logger.AddHook(newFileHook())
	if *conSoleRule == CONSOLE_FILE {
		logger.SetOutput(os.Stdout)
		logger.AddHook(NewLfsHook(time.Duration(1)*time.Hour, 24, dir))
	} else if *conSoleRule == CONSOLE {
		logger.SetOutput(os.Stdout)
	} else if *conSoleRule == FILE {
		logger.SetOutput(io.Discard)
		logger.AddHook(NewLfsHook(time.Duration(1)*time.Hour, 24, dir))
	} else if *conSoleRule == CLOSE {
		logger.SetOutput(io.Discard)
	} else {
		logger.SetOutput(os.Stdout)
	}
	return &Logger{
		logger,
		os.Getpid(),
	}
}

func NewLfsHook(rotationTime time.Duration, maxRemainNum uint, dir string) logrus.Hook {
	var formatter logrus.Formatter = &nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		HideKeys:        true,
		FieldsOrder:     []string{"FilePath"},
		TrimMessages:    true,
		NoFieldsSpace:   true,
		NoColors:        true,
	}
	formatter = UTCFormatter{formatter}
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: initRotateLogs(rotationTime, maxRemainNum, "all", dir),
		logrus.InfoLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", dir),
		logrus.WarnLevel:  initRotateLogs(rotationTime, maxRemainNum, "all", dir),
		logrus.ErrorLevel: initRotateLogs(rotationTime, maxRemainNum, "all", dir),
	}, formatter)
	return lfsHook
}

func initRotateLogs(rotationTime time.Duration, maxRemainNum uint, level string, dir string) *rotatelogs.RotateLogs {
	writer, err := rotatelogs.New(
		filepath.Join(dir, level+"."+"%Y-%m-%d-%H-%M"+".log"),
		rotatelogs.WithRotationTime(rotationTime),
		rotatelogs.WithRotationCount(maxRemainNum),
	)
	if err != nil {
		panic(err.Error())
	} else {
		return writer
	}
}

func Error(args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorln(args)
}
func Debug(args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Debugln(args)
}
func Debugf(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Debugf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.WithFields(logrus.Fields{}).Errorf(format, args...)
}

func WithError(err error, message ...string) *utils.Error {
	var u *utils.Error
	if !errors.As(err, &u) {
		u = utils.NewSysError(err)
	}
	if logger == nil {
		defaultLog()
	}
	logger.WithFields(logrus.Fields{}).Errorln(u, message)
	if !u.IsHasStack {
		if len(message) == 0 {
			err = errors.WithStack(u)
		} else {
			err = errors.Wrap(u, message[0])
		}
		tmpErr := &utils.Error{
			Code:       u.Code,
			Msg:        u.Msg,
			MsgZh:      u.MsgZh,
			IsHasStack: true,
		}
		u = tmpErr
	}
	return u
}
