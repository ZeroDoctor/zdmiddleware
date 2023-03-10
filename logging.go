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
	white = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F0F0F0"))

	yellow = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FDFDAA"))

	green = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ADFFAD"))

	magenta = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FDAAFD"))

	red = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA0A0"))

	blue = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#1080FF"))

	cyan = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#10E0FF"))
)

func statusCodeColor(code int) string {
	codeStr := fmt.Sprint(code)

	var style lipgloss.Style

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		style = green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		style = white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		style = yellow
	default:
		style = red
	}

	return style.Render(codeStr)
}

func methodColor(method string) string {
	var style lipgloss.Style

	switch method {
	case http.MethodGet:
		style = blue
	case http.MethodPost:
		style = cyan
	case http.MethodPut:
		style = yellow
	case http.MethodDelete:
		style = red
	case http.MethodPatch:
		style = green
	case http.MethodHead:
		style = magenta
	case http.MethodOptions:
		style = white
	default:
		style = lipgloss.NewStyle()
	}

	return style.Render(method)
}

func GinLogrus(logger *logrus.Logger) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		stop := time.Since(startTime)
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
			"query":      raw,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  clientUserAgent,
		})

		if len(ctx.Errors) > 0 {
			entry.Error(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {

			method := methodColor(ctx.Request.Method)
			status := statusCodeColor(statusCode)

			logger.Printf("%s [%s?] %s \"%s\" (%dms)",
				method, path, raw, status, clientUserAgent, latency,
			)

		}
	}
}
