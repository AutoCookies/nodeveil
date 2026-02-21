package main

import (
	"encoding/json"
	"os"
	"time"
)

type BenchResult struct {
	Bench     string `json:"bench"`
	Timestamp string `json:"timestamp"`
	OK        bool   `json:"ok"`
}

func main() {
	_ = json.NewEncoder(os.Stdout).Encode(BenchResult{
		Bench:     "phase0_smoke",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		OK:        true,
	})
}
