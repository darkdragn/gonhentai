package api

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

//NewClient will generate a basic api client for use.
//Valid optional args are: Int: set the limit for goroutines, bool: Flag artist in use for PrettyNames
func NewClient(args ...interface{}) (client Client) {
	limit := 3
	artist := false

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.Logger = nil

	for _, arg := range args {
		switch t := arg.(type) {
		case int:
			limit = t
		case bool:
			artist = t
		}
	}
	client = Client{
		Artist:  artist,
		BaseURL: "https://nhentai.net/api",
		Client:  retryClient.StandardClient(),
		Limit:   make(chan struct{}, limit),
	}
	return
}

//NewDoujin is used to generate a doujin instance with Image instances attached at APIImages
func (a *Client) NewDoujin(nnn int) Doujin {
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

//NewSearch will return a Search struct from the qurey and page information
func (a *Client) NewSearch(s *SearchOpts) Search {
	surl, err := url.Parse(a.BaseURL + "/galleries/search")
	catch(err)

	values := url.Values{}
	values.Add("query", s.Search)
	values.Add("page", strconv.Itoa(s.Page))
	if s.Sort != "" {
		values.Add("sort", s.Sort)
	}
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

func (a *Client) imageToZip(i Image, bufChan chan zipImage) {
	buf := new(bytes.Buffer)
	resp, err := a.Client.Get(i.URL)
	catch(err)

	defer resp.Body.Close()
	_, err = io.Copy(buf, resp.Body)
	catch(err)

	bufChan <- zipImage{Filename: i.Filename, Buf: *buf}
}

func (a *Client) handleZip(bufChan chan zipImage, zipFile *zip.Writer, completion chan bool) {
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
