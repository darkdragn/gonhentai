package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var limit int

func main() {
	var nnn = flag.Int("n", 1234, "The special sauce")
	flag.IntVar(&limit, "limit", 5, "Number of gofuncs to pull images with")
	flag.Parse()

	hold := NewDoujin(*nnn)
	fmt.Printf("%s\n", hold.Titles.English)

	downloadZip(hold)
}

func downloadZip(hold Doujin) {
	var wg sync.WaitGroup
	limit := make(chan struct{}, 5)
	images := hold.generateImages()

	filename := fmt.Sprintf("%s.cbz", hold.Titles.English)
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	bufChan := make(chan zipImage)
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	go func() {

		for run := range bufChan {
			fh := &zip.FileHeader{
				Name:     run.img.Filename,
				Modified: time.Now(),
				Method:   0,
			}
			f, err := zipFile.CreateHeader(fh)
			if err != nil {
				log.Println("At the header")
				log.Fatal(err)
			}
			_, err = io.Copy(f, &run.buf)
			if err != nil {
				log.Fatal(err)
			}
			run.wg.Done()
		}
	}()
	for _, img := range images {
		limit <- struct{}{}
		wg.Add(1)
		go func(img APIImage) { saveZip(img, &wg, bufChan); <-limit }(img)
	}
	wg.Wait()
	close(bufChan)
}

func saveZip(image APIImage, wg *sync.WaitGroup, bufChan chan zipImage) {
	buf := new(bytes.Buffer)
	resp, err := http.Get(image.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		log.Println("At the io.Copy")
		log.Fatal(err)
	}
	bufChan <- zipImage{img: image, buf: *buf, wg: wg}
}

func saveImage(image APIImage, wg *sync.WaitGroup) {
	out, err := os.Create(image.Filename)
	defer out.Close()
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Get(image.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
}
