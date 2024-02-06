package logger

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(prefix string) *log.Logger {
	logFileName := fmt.Sprintf("%s.log", prefix)
	// file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("Failed to open log file")
		log.Fatal(err)
	}
	logger := log.New(file, fmt.Sprintf("%s: ", prefix), log.LstdFlags|log.Lshortfile)
	logger.Println("logger init")
	return logger
}
