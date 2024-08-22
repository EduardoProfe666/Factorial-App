package utils

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
)

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveResultToFile(result *big.Int, number int, dir string) error {
	err := EnsureDir(dir)
	if err != nil {
		return err
	}

	fileName := dir + "Factorial of " + strconv.Itoa(number) + ".txt"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.WriteString(result.String())
	if err != nil {
		return err
	}

	return nil
}

func LogResult(number int, result *big.Int, logFileName string) error {
	err := EnsureDir("Logs/")
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {

		}
	}(logFile)

	logEntry := fmt.Sprintf("Factorial of %d is %s\n", number, result.String())
	_, err = logFile.WriteString(logEntry)
	if err != nil {
		return err
	}

	return nil
}
