package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/glebarez/sqlite"
)

var DB *sql.DB

func InitDB() {
	// Check if the database file exists
	if _, err := os.Stat("./database.db"); os.IsNotExist(err) {
		// Create the database file if it does not exist
		file, err := os.Create("./database.db")
		if err != nil {
			log.Fatal("Error creating database file:", err)
		}
		file.Close()
	}

	var err error
	DB, err = sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatal("Error opening database:", err)
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
		log.Fatal("Error creating table:", err)
	}

	fmt.Println("Database initialized successfully")
}

func SaveResult(number int, result string) error {
	stmt, err := DB.Prepare("INSERT INTO factorial(number, result) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(number, result)
	if err != nil {
		return err
	}

	return nil
}

func GetResults() ([]string, error) {
	rows, err := DB.Query("SELECT number, result FROM factorial")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var number int
		var result string
		err = rows.Scan(&number, &result)
		if err != nil {
			return nil, err
		}
		results = append(results, fmt.Sprintf("Factorial of %d is %s", number, result))
	}

	return results, nil
}

func GetFactorial(number int) (string, error) {
	var result string
	err := DB.QueryRow("SELECT result FROM factorial WHERE number = ?", number).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("factorial not found")
		}
		return "", err
	}
	return result, nil
}

func ClearResults() error {
	_, err := DB.Exec("DELETE FROM factorial")
	if err != nil {
		return err
	}
	return nil
}
