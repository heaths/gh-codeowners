package cmd

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

const (
	unknown = "unknown"
)

func LintCommand(globalOpts *GlobalOptions) *cobra.Command {
	var repo string
	opts := &lintOptions{}

	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Checks CODEOWNERS for errors",
		Long:  "Checks your CODEOWNERS files for errors as determined by GitHub.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			opts.GlobalOptions = globalOpts
			if repo != "" {
				opts.Repo, err = repository.Parse(repo)
				if err != nil {
					return
				}
			}

			err = opts.EnsureRepository()
			if err != nil {
				return
			}

			return lint(opts)
		},
	}

	cmd.Flags().StringVarP(&repo, "repo", "R", "", "Select another repository to use using the [HOST/]OWNER/REPO format")
	StringEnumVarP(cmd, &opts.filter, "filter", "f", "", []string{unknown}, "Show owners only for the given filter")

	return cmd
}

type lintOptions struct {
	*GlobalOptions

	filter string
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
				Errors []struct {
					Kind   string
					Line   int
					Column int
					Source string
					// Path   string
				}
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner": graphql.String(opts.Repo.Owner()),
		"repo":  graphql.String(opts.Repo.Name()),
	}
	err = client.Query("CodeownersErrors", &query, variables)
	if err != nil {
		return
	}

	missing := make(map[string]bool)
	for _, e := range query.Repository.Codeowners.Errors {
		if e.Kind == "Unknown owner" && e.Column > 0 {
			owner := e.Source[e.Column-1:]
			if idx := strings.IndexFunc(owner, func(r rune) bool {
				return unicode.IsSpace(r)
			}); idx > 0 {
				owner = owner[:idx]
			}
			missing[owner] = true
		}
	}

	if opts.filter == unknown {
		owners := make([]string, 0, len(missing))
		for k := range missing {
			owners = append(owners, k)
		}
		sort.Strings(owners)

		for _, owner := range owners {
			fmt.Fprintln(opts.Console.Stdout(), owner)
		}
	}

	return
}
