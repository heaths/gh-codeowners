package cmd

import (
	"fmt"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/heaths/gh-codeowners/internal/codeowners"
	"github.com/heaths/gh-codeowners/internal/git"
	"github.com/spf13/cobra"
)

const (
	indent = "  "
)

func LintCommand(globalOpts *GlobalOptions) *cobra.Command {
	opts := &lintOptions{
		GlobalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Checks CODEOWNERS for errors",
		Long:  "Checks your CODEOWNERS files for errors as determined by GitHub.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = opts.EnsureRepository()
			if err != nil {
				return
			}

			err = opts.IsAuthenticated()
			if err != nil {
				return
			}

			return lint(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.fix, "fix", false, "Fix errors in the CODEOWNERS file.")
	cmd.Flags().BoolVar(&opts.json, "json", false, "Show errors as JSON.")
	cmd.Flags().BoolVar(&opts.unknownOwners, "unknown-owners", false, "Only list unknown owners.")
	cmd.MarkFlagsMutuallyExclusive("fix", "json", "unknown-owners")

	return cmd
}

type lintOptions struct {
	*GlobalOptions

	fix           bool
	json          bool
	unknownOwners bool
}

func lint(opts *lintOptions) (err error) {
	clientOpts := &api.ClientOptions{
		Host:      opts.host,
		AuthToken: opts.authToken,
	}
	client, err := gh.GQLClient(clientOpts)
	if err != nil {
		return
	}

	refName, err := git.RefName()
	if err != nil {
		return
	}

	errors, err := codeowners.QueryErrors(client, opts.Repo, refName)
	if err != nil {
		return
	}

	if opts.json {
		return printJson(opts.GlobalOptions, errors)
	}

	if opts.unknownOwners {
		missing := errors.UnknownOwners()
		for _, owner := range missing {
			fmt.Fprintln(opts.Console.Stdout(), owner)
		}

		return
	}

	if opts.IsColorEnabled() {
		cs := opts.Console.ColorScheme()
		remove := cs.ColorFunc(opts.Color.Error)

		prettyPrint := func(e codeowners.Error) {
			for _, line := range strings.Split(e.Message, "\n") {
				if e.Kind == codeowners.ErrorKindUnknownOwner {
					line = strings.TrimSpace(line)
					if line == "^" {
						fmt.Fprintln(opts.Console.Stdout())
						return
					} else if line == strings.TrimSpace(e.Source) {
						owner := e.UnknownOwner()
						line = indent + strings.ReplaceAll(line, owner, remove(owner))
					}
				}

				fmt.Fprintln(opts.Console.Stdout(), line)
			}
		}

		for _, e := range errors {
			prettyPrint(e)
		}

		return
	}

	for _, e := range errors {
		fmt.Fprintln(opts.Console.Stdout(), e.Message)
	}

	return
}
