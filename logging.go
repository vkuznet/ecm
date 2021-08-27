package main

// logging module provides various logging methods
//
// Copyright (c) 2020 - Valentin Kuznetsov <vkuznet@gmail.com>
//

import (
	"fmt"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// custom rotate logger
type RotateLogWriter struct {
	RotateLogs *rotatelogs.RotateLogs
}

func (w RotateLogWriter) Write(data []byte) (int, error) {
	return w.RotateLogs.Write([]byte(data))
}

// custom logger
type LogWriter struct {
}

func (writer LogWriter) Write(data []byte) (int, error) {
	return fmt.Print(data)
}
