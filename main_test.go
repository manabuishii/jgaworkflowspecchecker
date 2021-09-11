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
