package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething3(t *testing.T) {
	md5value, _ := md5File("./test/testfile.txt")

	assert.Equal(t, "39a870a194a787550b6b5d1f49629236", md5value, "True is true!")

}
