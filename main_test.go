package main

import (
	"log"
	"testing"

	"./db"
)

func TestCalculatePPLNSPayout(t *testing.T) {
	block := db.Block{
		Hash:       "000000000000000002003f919f871052b52c0c2e0cef8fb7c2b7fef28a4ff3fb",
		Difficulty: 522724265323.4008,
		Reward:     12.50931102,
		Height:     547267,
		Coin:       "BCH",
	}

	err := db.CalculatePPLNSPayout(block)
	log.Print(err)
}

func TestReport(t *testing.T) {
	err := db.Report()
	log.Print(err)
}
