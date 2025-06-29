package cmd

import (
	"fincli/internal/csvstatement"
	"fincli/internal/iostreams"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type ConvertOptions struct {
	IO       *iostreams.IOStreams

	FilePath   string
	FromFormat string
	ToFormat   string
}

func NewCmdConvert(io *iostreams.IOStreams, runF func(*ConvertOptions) error) *cobra.Command {
	opts := &ConvertOptions{
		IO: io,
	}

	cmd := &cobra.Command{
		Use:   "convert [filepath]",
		Short: "Convert a CSV bank statement to a different format",
		Long: `Convert a CSV bank statement from one format to another.

		Provide the path to the CSV file as an argument. The argument supports glob patterns, but the pattern must match exactly one file.

		The file should be formatted according to the format specified by the required --from flag.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FilePath = args[0]
			fmt.Printf("convert called on file: %s\n", opts.FilePath)

			if opts.FromFormat == "" || opts.ToFormat == "" {
				msg := "required flags '--from' and '--to' must not be empty"
				return fmt.Errorf(msg)
			}

			if runF != nil {
				return runF(opts)
			}

			return convertRun(opts)
		},
	}

	cmd.Flags().StringVar(&opts.FromFormat, "from", "", "Name of input format (required)")
	cmd.MarkFlagRequired("from")
	cmd.Flags().StringVar(&opts.ToFormat, "to", "", "Name of output format (required)")
	cmd.MarkFlagRequired("to")

	cmd.Flags().Bool("auto-suggest", false, "Enable auto suggestion")
	return cmd
}

func convertRun(opts *ConvertOptions) error {
	file, err := os.Open(opts.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", opts.FilePath, err)
	}
	defer file.Close()
	fromFormat, err := csvstatement.GetFormat(opts.FromFormat)
	if err != nil {
		return fmt.Errorf("failed to get format '%s': %v", opts.FromFormat, err)
	}

	toFormat, err := csvstatement.GetFormat(opts.ToFormat)
	if err != nil {
		msg := fmt.Sprintf("failed to get format '%s': %v", opts.ToFormat, err)
		return fmt.Errorf(msg)
	}
	err = csvstatement.Convert(file, os.Stdout, fromFormat, toFormat)
	if err != nil {
		msg := fmt.Sprintf("failed to convert bank statement: %v", err)
		return fmt.Errorf(msg)
	}
	return nil
}
