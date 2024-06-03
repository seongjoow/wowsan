package logger

import (
	"fmt"
	"log"
	"os"
	"wowsan/pkg/model"

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

// func NewLoggerConsole() *logrus.Logger {
// 	logger := logrus.New()
// 	logger.SetFormatter(&logrus.JSONFormatter{
// 		TimestampFormat:   "2006-01-02 15:04:05",
// 		DisableHTMLEscape: true,
// 	})
// 	logger.SetLevel(logrus.InfoLevel)
// 	logger.SetOutput(os.Stdout)
// 	return logger
// }

func NewLogger(prefix, logDir string) (*logrus.Logger, error) {
	logFileName := fmt.Sprintf("%s/%s.json", logDir, prefix)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableHTMLEscape: true,
	})

	// check directory is exist, if not create directory
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(file)

	return logger, nil
}

type BrokerInfoLogger struct {
	dir string
}

func NewBrokerInfoLogger(dir string) *BrokerInfoLogger {
	// check if log directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &BrokerInfoLogger{
		dir: dir,
	}
}

func (l *BrokerInfoLogger) GetBrokerInfo(broker *model.Broker) {
	// save to log file (broker_id.log)
	// if file doesn't exist, create it, otherwise append to the file
	logFile, err := os.OpenFile(l.dir+"/"+broker.Port+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	brokerInfoString := "==Broker Info==\n"
	brokerInfoString += fmt.Sprintf("Id: %s\n", broker.Id)
	brokerInfoString += fmt.Sprintf("Ip: %s\n", broker.Ip)
	brokerInfoString += fmt.Sprintf("Port: %s\n", broker.Port)

	publishersString := "==Publishers==\n"
	for _, publisher := range broker.Publishers {
		publishersString += fmt.Sprintf("Id: %s\n", publisher.Id)
		publishersString += fmt.Sprintf("Ip: %s\n", publisher.Ip)
		publishersString += fmt.Sprintf("Port: %s\n", publisher.Port)
	}

	subscribersString := "==Subscribers==\n"
	for _, subscriber := range broker.Subscribers {
		subscribersString += fmt.Sprintf("Id: %s\n", subscriber.Id)
		subscribersString += fmt.Sprintf("Ip: %s\n", subscriber.Ip)
		subscribersString += fmt.Sprintf("Port: %s\n", subscriber.Port)
	}

	brokerString := "==Brokers==\n"
	for _, neighbor := range broker.Brokers {
		brokerString += fmt.Sprintf("Id: %s\n", neighbor.Id)
		brokerString += fmt.Sprintf("Ip: %s\n", neighbor.Ip)
		brokerString += fmt.Sprintf("Port: %s\n", neighbor.Port)
	}

	itemsSrt := "==SRT==\n"
	for _, item := range broker.SRT.Items {
		itemsSrt += fmt.Sprintf("Adv: %s %s %s (%s) | %s\n", item.Advertisement.Subject, item.Advertisement.Operator, item.Advertisement.Value, item.Identifier.MessageId, item.Identifier.SenderId)
		for i := 0; i < len(item.LastHop); i++ {
			itemsSrt += fmt.Sprintf("%s | %s | %d \n", item.LastHop[i].Id, item.LastHop[i].NodeType, item.HopCount)
		}
		itemsSrt += fmt.Sprintf("----------------------------\n")
	}

	itemsPrt := "==PRT==\n"
	for _, item := range broker.PRT.Items {
		itemsPrt += fmt.Sprintf("Sub: %s %s %s (%s) | %s\n", item.Subscription.Subject, item.Subscription.Operator, item.Subscription.Value, item.Identifier.MessageId, item.Identifier.SenderId)
		for i := 0; i < len(item.LastHop); i++ {
			itemsPrt += fmt.Sprintf("%s | %s\n", item.LastHop[i].Id, item.LastHop[i].NodeType)
		}
		itemsPrt += fmt.Sprintf("----------------------------\n")
	}

	log.Printf(brokerInfoString + publishersString + subscribersString + brokerString + itemsSrt + itemsPrt)
}
