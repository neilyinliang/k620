package global

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func LogLevel(level string) slog.Level {
	l := slog.LevelDebug
	switch strings.ToUpper(level) {
	case "DEBUG":
		l = slog.LevelDebug
	case "INFO":
		l = slog.LevelInfo
	case "WARN":
		l = slog.LevelWarn
	case "ERROR":
		l = slog.LevelError
	default:
		l = slog.LevelError
	}
	return l
}

func SetupLogger(path string, level string) (fd *os.File) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	slog.SetLogLoggerLevel(LogLevel(level))
	fd = os.Stdout
	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Println("Failed to open log file:", err, path)
		} else {
			fd = file
		}
	}
	log.SetOutput(fd)
	return fd
}
