package main

import (
	"bufio"
	"factorial/database"
	"factorial/ui"
	"os"
)

func main() {
	database.InitDB()

	reader := bufio.NewReader(os.Stdin)
	ui.RunUI(reader)
}
