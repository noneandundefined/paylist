package logger

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	logDir   string
	logQueue = make(chan logEntry, 2000)
	timeYkb  *time.Location
	dev      bool
)

type logEntry struct {
	filename string
	message  string
}

func LogDir() string {
	baseLogPath := os.Getenv("LOG_DIR")
	if baseLogPath == "" {
		baseLogPath = "./logs"
	}

	return strings.Trim(baseLogPath, `"' `)
}

func InitLogger() {
	var err error
	timeYkb, err = time.LoadLocation("Asia/Yekaterinburg")
	if err != nil {
		panic("Failed to load timezone: " + err.Error())
	}

	dev = os.Getenv("GO_ENV") == "DEV"

	ensureLogDir()
	if logDir == "" {
		fmt.Fprintf(os.Stderr, "Logger init failed: could not create log directory (LOG_DIR=%q)\n", LogDir())
	} else {
		fmt.Fprintf(os.Stderr, "Logger initialized: %s\n", logDir)
	}

	/* Workers */
	go logWorker()
	go compressWorker()
}

func logWorker() {
	for entry := range logQueue {
		writeLog(entry.filename, entry.message)
	}
}

func writeLog(filename, log string) {
	ensureLogDir()
	filePath := filepath.Join(logDir, filename)

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening log file %s: %v\n", filePath, err)
		return
	}
	defer f.Close()

	f.WriteString(log + "\n")
}

func writeDateTime(log *strings.Builder, t time.Time) {
	log.WriteString(fmt.Sprintf("%02d.%02d.%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute(), t.Second()))
}

func ensureLogDir() {
	now := time.Now().In(timeYkb)
	day := fmt.Sprintf("%02d", now.Day())
	month := fmt.Sprintf("%02d", now.Month())
	year := fmt.Sprintf("%d", now.Year())

	newLogDir := filepath.Join(LogDir(), fmt.Sprintf("log_%s%s%s", day, month, year))

	if newLogDir != logDir {
		if err := os.MkdirAll(newLogDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating log directory %s: %v\n", newLogDir, err)
			return
		}

		logDir = newLogDir
	}
}

func logWithLevel(level, filename, format string, args ...any) {
	now := time.Now().In(timeYkb)
	var log strings.Builder

	log.WriteString("[")
	writeDateTime(&log, now)
	log.WriteString("] [")
	log.WriteString(level)
	log.WriteString("] ")
	log.WriteString(fmt.Sprintf(format, args...))

	logMessage := log.String()

	/* If DEV write in console */
	if dev {
		fmt.Println(logMessage)
	}

	select {
	case logQueue <- logEntry{filename: filename, message: logMessage}:
	default:
		fmt.Fprintf(os.Stderr, "Log queue full, message dropped: %s\n", logMessage)
	}
}

func Info(format string, args ...any) {
	logWithLevel("INFO", "server.log", format, args...)
}

func Warning(format string, args ...any) {
	logWithLevel("WARN", "warning.log", format, args...)
}

func Session(format string, args ...any) {
	logWithLevel("SESS", "session.log", format, args...)
}

func Error(format string, args ...any) {
	logWithLevel("ERROR", "errors.log", format, args...)
}

func AI(format string, args ...any) {
	logWithLevel("AI", "ai.log", format, args...)
}

func Moderation(format string, args ...any) {
	logWithLevel("MOD", "moderation.log", format, args...)
}

func compressWorker() {
	compressOldDirs()

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		compressOldDirs()
	}
}

func compressOldDirs() {
	baseLogPath := filepath.Dir(logDir)

	entries, err := os.ReadDir(baseLogPath)
	if err != nil {
		fmt.Println("Error reading log directory:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(baseLogPath, entry.Name())
		if fullPath == logDir {
			continue
		}

		zipPath := fullPath + ".zip"
		if _, err := os.Stat(zipPath); err == nil {
			continue
		}

		if err := zipFolder(fullPath, zipPath); err != nil {
			fmt.Println("Error compressing folder:", fullPath, err)
		} else {
			os.RemoveAll(fullPath)
		}
	}
}

func zipFolder(src, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		fw, err := w.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, f)
		return err
	})
}
