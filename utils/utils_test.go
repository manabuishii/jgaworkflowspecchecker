package utils

import (
	"io/ioutil"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_md5File(t *testing.T) {
	md5value, _ := Md5File("../test/testfile.txt")

	assert.Equal(t, "39a870a194a787550b6b5d1f49629236", md5value, "md5 value is different. Contents updated ?")

}
func Test_getFileNameWithoutExtension_1ext(t *testing.T) {
	ext := getFileNameWithoutExtension("../test/testfile.txt")

	assert.Equal(t, "testfile", ext, "Remove directory and last extention.")

}

func Test_getFileNameWithoutExtension_2ext(t *testing.T) {
	ext := getFileNameWithoutExtension("../test/testfile.fa.txt")

	assert.Equal(t, "testfile.fa", ext, "Remove directory and last extention.")

}

func Test_CheckSecondaryFilesExists_all_files_exists(t *testing.T) {
	result, _ := CheckSecondaryFilesExists("../test/secondaryfile/case1/case1.fasta")

	assert.True(t, result, "Check all secondary file")

}

func Test_CheckSecondaryFilesExists_missing_pac_file(t *testing.T) {
	result, _ := CheckSecondaryFilesExists("../test/secondaryfile/case2/case2.fasta")

	assert.False(t, result, "pac file is missing so expected false")

}

func Test_CheckSecondaryFilesExists_missing_dict_file(t *testing.T) {
	result, _ := CheckSecondaryFilesExists("../test/secondaryfile/case3/case3.fasta")

	assert.False(t, result, "^.dict file is missing so expected false")

}
func Test_checkRunDataFile_success(t *testing.T) {
	fileExistsCheckFlag := true
	fileHashCheckFlag := true
	result, _ := CheckRunDataFile("../test/testfile.txt", "39a870a194a787550b6b5d1f49629236", fileExistsCheckFlag, fileHashCheckFlag)

	assert.True(t, result, "md5 match is expected")

}

func Test_checkRunDataFile_fail(t *testing.T) {
	fileExistsCheckFlag := true
	fileHashCheckFlag := true
	result, _ := CheckRunDataFile("../test/testfile.txt", "aa", fileExistsCheckFlag, fileHashCheckFlag)

	assert.False(t, result, "md5 not match is expected")

}

func Test_IsExistsAllResultFilesPrefixSampleId_success(t *testing.T) {
	result := IsExistsAllResultFilesPrefixSampleId("../test/resultfile/success", "XX00000")

	assert.True(t, result, "All files MUST be exists")

}

func Test_IsExistsAllResultFilesPrefixSampleId_missing_crai(t *testing.T) {
	// XX00000.cram.crai is missing
	result := IsExistsAllResultFilesPrefixSampleId("../test/resultfile/fail", "XX00000")

	assert.False(t, result, "XX00000.cram.crai is missing")

}

func Test_IsExistsAllResultFilesPrefixSampleId_filesize_zero(t *testing.T) {
	// XX00000.cram.crai is missing
	result := IsExistsAllResultFilesPrefixSampleId("../test/resultfile/filesizezero", "XX00000")

	assert.False(t, result, "XX00000.cram.crai is missing")

}

func Test_IsExistsAllResultFilesPrefixRunId_success(t *testing.T) {
	result := IsExistsAllResultFilesPrefixRunId("../test/resultfile/success/XX00000", "YYY0000000")

	assert.True(t, result, "All files MUST be exists")

}

func Test_IsExistsAllResultFilesPrefixRunId_fail_runid_bam(t *testing.T) {
	result := IsExistsAllResultFilesPrefixRunId("../test/resultfile/fail_runid_bam/XX00000", "YYY0000000")

	assert.False(t, result, "YYY0000000.bam is missing")

}

func Test_BuildVersionString_local(t *testing.T) {
	result := BuildVersionString("dev", "", "")

	assert.Equal(t, "Version: dev- (built at )\n", result, "version string just dev")

}

func Test_buildVersionString_full(t *testing.T) {
	result := BuildVersionString("0.9.0", "abcdefab", "2021-11-11T11:22:33")

	assert.Equal(t, "Version: 0.9.0-abcdefab (built at 2021-11-11T11:22:33)\n", result, "version string just version, commit id and date")

}

func Test_IsExistsWorkflowFile_valid(t *testing.T) {
	result := IsExistsWorkflowFile("../test/samplefiles/dummy.workflow.cwl")

	assert.True(t, result, "Workflow file is exists")
}

func Test_IsExistsWorkflowFile_fail(t *testing.T) {
	result := IsExistsWorkflowFile("../test/samplefiles/nosucha.workflow.cwl")

	assert.False(t, result, "Workflow file MUST be missing")
}

func Test_IsExistsWorkflowFile_startsWith_http(t *testing.T) {
	result := IsExistsWorkflowFile("http://example.com/test/samplefiles/nosucha.workflow.cwl")

	assert.True(t, result, "Workflow file is startsWith http://.")
}

func Test_IsExistsWorkflowFile_startsWith_https(t *testing.T) {
	result := IsExistsWorkflowFile("https://example.com/test/samplefiles/nosucha.workflow.cwl")

	assert.True(t, result, "Workflow file is startsWith https://.")
}

func Test_IsInCondaEnv_set(t *testing.T) {
	t.Setenv("CONDA_DEFAULT_ENV", "./dummy")
	result := IsInCondaEnv()

	assert.True(t, result, "CONDA_DEFAULT_ENV is set.")
}
func Test_IsInCondaEnv_empty(t *testing.T) {
	t.Setenv("CONDA_DEFAULT_ENV", "")
	result := IsInCondaEnv()

	assert.False(t, result, "CONDA_DEFAULT_ENV is empty so false is expected")
}

func Test_IsInPythonVirtualenv_set(t *testing.T) {
	t.Setenv("VIRTUAL_ENV", "./dummy")
	result := IsInPythonVirtualenv()

	assert.True(t, result, "VIRTUAL_ENV is set.")
}
func Test_IsInPythonVirtualenv_empty(t *testing.T) {
	t.Setenv("VIRTUAL_ENV", "")
	result := IsInPythonVirtualenv()

	assert.False(t, result, "VIRTUAL_ENV is empty so false is expected")
}

func Test_IsInVirtualenv_CONDA_DEFAULT_ENV_set(t *testing.T) {
	t.Setenv("CONDA_DEFAULT_ENV", "./dummy")
	result := IsInVirtualenv()

	assert.True(t, result, "CONDA_DEFAULT_ENV is set.")
}
func Test_IsInVirtualenv_CONDA_DEFAULT_ENV_empty(t *testing.T) {
	t.Setenv("CONDA_DEFAULT_ENV", "")
	result := IsInVirtualenv()

	assert.False(t, result, "CONDA_DEFAULT_ENV is empty so false is expected")
}

func Test_IsInVirtualenv_VIRTUAL_ENV_set(t *testing.T) {
	t.Setenv("VIRTUAL_ENV", "./dummy")
	result := IsInVirtualenv()

	assert.True(t, result, "VIRTUAL_ENV is set.")
}
func Test_IsInVirtualenv_VIRTUAL_ENV_empty(t *testing.T) {
	t.Setenv("VIRTUAL_ENV", "")
	result := IsInVirtualenv()

	assert.False(t, result, "VIRTUAL_ENV is empty so false is expected")
}

func Test_IsOnlyValidCharcterInFilepath(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc")
	assert.True(t, result, "/aaaa/bbb/ccc is all valid character")
}
func Test_IsOnlyValidCharcterInFilepath_lastcharacter_slash(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc/")
	assert.True(t, result, "/aaaa/bbb/ccc/ is last character is slash but all character is ")
}

func Test_IsOnlyValidCharcterInFilepath_period_is_valid_character(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc/a.txt")
	assert.True(t, result, "/aaaa/bbb/ccc/a.txt is all valid character")
}

func Test_IsOnlyValidCharcterInFilepath_capital_character_is_valid(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc/ABC")
	assert.True(t, result, "/aaaa/bbb/ccc/ABC is all valid character")
}
func Test_IsOnlyValidCharcterInFilepath_number_character_is_valid(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc/9a")
	assert.True(t, result, "/aaaa/bbb/ccc/9a is all valid character")
}
func Test_IsOnlyValidCharcterInFilepath_hyphen_is_valid(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc/a-b")
	assert.True(t, result, "/aaaa/bbb/ccc/a-b is all valid character")
}
func Test_IsOnlyValidCharcterInFilepath_underbar_is_valid(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/ccc/a_b")
	assert.True(t, result, "/aaaa/bbb/ccc/a_b is all valid character")
}

func Test_IsOnlyValidCharcterInFilepath_comma(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaa,a/bbb/ccc")
	assert.False(t, result, "/aaa,a/bbb/ccc has ',' invalid charcter")
}

func Test_IsOnlyValidCharcterInFilepath_dollmark(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/b$bb/ccc")
	assert.False(t, result, "/aaaa/b$bb/ccc has '$' invalid charcter")
}

func Test_IsOnlyValidCharcterInFilepath_colon(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/c:cc")
	assert.False(t, result, "/aaaa/bbb/c:cc has ':' invalid charcter")
}

func Test_IsOnlyValidCharcterInFilepath_semicolon(t *testing.T) {
	result := IsOnlyValidCharcterInFilepath("/aaaa/bbb/c;cc")
	assert.False(t, result, "/aaaa/bbb/c;cc has ';' invalid charcter")
}

func Test_SortByFileNameOrderDesc(t *testing.T) {
	files, _ := ioutil.ReadDir("../test/jobManager")
	// First sort by Asc Order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	// Next sort by Desc Order
	SortByFileNameOrderDesc(files)
	assert.Equal(t, files[0].Name(), "20211101145001")
	assert.Equal(t, files[1].Name(), "20211101143242")
	assert.Equal(t, files[2].Name(), "20211101142745")
	assert.Equal(t, files[3].Name(), "20211101141356")
	assert.Equal(t, files[4].Name(), "20211101140603")
	assert.Equal(t, files[5].Name(), "20211101124751")
}

func Test_RemoveStringFromSlice(t *testing.T) {
	sliceOrg := []string{"a", "b", "c", "d", "e"}
	slice := RemoveStringFromSlice(sliceOrg, "a")
	assert.Equal(t, len(slice), 4)
	assert.Equal(t, slice[0], "b")
	assert.Equal(t, slice[1], "c")
	assert.Equal(t, slice[2], "d")
	assert.Equal(t, slice[3], "e")
}

func Test_RemoveStringFromSlice_last(t *testing.T) {
	sliceOrg := []string{"a", "b", "c", "d", "e"}
	slice := RemoveStringFromSlice(sliceOrg, "e")
	assert.Equal(t, len(slice), 4)
	assert.Equal(t, slice[0], "a")
	assert.Equal(t, slice[1], "b")
	assert.Equal(t, slice[2], "c")
	assert.Equal(t, slice[3], "d")
}

func Test_RemoveStringFromSlice_center(t *testing.T) {
	sliceOrg := []string{"a", "b", "c", "d", "e"}
	slice := RemoveStringFromSlice(sliceOrg, "c")
	assert.Equal(t, len(slice), 4)
	assert.Equal(t, slice[0], "a")
	assert.Equal(t, slice[1], "b")
	assert.Equal(t, slice[2], "d")
	assert.Equal(t, slice[3], "e")
}
