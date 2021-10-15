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
	"io"
	"log"
	"os"

	"github.com/manabuishiii/jgaworkflowspecchecker/utils"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

//
var dryrunFlag bool
var fileExistsCheckFlag bool
var fileHashCheckFlag bool

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run workflow",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runmain(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().BoolVarP(&dryrunFlag, "dry-run", "n", false, "Dry-run, do not execute acutal command")
	runCmd.Flags().BoolVarP(&fileExistsCheckFlag, "file-exists-check", "", true, "Check file exists")
	runCmd.Flags().BoolVarP(&fileHashCheckFlag, "file-hash-check", "", true, "Check file hash value")

}

func runmain(args []string) {
	loadSampleSheetAndConfigFile(args)

	// files in sample sheet
	if !utils.CheckSampleSheetFiles(&ss, fileExistsCheckFlag, fileHashCheckFlag, displayMeesage) {
		fmt.Println("Some files in sample sheet are missing.")
		return
	}

	secondaryFilesCheck, _ := utils.CheckSecondaryFilesExists(rss.Reference.Path)
	if !secondaryFilesCheck {
		fmt.Println("Some secondary file is missing")
		return
	}

	// Set output directory path
	workflowFilePath := rss.WorkflowFile.Path

	// currently check local filesystem only
	if !utils.IsExistsWorkflowFile(workflowFilePath) {
		fmt.Printf("Missing workflow file [%s]\n", workflowFilePath)
		os.Exit(1)
	}
	// Set output directory path
	outputDirectoryPath := rss.OutputDirectory.Path
	// Create output directory
	// if not create , show error message and exit
	if err := os.MkdirAll(outputDirectoryPath, 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create output directory")
		return
	}

	if err := os.MkdirAll(outputDirectoryPath+"/toil-outputs", 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create toil outputs directory")
		return
	}
	if err := os.MkdirAll(outputDirectoryPath+"/logs", 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create logs directory")
		return
	}
	if err := os.MkdirAll(outputDirectoryPath+"/jobstores", 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create jobstores directory")
		return
	}

	// copy sample_sheet file
	original_sample_sheet, err := os.Open(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	defer original_sample_sheet.Close()
	copied_sample_sheet, err := os.Create(outputDirectoryPath + "/" + original_sample_sheet.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer copied_sample_sheet.Close()
	_, err = io.Copy(copied_sample_sheet, original_sample_sheet)
	if err != nil {
		log.Fatal(err)
	}

	// copy config file
	original_configfile, err := os.Open(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	defer original_configfile.Close()
	copied_configfile, err := os.Create(outputDirectoryPath + "/" + original_configfile.Name())
	if err != nil {
		log.Fatal(err)
	}
	defer copied_configfile.Close()
	_, err = io.Copy(copied_configfile, original_configfile)
	if err != nil {
		log.Fatal(err)
	}
	// check workflow file is exists
	if !utils.CheckAndDisplayFilesForExecute(&rss) {
		fmt.Println("Some files for workflow execution are missing.")
		os.Exit(1)
	}
	// create job file for CWL
	utils.CreateJobFile(&ss, &rss)

	// check toil-cwl-runner is exists or not
	foundToilCWLRunner := utils.IsExistsToilCWLRunner()

	// dry-run

	// exec and wait
	var eg errgroup.Group
	executeCount := 0
	for i, s := range ss.SampleList {
		isExecute := false
		// Check SampleId result directory is exist
		if _, err := os.Stat(outputDirectoryPath + "/" + s.SampleId); os.IsNotExist(err) {
			// SampleId result directory is missing
			// so this id must be executed
			isExecute = true
		} else {
			// check all result file is found or not
			// SampleId prefix files check
			check1 := utils.IsExistsAllResultFilesPrefixSampleId(outputDirectoryPath, s.SampleId)
			if !check1 {
				isExecute = true
			}
			// RunID prefix files check
			for _, r := range s.RunList {
				check2 := utils.IsExistsAllResultFilesPrefixRunId(outputDirectoryPath+"/"+s.SampleId, r.RunId)
				if !check2 {
					isExecute = true
				}
			}
			//fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, s.SampleId)
		}
		if isExecute {
			executeCount += 1
			sampleId := s.SampleId

			fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, sampleId)
			if !dryrunFlag {
				// TODO exec real
				if foundToilCWLRunner {
					// only exec when toil-cwl-runner is found
					eg.Go(func() error {
						utils.ExecCWL(outputDirectoryPath, workflowFilePath, sampleId)
						return nil
					})
				}
			}
		}
	}
	if dryrunFlag {
		fmt.Printf("[%d/%d] task will be executed.\n", executeCount, len(ss.SampleList))
	}
	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("fin")

}
