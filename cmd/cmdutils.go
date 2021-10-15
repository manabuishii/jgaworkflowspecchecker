package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/manabuishiii/jgaworkflowspecchecker/utils"
	"github.com/xeipuuv/gojsonschema"
)

//go:embed samplesheet_schema.json
var samplesheetfileBytes []byte

//go:embed configfile_schema.json
var configfileBytes []byte

/*
 * Behavior:
 *   All fine: true
 *   Something wrong: false
 */
func loadSampleSheetAndConfigFile(args []string) bool {
	if len(args) != 2 {
		fmt.Printf("Some required files are not specified.You pass [%d] file(s)\n", len(args))
		fmt.Println("samplesheet_data configfile_data")
		return false
	}
	path, _ := filepath.Abs("./")
	samplesheet_data_file := args[0]
	config_data_file := args[1]
	allfileexist := true
	if !utils.IsExistsFile(samplesheet_data_file) {
		fmt.Printf("[%s] is missing sample data file\n", samplesheet_data_file)
		allfileexist = false
	}
	if !utils.IsExistsFile(config_data_file) {
		fmt.Printf("[%s] is missing config data file\n", config_data_file)
		allfileexist = false
	}
	if !allfileexist {
		fmt.Println("Some required files are missing. So stop execute")
		return false
	}

	// MUST must be canonical
	schemaLoader := gojsonschema.NewStringLoader(string(samplesheetfileBytes))
	// MUST must be canonical
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + samplesheet_data_file)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}
	if result.Valid() {
		if displayMeesage {
			fmt.Printf("The sample sheet document is valid\n")
		}
	} else {
		fmt.Printf("The sample sheet document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	if displayMeesage {
		fmt.Println("Load sample sheet")
	}
	raw, err := ioutil.ReadFile(samplesheet_data_file)
	if err != nil {
		fmt.Println("Samplesheet has some problems")
		fmt.Println(err.Error())
		return false
	}

	json.Unmarshal(raw, &ss)
	if displayMeesage {
		fmt.Println("Load sample sheet end")
	}
	// configfile loader strings are embed variable
	rschemaLoader := gojsonschema.NewStringLoader(string(configfileBytes))
	// MUST must be canonical
	rdocumentLoader := gojsonschema.NewReferenceLoader("file://" + path + "/" + config_data_file)

	rresult, err := gojsonschema.Validate(rschemaLoader, rdocumentLoader)
	if err != nil {
		panic(err.Error())
	}

	if rresult.Valid() {
		if displayMeesage {
			fmt.Printf("The reference config document is valid\n")
		}
	} else {
		fmt.Printf("The reference config document is not valid. see errors :\n")
		for _, desc := range rresult.Errors() {
			fmt.Printf("- %s\n", desc)
		}
		return false
	}
	if displayMeesage {
		fmt.Println("Load config file")
	}
	rraw, err := ioutil.ReadFile(config_data_file)
	if err != nil {
		fmt.Println("Config file has problem")
		fmt.Println(err.Error())
		return false
	}

	json.Unmarshal(rraw, &rss)
	if displayMeesage {
		fmt.Println("Load config file end")
	}
	return true
}
