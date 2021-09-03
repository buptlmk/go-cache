package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

var Logger = logrus.New()

type format struct {
}

func (f *format) Format(entry *logrus.Entry) ([]byte, error) {

	str := fmt.Sprintf("[%s] [%s] [%s:%d]-->%v\n", entry.Time.Format("2006-01-02 15:04:05"), strings.ToUpper(entry.Level.String()), entry.Caller.File, entry.Caller.Line, entry.Message)
	return []byte(str), nil
}

func init() {

	Logger.SetLevel(logrus.InfoLevel)

	Logger.SetFormatter(&format{})
	Logger.SetReportCaller(true)
}
