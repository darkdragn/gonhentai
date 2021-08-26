package main

import (
	// "fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/darkdragn/gonhentai/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("gonhentai", "A command-line nhentai ripper.")

	pull       = app.Command("pull", "Pull a single gallery.")
	pullNumber = pull.Arg("index", "Index for the gallery.").Required().Int()

	search        = app.Command("search", "Search for hentai.")
	searchString  = search.Arg("search", "Search params (ref nhentai tag docs).").Required().Strings()
	searchAll     = search.Flag("all", "Download all results").Bool()
	searchArtist  = search.Flag("artist", "Store downloads in artist folder").Short('a').Bool()
	searchEnglish = search.Flag("english", "Add languages:english to the search string").Short('e').Bool()
	searchLong    = search.Flag("long", "Add pages:>50 to search string").Short('l').Bool()
	searchNumber  = search.Flag("number", "Pull index from the search").Short('n').String()
	searchPage    = search.Flag("page", "Add page number to query params").Short('p').Default("1").Int()
	searchPopular = search.Flag("popular", "Add popular query parameters to the search").Bool()
	searchSort    = search.Flag(
		"sort", "Add a sort string to the query params (popular, popular-week, popular-today").
		Short('s').
		Enum("popular", "popular-week", "popular-today")
	searchUncensored = search.Flag("uncensored", "Add tags:uncensored").Short('u').Bool()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case pull.FullCommand():
		index := *pullNumber
		client := api.NewClient()
		doujin := client.NewDoujin(index)
		doujin.DownloadZip()

	case search.FullCommand():
		text := strings.Join(*searchString, " ")
		client := api.NewClient(*searchArtist)
		sort := ""

		if *searchEnglish {
			text += " languages:english"
		}
		if *searchLong {
			text += " pages:>50"
		}
		if *searchUncensored {
			text += " tags:uncensored"
		}
		if *searchPopular {
			sort = "popular"
		} else if *searchSort != "" {
			sort = *searchSort
		}

		search := client.NewSearch(
			&api.SearchOpts{
				Search: text,
				Page:   *searchPage,
				Sort:   sort,
			},
		)

		if len(*searchNumber) > 0 {
			for _, n := range strings.Split(*searchNumber, ",") {
				num, err := strconv.Atoi(n)
				if err != nil {
					log.Panic(err)
				}
				d := search.ReturnDoujin(num)
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
