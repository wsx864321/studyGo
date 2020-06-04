package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type FileLogger struct {
	level        int
	logPath      string
	logName      string
	file         *os.File
	fileWarn     *os.File
	logChan      chan *LogData
	splitHour    int
	splitType    string
	splitSize    int64
	preSplitTime int64
}
//创建FileLogger实例
func NewFileLogger(config map[string]string) LogInterface {
	var(
		logPath   string
		logName   string
		levelStr  string
		level     int
		logger    FileLogger
		fileName  string
		file      *os.File
		err       error
		fileWarn  *os.File
		ok        bool
		splitType string
		splitHour int
		splitSize int64
		dateTime  string
		loc       *time.Location
		tmp       time.Time
		timeStamp int64
	)

	if logPath,ok = config["log_path"];!ok {
		panic("log path don`t set")
	}
	if logName,ok = config["log_name"];!ok {
		panic("log name don`t set")
	}
	if levelStr,ok = config["log_level"];!ok {
		panic("log level don`t set")
	}
	if splitType,ok = config["split_type"];!ok {
		panic("log split type don`t set")
	}

	if splitType == SPLIT_TYPE_HOUR {
		if _,ok = config["split_hour"];ok {
			splitHour,_ = strconv.Atoi(config["split_hour"])
		}else{
			splitHour = 1 //默认是1小时切割
		}
	}else{
		if _,ok = config["split_size"];ok {
			splitSize,_ = strconv.ParseInt(config["split_size"], 10, 64)
		}else{
			splitSize = 300 * 1024 * 1024 //默认是300M
		}
	}

	level = getLogLevel(levelStr)

	//打开普通日志
	fileName = fmt.Sprintf("%s/%s.log",logPath,logName)
	file,err = os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil{
		panic(fmt.Sprintf("open %s failed,err %v:", fileName, err))
	}

	//打开错误日志
	fileName = fmt.Sprintf("%s/%s.log.wf",logPath,logName)
	fileWarn,err = os.OpenFile(fileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil{
		panic(fmt.Sprintf("open %s failed,err %v:", fileName, err))
	}

	loc, _ = time.LoadLocation("Local")    //获取时区
	dateTime =  time.Now().Format("2006-01-02 15:00:00")
	tmp,_ =  time.ParseInLocation("2006-01-02 15:04:05", dateTime, loc)
	timeStamp = tmp.Unix()

	logger = FileLogger{
		level:        level,
		logPath:      logPath,
		logName:      logName,
		file:         file,
		fileWarn:     fileWarn,
		logChan:      make(chan *LogData, 1000), //暂时是1000吧，也不知道多少合适
		splitSize:    splitSize,
		splitType:    splitType,
		splitHour:    splitHour,
		preSplitTime: timeStamp,
	}

	go logger.writeLogBackground()

	return &logger
}

func (f *FileLogger)SetLevel(level int) {
	if level < LOG_LEVEL_DEBUG || level > LOG_LEVEL_FATAL{
		f.level = LOG_LEVEL_DEBUG
	}

	f.level = level
}

func (f *FileLogger)Debug(format string, arg... interface{}){
	if f.level > LOG_LEVEL_DEBUG {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_DEBUG, format, arg...)
	f.logChan <- log
}

func (f *FileLogger)Trace(format string, arg... interface{}){
	if f.level > LOG_LEVEL_TRACE {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_TRACE, format, arg...)
	f.logChan <- log
}

func (f *FileLogger)Info(format string, arg... interface{}){
	if f.level > LOG_LEVEL_INFO {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_INFO, format, arg...)
	f.logChan <- log
}

func (f *FileLogger)Warn(format string, arg... interface{}){
	if f.level > LOG_LEVEL_WARN {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_WARN, format, arg...)
	f.logChan <- log
}

func (f *FileLogger)Error(format string, arg... interface{}) {
	if f.level > LOG_LEVEL_ERROR {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_ERROR, format, arg...)
	f.logChan <- log
}

func (f *FileLogger)Fatal(format string, arg... interface{}){
	if f.level > LOG_LEVEL_FATAL {
		return
	}

	var(
		log *LogData
	)

	log = formatLog(LOG_LEVEL_FATAL, format, arg...)
	f.logChan <- log
}

func (f *FileLogger)Close(){
	_ = f.file.Close()
	_ = f.fileWarn.Close()
	close(f.logChan)
}

//异步写日志
func (f *FileLogger)writeLogBackground(){
	var(
		log *LogData
	)

	for log = range f.logChan {
		f.realWriteLog(log)
	}
}

func (f *FileLogger)realWriteLog(log *LogData){
	if f.splitType == SPLIT_TYPE_HOUR{
		f.splitFileByHour()
	}else{
		f.splitFileBySize(log)
	}

	if log.LogType == WF_LOG {
		_,_ = fmt.Fprintf(f.fileWarn, "%s 【%s】 (%s:%s:%d) " + log.Msg +
			"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
	}else{
		_,_ = fmt.Fprintf(f.file, "%s 【%s】 (%s:%s:%d) " + log.Msg +
			"\n",log.Date,log.LevelStr,log.FileName,log.FuncName,log.LineNo)
	}
}

//按大小分割日志
func (f *FileLogger)splitFileBySize(log *LogData){
	var(
		file          *os.File
		fileInfo      os.FileInfo
		err           error
		size          int64
		logPath       string
		backupLogPath string
		logDate       string
	)
	if log.LogType == WF_LOG {
		file = f.fileWarn
	}else{
		file = f.file
	}

	fileInfo,err = file.Stat()
	if err != nil {
		return
	}

	size = fileInfo.Size()

	if size < f.splitSize {
		return
	}

	logDate = time.Unix(f.preSplitTime,0).Format("2006010215")
	if log.LogType == WF_LOG {
		backupLogPath = fmt.Sprintf("%s/%s.log.wf_%s",f.logPath,f.logName, logDate)
		logPath = fmt.Sprintf("%s/%s.log.wf",f.logPath,f.logName)
	}else{
		backupLogPath = fmt.Sprintf("%s/%s.log_%s",f.logPath,f.logName, logDate)
		logPath = fmt.Sprintf("%s/%s.log",f.logPath,f.logName)
	}

	file.Close()
	os.Rename(logPath, backupLogPath)

	file,_ = os.OpenFile(logPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	if log.LogType == WF_LOG {
		f.fileWarn = file
	}else{
		f.file = file
	}
}
//按时间分割日志
func (f *FileLogger)splitFileByHour(){
	var(
		curTime         int64
		backupWfLogPath string
		backupLogPath   string
		logDate         string
		wfLogPath       string
		logPath         string
		dateTime        string
		loc             *time.Location
		tmp             time.Time
	)

	curTime = time.Now().Unix()
	if  curTime - f.preSplitTime < int64(f.splitHour * 3600)  {
		return
	}

	loc, _ = time.LoadLocation("Local")    //获取时区
	dateTime =  time.Now().Format("2006-01-02 15:00:00")
	tmp,_ =  time.ParseInLocation("2006-01-02 15:04:05", dateTime, loc)
	f.preSplitTime = tmp.Unix()

	logDate = time.Unix(f.preSplitTime,0).Format("2006010215")
	wfLogPath =fmt.Sprintf("%s/%s.log.wf",f.logPath,f.logName)
	logPath = fmt.Sprintf("%s/%s.log",f.logPath,f.logName)
	backupWfLogPath = fmt.Sprintf("%s/%s.log.wf_%s",f.logPath,f.logName, logDate)
	backupLogPath = fmt.Sprintf("%s/%s.log_%s",f.logPath,f.logName, logDate)

	f.file.Close()
	f.fileWarn.Close()

	os.Rename(logPath,backupLogPath)
	os.Rename(wfLogPath,backupWfLogPath)

	f.file,_ = os.OpenFile(logPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
	f.fileWarn,_ = os.OpenFile(wfLogPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0755)
}

