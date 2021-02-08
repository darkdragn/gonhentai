package api

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"time"

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

func (a *APIClient) imageToZip(i Image, bufChan chan zipImage) {
	buf := new(bytes.Buffer)
	resp, err := a.Client.Get(i.URL)
	catch(err)

	defer resp.Body.Close()
	_, err = io.Copy(buf, resp.Body)
	catch(err)

	bufChan <- zipImage{Filename: i.Filename, Buf: *buf}
}

func (a *APIClient) handleZip(bufChan chan zipImage, zipFile *zip.Writer, completion chan bool) {
	for run := range bufChan {
		fh := &zip.FileHeader{
			Name:     run.Filename,
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
