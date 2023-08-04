package main

import (
	"fmt"
	"math/rand"
	"time"
)

func chooseFile(files []string, history []string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no files")
	}

	fileset := make(map[string]struct{})
	for _, file := range files {
		fileset[file] = struct{}{}
	}

	// Number unique history items from oldest k to newest 1
	var uniqueHistory []string
	seen := make(map[string]struct{})
	for i := range history {
		entry := history[len(history)-1-i]
		if _, ok := fileset[entry]; !ok {
			continue
		}

		if _, ok := seen[entry]; !ok {
			uniqueHistory = append(uniqueHistory, entry)
			seen[entry] = struct{}{}
		}
	}

	historyIndex := make(map[string]int)
	for i, entry := range uniqueHistory {
		historyIndex[entry] = i
	}

	historyLen := len(uniqueHistory)

	// Generate weights
	weights := make([]float64, len(files))
	for i, entry := range files {
		if pos, ok := historyIndex[entry]; ok {
			weights[i] = float64(pos+1) / float64(historyLen+1)
		} else {
			weights[i] = 1.0
		}
	}

	if pos, ok := randIndexWeighted(weights); ok {
		return files[pos], nil
	}

	return "", fmt.Errorf("could not choose file")
}

func randIndexWeighted(weights []float64) (int, bool) {
	if len(weights) == 0 {
		return 0, false
	}

	s := 0.0
	for _, w := range weights {
		s += w
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tgt := rng.Float64() * s

	r := 0.0
	for i, w := range weights {
		r += w
		if r > tgt {
			return i, true
		}
	}

	return len(weights) - 1, true
}
