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

	"github.com/groveriffic/go-prevalent-colors/color"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: main urls.txt")
	}

	filename := os.Args[1]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	out := csv.NewWriter(os.Stdout)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		/* OPTIMIZE: This loop is probably IO bound on the first attempt */
		url := scanner.Text()
		log.Println(url)

		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		cc := color.ColorCounter{}
		cc.Image(img)
		record := []string{url}
		for _, rgb := range cc.TopThree() {
			record = append(record, rgb.String())
		}

		err = out.Write(record)
		if err != nil {
			log.Fatal(err)
		}
		out.Flush()
		if out.Error() != nil {
			log.Fatal(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Failed to read file: ", filename, err)
	}
}
