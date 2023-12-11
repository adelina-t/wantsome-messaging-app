package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"wantsome.ro/messagingapp/pkg/models"
)

type Logger struct {
	File            *os.File
	LoggerBroadcast chan models.LoggerMessage
	LoggerMutex     *sync.Mutex
	MaxLogLevel     models.LogLevel
	Method          string
}

func InitLogger() Logger {
	currentTime := time.Now()
	fileName := fmt.Sprintf("log_%d-%d-%d_%d_%d_%d.txt",
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second())

	log.Printf("%s \n", fileName)

	f, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating the file %s", err)
		return Logger{}
	}
	logger := Logger{}
	logger.LoggerMutex = &sync.Mutex{}
	logger.LoggerMutex.Lock()
	logger.File = f
	logger.LoggerBroadcast = make(chan models.LoggerMessage)
	logger.MaxLogLevel = models.LogLevel(2)
	logger.LoggerMutex.Unlock()
	log.Printf("Log created successfully")
	return logger
}

func (logger *Logger) DeInitLogger() {
	err := logger.File.Close()
	if err != nil {
		log.Printf("Error closing the file %s", err)
		return
	}
}

func (logger Logger) Log(msg string, method string, level models.LogLevel) {
	if level <= logger.MaxLogLevel {
		logger.LoggerBroadcast <- models.LoggerMessage{
			Message: msg,
			Method:  logger.Method,
			Level:   level,
		}
	}
}

func (logger Logger) WriteToFileService() {
	defer logger.File.Close()

	for {
		loggerMessage := <-logger.LoggerBroadcast
		logger.LoggerMutex.Lock()

		currentTime := time.Now()
		logToWrite := fmt.Sprintf("%d-%d-%d_%d:%d:%d",
			currentTime.Day(),
			currentTime.Month(),
			currentTime.Year(),
			currentTime.Hour(),
			currentTime.Minute(),
			currentTime.Second())

		logToWrite += " Level=" + strconv.Itoa(int(loggerMessage.Level))
		if loggerMessage.Message != "" {
			logToWrite += " Msg=" + loggerMessage.Message
		}

		if loggerMessage.Method != "" {
			logToWrite += "Method= " + loggerMessage.Method
		}

		logToWrite += "\n"
		_, err := logger.File.WriteString(logToWrite)
		if err != nil {
			log.Printf("Could not print the message: %s", err)
		}
		logger.LoggerMutex.Unlock()
	}
}
