package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BuildVersionString_local(t *testing.T) {
	result := BuildVersionString("dev", "", "")

	assert.Equal(t, "Version: dev- (built at )\n", result, "version string just dev")

}

func Test_buildVersionString_full(t *testing.T) {
	result := BuildVersionString("0.9.0", "abcdefab", "2021-11-11T11:22:33")

	assert.Equal(t, "Version: 0.9.0-abcdefab (built at 2021-11-11T11:22:33)\n", result, "version string just version, commit id and date")

}
