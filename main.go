package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	res, err := http.Get("http://nhentai.net/api/gallery/18511")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	hold := Doujin{}
	jsonErr := json.Unmarshal(body, &hold)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	fmt.Printf("%s", hold.Titles.English)

	var wg sync.WaitGroup
	limit := make(chan struct{}, 1)
	images := hold.generateImages()
	// for _, img := range images {
	// 	// fmt.Printf("%v\n", img.URL())
	// 	limit <- struct{}{}
	// 	wg.Add(1)
	// 	go func(img APIImage) { saveImage(img, &wg); <-limit }(img)
	// }

	file, err := os.Create("temp.cbz")
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	bufChan := make(chan bytes.Buffer)
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()
	for _, img := range images {
		limit <- struct{}{}
		wg.Add(1)
		fh := &zip.FileHeader{
			Name:     img.Filename,
			Modified: time.Now(),
			Method:   0,
		}
		f, err := zipFile.CreateHeader(fh)
		if err != nil {
			log.Println("At the header")
			log.Fatal(err)
		}
		go func(img APIImage) { saveZip(img, &wg, &f); <-limit }(img)
	}
	wg.Wait()
	// zipFile.Close()
}

func saveZip(image APIImage, wg *sync.WaitGroup, out *io.Writer) {
	resp, err := http.Get(image.URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(*out, resp.Body)
	if err != nil {
		log.Println("At the io.Copy")
		log.Fatal(err)
	}
	wg.Done()
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
