package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/inancgumus/screen"

	"factorial/database"
	. "factorial/logic"
	. "factorial/ui/colors"
	. "factorial/utils"
)

const logFileName = "Logs/factorial_results.log"
const resultsDir = "Results/"

func RunUI(reader *bufio.Reader) {
	banner := figure.NewFigure("Factorial App", "", true)
	banner.Print()

	err := EnsureDir(resultsDir)
	if err != nil {
		fmt.Println(Red("Error creating Results directory: " + err.Error()))
		pressAnyKeyToContinue(reader)
		return
	}

	err = EnsureDir("Logs/")
	if err != nil {
		fmt.Println(Red("Error creating Logs directory: " + err.Error()))
		pressAnyKeyToContinue(reader)
		return
	}

	for {
		clearScreen()
		fmt.Println(Yellow("Choose an option: 🚀"))
		fmt.Println("1. " + Cyan("Calculate factorial of a single number 🔢"))
		fmt.Println("2. " + Cyan("Calculate factorials of multiple numbers 🔢🔢"))
		fmt.Println("3. " + Cyan("Calculate factorials in a range 📊"))
		fmt.Println("4. " + Cyan("View stored results 📂"))
		fmt.Println("5. " + Cyan("Clear result files 🗑️"))
		fmt.Println("6. " + Cyan("Clear database results 🗑️"))
		fmt.Println("7. " + Cyan("Export results to CSV 📄"))
		fmt.Println("8. " + Cyan("Exit 🚪"))
		fmt.Print(Green("Enter option: "))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		option, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(Red("Error: Invalid input. Please enter a number. 🚫"))
			pressAnyKeyToContinue(reader)
			continue
		}

		switch option {
		case 1:
			calculateSingleFactorial(reader)
		case 2:
			calculateMultipleFactorials(reader)
		case 3:
			calculateRangeFactorials(reader)
		case 4:
			viewStoredResults(reader)
		case 5:
			clearResultFiles(reader)
		case 6:
			clearDatabaseResults(reader)
		case 7:
			exportResultsToCSV(reader)
		case 8:
			fmt.Println(Yellow("Exiting... 🚪"))
			return
		default:
			fmt.Println(Red("Invalid option. Please try again. 🚫"))
			pressAnyKeyToContinue(reader)
		}
	}
}

func calculateSingleFactorial(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Green("Enter a number to calculate its factorial: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	number, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println(Red("Error: Invalid input. Please enter an integer. 🚫"))
		pressAnyKeyToContinue(reader)
		return
	}

	result, err := database.GetFactorial(number)
	if err == nil {
		fmt.Printf("Factorial of %d is %s (loaded from database)\n", number, result)
	} else {
		loadingAnimation()
		result := Factorial(number)
		fmt.Printf("Factorial of %d is %s\n", number, result.String())

		err = SaveResultToFile(result, number, resultsDir)
		if err != nil {
			fmt.Println(Red("Error saving result to file: " + err.Error()))
		} else {
			fmt.Println(Green("Result saved to " + resultsDir + "Factorial of " + strconv.Itoa(number) + ".txt"))
		}

		err = LogResult(number, result, logFileName)
		if err != nil {
			fmt.Println(Red("Error logging result: " + err.Error()))
		}

		err = database.SaveResult(number, result.String())
		if err != nil {
			fmt.Println(Red("Error saving result to database: " + err.Error()))
		}
	}
	pressAnyKeyToContinue(reader)
}

func calculateMultipleFactorials(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Green("Enter numbers separated by commas to calculate their factorials: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	numbers := strings.Split(input, ",")

	var wg sync.WaitGroup
	for _, numStr := range numbers {
		num, err := strconv.Atoi(strings.TrimSpace(numStr))
		if err != nil {
			fmt.Printf(Red("Error: Invalid input '%s'. Please enter integers. 🚫\n"), numStr)
			continue
		}
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			result, err := database.GetFactorial(num)
			if err == nil {
				fmt.Printf("Factorial of %d is %s (loaded from database)\n", num, result)
			} else {
				loadingAnimation()
				result := Factorial(num)
				fmt.Printf("Factorial of %d is %s\n", num, result.String())

				err := SaveResultToFile(result, num, resultsDir)
				if err != nil {
					fmt.Println(Red("Error saving result to file: " + err.Error()))
				} else {
					fmt.Println(Green("Result saved to " + resultsDir + "Factorial of " + strconv.Itoa(num) + ".txt"))
				}

				err = LogResult(num, result, logFileName)
				if err != nil {
					fmt.Println(Red("Error logging result: " + err.Error()))
				}

				err = database.SaveResult(num, result.String())
				if err != nil {
					fmt.Println(Red("Error saving result to database: " + err.Error()))
				}
			}
		}(num)
	}
	wg.Wait()
	pressAnyKeyToContinue(reader)
}

