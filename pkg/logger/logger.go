package logger

import (
	"fmt"
	"os"

	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

// func NewLogger(prefix string) *log.Logger {
// 	logFileName := fmt.Sprintf("%s.log", prefix)
// 	// file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
// 	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
// 	if err != nil {
// 		fmt.Println("Failed to open log file")
// 		log.Fatal(err)
// 	}
// 	logger := log.New(file, fmt.Sprintf("%s: ", prefix), log.LstdFlags|log.Lshortfile)
// 	logger.Println("logger init")
// 	return logger
// }

func NewLogger(prefix string) (*logrus.Logger, error) {
	logFileName := fmt.Sprintf("%s.json", prefix)
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(file)
	return logger, nil
}

func Utilization() (float64, uint64) {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, 0
	}
	cpuPercent, _ := p.CPUPercent()
	memInfo, _ := p.MemoryInfo()
	return cpuPercent, memInfo.RSS
}
