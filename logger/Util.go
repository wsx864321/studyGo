package logger

import (
	"fmt"
	"path"
	"runtime"
	"time"
)

type LogData struct {
	Msg       string
	LevelStr  string
	LogType   int //是否存储到wf
	Date      string
	FileName  string
	FuncName  string
	LineNo    int
}


//获取函数行数信息
func getLineInfo() (string,string,int) {
	var(
		pc       uintptr
		file     string
		lineNo   int
		ok       bool
		funcName string
	)

	pc,file,lineNo,ok = runtime.Caller(3)
	if ok {
		funcName = runtime.FuncForPC(pc).Name()
	}

	return file,funcName,lineNo
}

//格式化日志数据
func formatLog(level int,format string,arg... interface{}) *LogData {
	var(
		msg      string
		date     string
		file     string
		funcName string
		LevelStr string
		lineNo   int
		logType  int
	)

	date = time.Now().Format("2006-01-02 15:04:05.999")
	file,funcName,lineNo = getLineInfo()
	file = path.Base(file)//去掉项目目录前面的地址信息
	funcName = path.Base(funcName)//去掉项目目录前面的地址信息
	msg = fmt.Sprintf(format,arg...)
	LevelStr = getLogLevelText(level)
	if level == LOG_LEVEL_ERROR || level == LOG_LEVEL_FATAL {
		logType = WF_LOG
	}else{
		logType = NOMAL_LOG
	}

	return &LogData{
		Msg:msg,
		LevelStr:LevelStr,
		LogType:logType,
		Date:date,
		FileName:file,
		FuncName:funcName,
		LineNo:lineNo,
	}

}
