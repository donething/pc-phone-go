package logger

import (
	"github.com/donething/utils-go/dolog"
	"log"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

const LogName = "run.log"

func init() {
	Info, Warn, Error = dolog.InitLog(LogName, dolog.DefaultFormat)
}
