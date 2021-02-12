//Package cmd is used to build the cli
/*
Copyright © 2021 Darkdragn

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/darkdragn/gonhentai/api"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [flags] SEARCH_QUERY",
	Short: "Run a search on nhentai! Pull a selected gallery with -n",
	Long: `The easiest way to rip a mountain of content from nhentai for
personal use! Search supports all tags and other. Checkout of nhentai's
search FAQ for details.

Example: This will pull all uncensored galleries from artist yamatogawa
that have been translated into english

gonhentai search "artist:yamatogawa tags:uncensored languages:english" -a --all

Or even easier:
gonhentai search "artist:yamatogawa" -e -u -a --all`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		artist, _ := cmd.Flags().GetBool("artist")
		page, _ := cmd.Flags().GetInt("page")

		client := api.NewClient(artist)
		searchString := args[0]

		checks := map[string]string{
			"uncensored": " tags:uncensored",
			"english":    " languages:english",
			"long":       " pages:>50",
		}

		for b, s := range checks {
			v, _ := cmd.Flags().GetBool(b)
			if v {
				searchString += s
			}
		}

		search := client.NewSearch(searchString, page)

		if cmd.Flags().Changed("number") {
			ns, _ := cmd.Flags().GetIntSlice("number")
			for _, n := range ns {
				d := search.ReturnDoujin(n)
				d.DownloadZip()
			}
		} else if all {
			for ind := range search.Result {
				doujin := search.ReturnDoujin(ind)
				doujin.DownloadZip()
			}
		} else {
			search.RenderTable(true, page)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().Bool("all", false, "Download all results")
	searchCmd.Flags().BoolP("artist", "a", false, "Store things in an artist directory")
	searchCmd.Flags().BoolP("english", "e", false, "Add languages:english to search string")
	searchCmd.Flags().BoolP("long", "l", false, "Add pages:>50 to search string")
	searchCmd.Flags().BoolP("uncensored", "u", false, "Add tags:uncensored to search string")
	searchCmd.Flags().IntSliceP("number", "n", []int{50}, "Pull")
	searchCmd.Flags().IntP("page", "p", 1, "Select search page")
}
