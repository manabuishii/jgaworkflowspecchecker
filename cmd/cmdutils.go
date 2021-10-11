package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/manabuishiii/jgaworkflowspecchecker/utils"
	"github.com/xeipuuv/gojsonschema"
)

//
func loadSampleSheetAndConfigFile(args []string) {
	path, _ := filepath.Abs("./")
	samplesheet_schema_file := args[0]
	samplesheet_data_file := args[1]
	config_schema_file := args[2]
	config_data_file := args[3]
	allfileexist := true
	if !utils.IsExistsFile(samplesheet_schema_file) {
		fmt.Printf("[%s] is missing\n", samplesheet_schema_file)
		allfileexist = false
	}
	if !utils.IsExistsFile(samplesheet_data_file) {
		fmt.Printf("[%s] is missing\n", samplesheet_data_file)
		allfileexist = false
	}
	if !utils.IsExistsFile(config_schema_file) {
		fmt.Printf("[%s] is missing\n", config_schema_file)
		allfileexist = false
	}
	if !utils.IsExistsFile(config_data_file) {
		fmt.Printf("[%s] is missing\n", config_data_file)
		allfileexist = false
	}
	if !allfileexist {
		fmt.Println("Some required files are missing. So stop execute")
		return
	}

	// MUST must be canonical
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + samplesheet_schema_file)
	// MUST must be canonical
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + samplesheet_data_file)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if result.Valid() {
		fmt.Printf("The sample sheet document is valid\n")
	} else {
		fmt.Printf("The sample sheet document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	fmt.Println("Load sample sheet")
	raw, err := ioutil.ReadFile(samplesheet_data_file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &ss)
	fmt.Println("Load sample sheet end")

	// validate
	/*
		checkResult := false
		for _, s := range ss.SampleList {
			//fmt.Printf("Check index: %d, SampleId: %s\n", i, s.SampleId)
			for j, t := range s.RunList {
				r1, _ := utils.CheckRunData(&t.RunData, fileExistsCheckFlag, fileHashCheckFlag)
				checkResult = checkResult || r1
				if !r1 {
					fmt.Println("At sample sheet check. Some error found. Sample Not exist or Hash value error")
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
	*/
	// reference config validate

	// MUST must be canonical
	rschemaLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + config_schema_file)
	// MUST must be canonical
	rdocumentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + config_data_file)

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
	fmt.Println("Load config file")
	rraw, err := ioutil.ReadFile(config_data_file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	json.Unmarshal(rraw, &rss)
	fmt.Println("Load config file end")
}
