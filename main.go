package main

import (
	"flag"

	"github.com/darkdragn/gonhentai/v2/api"
)

var limitRAW int

func main() {
	var nnn = flag.Int("n", 1234, "The special sauce")
	var query = flag.String("query", "", "Run a query on nhentai")
	var page = flag.Int("page", 1, "Select search page number.")
	var pretty = flag.Bool("p", false, "Use pretty name")
	flag.IntVar(&limitRAW, "limit", 5, "Number of gofuncs to pull images with")
	flag.Parse()

	if f := flag.CommandLine.Lookup("query"); *query != f.DefValue {
		search := api.NewSearch(*query, *page)

		if *nnn < 25 {
			d := search.ReturnDoujin(*nnn)
			d.DownloadZip(limitRAW, *pretty)
		} else {
			search.RenderTable(*pretty, *page)
		}
	} else {
		hold := api.NewDoujin(*nnn)
		hold.DownloadZip(limitRAW, *pretty)
	}
}
