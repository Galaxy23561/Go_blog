package middleware

import (
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"os"
	"fmt"
	"github.com/rifflock/lfshook"
	retalog "github.com/lestrrat-go/file-rotatelogs"
	"time"
)

func Logger() gin.HandlerFunc{
	filePath:="log/log"

	src,err:=os.OpenFile(filePath,os.O_RDWR|os.O_CREATE,0755)
	if err!=nil{
		fmt.Println("err:",err)
	}
	logger:=logrus.New()

	logger.Out=src

	logger.SetLevel(logrus.DebugLevel)

	logWriter, _:=retalog.New(
		filePath+"%Y%m%d.log",
		retalog.WithMaxAge(7*24*time.Hour),
		retalog.WithRotationTime(24*time.Hour),
	)

	writeMap:=lfshook.WriterMap{
		logrus.InfoLevel:logWriter,
		logrus.FatalLevel:logWriter,
		logrus.DebugLevel:logWriter,
		logrus.ErrorLevel:logWriter,
		logrus.PanicLevel:logWriter,
		logrus.WarnLevel:logWriter,
	}
	Hook:=lfshook.NewHook(writeMap,&logrus.TextFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})

	logger.AddHook(Hook)

	return func(c *gin.Context){
		startTime:=time.Now()
		c.Next()
		stopTime:=time.Since(startTime)
		spendTime:=fmt.Sprintf("%d ms",stopTime)
		hostName,err:=os.Hostname()
		if err!=nil{
			hostName="unknow"
		}
		statusCode:=c.Writer.Status()
		clientIp:=c.ClientIP()
		userAgent:=c.Request.UserAgent()
		dataSize:=c.Writer.Size()
		if dataSize<0{
			dataSize=0
		}
		method:=c.Request.Method
		path:=c.Request.RequestURI

		entry:=logger.WithFields(logrus.Fields{
			"hostname":hostName,
			"status_code":statusCode,
			"spend_time":spendTime,
			"client_ip":clientIp,
			"user_agent":userAgent,
			"data_size":dataSize,
			"method":method,
			"path":path,
		})
		if len(c.Errors)>0{
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode>=500{
			entry.Error()
		} else if statusCode>=400{
			entry.Warn()
		} else {
			entry.Info()
		}
	}
}