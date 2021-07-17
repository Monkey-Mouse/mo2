package basiclog

import "log"

var Logger *log.Logger
var ErrLogger *log.Logger

func SetLoggeer(logger *log.Logger, errLogger *log.Logger) {
	Logger = logger
	ErrLogger = errLogger
}
