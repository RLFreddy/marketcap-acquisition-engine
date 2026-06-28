package logger

import (
	"os"

	"github.com/fatih/color"
)

var (
	infoColor    = color.New(color.FgCyan, color.Bold)
	successColor = color.New(color.FgGreen, color.Bold)
	warnColor    = color.New(color.FgYellow, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	traceColor   = color.New(color.FgHiBlack)
)

// Info logs an informational message
func Info(format string, a ...interface{}) {
	infoColor.Printf(format+"\n", a...)
}

// Success logs a successful operation
func Success(format string, a ...interface{}) {
	successColor.Printf(format+"\n", a...)
}

// Warn logs a warning message
func Warn(format string, a ...interface{}) {
	warnColor.Printf(format+"\n", a...)
}

// Error logs an error message
func Error(format string, a ...interface{}) {
	errorColor.Printf(format+"\n", a...)
}

// Trace logs verbose/progress messages
func Trace(format string, a ...interface{}) {
	traceColor.Printf(format+"\n", a...)
}

// Fatal logs an error message and terminates the program
func Fatal(format string, a ...interface{}) {
	errorColor.Printf(format+"\n", a...)
	os.Stdout.Sync()
	os.Stderr.Sync()
	os.Exit(1)
}
