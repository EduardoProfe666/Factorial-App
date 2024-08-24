package logic

import (
	"math/big"
	"sync"
)

func Factorial(number int) *big.Int {
	if number <= 1 {
		return big.NewInt(1)
	}

	numGoroutines := 4
	if number < numGoroutines {
		numGoroutines = number
	}

	results := make(chan *big.Int, numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	chunkSize := number / numGoroutines
	for i := 0; i < numGoroutines; i++ {
		start := i*chunkSize + 1
		end := (i + 1) * chunkSize
		if i == numGoroutines-1 {
			end = number
		}
		go func(start, end int) {
			defer wg.Done()
			partial := big.NewInt(1)
			for j := start; j <= end; j++ {
				partial.Mul(partial, big.NewInt(int64(j)))
			}
			results <- partial
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	result := big.NewInt(1)
	for partial := range results {
		result.Mul(result, partial)
	}

	return result
}
