package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/groveriffic/go-prevalent-colors/color"
)

var cpuprofile = flag.String("cpu", "", "write cpu profile to file")
var workers = flag.Int("n", 5, "number of concurrent workers to run")
var help = flag.Bool("help", false, "display help message")
var input = flag.String("input", "", "(required) input file with image urls")
var outputFile = flag.String("output", "", "output file for CSV")
var logFile = flag.String("log", "", "log to file")

func main() {
	flag.Parse()

	if *input == "" || *help {
		fmt.Fprintln(os.Stderr, "Usage: ./go-prevalent-colors -in input.txt")
		flag.PrintDefaults()
		return
	}

	out := os.Stdout
	if *outputFile != "" {
		outF, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal("failed to create output file", err)
		}
		defer outF.Close()
		out = outF
	}

	if *logFile != "" {
		logF, err := os.Create(*logFile)
		if err != nil {
			log.Fatal("failed to create log file", err)
		}
		log.SetOutput(logF)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	log.Println("Input file:", *input)
	log.Println("Workers:", *workers)

	lines := generateLines(*input)

	records := processURLs(lines, *workers)

	writeCSV(records, out)
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
		defer f.Close()

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

func writeCSV(records chan []string, out io.Writer) {
	w := csv.NewWriter(out)
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
