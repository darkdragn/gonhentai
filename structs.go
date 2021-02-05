package main

import "fmt"

type Doujin struct {
	ID       int    `json:"id"`
	MediaID  string `json:"media_id"`
	NumPages int    `json:"num_pages"`
	Images   images `json:"images"`
	Titles   title  `json:"title"`
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

type APIImage struct {
	MediaID  string
	Index    int
	Type     imageType
	Filename string
	URL      string
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

func (d *Doujin) generateImage(i int, t imageType) APIImage {
	image := APIImage{Index: i, MediaID: d.MediaID, Type: t}
	image.Filename = image.filename()
	image.URL = image.generateURL()
	return image
}

func (d *Doujin) generateImages() []APIImage {
	images := make([]APIImage, len(d.Images.Pages))

	for index, img := range d.Images.Pages {
		images[index] = d.generateImage(index+1, img.Type)
	}
	return images
}

// func NewAPIImage() APIImage {
// 	image := APIImage{}
// }

func (i *APIImage) filename() string {
	return fmt.Sprintf("%d.%s", i.Index, i.Type.extension())
}

func (i *APIImage) generateURL() string {
	const ImageBase = "https://i.nhentai.net"
	return fmt.Sprintf("%s/galleries/%s/%s", ImageBase, i.MediaID, i.filename())
	// return IMAGE_BASE + "/galleries/" + i.MediaID + "/" + i.Index + "." + i.Type.extension()
}
