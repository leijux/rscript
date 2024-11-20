package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leijux/rscript/internal/pkg/parser"
)

func TestParseRScriptFile(t *testing.T) {
	script, err := parser.ParseWithPath("testdata/script.yaml")
	require.NoError(t, err)

	assert.Equal(t, uint(1), script.SchemaVersion)
	assert.Equal(t, "example.zip", script.Variables["update_package_name"])
	assert.Equal(t, "rm -rf ./example.zip", script.Commands[5])
	assert.Equal(t, "192.168.0.1:22", script.Remotes[0].AddrPort.String())

	assert.Equal(t, "root", script.Remotes[0].Username)
	assert.Equal(t, "123456", script.Remotes[0].Password)

}

func TestParseEmptyConfigFile(t *testing.T) {
	script, err := parser.ParseWithPath("testdata/script_empty.yaml")
	require.Error(t, err)

	require.Nil(t, script)
}

func TestParseDefault(t *testing.T) {
	script, err := parser.ParseWithPath("testdata/script_default.yaml")
	require.NoError(t, err)

	assert.Equal(t, "192.168.0.1:22", script.Remotes[0].AddrPort.String())
	assert.Equal(t, "root", script.Remotes[0].Username)
	assert.Equal(t, "123456", script.Remotes[0].Password)
}
