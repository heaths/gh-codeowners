package cmd

import (
	"fmt"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/heaths/gh-codeowners/internal/codeowners"
	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
)

func PrCommand(globalOpts *GlobalOptions) *cobra.Command {
	opts := &prOptions{
		GlobalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:   "pr",
		Short: "Views the owners for a list of files in a pull request",
		Long:  "Shows the owners for each file in a pull request. You must be in a repository to find the CODEOWNERS file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = opts.EnsureRepository()
			if err != nil {
				return
			}

			err = opts.IsAuthenticated()
			if err != nil {
				return
			}

			opts.number, err = parseNumberRef(args[0])
			if err != nil {
				return
			}

			return pr(opts)
		},
	}

	return cmd
}

type prOptions struct {
	*GlobalOptions

	number int
}

func pr(opts *prOptions) (err error) {
	clientOpts := &api.ClientOptions{
		Host:      opts.host,
		AuthToken: opts.authToken,
	}
	client, err := gh.GQLClient(clientOpts)
	if err != nil {
		return
	}

	fs := opts.RootFS()
	path := codeowners.Find(fs)
	if path == "" {
		return fmt.Errorf("CODEOWNERS not found")
	}

	var c *codeowners.Codeowners
	c, err = codeowners.Open(fs, path)
	if err != nil {
		return err
	}

	var query struct {
		Repository struct {
			PullRequest struct {
				Files struct {
					Nodes []struct {
						Path       string
						ChangeType string
					}
					PageInfo struct {
						HasNextPage bool
						EndCursor   string
					}
				} `graphql:"files(first: 100, after: $endCursor)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	variables := map[string]interface{}{
		"owner":     graphql.String(opts.Repo.Owner()),
		"repo":      graphql.String(opts.Repo.Name()),
		"number":    graphql.Int(opts.number),
		"endCursor": graphql.String(""),
	}

	var files []file
	for {
		err = client.Query("PullRequestFiles", &query, variables)
		if err != nil {
			return
		}

		for _, node := range query.Repository.PullRequest.Files.Nodes {
			files = append(files, file{
				Path:       node.Path,
				ChangeType: node.ChangeType,
				Owners:     c.Owners(node.Path),
			})
		}

		if query.Repository.PullRequest.Files.PageInfo.HasNextPage {
			variables["endCursor"] = graphql.String(query.Repository.PullRequest.Files.PageInfo.EndCursor)
		} else {
			break
		}
	}

	return printJson(opts.GlobalOptions, files)
}

type file struct {
	Path       string   `json:"path"`
	ChangeType string   `json:"changeType"`
	Owners     []string `json:"owners"`
}
