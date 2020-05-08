package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/google/go-github/v31/github"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

func action(c *cli.Context) error {
	if c.NArg() != 3 {
		if err := cli.ShowSubcommandHelp(c); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return cli.NewExitError("", 1)
	}
	milestoneTitle := c.Args().Get(0)
	prID, err := strconv.Atoi(c.Args().Get(1))
	if err != nil {
		return cli.NewExitError("The PR ID must be an integer", 1)
	}
	token := c.Args().Get(2)

	log.Info().Msgf("Verifying cherry-pick PR %d for milestone %q...", prID, milestoneTitle)

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	repoSvc := client.Repositories

	ms, err := getMilestone(milestoneTitle, repoSvc)
	if err != nil {
		return err
	}

	return nil
}

type milestone struct {
}

func getMilestone(title string, repoSvc *github.RepositoriesService) (milestone, error) {
	return milestone{}, nil
}

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		PartsOrder: []string{"message"},
		FormatMessage: func(msg interface{}) string {
			return fmt.Sprintf("* %s", msg)
		},
	}
	log.Logger = log.Output(output)

	app := &cli.App{
		Name:      "verify-cherrypicks",
		Usage:     "Tool to verify a Grafana cherry-picks PR",
		ArgsUsage: "<milestone-title> <pr-id> <github-token>",
		Version:   "0.1.0",
		Action:    action,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("An unexpected error occurred")
	}
}
