package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TempStats struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <filename>", os.Args[0])
	}

	startTime := time.Now()

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	stats := make(map[string]*TempStats)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		if len(parts) != 2 {
			continue
		}

		city := parts[0]
		temp, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}

		if _, exists := stats[city]; !exists {
			stats[city] = &TempStats{Min: temp, Max: temp, Sum: temp, Count: 1}
		} else {
			if temp < stats[city].Min {
				stats[city].Min = temp
			}
			if temp > stats[city].Max {
				stats[city].Max = temp
			}
			stats[city].Sum += temp
			stats[city].Count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	var cities []string
	for city := range stats {
		cities = append(cities, city)
	}
	sort.Strings(cities)

	result := "{"
	first := true
	for _, city := range cities {
		stat := stats[city]
		mean := stat.Sum / float64(stat.Count)
		if !first {
			result += ", "
		}
		result += fmt.Sprintf("%s=%.1f/%.1f/%.1f", city, stat.Min, mean, stat.Max)
		first = false
	}
	result += "}"

	fmt.Println(result)
	duration := time.Since(startTime)
	fmt.Println("Tempo de execução:", duration)
}
