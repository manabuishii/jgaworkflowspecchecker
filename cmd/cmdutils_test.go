package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadSampleSheetAndConfigFile_PE_success(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json", "../test/datafiles/configfile_1run-test.json"})

	assert.True(t, result, "All files MUST be exists")
}
func Test_loadSampleSheetAndConfigFile_SE_success(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
}
func Test_loadSampleSheetAndConfigFile_SampleSheet_Filename_has_invalid_char(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1ru;n-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.False(t, result, "SampleSheet filename has invalid char")
}
func Test_loadSampleSheetAndConfigFile_Configfile_Filename_has_invalid_char(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run;-test.json"})
	assert.False(t, result, "Config file filename has invalid char")
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

func Test_loadSampleSheetAndConfigFile_samplesheet_PE_fq1_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	ss.SampleList[0].RunList[0].FQ1 = "aaa;aaa"
	ssresult := IsAllSamplesheetFilepathHasValidchar(&ss)

	assert.False(t, ssresult, "ss.SampleList[0].RunList[0].FQ1 has invalid character")
}
func Test_loadSampleSheetAndConfigFile_samplesheet_PE_fq2_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	ss.SampleList[0].RunList[0].FQ2 = "aaa;aaa"
	ssresult := IsAllSamplesheetFilepathHasValidchar(&ss)

	assert.False(t, ssresult, "ss.SampleList[0].RunList[0].FQ2 has invalid character")
}
func Test_loadSampleSheetAndConfigFile_samplesheet_SE_fq1_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	ss.SampleList[0].RunList[0].FQ1 = "aaa;aaa"
	ssresult := IsAllSamplesheetFilepathHasValidchar(&ss)

	assert.False(t, ssresult, "ss.SampleList[0].RunList[0].FQ1 has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_workflow_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.WorkflowFile.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.WorkflowFile.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_outputdirectory_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.OutputDirectory.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.OutputDirectory.Path has invalid character")
}
func Test_loadSampleSheetAndConfigFile_configfile_container_cache_directgory_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.ContainerCacheDirectory.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.ContainerCacheDirectory.Path has invalid character")
}
func Test_loadSampleSheetAndConfigFile_configfile_reference_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.Reference.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.Reference.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_dbsnp_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.Dbsnp.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.Dbsnp.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_mills_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.Mills.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.Mills.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_known_indels_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.KnownIndels.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.KnownIndels.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_haplotypecaller_autosome_PAR_interval_bed_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.HaplotypecallerAutosomePARIntervalBed.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.HaplotypecallerAutosomePARIntervalBed.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_haplotypecaller_autosome_PAR_interval_list_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.HaplotypecallerAutosomePARIntervalList.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.HaplotypecallerAutosomePARIntervalList.Path has invalid character")
}

