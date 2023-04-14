package cmd

import (
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/heaths/gh-codeowners/internal/codeowners"
	"github.com/heaths/gh-codeowners/internal/git"
	"github.com/spf13/cobra"
)

func ViewCommand(globalOpts *GlobalOptions) *cobra.Command {
	opts := &viewOptions{
		GlobalOptions: globalOpts,
	}

	cmd := &cobra.Command{
		Use:   "view",
		Short: "Views the CODEOWNERS file with errors highlighted",
		Long:  "Checks your CODEOWNERS files for errors as determined by GitHub and renders the CODEOWNERS file.",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			err = opts.EnsureRepository()
			if err != nil {
				return
			}

			err = opts.IsAuthenticated()
			if err != nil {
				return
			}

			return view(opts)
		},
	}

	return cmd
}

type viewOptions struct {
	*GlobalOptions
}

func view(opts *viewOptions) (err error) {
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

	renderOpts := codeowners.RenderOptions{
		Console: opts.Console,
		Color:   opts.Color,
	}

	root, err := git.RootFS()
	if err != nil {
		return
	}

	return codeowners.Render(root, errors, renderOpts)
}
