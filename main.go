package main

import (
	"flag"
	"fmt"

	"github.com/darkdragn/gonhentai/v2/api"
)

var limitRAW int

func main() {
	var nnn = flag.Int("n", 1234, "The special sauce")
	flag.IntVar(&limitRAW, "limit", 5, "Number of gofuncs to pull images with")
	flag.Parse()

	hold := api.NewDoujin(*nnn)
	fmt.Printf("%s\n", hold.Titles.English)

	// downloadZip(hold)
	hold.DownloadZip(limitRAW)
}
