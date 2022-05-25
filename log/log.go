package log

import (
	"path"
	"runtime"
	"strconv"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
)

const (
	PanicLevel int = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

// CreateLogrusInstance。
// 获取输出日志到文件的 logrus 实例。
// filePath：日志输出文件路径。
// size：日志文件最大 size, 单位是 MB。
// backups：最大过期日志保留的个数。
// local：控制是否打印时间。
// level：日志级别。
func CreateLogrusInstance(filePath string, size, backups, level int, local bool) *log.Logger {
	var instance = log.New()
	logger := &lumberjack.Logger{
		Filename:   filePath, // 日志输出文件路径
		MaxSize:    size,     // 日志文件最大 size, 单位是 MB
		MaxBackups: backups,  // 最大过期日志保留的个数
		MaxAge:     30,       // 保留过期文件的最大时间间隔,单位是天
		Compress:   false,    // disabled by default
		LocalTime:  local,
	}
	instance.SetOutput(logger)
	instance.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			functionName := path.Base(frame.Function)
			return functionName, fileName
		},
	})
	instance.SetReportCaller(true)
	instance.SetLevel(log.Level(level))
	return instance
}
