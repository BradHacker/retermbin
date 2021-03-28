package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Stats struct {
	currentWorkers int
	urls []string
	successes int
	failures int
	startTime time.Time
	endTime time.Time
} 

const (
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

func testUrl(stats *Stats) {
	(*stats).currentWorkers++
	url := "https://termbin.com/" + generateSlug(4)
	fmt.Printf("  | %s\r", url)
	resp, _ := http.Get(url)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		(*stats).successes++
		(*stats).urls = append((*stats).urls, url)
	} else {
		(*stats).failures++
	}
	(*stats).currentWorkers--
}

func writeOutput(urls []string, outFile string) {
	println("Saving url's to " + outFile)
	data := []byte(strings.Join(urls, "\n"))
	err := ioutil.WriteFile(outFile, data, 0644)
	if err != nil {
		println(fmt.Errorf("error writing to output: %v", err))
		os.Exit(1)
	}
}

func printStats(stats Stats) {
	println("Stats")
	println("-----")
	fmt.Printf("Valid Urls:\t%d\n", stats.successes)
	fmt.Printf("Invalid Urls:\t%d\n", stats.failures)
	println("-----")
	fmt.Printf("Total Time:\t\t%s\n", time.Time.Sub(stats.endTime, stats.startTime).String())
	fmt.Printf("Avg Request Time:\t%0.2fs\n", time.Time.Sub(stats.endTime, stats.startTime).Seconds() / (float64(stats.successes + stats.failures)))
}

func main() {
	if len(os.Args) < 2 {
		println("Usage: retermbin [max_workers] (output_file)")
		os.Exit(1)
	}

	MAX_CONCURRENT_WORKERS, err := strconv.Atoi(os.Args[1])
	if err != nil {
		println("Max workers needs to be an integer")
		os.Exit(1)
	}
	OUTPUT_FILE := "output/retermbin_out_" + fmt.Sprint(time.Now().UnixNano()) + ".txt"
	if len(os.Args) > 2 {
		OUTPUT_FILE = os.Args[2]
	}

	println("ReTermbin v0.1")
	println("^C to stop testing")
	stats := Stats{
		currentWorkers: 0,
		successes: 0,
		failures: 0,
		startTime: time.Now(),
		urls: []string{},
	}

	sigs := make(chan os.Signal, 1)
	done := false

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		stats.endTime = time.Now()
		writeOutput(stats.urls, OUTPUT_FILE)
		printStats(stats)
		done = true
	}()

	for !done {
		if stats.currentWorkers < MAX_CONCURRENT_WORKERS {
			go testUrl(&stats)
		}
	}
}
