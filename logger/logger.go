package logger

import "errors"

var(
	G_Logger LogInterface
)

func InitLogger(logType int,config map[string]string) error {
	switch logType {
	case FILE_TYPE:
		G_Logger = NewFileLogger(config)
	case CONSOLE_TYPE:
		G_Logger = NewConsoleLogger(config)
	default:
		return errors.New("不存在的日志打印类型")
	}

	return nil
}
