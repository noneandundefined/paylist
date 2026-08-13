package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"paylist.server/infra/logger"
)

type Logger struct {
	currentLogDir string
	loc           *time.Location
}

func NewLogger() *Logger {
	loc, _ := time.LoadLocation("Asia/Yekaterinburg")

	return &Logger{loc: loc}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int64
	wroteHeader bool
}

func (l *Logger) getCurrentLogDir() string {
	now := time.Now().In(l.loc)

	day := fmt.Sprintf("%02d", now.Day())
	month := fmt.Sprintf("%02d", now.Month())
	year := fmt.Sprintf("%d", now.Year())

	dirName := fmt.Sprintf("log_%s%s%s", day, month, year)

	dir := filepath.Join(logger.LogDir(), dirName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create log directory: %s", err.Error())
	}

	return dir
}

// Обновление даты в папке
func (l *Logger) updateLogDirIfNeeded() {
	newLogDir := l.getCurrentLogDir()
	if newLogDir != l.currentLogDir {
		l.currentLogDir = newLogDir
	}
}

func (l *Logger) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         200,
		}

		next.ServeHTTP(wrapped, r)

		entry := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d \"%s\" %v\n",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			getHTTPVersion(r),
			r.Method,
			r.URL.Path,
			wrapped.status,
			wrapped.size,
			r.UserAgent(),
			time.Since(start),
		)

		l.updateLogDirIfNeeded()
		l.productionLogging(entry)
	})
}

func getHTTPVersion(r *http.Request) string {
	switch r.ProtoMajor {
	case 1:
		return "HTTP/1.1"
	case 2:
		return "HTTP/2.0"
	case 3:
		return "HTTP/3.0"
	default:
		return fmt.Sprintf("HTTP/%d.%d", r.ProtoMajor, r.ProtoMinor)
	}
}

func (l *Logger) productionLogging(log string) {
	APP_DIR := filepath.Join(l.currentLogDir, "server.log")

	appLogFile, err := os.OpenFile(APP_DIR, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("failed to open app log file: %s", err.Error())
		return
	}
	defer appLogFile.Close()

	_, err = appLogFile.WriteString(log)
	if err != nil {
		logger.Error("error writing to log file: %s", err.Error())
		return
	}
}

func (rw *responseWriter) WriteHeader(status int) {
	if rw.wroteHeader {
		return
	}

	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
	rw.wroteHeader = true
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}

	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

func (rw *responseWriter) Flush() {
	if fl, ok := rw.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacker not supported")
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	if cn, ok := rw.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	return make(chan bool)
}
