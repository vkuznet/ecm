package main

// logging module provides various logging methods
//
// Copyright (c) 2020 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"fmt"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// RotateLogWriter represents rorate log writer
type RotateLogWriter struct {
	RotateLogs *rotatelogs.RotateLogs
}

// Write method of our rotate log writer
func (w RotateLogWriter) Write(data []byte) (int, error) {
	return w.RotateLogs.Write([]byte(data))
}

// custom logger
type LogWriter struct {
}

// Write method of our log writer
func (writer LogWriter) Write(data []byte) (int, error) {
	return fmt.Print(data)
}