func Test_loadSampleSheetAndConfigFile_configfile_haplotypecaller_chrX_nonPAR_interval_bed_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.HaplotypecallerChrXNonPARIntervalBed.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.HaplotypecallerChrXNonPARIntervalBed.Path has invalid character")
}
func Test_loadSampleSheetAndConfigFile_configfile_haplotypecaller_chrX_nonPAR_interval_list_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.HaplotypecallerChrXNonPARIntervalList.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.HaplotypecallerChrXNonPARIntervalList.Path has invalid character")
}
func Test_loadSampleSheetAndConfigFile_configfile_haplotypecaller_chrY_nonPAR_interval_bed_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.HaplotypecallerChrYNonPARIntervalBed.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.HaplotypecallerChrYNonPARIntervalBed.Path has invalid character")
}
func Test_loadSampleSheetAndConfigFile_configfile_haplotypecaller_chrY_nonPAR_interval_list_has_invalid_character(t *testing.T) {
	result := loadSampleSheetAndConfigFile([]string{"../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/configfile_1run-test.json"})
	assert.True(t, result, "All files MUST be exists")
	rss.HaplotypecallerChrYNonPARIntervalList.Path = "aaa;aaa"
	rssresult := IsAllFilepathInConfigFileHasValidchar(&rss)

	assert.False(t, rssresult, "rss.HaplotypecallerChrYNonPARIntervalList.Path has invalid character")
}
func Test_IsSameFilePath(t *testing.T) {
	result := IsSameFilePath("../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/samplesheet_1run-SE-test.json")
	assert.True(t, result, "The same file")
}
func Test_IsSameFilePath_different_file(t *testing.T) {
	result := IsSameFilePath("../test/datafiles/samplesheet_1run-SE-test.json", "../test/datafiles/samplesheet_1run-SE-test.json.bak")
	assert.False(t, result, "The different file")
}
func Test_IsSameFilePath_symlink(t *testing.T) {
	// create temp directory
	tmpDir := t.TempDir()
	// create temporary file
	tmpFile := filepath.Join(tmpDir, "tmpFile")
	err := ioutil.WriteFile(tmpFile, []byte("test"), 0644)
	assert.NoError(t, err, "create temporary file")

	// create symbolic link
	symlink := filepath.Join(tmpDir, "symlink")
	err = os.Symlink(tmpFile, symlink)
	assert.NoError(t, err, "create symbolic link")
	result := IsSameFilePath(symlink, tmpFile)
	assert.True(t, result, "The same file")
}
func Test_IsSameFilePath_symlink_different_file(t *testing.T) {
	// create temp directory
	tmpDir := t.TempDir()
	// create temporary file
	tmpFile := filepath.Join(tmpDir, "tmpFile")
	err := ioutil.WriteFile(tmpFile, []byte("test"), 0644)
	assert.NoError(t, err, "create temporary file")
	// create symbolic link
	symlink := filepath.Join(tmpDir, "symlink")
	err = os.Symlink(tmpFile, symlink)
	assert.NoError(t, err, "create symbolic link")
	result := IsSameFilePath(symlink, tmpFile+".bak")
	assert.False(t, result, "The different file")
}
func Test_IsSameFilePath_symlink_symlink(t *testing.T) {
	// create temp directory
	tmpDir := t.TempDir()
	// create temporary file
	tmpFile := filepath.Join(tmpDir, "tmpFile")
	err := ioutil.WriteFile(tmpFile, []byte("test"), 0644)
	assert.NoError(t, err, "create temporary file")
	// create symbolic link
	symlink := filepath.Join(tmpDir, "symlink")
	err = os.Symlink(tmpFile, symlink)
	assert.NoError(t, err, "create symbolic link")
	// create symbolic link
	symlink2 := filepath.Join(tmpDir, "symlink2")
	err = os.Symlink(symlink, symlink2)
	assert.NoError(t, err, "create symbolic link")
	result := IsSameFilePath(symlink2, symlink)
	assert.True(t, result, "The same file")
}

func Test_IsSameFilePath_symlink_symlink_different_file(t *testing.T) {
	// create temp directory
	tmpDir := t.TempDir()
	// create temporary file
	tmpFile := filepath.Join(tmpDir, "tmpFile")
	err := ioutil.WriteFile(tmpFile, []byte("test"), 0644)
	assert.NoError(t, err, "create temporary file")
	// create symbolic link
	symlink := filepath.Join(tmpDir, "symlink")
	err = os.Symlink(tmpFile, symlink)
	assert.NoError(t, err, "create symbolic link")
	// create temporary file
	tmpFile2 := filepath.Join(tmpDir, "tmpFile2")
	err = ioutil.WriteFile(tmpFile2, []byte("test"), 0644)
	assert.NoError(t, err, "create temporary file")
	// create symbolic link
	symlink2 := filepath.Join(tmpDir, "symlink2")
	err = os.Symlink(tmpFile2, symlink2)
	assert.NoError(t, err, "create symbolic link")
	result := IsSameFilePath(symlink2, symlink)
	assert.False(t, result, "The different file")

}

func Test_getExitCodeContent(t *testing.T) {
	result := getExitCodeContent("../test/jobManager/20211101145001/XX00001/toil.exitcode.txt")
	assert.Equal(t, "1", result, "Exit code 1")
}

func Test_getExitCodeContent_otherfile(t *testing.T) {
	result := getExitCodeContent("../test/jobManager/20211101143242/XX00002/toil.exitcode.txt")
	assert.Equal(t, "1", result, "Exit code 1")
}

func Test_getExitCodeContent_exitcode_0(t *testing.T) {
	result := getExitCodeContent("../test/jobManager/20211101143242/XX00003/toil.exitcode.txt")
	assert.Equal(t, "0", result, "Exit code 0")
}
