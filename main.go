// Copyright 2023 Heath Stewart.
// Licensed under the MIT License. See LICENSE.txt in the project root for license information.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/heaths/gh-codeowners/internal/cmd"
	"github.com/heaths/go-console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultColorComment = "#6A9955"
	defaultColorError   = "#F44747"
)

var (
	colorRE = regexp.MustCompile("#[0-9A-Fa-f]{6}")
)

func main() {
	con := console.System()
	log := log.New(con.Stderr(), "", log.Ltime)

	opts := &cmd.GlobalOptions{
		Console: con,
		Log:     log,
	}

	v := viper.New()
	if dir, err := os.UserHomeDir(); err == nil {
		dir = filepath.Join(dir, ".config", "gh-codeowners")
		v.SetConfigType("yml")
		v.AddConfigPath(dir)
	}

	loadColorConfig := func(key string, field *string) {
		val := v.Get(key)
		if s, ok := val.(string); ok && colorRE.MatchString(s) {
			*field = s
			return
		}

		if opts.Verbose {
			log.Printf("config %q is not a valid color matching #RRGGBB; skipping...", key)
		}
	}

	rootCmd := cobra.Command{
		Use:   "codeowners",
		Short: "Check CODEOWNERS file",
		Long:  "GitHub CLI extension to check your CODEOWNERS file.",
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			opts.Color = cmd.ColorOptions{
				Comment: defaultColorComment,
				Error:   defaultColorError,
			}

			if err := v.ReadInConfig(); err != nil && opts.Verbose {
				log.Printf("failed to load config: %q, skipping...", err)
				return
			}

			loadColorConfig("color.comment", &opts.Color.Comment)
			loadColorConfig("color.error", &opts.Color.Error)
		},
		SilenceUsage: true,
	}

	// Output options
	rootCmd.SetOut(con.Stdout())
	rootCmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Log verbose output.")

	// Colors options
	rootCmd.PersistentFlags().String("color-comment", defaultColorComment, fmt.Sprintf("Hex RGB color code for comments e.g., %q.", defaultColorComment))
	rootCmd.PersistentFlags().String("color-error", defaultColorError, fmt.Sprintf("Hex RGB color code for errors e.g., %q.", defaultColorError))

	_ = v.BindPFlag("color.comment", rootCmd.PersistentFlags().Lookup("color-comment"))
	_ = v.BindPFlag("color.error", rootCmd.PersistentFlags().Lookup("color-error"))

	// Subcommands
	rootCmd.AddCommand(cmd.LintCommand(opts))
	rootCmd.AddCommand(cmd.ViewCommand(opts))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
