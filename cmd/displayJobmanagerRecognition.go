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

	"github.com/manabuishiii/jgaworkflowspecchecker/utils"
	"github.com/spf13/cobra"
)

// displayJobmanagerRecognitionCmd represents the displayJobmanagerRecognition command
var displayJobmanagerRecognitionCmd = &cobra.Command{
	Use:   "display-jobmanager-recognition",
	Short: "Display JobManager Recognition",
	Long: `Display JobManager Recognition
Virutlenv state
Singularity command
Slurm command
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("displayJobmanagerRecognition called")
		// TODO resolv workflow path
		loadSampleSheetAndConfigFile(args)
		utils.DisplayJobManagerRecoginition(&rss)
		utils.CheckSampleSheetFiles(&ss, fileExistsCheckFlag, fileHashCheckFlag)
	},
}

func init() {
	rootCmd.AddCommand(displayJobmanagerRecognitionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// displayJobmanagerRecognitionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// displayJobmanagerRecognitionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
