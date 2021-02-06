package api

import (
	"bytes"
	"sync"
)

//Doujin a quick struct for unpacking the response from the nhentai API.
// Used for responses from nhentai.net/api/galleries/:magicNumber
type Doujin struct {
	ID        int    `json:"id"`
	MediaID   string `json:"media_id"`
	NumPages  int    `json:"num_pages"`
	Images    images `json:"images"`
	Titles    title  `json:"title"`
	APIImages []Image
}

type title struct {
	English  string `json:"english"`
	Japanese string `json:"japanese"`
	Pretty   string `json:"pretty"`
}

type page struct {
	Type string
}

type image struct {
	Type   imageType `json:"t"`
	Width  int       `json:"w"`
	Height int       `json:"h"`
}

type images struct {
	Cover     image   `json:"cover"`
	Thumbnail image   `json:"thumbnail"`
	Pages     []image `json:"pages"`
}

type imageType string

const (
	jpeg imageType = "j"
	png            = "p"
	gif            = "g"
)

//Image is a quick build for generating a URL and filename from the
//nhentai api resonse.
type Image struct {
	MediaID  string
	Index    int
	Type     imageType
	Filename string
	URL      string
}

//ZipImage is used for passing an image struct and a buffer between goroutines
type zipImage struct {
	Img Image
	Buf bytes.Buffer
	Wg  *sync.WaitGroup
}
