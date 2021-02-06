package api

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

	"github.com/schollz/progressbar/v3"
)

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//DownloadZip is used to pull the Doujin to disk as {doujin.Title.English}.cbz
func (d Doujin) DownloadZip(limitRAW int) {
	var wg sync.WaitGroup
	limit := make(chan struct{}, limitRAW)
	images := d.APIImages

	bar := progressbar.Default(int64(len(images)))
	filename := fmt.Sprintf("%s.cbz", d.Titles.English)
	file, err := os.Create(filename)
	defer file.Close()
	catch(err)

	bufChan := make(chan zipImage)
	zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	completion := make(chan bool)
	go HandleZip(bufChan, zipFile, completion)
	go func() {
		for range completion {
			wg.Done()
			bar.Add(1)
		}
	}()
	for _, img := range images {
		limit <- struct{}{}
		wg.Add(1)
		go func(img Image) { img.ImageToZip(bufChan); <-limit }(img)
	}

	wg.Wait()
	close(bufChan)
}

func (d *Doujin) generateImage(i int, t imageType) Image {
	image := Image{Index: i, MediaID: d.MediaID, Type: t}
	image.Filename = image.filename()
	image.URL = image.generateURL()
	return image
}

func (d *Doujin) generateImages() []Image {
	images := make([]Image, len(d.Images.Pages))

	for index, img := range d.Images.Pages {
		images[index] = d.generateImage(index+1, img.Type)
	}
	return images
}

func HandleZip(bufChan chan zipImage, zipFile *zip.Writer, completion chan bool) {

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
		completion <- true
	}
}

func (i Image) ImageToZip(bufChan chan zipImage) {
	buf := new(bytes.Buffer)
	resp, err := http.Get(i.URL)
	catch(err)
	defer resp.Body.Close()
	_, err = io.Copy(buf, resp.Body)
	catch(err)
	bufChan <- zipImage{Img: i, Buf: *buf}
}

func (i *Image) filename() string {
	return fmt.Sprintf("%d.%s", i.Index, i.Type.extension())
}

func (i *Image) generateURL() string {
	const ImageBase = "https://i.nhentai.net"
	return fmt.Sprintf("%s/galleries/%s/%s", ImageBase, i.MediaID, i.filename())
}

func (it *imageType) extension() (ext string) {

	switch *it {
	case jpeg:
		ext = "jpg"
	case png:
		ext = "png"
	case gif:
		ext = "gif"
	}
	return
}

//NewDoujin is used to generate a doujin instance with Image instances attached at APIImages
func NewDoujin(nnn int) Doujin {
	url := fmt.Sprintf("http://nhentai.net/api/gallery/%d", nnn)
	res, err := http.Get(url)
	catch(err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	catch(err)

	hold := Doujin{}
	jsonErr := json.Unmarshal(body, &hold)
	catch(jsonErr)

	hold.APIImages = hold.generateImages()
	return hold
}
