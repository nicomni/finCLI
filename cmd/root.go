// Package cmd
package cmd

import (
	"fincli/internal/iostreams"

	"github.com/spf13/cobra"
)

var cfgFile string

func NewCmdRoot(io *iostreams.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fincli",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fincli.yaml)")

	cmd.AddCommand(NewCmdConvert(io, nil))

	return cmd
}
