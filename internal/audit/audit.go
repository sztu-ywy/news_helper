package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/signal"
	"time"

	"git.uozi.org/uozi/burn-api/model"
	"git.uozi.org/uozi/burn-api/settings"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cast"
	"github.com/uozi-tech/cosy/logger"
)

const Topic = "audit"

var producerInstance *producer.Producer

func Init(ctx context.Context) {
	if !settings.SLSSettings.Enable() {
		return
	}
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.Logger = &ZapLogger{
		logger: logger.GetLogger(),
	}
	producerConfig.Endpoint = settings.SLSSettings.EndPoint
	provider := getCredentialsProvider()
	producerConfig.CredentialsProvider = provider
	// if you want to use log context, set the GeneratePackId to true
	producerConfig.GeneratePackId = true
	producerConfig.LogTags = []*sls.LogTag{
		{
			Key:   proto.String("type"),
			Value: proto.String(Topic),
		},
	}
	producerInstance = producer.InitProducer(producerConfig)
	producerInstance.Start()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	defer producerInstance.SafeClose()

	<-ch
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !settings.SLSSettings.Enable() {
			c.Next()
			return
		}
		startTime := time.Now()
		ip := c.ClientIP()
		reqURL := c.Request.URL.String()
		reqHeader := c.Request.Header
		reqMethod := c.Request.Method
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = c.GetRawData()
		}
		reqBody := string(reqBodyBytes)
		// re-assigned the request body to the original one, to prevent the request body from being consumed
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))

		responseBodyWriter := &responseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = responseBodyWriter

		// continue the request
		c.Next()

		// get the response meta
		respStatusCode := cast.ToString(c.Writer.Status())
		respHeader := c.Writer.Header()
		respBody := responseBodyWriter.body.String()
		latency := time.Since(startTime).String()

		var userId uint64
		if user, ok := c.Get("user"); ok {
			userId = user.(*model.User).ID
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Error(r)
				}
			}()
			reqHeaderBytes, _ := json.Marshal(reqHeader)
			respHeaderBytes, _ := json.Marshal(respHeader)
			log := producer.GenerateLog(uint32(time.Now().Unix()), map[string]string{
				"ip":               ip,
				"user_id":          cast.ToString(userId),
				"req_url":          reqURL,
				"req_method":       reqMethod,
				"req_header":       string(reqHeaderBytes),
				"req_body":         reqBody,
				"resp_header":      string(respHeaderBytes),
				"resp_status_code": respStatusCode,
				"resp_body":        respBody,
				"latency":          latency,
			})
			err := producerInstance.SendLog(settings.SLSSettings.ProjectName,
				settings.SLSSettings.LogStoreName, Topic, settings.SLSSettings.Source, log)
			if err != nil {
				logger.Error(err)
			}
		}()
	}
}
