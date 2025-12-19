package logger

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

func New() *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	// log.SetFormatter(&logrus.JSONFormatter{}) // atau TextFormatter

	log.SetFormatter(&ColorFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		EnableColors:    true, // aktifkan warna
	})

	log.SetLevel(logrus.InfoLevel)

	return log
}

type ColorFormatter struct {
	TimestampFormat string
	EnableColors    bool
}

// ANSI colors
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

func (f *ColorFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// auto-detect TTY (disable color in file / docker logs pipe)
	enableColor := f.EnableColors && term.IsTerminal(int(os.Stdout.Fd()))

	// timestamp
	tsFormat := f.TimestampFormat
	if tsFormat == "" {
		tsFormat = time.RFC3339
	}

	fmt.Fprintf(
		&b,
		"%s [%s] %s",
		entry.Time.Format(tsFormat),
		entry.Level.String(),
		entry.Message,
	)

	// sort fields biar konsisten
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// append fields
	for _, k := range keys {
		val := entry.Data[k]

		// 🎨 warna khusus untuk status code
		if k == "status" && enableColor {
			if code, ok := val.(int); ok {
				fmt.Fprintf(&b, " %s=%s%d%s", k, statusColor(code), code, colorReset)
				continue
			}
		}

		fmt.Fprintf(&b, " %s=%v", k, val)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen
	case code >= 400 && code < 500:
		return colorYellow
	case code >= 500:
		return colorRed
	default:
		return ""
	}
}
