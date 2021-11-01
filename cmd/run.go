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
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
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
	Long: `Run workflow. When to actually to execute workflow, create output directory.
	And stdout and sterr is redirected to output directory.
	If '--dry-run' flag is set, it only display information do not create directory`,
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
	runCmd.Flags().BoolVarP(&dryrunFlag, "dry-run", "n", false, "Dry-run, do not execute acutal command")
	runCmd.Flags().BoolVarP(&fileExistsCheckFlag, "file-exists-check", "", true, "Check file exists")
	runCmd.Flags().BoolVarP(&fileHashCheckFlag, "file-hash-check", "", true, "Check file hash value")

}
func copyFiles(outputDirectoryPath string, samplesheet_data_file string, config_data_file string) bool {
	// copy sample_sheet file
	sampleSheetDst := outputDirectoryPath + "/" + filepath.Base(samplesheet_data_file)
	// check sample sheet file is same
	if IsSameFilePath(samplesheet_data_file, sampleSheetDst) {
		fmt.Println("Sample sheet file is in output directory backup file. So do not copy")
	} else {

		original_sample_sheet, err := os.Open(samplesheet_data_file)
		if err != nil {
			fmt.Println(err)
		}
		defer original_sample_sheet.Close()
		copied_sample_sheet, err := os.Create(sampleSheetDst)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer copied_sample_sheet.Close()
		_, err = io.Copy(copied_sample_sheet, original_sample_sheet)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	// copy config file
	configFileDst := outputDirectoryPath + "/" + filepath.Base(config_data_file)
	// check config file is same
	if IsSameFilePath(config_data_file, configFileDst) {
		fmt.Println("Config file is in output directory backup file. So do not copy")
	} else {
		original_config_file, err := os.Open(config_data_file)
		if err != nil {
			fmt.Println(err)
		}
		defer original_config_file.Close()
		copied_config_file, err := os.Create(configFileDst)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer copied_config_file.Close()
		_, err = io.Copy(copied_config_file, original_config_file)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}

	return true
}

func createDirectory(outputDirectoryPath string) bool {
	// Create output directory
	// if not create , show error message and exit
	if err := os.MkdirAll(outputDirectoryPath, 0755); err != nil {
		fmt.Println(err)
		fmt.Println("cannot create output directory")
		return false
	}
	return true
}

func checkSampleSheet(ss *utils.SimpleSchema) bool {
	if !utils.CheckSampleSheetFiles(ss, fileExistsCheckFlag, fileHashCheckFlag, displayMeesage) {
		fmt.Println("Some files in sample sheet are missing.")
		return false
	}
	return true
}
func checkConfigFile(rss *utils.ReferenceSchema) bool {
	secondaryFilesCheck, _ := utils.CheckSecondaryFilesExists(rss.Reference.Path)
	if !secondaryFilesCheck {
		fmt.Println("Some secondary file is missing")
		return false
	}

	// Set output directory path
	workflowFilePath := rss.WorkflowFile.Path

	// currently check local filesystem only
	if !utils.IsExistsWorkflowFile(workflowFilePath) {
		fmt.Printf("Missing workflow file [%s]\n", workflowFilePath)
		return false
	}
	// check workflow file is exists
	if !utils.CheckAndDisplayFilesForExecute(rss) {
		fmt.Println("Some files for workflow execution are missing.")
		return false
	}
	return true
}
func runmain(args []string) {
	loadSampleSheetAndConfigFile(args)
	// check in sample sheet data
	if !checkSampleSheet(&ss) {
		return
	}
	// check in config data
	if !checkConfigFile(&rss) {
		return
	}
	// Get Current Time for output directory
	currentTime := utils.GetCurrentTime()
	// Create JobManager Top level directory
	// Setup stdout and stderr to Console and Files (JobManager Top Level Directory).
	if !dryrunFlag {
		// JobManager Top Direcotry
		// This directory is used to store Samples and jobmanager it self logs
		jobManagerTopDirectory := rss.OutputDirectory.Path + "/jobManager/" + currentTime
		if err := os.MkdirAll(jobManagerTopDirectory, 0755); err != nil {
			fmt.Println(err)
			fmt.Println("cannot create JobManager Top Directory")
			return
		}
		// for jobmanager log files
		// Setup log files
		// Stdout and Stderr are redirected to log files
		// if not create , show error message and exit
		jobmanagerstdoutfile := utils.LogStdout(jobManagerTopDirectory + "/jobmanager-stdout.log")
		defer jobmanagerstdoutfile()
		jobmanagerstderrfile := utils.LogStderr(jobManagerTopDirectory + "/jobmanager-stderr.log")
		defer jobmanagerstderrfile()
		// Display Version
		fmt.Print(utils.BuildVersionString(Version, Revision, Date))
		// Display JobManager Top Directory
		fmt.Println("JobManager Top Directory: " + jobManagerTopDirectory)
	} else {
		// Display Version
		fmt.Print(utils.BuildVersionString(Version, Revision, Date))
	}
	//
	foundToilCWLRunner := utils.IsExistsToilCWLRunner()

	// Setup output directory
	outputDirectoryPath := rss.OutputDirectory.Path
	if !dryrunFlag {
		//
		if !foundToilCWLRunner {
			fmt.Println("toil-cwl-runner not found, so can not execute anything.")
			fmt.Println("To ckeck execution environment using `display-jobmanager-recognition`")
			return
		}
		// create output directory
		isDirecotryCreate := createDirectory(outputDirectoryPath)
		if !isDirecotryCreate {
			fmt.Println("Can not create output direcoty")
			os.Exit(1)
		}
		//
		samplesheet_data_file := args[0]
		config_data_file := args[1]
		// copy samplesheet and config file
		isCopyFiles := copyFiles(outputDirectoryPath, samplesheet_data_file, config_data_file)
		if !isCopyFiles {
			fmt.Println("Can not copy files to output direcoty")
			os.Exit(1)
		}
		// Generate sample id list
		utils.GenerateSampleList(&ss, &rss)
	}

	// exec and wait
	var eg errgroup.Group
	executeCount := 0
	for i, s := range ss.SampleList {
		// sample id has something missing. sample id executes
		isExecute := !utils.CheckAllResultFiles(outputDirectoryPath, s)
		if isExecute {
			executeCount += 1
			fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, s.SampleId)
			if !dryrunFlag {
				// TODO exec real
				// check toil-cwl-runner is exists or not
				if foundToilCWLRunner {
					// only exec when toil-cwl-runner is found
					var sampleForExecCWL utils.Sample
					copier.Copy(&sampleForExecCWL, &s)
					eg.Go(func() error {
						utils.ExecCWL(&sampleForExecCWL, &rss, currentTime)
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
