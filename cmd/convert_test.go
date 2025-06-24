package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: What should happen if run in non-interactive mode?
// TODO: What should happen if stdin is not TTY?
// TODO: What should happen if stdout is not TTY?
// Are there situations when stdin/out is TTY, but we still can't prompt the user?

func TestNewCmdConvert(t *testing.T) {
	tests := []struct {
		name        string
		cli         string
		wantsOpts   ConvertOptions
		wantsErr    bool
		wantsErrMsg string
	}{
		{
			name:     "convert file",
			cli:      "path/to/file --from FROM_FORMAT --to TO_FORMAT",
			wantsErr: false,
			wantsOpts: ConvertOptions{
				FilePath:   "path/to/file",
				FromFormat: "FROM_FORMAT",
				ToFormat:   "TO_FORMAT",
			},
		},
		{
			name:     "no args",
			cli:      "",
			wantsErr: true,
		},
		{
			name:     "too many args",
			cli:      "arg1 arg2",
			wantsErr: true,
		},
		{
			name:     "missing required --from flag",
			cli:      "file1 --to format",
			wantsErr: true,
		},
		{
			name:     "missing required --to flag",
			cli:      "file1 --from format",
			wantsErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts *ConvertOptions
			cmd := NewCmdConvert(func(o *ConvertOptions) error {
				opts = o
				return nil
			})

			// TODO: consider using github.com/google/shlex.Split()
			var args []string
			if tt.cli == "" {
				args = nil
			} else {
				args = strings.Split(tt.cli, " ")
			}
			cmd.SetArgs(args)

			cmd.SetIn(new(bytes.Buffer))
			cmd.SetOut(new(bytes.Buffer))
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.Execute()
			if tt.wantsErr {
				assert.Error(t, err)
				if tt.wantsErrMsg != "" {
					assert.Equal(t, tt.wantsErrMsg, err.Error())
				}
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.wantsOpts.FilePath, opts.FilePath)
			assert.Equal(t, tt.wantsOpts.FromFormat, opts.FromFormat)
			assert.Equal(t, tt.wantsOpts.ToFormat, opts.ToFormat)
		})
	}
}

func Test_convertRun(t *testing.T) {
	assert.True(t, false)
}
