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
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
)

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//DownloadZip is used to pull the Doujin to disk as {doujin.Title.English}.cbz
func (d Doujin) DownloadZip(limitRAW int, pretty bool, artist bool) {
	var wg sync.WaitGroup
	limit := make(chan struct{}, limitRAW)
	images := d.APIImages

	var filename string
	if artist {
		pretty = true
		filename = d.Artist() + "/"
		_ = os.Mkdir(filename, 0755)
		// catch(err)
	}

	if pretty {
		// filename = fmt.Sprintf("%s.cbz", d.Titles.Pretty)
		filename += d.Titles.Pretty + ".cbz"
	} else {
		filename = fmt.Sprintf("%s.cbz", d.Titles.English)
	}
	// fmt.Printf("%s\n", filename)
	if _, err := os.Stat(filename); err == nil {
		return
	}
	bar := progressbar.DefaultBytes(-1, d.Titles.Pretty)
	// bar := progressbar.Default(int64(len(images)))

	file, err := os.Create(filename)
	defer file.Close()
	catch(err)

	bufChan := make(chan zipImage)
	zipFile := zip.NewWriter(io.MultiWriter(file, bar))
	// zipFile := zip.NewWriter(file)
	defer zipFile.Close()

	completion := make(chan bool)
	go handleZip(bufChan, zipFile, completion)
	go func() {
		for range completion {
			wg.Done()
			bar.Add(1)
		}
	}()
	for _, img := range images {
		limit <- struct{}{}
		wg.Add(1)
		go func(img Image) { img.imageToZip(bufChan); <-limit }(img)
	}

	wg.Wait()
	close(bufChan)
	bar.Finish()
}

func (d *Doujin) Artist() string {
	for _, tag := range d.Tags {
		if tag.Type == "artist" {
			return tag.Name
		}
	}
	return "Not Found"
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

func (i Image) imageToZip(bufChan chan zipImage) {
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

func (s Search) ReturnDoujin(index int) Doujin {
	magicNumber, _ := s.Result[index].ID.Int64()
	return NewDoujin(int(magicNumber))
}

func (s Search) RenderTable(pretty bool, page int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	for ind, d := range s.Result {
		title := d.Titles.English
		if pretty {
			title = d.Titles.Pretty
		}
		t.AppendRow([]interface{}{ind, d.ID, title})
	}
	t.Render()
	fmt.Printf("Page %d/%d\n", page, s.Pages)
}

func handleZip(bufChan chan zipImage, zipFile *zip.Writer, completion chan bool) {

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

func NewSearch(query string, page int) Search {
	surl, err := url.Parse("https://nhentai.net/api/galleries/search")
	catch(err)

	values := url.Values{}
	values.Add("query", query)
	values.Add("page", strconv.Itoa(page))
	surl.RawQuery = values.Encode()
	resp, err := http.Get(surl.String())
	catch(err)

	search := Search{}
	body, err := ioutil.ReadAll(resp.Body)
	catch(err)

	jsonErr := json.Unmarshal(body, &search)
	catch(jsonErr)
	return search
}
