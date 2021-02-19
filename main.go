package main

import (
	// "fmt"
	"fmt"
	"os"
	"strings"

	"github.com/darkdragn/gonhentai/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("gonhentai", "A command-line nhentai ripper.")

	pull = app.Command("pull", "Pull a single gallery.")
	//   registerNick = register.Arg("nick", "Nickname for user.").Required().String()
	//   registerName = register.Arg("name", "Name of user.").Required().String()

	search           = app.Command("search", "Search for hentai.")
	searchString     = search.Arg("search", "Search params (ref nhentai tag docs).").Strings()
	searchAll        = search.Flag("all", "Download all results").Bool()
	searchArtist     = search.Flag("artist", "Store downloads in artist folder").Short('a').Bool()
	searchEnglish    = search.Flag("english", "Add languages:english to the search string").Short('e').Bool()
	searchLong       = search.Flag("long", "Add pages:>50 to search string").Short('l').Bool()
	searchNumber     = search.Flag("number", "Pull index from the search").Short('n').String()
	searchPage       = search.Flag("page", "Add page number to query params").Short('p').Default("1").Int()
	searchUncensored = search.Flag("uncensored", "Add tags:uncensored").Short('u').Bool()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case search.FullCommand():
		text := strings.Join(*searchString, " ")
		client := api.NewClient(*searchArtist)

		if *searchEnglish {
			text += " languages:english"
		}
		if *searchLong {
			text += " pages:>50"
		}
		if *searchUncensored {
			text += " tags:uncensored"
		}
		search := client.NewSearch(text, *searchPage)

		fmt.Printf("%v", *searchNumber)
		if len(*searchNumber) > 0 {
			for n := range strings.Split(*searchNumber, ",") {
				fmt.Printf("%d", n)
				d := search.ReturnDoujin(int(n))
				d.DownloadZip()
			}
		} else if *searchAll {
			for ind := range search.Result {
				doujin := search.ReturnDoujin(ind)
				doujin.DownloadZip()
			}
		} else {
			search.RenderTable(true, *searchPage)
		}
	}
}
