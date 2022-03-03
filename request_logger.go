package gin_request_logger

import (
	"bytes"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handler struct {
	logger *zap.Logger
}

func New(options Options) gin.HandlerFunc {
	handler := handler{
		logger: zap.NewExample(zap.Development()),
	}

	if options.Logger != nil {
		handler.logger = options.Logger
	}

	return handler.handle(options)
}

func (h *handler) handle(options Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-ID")
		if requestId == "" {
			requestId = uuid.New().String()
		}

		c.Set(RequestContextUUIDTag, requestId)

		c.Header("X-Request-ID", requestId)

		start := time.Now().UTC()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				h.logger.Error(e)
			}
		} else {
			fields := []zapcore.Field{
				zap.String("request-id", requestId),
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
				zap.String("time", end.Format(time.RFC3339)),
			}
			if c.Request.Method == "POST" ||
				c.Request.Method == "PUT" ||
				c.Request.Method == "PATCH" {
				bufs, err := ioutil.ReadAll(c.Request.Body)
				if err != nil {
					h.logger.Error("error while reading request body", zap.Error(err))
				}
				fields = append(fields, zap.String("request-body", string(bufs)))
			}
			if options.LogResponse {
				fields = append(fields, zap.String("response-body", blw.body.String()))
			}
			h.logger.Info(path, fields...)
		}

	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (h *handler) readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		h.logger.Error("read body", zap.Error(err))
	}

	s := buf.String()
	return s
}
