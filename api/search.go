package api

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

//ReturnDoujin will take an index number and return a doujin from the indexed search `Result`
func (s Search) ReturnDoujin(index int) Doujin {
	magicNumber, _ := s.Result[index].ID.Int64()
	return s.Client.NewDoujin(int(magicNumber))
}

//RenderTable will provide a pretty view of the search results
func (s Search) RenderTable(pretty bool, page int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	for ind, d := range s.Result {
		title := d.Titles.English
		artist := d.Artist()
		if pretty {
			title = d.Titles.Pretty
		}
		if len(title) > 75 {
			title = title[0:75]
		}
		t.AppendRow([]interface{}{ind, d.ID, artist, title})
	}
	t.Render()
	fmt.Printf("Page %d/%d\n", page, s.Pages)
}
