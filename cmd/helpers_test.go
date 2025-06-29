package cmd

import (
	"testing"

	"github.com/google/shlex"
	"github.com/stretchr/testify/require"
)

// parseArgs is a helper function to parse command line arguments
func parseArgs(t *testing.T, args string) []string {
	t.Helper()
	parsedArgs, err := shlex.Split(args)
	require.NoError(t, err)
	return parsedArgs
}
