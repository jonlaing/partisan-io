package logger

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Logger is a wrapper for the standard logger to add some extra functionality to it

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func init() {
	traceHandle := ioutil.Discard
	infoHandle := os.Stdout
	warningHandle := os.Stdout
	errorHandle := os.Stderr

	Trace = log.New(traceHandle,
		withColor("TRACE", white),
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		withColor("INFO", cyan),
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		withColor("WARNING", yellow),
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		withColor("ERROR", red),
		log.Ldate|log.Ltime|log.Lshortfile)

}

func withColor(label, color string) string {
	return fmt.Sprintf("%v%s:%v ", color, label, reset)
}
