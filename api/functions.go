package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

func catch(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func NewClient(args ...interface{}) (client APIClient) {
	limit := 5
	client = APIClient{
		BaseURL: "https://nhentai.net/api",
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		// Limit: 3,
	}
	for _, arg := range args {
		switch t := arg.(type) {
		case int:
			limit = t
		case bool:
			client.Artist = t
		}
	}
	client.Limit = make(chan struct{}, limit)
	// client.Buffer = make(chan zipImage)
	return
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
	return s.Client.NewDoujin(int(magicNumber))
}

func (s Search) RenderTable(pretty bool, page int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	for ind, d := range s.Result {
		title := d.Titles.English
		artist := d.Artist()
		if pretty {
			title = d.Titles.Pretty
		}
		t.AppendRow([]interface{}{ind, d.ID, artist, title})
	}
	t.Render()
	fmt.Printf("Page %d/%d\n", page, s.Pages)
}

//NewDoujin is used to generate a doujin instance with Image instances attached at APIImages
func (a *APIClient) NewDoujin(nnn int) Doujin {
	url := fmt.Sprintf("%s/gallery/%d", a.BaseURL, nnn)
	res, err := a.Client.Get(url)
	catch(err)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	catch(err)

	hold := Doujin{}
	jsonErr := json.Unmarshal(body, &hold)
	catch(jsonErr)

	hold.APIImages = hold.generateImages()
	hold.Client = a
	return hold
}

func (a *APIClient) NewSearch(query string, page int) Search {
	surl, err := url.Parse(a.BaseURL + "/galleries/search")
	catch(err)

	values := url.Values{}
	values.Add("query", query)
	values.Add("page", strconv.Itoa(page))
	surl.RawQuery = values.Encode()
	resp, err := a.Client.Get(surl.String())
	catch(err)

	search := Search{}
	body, err := ioutil.ReadAll(resp.Body)
	catch(err)

	jsonErr := json.Unmarshal(body, &search)
	catch(jsonErr)

	search.Client = a
	return search
}
