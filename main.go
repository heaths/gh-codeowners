// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/heaths/gh-codeowners/internal/cmd"
	"github.com/heaths/go-console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultColorComment = "#6A9955"
	defaultColorError   = "#F44747"
)

func main() {
	con := console.System()
	opts := &cmd.GlobalOptions{
		Console: con,
		Log:     log.New(con.Stdout(), "", log.Ltime),
	}

	v := viper.New()
	if dir, err := os.UserHomeDir(); err == nil {
		dir = filepath.Join(dir, ".config", "gh-codeowners")
		v.SetConfigType("yml")
		v.AddConfigPath(dir)
	}

	rootCmd := cobra.Command{
		Use:          "codeowners",
		Short:        "Check CODEOWNERS file",
		Long:         "GitHub CLI extension to check your CODEOWNERS file.",
		SilenceUsage: true,
	}

	// Output options
	rootCmd.SetOut(con.Stdout())
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Log verbose output.")

	// Colors options
	rootCmd.PersistentFlags().StringVar(&opts.Color.Comment, "color-comment", defaultColorComment, fmt.Sprintf("Hex RGB color code for comments e.g., %q.", defaultColorComment))
	rootCmd.PersistentFlags().StringVar(&opts.Color.Error, "color-error", defaultColorError, fmt.Sprintf("Hex RGB color code for errors e.g., %q.", defaultColorError))

	// BUGBUG: https://github.com/spf13/viper/issues/1537
	_ = v.BindPFlag("color.comment", rootCmd.PersistentFlags().Lookup("color-comment"))
	_ = v.BindPFlag("color.error", rootCmd.PersistentFlags().Lookup("color-error"))

	// Subcommands
	rootCmd.AddCommand(cmd.LintCommand(opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
