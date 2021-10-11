package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_md5File(t *testing.T) {
	md5value, _ := md5File("./test/testfile.txt")

	assert.Equal(t, "39a870a194a787550b6b5d1f49629236", md5value, "md5 value is different. Contents updated ?")

}

func Test_getFileNameWithoutExtension_1ext(t *testing.T) {
	ext := getFileNameWithoutExtension("./test/testfile.txt")

	assert.Equal(t, "testfile", ext, "Remove directory and last extention.")

}

func Test_getFileNameWithoutExtension_2ext(t *testing.T) {
	ext := getFileNameWithoutExtension("./test/testfile.fa.txt")

	assert.Equal(t, "testfile.fa", ext, "Remove directory and last extention.")

}

func Test_checkSecondaryFilesExists_all_files_exists(t *testing.T) {
	result, _ := checkSecondaryFilesExists("./test/secondaryfile/case1/case1.fasta")

	assert.True(t, result, "Check all secondary file")

}

func Test_checkSecondaryFilesExists_missing_pac_file(t *testing.T) {
	result, _ := checkSecondaryFilesExists("./test/secondaryfile/case2/case2.fasta")

	assert.False(t, result, "pac file is missing so expected false")

}

func Test_checkSecondaryFilesExists_missing_dict_file(t *testing.T) {
	result, _ := checkSecondaryFilesExists("./test/secondaryfile/case3/case3.fasta")

	assert.False(t, result, "^.dict file is missing so expected false")

}

func Test_checkRunDataFile_success(t *testing.T) {
	fileExistsCheckFlag = true
	fileHashCheckFlag = true
	result, _ := checkRunDataFile("./test/testfile.txt", "39a870a194a787550b6b5d1f49629236")

	assert.True(t, result, "md5 match is expected")

}

func Test_checkRunDataFile_fail(t *testing.T) {
	fileExistsCheckFlag = true
	fileHashCheckFlag = true
	result, _ := checkRunDataFile("./test/testfile.txt", "aa")

	assert.False(t, result, "md5 not match is expected")

}

func Test_isExistsAllResultFilesPrefixSampleId_success(t *testing.T) {
	result := isExistsAllResultFilesPrefixSampleId("test/resultfile/success", "XX00000")

	assert.True(t, result, "All files MUST be exists")

}

func Test_isExistsAllResultFilesPrefixSampleId_missing_crai(t *testing.T) {
	// XX00000.cram.crai is missing
	result := isExistsAllResultFilesPrefixSampleId("test/resultfile/fail", "XX00000")

	assert.False(t, result, "XX00000.cram.crai is missing")

}

func Test_isExistsAllResultFilesPrefixSampleId_filesize_zero(t *testing.T) {
	// XX00000.cram.crai is missing
	result := isExistsAllResultFilesPrefixSampleId("test/resultfile/filesizezero", "XX00000")

	assert.False(t, result, "XX00000.cram.crai is missing")

}

func Test_isExistsAllResultFilesPrefixRunId_success(t *testing.T) {
	result := isExistsAllResultFilesPrefixRunId("test/resultfile/success/XX00000", "YYY0000000")

	assert.True(t, result, "All files MUST be exists")

}

func Test_isExistsAllResultFilesPrefixRunId_fail_runid_bam(t *testing.T) {
	result := isExistsAllResultFilesPrefixRunId("test/resultfile/fail_runid_bam/XX00000", "YYY0000000")

	assert.False(t, result, "YYY0000000.bam is missing")

}
