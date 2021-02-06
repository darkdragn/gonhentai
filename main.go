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

	. "github.com/darkdragn/gonhentai/v2/api"
	"github.com/schollz/progressbar/v3"
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

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func downloadZip(hold Doujin) {
	var wg sync.WaitGroup
	limit := make(chan struct{}, 5)
	images := hold.GenerateImages()

	bar := progressbar.Default(int64(len(images)))
	filename := fmt.Sprintf("%s.cbz", hold.Titles.English)
	file, err := os.Create(filename)
	defer file.Close()
	catch(err)

	bufChan := make(chan ZipImage)
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	go func() {

		for run := range bufChan {
			fh := &zip.FileHeader{
				Name:     run.Img.Filename,
				Modified: time.Now(),
				Method:   0,
			}
			f, err := zipFile.CreateHeader(fh)
			catch(err)
			_, err = io.Copy(f, &run.Buf)
			catch(err)
			bar.Add(1)
			run.Wg.Done()
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

func saveZip(image APIImage, wg *sync.WaitGroup, bufChan chan ZipImage) {
	buf := new(bytes.Buffer)
	resp, err := http.Get(image.URL)
	catch(err)
	defer resp.Body.Close()
	_, err = io.Copy(buf, resp.Body)
	catch(err)
	bufChan <- ZipImage{Img: image, Buf: *buf, Wg: wg}
}

func saveImage(image APIImage, wg *sync.WaitGroup) {
	out, err := os.Create(image.Filename)
	defer out.Close()
	catch(err)
	resp, err := http.Get(image.URL)
	catch(err)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	catch(err)
	wg.Done()
}
