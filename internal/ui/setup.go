package ui

import (
	"factorial/internal/database"
	"factorial/internal/utils"
)

func initDatabase() {
	database.InitDB()
}

func initLogging() {
	utils.LogInfo("Factorial App Started")
}

func closeApp() {
	utils.LogInfo("Factorial App Exited")
}
