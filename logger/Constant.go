package logger

//日志级别配置
const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_TRACE
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
)
//普通日志和错误日志
const (
	NOMAL_LOG = iota
	WF_LOG
)
//日志打印类型
const (
	FILE_TYPE = iota
	CONSOLE_TYPE
)
//日志拆分类型
const (
	SPLIT_TYPE_SIZE = "size"
	SPLIT_TYPE_HOUR = "hour"
)

func getLogLevelText(levle int) string {
	switch levle {
	case LOG_LEVEL_DEBUG:
		return "DEBUG"
	case LOG_LEVEL_TRACE:
		return "TRACE"
	case LOG_LEVEL_INFO:
		return "INFO"
	case LOG_LEVEL_WARN:
		return "WARN"
	case LOG_LEVEL_ERROR:
		return "ERROR"
	case LOG_LEVEL_FATAL:
		return "FATAL"
	}

	return "UNKNOW"
}

func getLogLevel(str string) int {
	switch str {
	case "debug":
		return LOG_LEVEL_DEBUG
	case "trace":
		return LOG_LEVEL_TRACE
	case "info":
		return LOG_LEVEL_INFO
	case "warn":
		return LOG_LEVEL_WARN
	case "error":
		return LOG_LEVEL_ERROR
	case "fatal":
		return LOG_LEVEL_FATAL
	}

	return LOG_LEVEL_DEBUG
}