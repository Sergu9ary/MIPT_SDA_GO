package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var counts = make(map[string]int)
	for _, filename := range os.Args[1:] {
		file, _ := os.OpenFile(filename, os.O_RDONLY, 0666)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			counts[str]++
		}
	}
	for str, count := range counts {
		if count >= 2 {
			fmt.Printf("%d\t%s\n", count, str)
		}
	}
}
