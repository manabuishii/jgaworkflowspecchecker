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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
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
		fmt.Println("run called")
		fmt.Printf("[%t]\n", dryrunFlag)
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

func runmain() {
	path, err := filepath.Abs("./")

	// MUST must be canonical
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[1])
	// MUST must be canonical
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[2])

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}

}

/*

	fmt.Println("Load sample")
	raw, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var ss simpleSchema

	json.Unmarshal(raw, &ss)
	fmt.Println("Load end")

	// validate
	checkResult := false
	for _, s := range ss.SampleList {
		//fmt.Printf("Check index: %d, SampleId: %s\n", i, s.SampleId)
		for j, t := range s.RunList {
			r1, _ := checkRunData(&t.RunData)
			checkResult = checkResult || r1
			if !r1 {
				fmt.Println("Some error found. Not exist or Hash value error")
				fmt.Printf("Check index: %d, RunId: %s\n", j, t.RunId)
				fmt.Printf("pe or se: [%s]\n", t.RunData.PEOrSE)
				fmt.Printf("fq1: [%s]\n", t.RunData.FQ1)
				fmt.Printf("fq2: [%s]\n", t.RunData.FQ2)
				fmt.Printf("result=%t\n", r1)
			}
		}
	}
	if !checkResult {
		fmt.Println("some thing wrong. do not execute")
		return
	}
	// reference config validate

	// MUST must be canonical
	rschemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[3])
	// MUST must be canonical
	rdocumentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + os.Args[4])

	rresult, err := gojsonschema.Validate(rschemaLoader, rdocumentLoader)
	if err != nil {
		panic(err.Error())
	}

	if rresult.Valid() {
		fmt.Printf("The reference config document is valid\n")
	} else {
		fmt.Printf("The reference config document is not valid. see errors :\n")
		for _, desc := range rresult.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return
	}
	fmt.Println("Load sample")
	rraw, err := ioutil.ReadFile(os.Args[4])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var rss referenceSchema

	json.Unmarshal(rraw, &rss)
	fmt.Println("Load end")

	// if dryrunFlag {
	// 	fmt.Println("Dry-run flag is set")
	// 	return
	// }
	//
	secondaryFilesCheck, err := checkSecondaryFilesExists(rss.Reference.Path)
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

	// create job file for CWL
	createJobFile(&ss, &rss)

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
			check1 := isExistsAllResultFilesPrefixSampleId(outputDirectoryPath, s.SampleId)
			if !check1 {
				isExecute = true
			}
			// RunID prefix files check
			for _, r := range s.RunList {
				check2 := isExistsAllResultFilesPrefixRunId(outputDirectoryPath+"/"+s.SampleId, r.RunId)
				if !check2 {
					isExecute = true
				}
			}
			//fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, s.SampleId)
		}
		if isExecute {
			executeCount += 1

			fmt.Printf("index: %d, SampleId: %s will be Execute new.\n", i, s.SampleId)
			sampleId := s.SampleId
			if !dryrunFlag {
				// TODO exec real
				// if foundToilCWLRunner {
				// 	// only exec when toil-cwl-runner is found
				// 	eg.Go(func() error {
				// 		execCWL(outputDirectoryPath, workflowFilePath, sampleId)
				// 		return nil
				// 	})
				// }
			}
		}
	}
	if dryrunFlag {
		fmt.Printf("[%d/%d] task will be executed.\n", executeCount, len(ss.SampleList))
	}

	// for i, s := range ss.SampleList {
	// 	fmt.Printf("index: %d, SampleId: %s\n", i, s.SampleId)
	// 	sampleId := s.SampleId
	// 	eg.Go(func() error {
	// 		// TODO check exit status
	// 		execCWL(outputDirectoryPath, workflowFilePath, sampleId)
	// 		return nil
	// 	})
	// }

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	if foundToilCWLRunner {
		fmt.Println("Can not find toil-cwl-runner")
	}

	fmt.Println("fin")

}
*/
