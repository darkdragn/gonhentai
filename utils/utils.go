package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/darkdragn/gonhentai/v2/api"
)

func catch(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func saveImage(image api.Image, wg *sync.WaitGroup) {
	out, err := os.Create(image.Filename)
	defer out.Close()
	catch(err)
	resp, err := http.Get(image.URL)
	catch(err)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	catch(err)
	wg.Done()
}
