package logger

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
	"os"
)

const ReqIDKey = "req_id"

func RequestIDSetter() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set(ReqIDKey, reqID)
	}
}

func RequestLogger(lg *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		lg := lg.WithFields(map[string]interface{}{"req_id": c.Value(ReqIDKey)})
		lg.WithFields(map[string]interface{}{"client": c.ClientIP()}).Info()
	}
}

type Logger struct {
	*logrus.Logger
	file *os.File
}

func (lg *Logger) Close() {
	lg.file.Close()
}
func New(fileName string, logLevel string) *Logger {
	lg := logrus.New()
	var file *os.File = nil
	if fileName == "" {
		lg.Out = os.Stdout
	} else {
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		lg.Out = f
		file = f
	}
	if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return nil
		}
		lg.SetLevel(level)
	}
	lg.SetReportCaller(true)

	formatter := &prefixed.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	}
	lg.SetFormatter(formatter)

	return &Logger{lg, file}
}
