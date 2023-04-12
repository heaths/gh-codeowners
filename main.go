// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"log"
	"os"

	"github.com/heaths/gh-codeowners/internal/cmd"
	"github.com/heaths/go-console"
	"github.com/spf13/cobra"
)

func main() {
	con := console.System()
	opts := &cmd.GlobalOptions{
		Console: con,
		Log:     log.New(con.Stdout(), "", log.Ltime),
	}

	rootCmd := cobra.Command{
		Use:          "codeowners",
		Short:        "Check CODEOWNERS file",
		Long:         "GitHub CLI extension to check your CODEOWNERS file.",
		SilenceUsage: true,
	}

	rootCmd.SetOut(con.Stdout())

	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Log verbose output.")
	rootCmd.AddCommand(cmd.LintCommand(opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
