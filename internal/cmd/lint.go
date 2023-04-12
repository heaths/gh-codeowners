package cmd

import (
	"fmt"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/heaths/gh-codeowners/internal/codeowners"
	"github.com/heaths/gh-codeowners/internal/git"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func LintCommand(globalOpts *GlobalOptions) *cobra.Command {
	opts := &lintOptions{}

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Checks CODEOWNERS for errors",
		Long:  "Checks your CODEOWNERS files for errors as determined by GitHub.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.GlobalOptions = globalOpts

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
	cmd.Flags().BoolVar(&opts.unknownOwners, "unknown-owners", false, "Only list unknown owners.")
	cmd.MarkFlagsMutuallyExclusive("fix", "unknown-owners")

	return cmd
}

type lintOptions struct {
	*GlobalOptions

	fix           bool
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

	var query struct {
		Repository struct {
			Codeowners struct {
				Errors codeowners.Errors
			} `graphql:"codeowners(refName: $branch)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	branch, err := git.BranchRef()
	if err != nil {
		return
	}

	variables := map[string]interface{}{
		"owner":  graphql.String(opts.Repo.Owner()),
		"repo":   graphql.String(opts.Repo.Name()),
		"branch": graphql.String(branch),
	}
	err = client.Query("CodeownersErrors", &query, variables)
	if err != nil {
		return
	}

	if opts.unknownOwners {
		missing := query.Repository.Codeowners.Errors.UnknownOwners()
		for _, owner := range missing {
			fmt.Fprintln(opts.Console.Stdout(), owner)
		}

		return
	}

	linterOpts := codeowners.LintOptions{
		Console: opts.Console,
		Fix:     opts.fix,
	}

	root, err := git.RootFS()
	if err != nil {
		return
	}

	return codeowners.Lint(root, query.Repository.Codeowners.Errors, linterOpts)
}
