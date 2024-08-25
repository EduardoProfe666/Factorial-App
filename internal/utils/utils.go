package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

const logFileName = "logs.txt"

var logger *log.Logger

func init() {
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
	logger.SetFlags(0)
}

func logMessage(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logger.Printf("%s : %s : %s", timestamp, level, message)
	fmt.Printf("%s : %s : %s\n", timestamp, level, message)
}

func LogInfo(message string) {
	logMessage("ℹ️INFOℹ️", message)
}

func LogWarning(message string) {
	logMessage("⚠️WARNING⚠️", message)
}

func LogError(message string) {
	logMessage("❌ERROR❌", message)
}

func LogResult(number int, fromDatabase bool) {
	if fromDatabase {
		LogInfo(fmt.Sprintf("Factorial of %d retrieved from database", number))
	} else {
		LogInfo(fmt.Sprintf("Factorial of %d calculated", number))
	}
}
