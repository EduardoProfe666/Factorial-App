package database

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"factorial/internal/utils"
	"os"
	"strconv"
	"time"

	_ "github.com/glebarez/sqlite"
)

var DB *sql.DB

func InitDB() {
	if _, err := os.Stat("./database.db"); os.IsNotExist(err) {
		file, err := os.Create("./database.db")
		if err != nil {
			utils.LogError("Error creating database file: " + err.Error())
			return
		}
		file.Close()
	}

	var err error
	DB, err = sql.Open("sqlite", "file:./database.db?cache=shared&mode=rwc&_busy_timeout=5000")
	if err != nil {
		utils.LogError("Error opening database: " + err.Error())
		return
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS factorial (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		number INTEGER,
		result TEXT
	);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		utils.LogError("Error creating table: " + err.Error())
		return
	}

	utils.LogInfo("Database initialized successfully")
}

func SaveResult(number int, result string) error {
	maxRetries := 3
	retryDelay := 500 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := saveResultWithRetry(number, result)
		if err == nil {
			return nil
		}

		if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone) || errors.Is(err, sql.ErrNoRows) {
			utils.LogError("Operation failed, retrying...")
			time.Sleep(retryDelay)
			continue
		}

		return err
	}

	return errors.New("failed to save result after multiple retries")
}

func saveResultWithRetry(number int, result string) error {
	stmt, err := DB.Prepare("INSERT INTO factorial(number, result) VALUES(?, ?)")
	if err != nil {
		utils.LogError("Error preparing statement: " + err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(number, result)
	if err != nil {
		utils.LogError("Error executing statement: " + err.Error())
		return err
	}

	utils.LogInfo("Result saved to database for number: " + strconv.Itoa(number))
	return nil
}

type FactorialResult struct {
	Number int
	Result string
}

func GetResults() ([]FactorialResult, error) {
	rows, err := DB.Query("SELECT number, result FROM factorial")
	if err != nil {
		utils.LogError("Error querying database: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var results []FactorialResult
	for rows.Next() {
		var number int
		var result string
		err = rows.Scan(&number, &result)
		if err != nil {
			utils.LogError("Error scanning row: " + err.Error())
			return nil, err
		}
		results = append(results, FactorialResult{Number: number, Result: result})
	}

	return results, nil
}

func GetFactorial(number int) (string, error) {
	var result string
	err := DB.QueryRow("SELECT result FROM factorial WHERE number = ?", number).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.LogWarning("Factorial not found for number: " + strconv.Itoa(number))
			return "", errors.New("factorial not found")
		}
		utils.LogError("Error querying database: " + err.Error())
		return "", err
	}

	utils.LogInfo("Factorial retrieved from database for number: " + strconv.Itoa(number))
	return result, nil
}

func ClearResults() error {
	_, err := DB.Exec("DELETE FROM factorial")
	if err != nil {
		utils.LogError("Error clearing results: " + err.Error())
		return err
	}

	utils.LogInfo("Database results cleared successfully")
	return nil
}

func ExportToCSV(filename string) error {
	results, err := GetResults()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Number", "Result"})

	for _, result := range results {
		writer.Write([]string{strconv.Itoa(result.Number), result.Result})
	}

	return nil
}
