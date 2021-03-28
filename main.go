package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	MAX_CONCURRENT_WORKERS = 10
	SLUG_SYMBOLS = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func generateSlug(length int) string {
	rand.Seed(time.Now().UnixNano())
	slug := ""
	for i := 0; i < length; i++ {
		slug += string(SLUG_SYMBOLS[rand.Intn(len(SLUG_SYMBOLS))])
	}
	return slug
}

func testUrl(currentRequests *int) {
	*currentRequests++
	url := "https://termbin.com/" + generateSlug(4)
	resp, _ := http.Get(url)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("Valid URL | %s\n", url)
	}
	*currentRequests--
}

func main() {
	println("ReTermbin v0.1")
	currentWorkers := 0

	for {
		if currentWorkers < MAX_CONCURRENT_WORKERS {
			go testUrl(&currentWorkers)
		}
	}
}
