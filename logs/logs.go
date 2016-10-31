package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

const DefaultLogName = "rtb.log"

var (
	curLogger    atomic.Value
	debugLogger  atomic.Value
	MainLogname  string
	logger_debug *log.Logger = log.New(os.Stderr, ``, log.Ldate|log.Ltime)
)

func init() {
	logger := log.New(os.Stderr, ``, log.Ldate|log.Ltime)
	curLogger.Store(logger)
	debugLogger.Store(logger)
}

// ==

// rtb_2016-04-05_15.log

func UpdateCurrentLogger(filesufix string) {
	if newLogger := NewLogger(filesufix); newLogger != nil {
		newLogger.SetFlags(0)
		curLogger.Store(newLogger)
	}
}

func UpdateDebugLogger(filesufix string) {
	if newLogger := NewDebugLogger("_debug" + filesufix); newLogger != nil {
		debugLogger.Store(newLogger)
	}
}

func Init(logname string) {

	if logname == "" {
		MainLogname = DefaultLogName
	} else {
		MainLogname = logname
	}

	tm := time.Now()
	currentLogSufix := fmt.Sprintf("_%04d-%02d-%02d_%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour())

	UpdateCurrentLogger(currentLogSufix)
	UpdateDebugLogger(currentLogSufix)

	redirectStderr()

	go LogSwitcher()
}

// ==

func leftTime() time.Duration {

	tm := time.Now()
	return time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour()+1, 0, 0, 0, tm.Location()).Sub(tm)
}

func LogSwitcher() {
	for {

		time.Sleep(leftTime() + time.Second)

		tm := time.Now()
		currentLogSufix := fmt.Sprintf("_%04d-%02d-%02d_%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour())

		UpdateCurrentLogger(currentLogSufix)
		UpdateDebugLogger(currentLogSufix)
	}
}

// ==

func Report(typ string, test, seatId int, requestId, data string) {

	logger := curLogger.Load().(*log.Logger)

	unixTime := time.Now().UnixNano() / 1000

	logger.Print(unixTime, "\t", typ, "\t", test, "\t", seatId, "\t", requestId, "\t", data)
}

func Write(arg ...interface{}) {

	logger := curLogger.Load().(*log.Logger)
	logger.Println(arg...)
}

func Critical(arg ...interface{}) {

	args := make([]interface{}, 0, len(arg)+1)
	args = append(args, "Critical:")
	args = append(args, arg...)

	Debug(args...)
}

func Debug(arg ...interface{}) {

	logger := debugLogger.Load().(*log.Logger)
	logger.Println(arg...)
}

func Recover() {

	if err := recover(); err != nil {

		pc, file, line, ok := runtime.Caller(4)

		if !ok {
			file = "?"
			line = 0
		}

		fn_name := ""
		fn := runtime.FuncForPC(pc)

		if fn == nil {
			fn_name = "?()"
		} else {
			dot_name := filepath.Ext(fn.Name())
			fn_name = strings.TrimLeft(dot_name, ".") + "()"
		}

		var buf [10240]byte
		number := runtime.Stack(buf[:], false)

		debugStr := fmt.Sprintf("%s:%d %s: %s\n\n%s", file, line, fn_name, err, buf[:number])

		Critical(debugStr) // test
	}
}

//

func getLogName(logname, sufix string) string {

	return strings.TrimSuffix(logname, ".log") + sufix + ".log"
}

func NewDebugLogger(sufix string) *log.Logger {

	logname := getLogName("./bidder.log", sufix)

	file, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Critical("log file \"%s\" cannot be written", logname)
		return nil
	}

	return log.New(file, ``, log.Ldate|log.Ltime)
}

func NewLogger(sufix string) *log.Logger {

	logname := getLogName(MainLogname, sufix)

	file, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		Critical("log file \"%s\" cannot be written", logname)
		return nil
	}

	return log.New(file, ``, log.Ldate|log.Ltime)
}

func redirectStderr() {

	logname := "bidder_crash.log"

	file, err := os.OpenFile(logname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to redirect stderr to file: %v", err)
	}

	if err := syscall.Dup2(int(file.Fd()), 2); err != nil {
		fmt.Printf("Failed to redirect stderr to file: %v", err)
	}
}
