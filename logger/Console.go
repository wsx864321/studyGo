package logger

import (
	"fmt"
)

type ConsoleLogger struct {
	level    int
}

//创建logger实例
func NewConsoleLogger(config map[string]string) LogInterface {
	var(
		level    int
		levelStr string
		logger   ConsoleLogger
		ok       bool
	)

	if levelStr,ok = config["log_level"];!ok {
		panic("log level don`t set")
	}
	level = getLogLevel(levelStr)

	logger.SetLevel(level)

	return &logger
}

func (f *ConsoleLogger)SetLevel(level int) {
	if level < LOG_LEVEL_DEBUG || level > LOG_LEVEL_FATAL{
		f.level = LOG_LEVEL_DEBUG
	}

	f.level = level
}

func (c *ConsoleLogger)Debug(format string, arg... interface{}){
	if c.level > LOG_LEVEL_DEBUG {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_DEBUG, format, arg...)
	fmt.Printf("%s 【%s】 (%s:%s:%d) " + log.Msg +
		"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
}

func (c *ConsoleLogger)Trace(format string, arg... interface{}){
	if c.level > LOG_LEVEL_TRACE {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_TRACE, format, arg...)
	fmt.Printf("%s 【%s】 (%s:%s:%d) " + log.Msg +
		"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
}

func (c *ConsoleLogger)Info(format string, arg... interface{}){
	if c.level > LOG_LEVEL_INFO {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_INFO, format, arg...)
	fmt.Printf("%s 【%s】 (%s:%s:%d) " + log.Msg +
		"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
}

func (c *ConsoleLogger)Warn(format string, arg... interface{}){
	if c.level > LOG_LEVEL_WARN {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_WARN, format, arg...)
	fmt.Printf("%s 【%s】 (%s:%s:%d) " + log.Msg +
		"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
}

func (c *ConsoleLogger)Error(format string, arg... interface{}) {
	if c.level > LOG_LEVEL_ERROR {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_ERROR, format, arg...)
	fmt.Printf("%s 【%s】 (%s:%s:%d) " + log.Msg +
		"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
}

func (c *ConsoleLogger)Fatal(format string, arg... interface{}) {
	if c.level > LOG_LEVEL_FATAL {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_FATAL, format, arg...)
	fmt.Printf("%s 【%s】 (%s:%s:%d) " + log.Msg +
		"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
}

func (c *ConsoleLogger)Close(){

}


