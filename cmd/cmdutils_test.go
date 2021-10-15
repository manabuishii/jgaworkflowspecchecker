package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadSampleSheetAndConfigFile_1file_success(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json", "../test/datafiles/configfile_1run-test.json"})

	assert.True(t, result, "All files MUST be exists")
}

func Test_loadSampleSheetAndConfigFile_pass_1file_fail(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json"})

	assert.False(t, result, "input file MUST be exactly 2")
}

func Test_loadSampleSheetAndConfigFile_pass_3file_fail(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json", "../test/datafiles/configfile_1run-test.json", "3rd.file"})

	assert.False(t, result, "input file MUST be exactly 2")
}

func Test_loadSampleSheetAndConfigFile_1file_fq1_missing_success_as_configfile(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test-fail.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "fq1 missing but Valid as samplesheet file and config file")
}

func Test_loadSampleSheetAndConfigFile_sample_data_file_missing(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/nosuchasamplesheet.json", "../test/datafiles/configfile_1run-test.json"})

	assert.False(t, result, "nosuchasamplesheet.json MUST be missing")
}

func Test_loadSampleSheetAndConfigFile_config_data_file_missing(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test-fail.json", "../test/datafiles/nosuchaconfigfile.json"})

	assert.False(t, result, "nosuchaconfigfile MUST be missing")
}

func Test_loadSampleSheetAndConfigFile_samplesheet_data_and_config_data_file_missing(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/nosuchasamplesheet.json", "../test/datafiles/nosuchaconfigfile.json"})

	assert.False(t, result, "nosuchasamplesheet.json and nosuchaconfigfile MUST be missing")

}

func Test_loadSampleSheetAndConfigFile_1file_fq1_missing_fail(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test-fail.json", "../test/datafiles/configfile_1run-test.json"})

	assert.True(t, result, "All files MUST be exists")
}

func Test_loadSampleSheetAndConfigFile_samplesheet_is_valid_configfile_data_is_invalid_fail(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json", "../test/datafiles/invalid_configfile_data.json"})

	assert.False(t, result, "Samplesheet data is valid But Configfile is invalid")
}

func Test_loadSampleSheetAndConfigFile_samplesheet_is_invalid_configfile_data_is_valid_fail(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/invalid_samplesheet_data.json", "../test/datafiles/configfile_1run-test.json"})

	assert.False(t, result, "Samplesheet data is invalid But Configfile is valid")
}
