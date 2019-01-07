package main

import (
	"bufio"
	"encoding/csv"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/groveriffic/go-prevalent-colors/color"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: main urls.txt")
	}

	filename := os.Args[1]
	log.Println("Input file:", filename)

	lines := generateLines(filename)

	workers := 2
	records := processURLs(lines, workers)

	writeCSV(records)
	log.Println("Done")
}

func generateLines(filename string) (lines chan string) {
	lines = make(chan string, 1)

	go func() {
		defer close(lines)

		f, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			return
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		err = scanner.Err()
		if err != nil {
			log.Println(err)
		}
	}()
	return
}

func processURL(url string) (record []string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return
	}

	cc := color.Counter{}
	cc.Image(img)
	record = []string{url}
	for _, rgb := range cc.TopThree() {
		record = append(record, rgb.String())
	}
	return
}

func writeCSV(records chan []string) {
	w := csv.NewWriter(os.Stdout)
	for record := range records {
		err := w.Write(record)
		if err != nil {
			log.Fatal(err)
		}
		w.Flush()
		err = w.Error()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func processURLs(urls chan string, workers int) chan []string {
	records := make(chan []string)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for url := range urls {
				log.Println("Processing:", url)
				record, err := processURL(url)
				if err != nil {
					log.Println(err, url)
				}
				records <- record
			}
		}()
	}

	go func() {
		wg.Wait()
		close(records)
	}()

	return records
}
