package utils

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Debugf(format string, args ...interface{}) {
	if gin.Mode() == gin.DebugMode {
		log.Printf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}
