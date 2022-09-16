package zdmiddleware

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	statusOKStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ADFAAD"))

	statusErrStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAADAD"))

	statusWarnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAAD"))
)

func Logrus(logger *logrus.Logger) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		start := time.Now()
		ctx.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()
		clientUserAgent := ctx.Request.UserAgent()
		referer := ctx.Request.Referer()
		dataLength := ctx.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		entry := logger.WithFields(logrus.Fields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency, // time to process
			"clientIP":   clientIP,
			"method":     ctx.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(ctx.Errors) > 0 {
			entry.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			statusCodeStr := fmt.Sprint(statusCode)

			var lg func(...interface{})
			if statusCode >= http.StatusInternalServerError {
				lg = entry.Error
				statusCodeStr = statusErrStyle.Render(statusCodeStr)
			} else if statusCode >= http.StatusBadRequest {
				lg = entry.Warn
				statusCodeStr = statusWarnStyle.Render(statusCodeStr)
			} else {
				lg = entry.Info
				statusCodeStr = statusOKStyle.Render(statusCodeStr)
			}

			lg(fmt.Sprintf("[%s %s] %s \"%s\" (%dms)",
				ctx.Request.Method, path, statusCodeStr, clientUserAgent, latency,
			))
		}
	}
}
