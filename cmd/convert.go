/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fincli/internal/csvstatement"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert [filepath]",
	Short: "Convert a CSV bank statement to a different format",
	Long: `Convert a CSV bank statement from one format to another.

Provide the path to the CSV file as an argument. The argument supports glob patterns, but the pattern must match exactly one file.

The file should be formatted according to the format specified by the required --from flag.`,
	Args: validateArgs,
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		fmt.Printf("convert called on file: %s\n", filepath)

		file, err := os.Open(filepath)
		if err != nil {
			slog.Error("failed to open file", slog.String("filepath", filepath), slog.Any("error", err))
			fmt.Printf("failed to open file %s: %v\n", filepath, err)
			os.Exit(1)
		}
		defer file.Close()

		fromFormatName, err := cmd.Flags().GetString("from")
		if err != nil {
			fmt.Printf("failed to get --from flag: %v", err)
			slog.Error(fmt.Sprintf("failed to get '--from' flag: %v", err))
			os.Exit(1)
		}

		toFormatName, err := cmd.Flags().GetString("to")
		if err != nil {
			fmt.Printf("failed to get --to flag: %v\n", err)
			slog.Error(fmt.Sprintf("failed to get '--to' flag: %v", err))
			os.Exit(1)
		}

		if fromFormatName == "" || toFormatName == "" {
			msg := "required flags '--from' and '--to' must not be empty"
			slog.Error(msg)
			fmt.Println("Error: " + msg)
			os.Exit(1)
		}
		fromFormat, err := csvstatement.GetFormat(fromFormatName)
		if err != nil {
			fmt.Printf("failed to get format '%s': %v\n", fromFormatName, err)
			slog.Error(fmt.Sprintf("failed to get format '%s': %v", fromFormatName, err))
			os.Exit(1)
		}

		toFormat, err := csvstatement.GetFormat(toFormatName)
		if err != nil {
			msg := fmt.Sprintf("failed to get format '%s': %v", toFormatName, err)
			fmt.Println(msg)
			slog.Error(msg)
			os.Exit(1)
		}

		err = csvstatement.Convert(file, os.Stdout, fromFormat, toFormat)
		if err != nil {
			msg := fmt.Sprintf("failed to convert bank statement: %v", err)
			fmt.Println(msg)
			slog.Error(msg)
		}
	},
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		slog.Error("expected exactly one argument, but received multiple", slog.Any("args", args))
		return fmt.Errorf("requires exactly one argument: the path to the CSV file")
	}
	filepathArg := args[0]

	if _, err := os.Stat(filepathArg); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Here you will define your flags and configuration settings.
	convertCmd.Flags().String("from", "", "Name of input format (required)")
	convertCmd.Flags().String("to", "", "Name of output format (required)")

	convertCmd.MarkFlagRequired("from")
	convertCmd.MarkFlagRequired("to")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// convertCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// convertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
