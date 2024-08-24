package api

import (
	"encoding/csv"
	"encoding/json"
	"factorial/internal/database"
	"factorial/internal/logic"
	"factorial/internal/utils"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func StartAPI() {
	http.HandleFunc("/results", getResults)
	http.HandleFunc("/delete", deleteResults)
	http.HandleFunc("/export", exportCSV)
	http.HandleFunc("/factorial", getFactorial)
	http.HandleFunc("/range", calculateRange)

	utils.LogInfo("Starting API server on localhost:2999")
	http.ListenAndServe(":2999", nil)
}

func getResults(w http.ResponseWriter, r *http.Request) {
	results, err := database.GetResults()
	if err != nil {
		utils.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonData, err := json.Marshal(results)
	if err != nil {
		utils.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func deleteResults(w http.ResponseWriter, r *http.Request) {
	err := database.ClearResults()
	if err != nil {
		utils.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Results deleted successfully"))
}

func exportCSV(w http.ResponseWriter, r *http.Request) {
	results, err := database.GetResults()
	if err != nil {
		utils.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Create("results.csv")
	if err != nil {
		utils.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, result := range results {
		record := []string{strconv.Itoa(result.Number), result.Result}
		writer.Write(record)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Results exported to results.csv"))
}

func getFactorial(w http.ResponseWriter, r *http.Request) {
	numberStr := r.URL.Query().Get("number")
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		utils.LogWarning("Invalid number")
		http.Error(w, "Invalid number", http.StatusBadRequest)
		return
	}

	result, err := database.GetFactorial(number)
	if err != nil {
		utils.LogError(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func calculateRange(w http.ResponseWriter, r *http.Request) {
	lowerStr := r.URL.Query().Get("lower")
	upperStr := r.URL.Query().Get("upper")
	lower, err1 := strconv.Atoi(lowerStr)
	upper, err2 := strconv.Atoi(upperStr)
	if err1 != nil || err2 != nil || lower > upper {
		utils.LogWarning("Invalid range")
		http.Error(w, "Invalid range", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	for i := lower; i <= upper; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			_, err := database.GetFactorial(num)
			if err != nil {
				result := logic.Factorial(num)
				err = database.SaveResult(num, result.String())
				if err != nil {
					utils.LogError(err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}(i)
	}
	wg.Wait()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Range factorial calculation completed"))
}
