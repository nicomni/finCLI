package cmd

import (
	"fincli/internal/iostreams"
	"fmt"

	"github.com/spf13/cobra"
)

type EditOpts struct {
	IO       *iostreams.IOStreams
	Format   string
	FilePath string
}

func NewCmdEdit(io *iostreams.IOStreams, runFunc func(o *EditOpts) error) *cobra.Command {
	opts := &EditOpts{
		IO: io,
	}
	cmd := &cobra.Command{
		Use:  "edit [filepath]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FilePath = args[0]

			if runFunc != nil {
				return runFunc(opts)
			}

			return editRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.Format, "format", "", "Format of the file to edit (required)")
	cmd.MarkFlagRequired("format")

	return cmd
}

func editRun(opts *EditOpts) error {
	stdout := opts.IO.Out
	fmt.Fprintf(stdout, "hello edit command")
	return nil
}