func calculateRangeFactorials(reader *bufio.Reader) {
	clearScreen()
	fmt.Print(Green("Enter the start and end of the range separated by a comma: "))
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	rangeParts := strings.Split(input, ",")
	if len(rangeParts) != 2 {
		fmt.Println(Red("Error: Invalid input. Please enter two integers separated by a comma. 🚫"))
		pressAnyKeyToContinue(reader)
		return
	}
	start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
	if err != nil {
		fmt.Println(Red("Error: Invalid start of range. Please enter an integer. 🚫"))
		pressAnyKeyToContinue(reader)
		return
	}
	end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
	if err != nil {
		fmt.Println(Red("Error: Invalid end of range. Please enter an integer. 🚫"))
		pressAnyKeyToContinue(reader)
		return
	}

	var wg sync.WaitGroup
	for i := start; i <= end; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			result, err := database.GetFactorial(num)
			if err == nil {
				fmt.Printf("Factorial of %d is %s (loaded from database)\n", num, result)
			} else {
				loadingAnimation()
				result := Factorial(num)
				fmt.Printf("Factorial of %d is %s\n", num, result.String())

				err := SaveResultToFile(result, num, resultsDir)
				if err != nil {
					fmt.Println(Red("Error saving result to file: " + err.Error()))
				} else {
					fmt.Println(Green("Result saved to " + resultsDir + "Factorial of " + strconv.Itoa(num) + ".txt"))
				}

				err = LogResult(num, result, logFileName)
				if err != nil {
					fmt.Println(Red("Error logging result: " + err.Error()))
				}

				err = database.SaveResult(num, result.String())
				if err != nil {
					fmt.Println(Red("Error saving result to database: " + err.Error()))
				}
			}
		}(i)
	}
	wg.Wait()
	pressAnyKeyToContinue(reader)
}

func viewStoredResults(reader *bufio.Reader) {
	clearScreen()
	loadingAnimation()
	results, err := database.GetResults()
	if err != nil {
		fmt.Println(Red("Error retrieving results from database: " + err.Error()))
		pressAnyKeyToContinue(reader)
		return
	}
	fmt.Println(Yellow("Stored results: 📂"))
	for _, result := range results {
		fmt.Println(result)
	}
	pressAnyKeyToContinue(reader)
}

func clearResultFiles(reader *bufio.Reader) {
	clearScreen()
	loadingAnimation()
	files, err := os.ReadDir(resultsDir)
	if err != nil {
		fmt.Println(Red("Error reading directory: " + err.Error()))
		pressAnyKeyToContinue(reader)
		return
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "Factorial of ") && strings.HasSuffix(file.Name(), ".txt") {
			err := os.Remove(resultsDir + file.Name())
			if err != nil {
				fmt.Println(Red("Error deleting file: " + err.Error()))
			} else {
				fmt.Println(Green("Deleted file: " + file.Name()))
			}
		}
	}
	err = os.Remove(logFileName)
	if err != nil {
		fmt.Println(Red("Error deleting log file: " + err.Error()))
	} else {
		fmt.Println(Green("Deleted log file: " + logFileName))
	}
	pressAnyKeyToContinue(reader)
}

func exportResultsToCSV(reader *bufio.Reader) {
	clearScreen()
	loadingAnimation()
	results, err := database.GetResults()
	if err != nil {
		fmt.Println(Red("Error retrieving results from database: " + err.Error()))
		pressAnyKeyToContinue(reader)
		return
	}

	csvFile, err := os.Create(resultsDir + "results.csv")
	if err != nil {
		fmt.Println(Red("Error creating CSV file: " + err.Error()))
		pressAnyKeyToContinue(reader)
		return
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {

		}
	}(csvFile)

	csvFile.WriteString("Number,Factorial\n")
	for _, result := range results {
		parts := strings.Split(result, " is ")
		number := strings.TrimPrefix(parts[0], "Factorial of ")
		factorial := parts[1]
		csvFile.WriteString(fmt.Sprintf("%s,%s\n", number, factorial))
	}

	fmt.Println(Green("Results exported to " + resultsDir + "results.csv 📄"))
	pressAnyKeyToContinue(reader)
}

func clearDatabaseResults(reader *bufio.Reader) {
	clearScreen()
	loadingAnimation()
	err := database.ClearResults()
	if err != nil {
		fmt.Println(Red("Error clearing results from database: " + err.Error()))
	} else {
		fmt.Println(Green("Database results cleared successfully"))
	}
	pressAnyKeyToContinue(reader)
}

func loadingAnimation() {
	dots := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	for i := 0; i < 10; i++ {
		fmt.Printf("\rLoading %s", dots[i%len(dots)])
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("\rDone!          ")
}

func clearScreen() {
	screen.Clear()
	screen.MoveTopLeft()
}

func pressAnyKeyToContinue(reader *bufio.Reader) {
	fmt.Println(Yellow("Press enter to continue..."))
	reader.ReadString('\n')
}
