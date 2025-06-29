package cmd

import (
	"bytes"
	"fincli/internal/iostreams"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCmdEdit(t *testing.T) {
	tests := []struct {
		name        string
		cli         string
		wantsOpts   EditOpts
		wantsErr    bool
		wantsErrMsg string
	}{
		{
			name:      "no args or flags",
			cli:       "",
			wantsOpts: EditOpts{},
			wantsErr:  true,
		},
		{
			name:     "called without required flag",
			cli:      "filename.ext",
			wantsErr: true,
		},
		{
			name:     "no arg, only flag",
			cli:      "--format test_format",
			wantsErr: true,
		},
		{
			name:     "two args, with req flag",
			cli:      "file1 file2 --format test_format",
			wantsErr: true,
		},

		{
			name: "correct command with args",
			cli:  "statement.csv --format test_format",
			wantsOpts: EditOpts{
				FilePath: "statement.csv",
				Format:   "test_format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			stdin := new(bytes.Buffer)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			io := &iostreams.IOStreams{
				In:  stdin,
				Out: stdout,
				Err: stderr,
			}
			var opts *EditOpts
			cmd := NewCmdEdit(io, func(o *EditOpts) error {
				opts = o
				return nil
			})
			require.NotNil(cmd)

			args := parseArgs(t, tt.cli)
			cmd.SetArgs(args)

			err := cmd.Execute()
			if tt.wantsErr {
				require.Error(err)
				assert.Contains(err.Error(), tt.wantsErrMsg)
				return
			}
			require.NoError(err)
			assert.NotNil(opts)
			assert.Same(io, opts.IO)
			assert.Empty(stdin)
			assert.Empty(stdout)
			assert.Empty(stderr)

			assert.Equal(tt.wantsOpts.FilePath, opts.FilePath)
			assert.Equal(tt.wantsOpts.Format, opts.Format)
		})
	}
}

func Test_editRun(t *testing.T) {
	tests := []struct {
		name       string
		opts       *EditOpts
		wantStdout string
	}{
		{
			name:       "initial test",
			opts:       &EditOpts{},
			wantStdout: "hello edit command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin := new(bytes.Buffer)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)
			io := &iostreams.IOStreams{
				In:  stdin,
				Out: stdout,
				Err: stderr,
			}
			tt.opts.IO = io
			err := editRun(tt.opts)

			assert.NoError(t, err)
			assert.Empty(t, stderr.String())
			assert.Equal(t, tt.wantStdout, stdout.String())
		})
	}
}
