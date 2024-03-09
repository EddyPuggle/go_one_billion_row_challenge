package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const WEATHER_STATION_PATH = "./weather_stations.csv"

// runtime.GOMAXPROCS(0) for maximum power
var NUM_WORKERS = 2000

func constructResultRowString(min float64, mean float64, max float64) string {
	return fmt.Sprintf("%f/%f/%f", round(min), round(mean), round(max))
}

func processLineWorker(rows chan string, wg *sync.WaitGroup, resultMap *SafeMap) {
	line := <-rows
	lineSplit := strings.Split(line, ";")
	key := lineSplit[0]
	value, err := strconv.ParseFloat(lineSplit[1], 64)
	assertError(err)

	resultMap.Set(key, value)
	wg.Done()
}

func readLinesWorker(rows chan string, scanner *bufio.Scanner) {
	defer close(rows)
	for scanner.Scan() {
		rows <- scanner.Text()
	}
	assertError(scanner.Err())
}

func main() {
	startTime := time.Now()
	resultMap := SafeMap{val: make(map[string]float64)}
	// open file
	file, err := os.Open(WEATHER_STATION_PATH)
	assertError(err)
	defer file.Close()
	// waitgroup to wait for
	var wg sync.WaitGroup
	// channel to transport lines to the different workers
	rowsChan := make(chan string)
	// prepare workers to process incoming lines; 100 for now; use runtime.GOMAXPROCS(0) for max speed
	wg.Add(NUM_WORKERS)
	for range NUM_WORKERS {
		go processLineWorker(rowsChan, &wg, &resultMap)
	}
	// prepare worker to read and send lines to "rowsChan"
	scanner := bufio.NewScanner(file)
	go readLinesWorker(rowsChan, scanner)
	// wait for completion
	wg.Wait()
	// print map
	for key, val := range resultMap.val {
		fmt.Printf("%s: %f\n", key, val)
	}
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime).Seconds())
}
