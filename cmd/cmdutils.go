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

//go:embed create_docker_image.sh
var createDockerImageScript []byte

//go:embed create_singularity_image.sh
var createSingularityImageScript []byte

/*
 * Behavior:
 *   All fine: true
 *   Something wrong: false
 */
func loadSampleSheetAndConfigFile(args []string) bool {
	if len(args) != 2 {
		fmt.Printf("Some required files are not specified. You pass [%d] file(s)\n", len(args))
		fmt.Println("samplesheet_data configfile_data")
		return false
	}
	path, _ := filepath.Abs("./")
	samplesheet_data_file := args[0]
	config_data_file := args[1]
	// Check sample sheet filename and config filename has invalid character.
	allfilepathisvalidchar := true
	if !utils.IsOnlyValidCharcterInFilepath(samplesheet_data_file) {
		fmt.Printf("[%s] has invalid character.\n", samplesheet_data_file)
		allfilepathisvalidchar = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(config_data_file) {
		fmt.Printf("[%s] has invalid character.\n", config_data_file)
		allfilepathisvalidchar = false
	}
	if !allfilepathisvalidchar {
		fmt.Println("Some required files has invalid character. So stop execute")
		return false
	}
	// Check sample sheet filename and config filename is exist.
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
	// files are provided. check both files contents.
	validateisfine := true
	// sample sheet schema provided by embed.
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
		validateisfine = false
	}
	if displayMeesage {
		fmt.Println("Load sample sheet")
	}
	raw, err := ioutil.ReadFile(samplesheet_data_file)
	if err != nil {
		fmt.Println("Samplesheet has some problems")
		fmt.Println(err.Error())
		validateisfine = false
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
		validateisfine = false
	}
	if displayMeesage {
		fmt.Println("Load config file")
	}
	rraw, err := ioutil.ReadFile(config_data_file)
	if err != nil {
		fmt.Println("Config file has problem")
		fmt.Println(err.Error())
		validateisfine = false
	}

	json.Unmarshal(rraw, &rss)
	if displayMeesage {
		fmt.Println("Load config file end")
	}
	validateisfine = validateisfine && IsAllSamplesheetFilepathHasValidchar(&ss)
	validateisfine = validateisfine && IsAllFilepathInConfigFileHasValidchar(&rss)
	return validateisfine
}

func IsAllSamplesheetFilepathHasValidchar(samplesheet *utils.SimpleSchema) bool {
	result := true
	for _, s := range ss.SampleList {
		for _, r := range s.RunList {
			if r.PEOrSE == "PE" {
				// PE
				if !utils.IsOnlyValidCharcterInFilepath(r.FQ1) {
					fmt.Printf("In SampleID[%s] RunID[%s] [%s] has invalid character in filepath\n", s.SampleId, r.RunId, r.FQ1)
					result = false
				}
				if !utils.IsOnlyValidCharcterInFilepath(r.FQ2) {
					fmt.Printf("In SampleID[%s] RunID[%s] [%s] has invalid character in filepath\n", s.SampleId, r.RunId, r.FQ2)
					result = false
				}
			} else {
				// SE
				if !utils.IsOnlyValidCharcterInFilepath(r.FQ1) {
					fmt.Printf("In SampleID[%s] RunID[%s] [%s] has invalid character in filepath\n", s.SampleId, r.RunId, r.FQ1)
					result = false
				}
			}
		}
	}
	return result
}

func IsAllFilepathInConfigFileHasValidchar(rss *utils.ReferenceSchema) bool {
	result := true
	if !utils.IsOnlyValidCharcterInFilepath(rss.WorkflowFile.Path) {
		fmt.Printf("In config file, `workflow_file` path [%s] has invalid character.\n", rss.WorkflowFile.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.OutputDirectory.Path) {
		fmt.Printf("In config file, `output_directory` path [%s] has invalid character.\n", rss.OutputDirectory.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.ContainerCacheDirectory.Path) {
		fmt.Printf("In config file, `container_cache_directory` path [%s] has invalid character.\n", rss.ContainerCacheDirectory.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.Reference.Path) {
		fmt.Printf("In config file, `reference` path [%s] has invalid character.\n", rss.Reference.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.Dbsnp.Path) {
		fmt.Printf("In config file, `dnsnp` path [%s] has invalid character.\n", rss.Dbsnp.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.Mills.Path) {
		fmt.Printf("In config file, `mills` path [%s] has invalid character.\n", rss.Mills.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.KnownIndels.Path) {
		fmt.Printf("In config file, `known_indels` path [%s] has invalid character.\n", rss.KnownIndels.Path)
		result = false
	}
	// Autosome PAR
	if !utils.IsOnlyValidCharcterInFilepath(rss.HaplotypecallerAutosomePARIntervalBed.Path) {
		fmt.Printf("In config file, `haplotypecaller_autosome_PAR_interval_bed` path [%s] has invalid character.\n", rss.HaplotypecallerAutosomePARIntervalBed.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.HaplotypecallerAutosomePARIntervalList.Path) {
		fmt.Printf("In config file, `haplotypecaller_autosome_PAR_interval_list` path [%s] has invalid character.\n", rss.HaplotypecallerAutosomePARIntervalList.Path)
		result = false
	}
	// ChrX NonPAR
	if !utils.IsOnlyValidCharcterInFilepath(rss.HaplotypecallerChrXNonPARIntervalBed.Path) {
		fmt.Printf("In config file, `haplotypecaller_chrX_nonPAR_interval_bed` path [%s] has invalid character.\n", rss.HaplotypecallerChrXNonPARIntervalBed.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.HaplotypecallerChrXNonPARIntervalList.Path) {
		fmt.Printf("In config file, `haplotypecaller_chrX_nonPAR_interval_list` path [%s] has invalid character.\n", rss.HaplotypecallerChrXNonPARIntervalList.Path)
		result = false
	}
	// ChrY NonPar
	if !utils.IsOnlyValidCharcterInFilepath(rss.HaplotypecallerChrYNonPARIntervalBed.Path) {
		fmt.Printf("In config file, `haplotypecaller_chrY_nonPAR_interval_bed` path [%s] has invalid character.\n", rss.HaplotypecallerChrYNonPARIntervalBed.Path)
		result = false
	}
	if !utils.IsOnlyValidCharcterInFilepath(rss.HaplotypecallerChrYNonPARIntervalList.Path) {
		fmt.Printf("In config file, `haplotypecaller_chrY_nonPAR_interval_list` path [%s] has invalid character.\n", rss.HaplotypecallerChrYNonPARIntervalList.Path)
		result = false
	}
	return result
}
