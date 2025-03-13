package utils

import (
	"github.com/gin-gonic/gin"
	"log"
)

// Debugf logs a debug message only if the server is in debug mode
func Debugf(format string, args ...interface{}) {
	if gin.Mode() == gin.DebugMode {
		log.Printf(format, args...)
	}
}

// Infof logs an info message
func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}
