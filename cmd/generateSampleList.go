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
	"os"

	"github.com/manabuishiii/jgaworkflowspecchecker/utils"
	"github.com/spf13/cobra"
)

// generateSampleListCmd represents the generateSampleList command
var generateSampleListCmd = &cobra.Command{
	Use:   "generate-sample-list",
	Short: "Generate sample list",
	Long:  `Generate sample list from samplesheet file.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateSampleListMain(args)
	},
}

func init() {
	rootCmd.AddCommand(generateSampleListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateSampleListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateSampleListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func generateSampleListMain(args []string) bool {
	// load SampleSheet and ConfigFile
	loadSuccess := loadSampleSheetAndConfigFile(args)
	if !loadSuccess {
		os.Exit(1)
	}
	result := utils.GenerateSampleList(&ss, &rss)
	return result
}
