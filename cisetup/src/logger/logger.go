package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
	"syscall"
)

func file_line() string {
	_, fileName, fileLine, ok := runtime.Caller(2)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", fileName, fileLine)
	} else {
		s = ""
	}
	return s
}

func Log(format string, a ...interface{}) {
	strings.Replace(format, "\n", "", -1)
	t := time.Now()
	tm := fmt.Sprintf("%s", t.Format("2006.01.02 15:04:05"))
	buf := fmt.Sprintf(format, a...)
	f, err := os.OpenFile("/var/log/citool/citool.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("(%d)[%s]%s: %s\n", syscall.Getpid(), tm, file_line(), buf)
}

func LogString(a string) {
	Log("%s", a)
}
