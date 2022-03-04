package utils

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	rotatelogs "github.com/kelofox/golang-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	red                           = "31"
	green                         = "32"
	yellow                        = "33"
	blue                          = "34"
	magenta                       = "35"
	cyan                          = "36"
	white                         = "37"
	LogDefaultDir                 = "logs/"
	LogDefaultLogFormat           = "[%time%] [%lvl%] %f% - %msg%"
	LogDefaultTimestampFormat     = "2006-01-02 15:04:05.000"
	LogDefaultTimeUnit            = "hour"
	LogDefaultRotationType        = "time"
	LogDefaultWithMaxAge          = 120
	LogDefaultWithRotationTime    = 24
	LogDefaultWithRotationCount   = 5
	LogDefaultWithRotationSize_MB = 100
)

type MyFormatter struct {
	logrus.TextFormatter
	LogFormat  string
	LoggerName string
}

type LogSetting struct {
	TimestampFormat     string
	TimeUnit            string
	RotationType        string
	WithMaxAge          int
	WithRotationTime    int
	WithRotationCount   uint
	WithRotationSize_MB int64
	LogDir              string
}

var LogSettingMap map[string]LogSetting

func init() {
	LogSettingMap, _ = getLogSetting()
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	level := strings.ToUpper(entry.Level.String())
	var levelColor string
	switch level {
	case "TRACE":
		levelColor = white
	case "DEBUG":
		levelColor = cyan
	case "INFO":
		levelColor = green
	case "WARNING":
		levelColor = yellow
	case "ERROR":
		levelColor = red
	default:
		levelColor = blue
	}
	msg := entry.Message
	if !f.DisableColors {
		output = "\x1b[" + levelColor + "m" + output
		msg = "\x1b[0m" + msg
	}
	output = strings.Replace(output, "%lvl%", level, 1)
	output = strings.Replace(output, "%f%", f.LoggerName, 1)
	timestampFormat := f.TimestampFormat
	// if timestampFormat == "" {
	// 	timestampFormat = defaultTimestampFormat
	// }
	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)
	output = strings.Replace(output, "%msg%", msg, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			output = strings.Replace(output, "%"+k+"%", v, 1)
		case int:
			s := strconv.Itoa(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		case bool:
			s := strconv.FormatBool(v)
			output = strings.Replace(output, "%"+k+"%", s, 1)
		}
	}
	output += "\n"
	return []byte(output), nil
}

type MyHook struct {
	Formatter logrus.Formatter
	logWriter io.Writer
}

// Levels 只定义 error 和 panic 等级的日志,其他日志等级不会触发 hook
func (h *MyHook) Levels() []logrus.Level {
	// return logrus.AllLevels
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.PanicLevel,
		logrus.InfoLevel,
	}
}

// Fire 将异常日志写入到指定日志文件中
func (h *MyHook) Fire(entry *logrus.Entry) error {
	b, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}
	h.logWriter.Write(b)
	return nil
	// f, err := os.OpenFile("err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 	return err
	// }
	// if _, err := f.Write([]byte(entry.Message)); err != nil {
	// 	return err
	// }
	// return nil
}

// @title GetLogger
// @description 取得logger
// @param loggerName string logger名稱
// @param enableColor bool 是否使用色碼
// @return *logrus.Logger
func GetLogger(loggerName string, enableColor bool) *logrus.Logger {
	// 取得設定檔
	var setting LogSetting
	if val, found := LogSettingMap[loggerName]; !found {
		setting = LogSetting{
			TimestampFormat:     LogDefaultTimestampFormat,
			TimeUnit:            LogDefaultTimeUnit,
			RotationType:        LogDefaultRotationType,
			WithMaxAge:          LogDefaultWithMaxAge,
			WithRotationTime:    LogDefaultWithRotationTime,
			WithRotationCount:   LogDefaultWithRotationCount,
			WithRotationSize_MB: LogDefaultWithRotationSize_MB,
			LogDir:              LogDefaultDir,
		}
	} else {
		setting = val
	}

	myFormatter := new(MyFormatter)
	myFormatter.LoggerName = loggerName
	myFormatter.TimestampFormat = setting.TimestampFormat
	myFormatter.LogFormat = LogDefaultLogFormat
	myFormatter.DisableColors = !enableColor

	myFileFormatter := new(MyFormatter)
	myFileFormatter.LoggerName = loggerName
	myFileFormatter.TimestampFormat = setting.TimestampFormat
	myFileFormatter.LogFormat = LogDefaultLogFormat
	myFileFormatter.DisableColors = true

	var tUnit time.Duration
	var fileTimeFormat string
	switch setting.TimeUnit {
	case "second":
		tUnit = time.Second
		fileTimeFormat = "%Y_%m%d_%H%M%S"
	case "minute":
		tUnit = time.Minute
		fileTimeFormat = "%Y_%m%d_%H%M"
	case "hour":
		tUnit = time.Hour
		fileTimeFormat = "%Y_%m%d_%H"
	case "day":
		tUnit = time.Hour * 24
		fileTimeFormat = "%Y_%m%d"
	default:
		tUnit = time.Hour
		fileTimeFormat = "%Y_%m%d_%H"
	}

	// 確認檔案分割方式
	rotationType := setting.RotationType
	if rotationType == "" {
		rotationType = LogDefaultRotationType
	}

	var logFile *rotatelogs.RotateLogs
	logDir := setting.LogDir
	if logDir == "" {
		logDir = LogDefaultDir
	} else {
		// 檔案路徑補加/
		str := []rune(logDir)
		lastChar := string(str[len(str)-1])
		if lastChar != "/" {
			logDir += "/"
		}
	}
	if rotationType == LogDefaultRotationType {
		logFile, _ = rotatelogs.New(
			logDir+loggerName+"_"+fileTimeFormat+".log", // 檔名時間格式化
			rotatelogs.WithLinkName(logDir+loggerName+".log"),
			rotatelogs.WithMaxAge(tUnit*time.Duration(setting.WithMaxAge)),             // 存活時限
			rotatelogs.WithRotationTime(tUnit*time.Duration(setting.WithRotationTime)), // 分割時間間隔
		)
	} else {
		fileTimeFormat = "%Y_%m%d"
		logFile, _ = rotatelogs.New(
			logDir+loggerName+".log", // 檔名時間格式化
			rotatelogs.WithRotationCount(setting.WithRotationCount),            // 最高留存file數量(不可與WithMaxAge混用)
			rotatelogs.WithRotationSize(setting.WithRotationSize_MB*1024*1024), // 分割size，單位byte(不可與WithRotateTime混用)
		)
	}

	logger := &logrus.Logger{
		Out:       os.Stdout,
		Level:     logrus.TraceLevel,
		Formatter: myFormatter,
		Hooks:     make(logrus.LevelHooks),
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: logFile, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  logFile,
		logrus.WarnLevel:  logFile,
		logrus.ErrorLevel: logFile,
		logrus.FatalLevel: logFile,
		logrus.PanicLevel: logFile,
	}, myFileFormatter)

	logger.AddHook(lfHook)
	return logger
}

// 讀取log設定檔
func getLogSetting() (sMap map[string]LogSetting, err error) {
	sMap = make(map[string]LogSetting)
	file, err := ioutil.ReadFile("config/logger.json")
	if err != nil {
		// fmt.Println("讀取log設定檔失敗 : ", err)
		return sMap, err
	}
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	_ = json.Unmarshal([]byte(file), &sMap)
	return sMap, nil
}
