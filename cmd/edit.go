package cmd

import (
	"fincli/internal/iostreams"
	"fmt"

	"github.com/spf13/cobra"
)

type EditOpts struct {
	IO *iostreams.IOStreams
}

func NewCmdEdit(io *iostreams.IOStreams, runFunc func(o *EditOpts) error) *cobra.Command {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := &EditOpts{
				IO: io,
			}
			if runFunc != nil {
				return runFunc(opts)
			}

			return editRun(opts)
		},
	}

	return cmd
}

func editRun(opts *EditOpts) error {
	stdout := opts.IO.Out
	fmt.Fprintf(stdout, "hello edit command")
	return nil
}
