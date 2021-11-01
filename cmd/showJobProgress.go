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

// showJobProgressCmd represents the showJobProgress command
var showJobProgressCmd = &cobra.Command{
	Use:   "show-job-progress",
	Short: "Show Job Progress",
	Long:  `Show Job Progress`,
	Run: func(cmd *cobra.Command, args []string) {
		showJobProgress(args)
	},
}

var onlynew bool
var onlyfinish bool

func init() {
	rootCmd.AddCommand(showJobProgressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showJobProgressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showJobProgressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	showJobProgressCmd.Flags().BoolVarP(&onlynew, "only-new", "", false, "Show newly execute sample id only")
	showJobProgressCmd.Flags().BoolVarP(&onlyfinish, "only-finish", "", false, "Show finished sample id only")
}
func contains(sampleIdList []string, sampleId string) bool {
	for _, v := range sampleIdList {
		if sampleId == v {
			return true
		}
	}
	return false
}
func showJobProgress(args []string) {
	loadSampleSheetAndConfigFile(args)
	// check in sample sheet data
	if !checkSampleSheet(&ss) {
		return
	}
	// check in config data
	if !checkConfigFile(&rss) {
		return
	}
	displayfinish := true
	if onlynew {
		displayfinish = false
	}
	displaynew := true
	if onlyfinish {
		displaynew = false
	}

	// Setup output directory
	outputDirectoryPath := rss.OutputDirectory.Path
	// Create Sample id list will be executed
	execSampleIdList := utils.CreateExecuteSampleIDList(outputDirectoryPath, &ss)
	if displayfinish {
		//
		for _, s := range ss.SampleList {
			if !contains(execSampleIdList, s.SampleId) {
				fmt.Printf("%s is finished.\n", s.SampleId)
			}
		}
	}
	// TODO display execute information
	DisplayJobInfo(outputDirectoryPath, execSampleIdList)
	if displaynew {
		for _, s := range execSampleIdList {
			fmt.Printf("%s will be Execute new.\n", s)
		}

	}
	fmt.Printf("%d / %d SampleID are finished.\n", len(ss.SampleList)-len(execSampleIdList), len(ss.SampleList))
	fmt.Printf("%d will be executed new.\n", len(execSampleIdList))

}
