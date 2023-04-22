package logger

import (
	"fmt"
	"os"
	"path"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zput/zxcTool/ztLog/zt_formatter"
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
}

func New(fileName string) *Logger {
	lg := logrus.New()
	if fileName == "" {
		lg.Out = os.Stdout
	} else {
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		lg.Out = f
	}
	lg.SetReportCaller(true)

	formatter := &zt_formatter.ZtFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		Formatter: nested.Formatter{
			//HideKeys: true,
			TimestampFormat: "2006-01-02T15:04:05Z07:00",
			FieldsOrder:     []string{"req_id"},
		},
	}

	lg.SetFormatter(formatter)

	return &Logger{lg}
}
