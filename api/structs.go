package api

import (
	"bytes"
	"encoding/json"
	"sync"
)

//Doujin a quick struct for unpacking the response from the nhentai API.
// Used for responses from nhentai.net/api/galleries/:magicNumber
type Doujin struct {
	ID       json.Number `json:"id"`
	MediaID  string      `json:"media_id"`
	NumPages int         `json:"num_pages"`
	Images   struct {
		Cover     image   `json:"cover"`
		Thumbnail image   `json:"thumbnail"`
		Pages     []image `json:"pages"`
	} `json:"images"`
	Titles struct {
		English  string `json:"english"`
		Japanese string `json:"japanese"`
		Pretty   string `json:"pretty"`
	} `json:"title"`
	APIImages []Image
}

type Search struct {
	Result []Doujin `json:"result"`
	Pages  int      `json:"num_pages"`
}

type image struct {
	Type   imageType `json:"t"`
	Width  int       `json:"w"`
	Height int       `json:"h"`
}

//Image is a quick build for generating a URL and filename from the
//nhentai api resonse.
type Image struct {
	MediaID  string
	Index    int
	Type     imageType
	Filename string
	URL      string
}

type imageType string

const (
	jpeg imageType = "j"
	png            = "p"
	gif            = "g"
)

type zipImage struct {
	Img Image
	Buf bytes.Buffer
	Wg  *sync.WaitGroup
}
