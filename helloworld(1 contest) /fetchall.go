package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	results := make(chan string)
	for _, url := range os.Args[1:] {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			startTime := time.Now()
			resp, err := http.Get(url)
			if err != nil {
				results <- fmt.Sprintf("Error fetching %s: %v\n", url, err)
				return
			}
			defer resp.Body.Close()
			_, err = io.ReadAll(resp.Body)
			if err != nil {
				results <- fmt.Sprintf("Error reading response body for %s: %v\n", url, err)
				return
			}
			elapsed := time.Since(startTime).Seconds()
			results <- fmt.Sprintf("%.2fs elapsed %s\n", elapsed, url)
		}(url)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Print(r)
	}

	fmt.Printf("%.2fs total elapsed\n", time.Since(start).Seconds())
}
