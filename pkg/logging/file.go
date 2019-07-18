package logging

import (
	"caspar/gin-blog/pkg/file"
	"caspar/gin-blog/pkg/setting"
	"fmt"
	"os"
	"time"
)

/*var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060102"
)*/

func getLogFilePath() string {
	return fmt.Sprintf("%s", setting.AppSetting.LogSavePath)
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
}

/*func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat), setting.AppSetting.LogFileExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}*/

func openLogFile(filepath string, filename string) (*os.File, error){
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err:%v", err)
	}
	src := dir + "/" + filepath
	perm := file.CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src:", src)
	}
	err = file.IsNotExisMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExisMkDir src:%s, err:%v", src, err)
	}

	f, err := file.Open(src+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, fmt.Errorf("fail to openFile:%v", err)
	}
	return f, nil
}


