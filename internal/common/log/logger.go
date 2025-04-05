package log

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var logger *log.Logger

func getLogger() *log.Logger {
	if logger == nil {
		styles := log.DefaultStyles()
		styles.Levels[log.FatalLevel] = lipgloss.NewStyle().SetString("FATAL")
		styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().SetString("ERROR")
		styles.Levels[log.WarnLevel] = lipgloss.NewStyle().SetString("WARN")
		styles.Levels[log.InfoLevel] = lipgloss.NewStyle().SetString("INFO")
		logger = log.New(os.Stderr)
		logger.SetStyles(styles)
		logger.SetReportTimestamp(false)
	}
	return logger
}

func Info(msg interface{}, keyvals ...interface{}) {
	getLogger().Info(msg, keyvals...)
}

func Infof(format string, args ...any) {
	getLogger().Infof(format, args...)
}

func Error(msg interface{}, keyvals ...interface{}) {
	getLogger().Error(msg, keyvals...)
}

func Errorf(format string, args ...any) {
	getLogger().Errorf(format, args...)
}

func Warn(msg interface{}, keyvals ...interface{}) {
	getLogger().Warn(msg, keyvals...)
}

func Warnf(format string, args ...any) {
	getLogger().Warnf(format, args...)
}

func Fatal(msg interface{}, keyvals ...interface{}) {
	getLogger().Fatal(msg, keyvals...)
}

func Fatalf(format string, args ...any) {
	getLogger().Fatalf(format, args...)
}
