package logger

import (
	"testing"
	"time"
)

func TestFileLogger(t *testing.T){
	_ = InitLogger(FILE_TYPE, map[string]string{
		"log_level":"debug",
		"log_path":"/usr/wsx/www/go/src/studyGo/logger",
		"log_name":"application",
		"split_type":"size",
		"split_size":"500",
	})
	//G_Logger.Debug("debug %s","====")
	//G_Logger.Info("Info %s","====")
	G_Logger.Error("error %s","====")
	time.Sleep(1 * time.Second)
	G_Logger.Close()
}

func TestConsoleLogger(t *testing.T) {
	_ = InitLogger(CONSOLE_TYPE, map[string]string{"log_level":"debug",})
	G_Logger.Debug("debug %s","====")
	G_Logger.Info("Info %s","====")
	G_Logger.Error("error %s","====")
	G_Logger.Close()
}