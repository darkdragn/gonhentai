/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"fmt"
	"strconv"

	"github.com/darkdragn/gonhentai/api"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull MAGICNUMBER",
	Short: "Pull a single gallery from nhentai, use the Magic Number",
	Long: `Pull a single galler from nhentai, use the Magic Number...
	
If you're not sure what that is, look at the gallery url:
https://nhentai.net/g/18511/

In this case it's 18511... you can thank me later. ;)`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ind, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Please provide and integer as the argument.")
			fmt.Printf("Argument provided: %v\n", args[0])
		}
		client := api.NewClient()
		doujin := client.NewDoujin(ind)
		doujin.DownloadZip()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
