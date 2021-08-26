package api

import (
	"archive/zip"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

//DownloadZip is used to pull the Doujin to disk as {doujin.Title.English}.cbz
func (d Doujin) DownloadZip() {
	var wg sync.WaitGroup
	bufChan := make(chan zipImage)
	// d.Client.Buffer = make(chan zipImage)

	filename := d.determineFilename()
	if _, err := os.Stat(filename); err == nil {
		return
	}

	bar := progressbar.DefaultBytes(-1, d.Titles.Pretty)

	file, err := os.Create(filename)
	defer file.Close()
	catch(err)

	zipFile := zip.NewWriter(io.MultiWriter(file, bar))
	defer zipFile.Close()

	completion := make(chan bool)
	go d.Client.handleZip(bufChan, zipFile, completion)
	go func() {
		for range completion {
			wg.Done()
			// bar.Add(1)
		}
	}()

	for _, img := range d.APIImages {
		d.Client.Limit <- struct{}{}
		wg.Add(1)
		go func(img Image) { d.Client.imageToZip(img, bufChan); <-d.Client.Limit }(img)
	}

	wg.Wait()
	close(bufChan)
	bar.Finish()
}

//Artist will walk tags to discover the first artist tag for the doujin
func (d *Doujin) Artist() string {
	for _, tag := range d.Tags {
		if tag.Type == "artist" {
			return tag.Name
		}
	}
	return "Not Found"
}

func (d *Doujin) determineFilename() (filename string) {
	var title string
	if d.Client.Artist {
		filename = d.Artist() + "/"
		title = d.Titles.Pretty
		_ = os.Mkdir(filename, 0755)
	} else {
		title = d.Titles.English
	}

	title = strings.ReplaceAll(title, "/", "")
	filename += title + ".cbz"
	return
}

func (d *Doujin) generateImage(i int, t imageType) Image {
	image := Image{Index: i, MediaID: d.MediaID, Type: t}
	image.Filename = image.zfilename()
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
