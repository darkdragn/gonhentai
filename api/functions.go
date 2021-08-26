package api

import (
	"fmt"
	"log"
)

func catch(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func (i *Image) filename() string {
	return fmt.Sprintf("%d.%s", i.Index, i.Type.extension())
}

func (i *Image) zfilename() string {
	return fmt.Sprintf("%03d.%s", i.Index, i.Type.extension())
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
