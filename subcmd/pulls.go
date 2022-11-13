package subcmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/go-github/v48/github"
	"github.com/winebarrel/frk"
	"github.com/winebarrel/frk/utils"
)

type PullsCmd struct {
	StartDate string `short:"s"`
	EndDate   string `short:"e"`
}

func (cmd *PullsCmd) Run(binds *frk.Binds) error {
	startDate := time.Time{}
	endDate := time.Time{}

	if cmd.StartDate != "" {
		var err error
		startDate, err = dateparse.ParseAny(cmd.StartDate)

		if err != nil {
			return fmt.Errorf("failed to parse '--start-date': %s", cmd.StartDate)
		}
	}

	if cmd.EndDate != "" {
		var err error
		endDate, err = dateparse.ParseAny(cmd.EndDate)

		if err != nil {
			return fmt.Errorf("failed to parse '--end-date': %s", cmd.EndDate)
		}
	}

	if !startDate.IsZero() && !endDate.IsZero() {
		if startDate.After(endDate) {
			return errors.New("'start-date' is greater than 'end-date'")
		}
	}

	prsByRepo, err := binds.Github.PullRequests()

	if err != nil {
		return fmt.Errorf("failed to get pull requests: %w", err)
	}

	fmt.Print("Pull Requests")

	if !startDate.IsZero() && !endDate.IsZero() {
		fmt.Printf(" (%s...%s)", startDate.Local().Format("2006-01-02"), endDate.Local().Format("2006-01-02"))
	}

	fmt.Printf("\n-\n\n")

	for _, i := range prsByRepo.Itor() {
		pullRequests := []*github.PullRequest{}

		for _, pr := range i.PullRequests {
			if !startDate.IsZero() && pr.CreatedAt.Before(startDate) {
				continue
			}

			if !endDate.IsZero() && pr.CreatedAt.After(endDate) {
				continue
			}

			pullRequests = append(pullRequests, pr)
		}

		if len(pullRequests) == 0 {
			continue
		}

		fmt.Printf("### %s\n", i.Repo)

		for _, pr := range pullRequests {
			fmt.Printf("- [%s](%s):", *pr.Title, *pr.HTMLURL)

			if pr.Body != nil && *pr.Body != "" {
				fmt.Printf(" (%s)", utils.Str(*pr.Body).Plain().Cut())
			}

			fmt.Printf(" %s\n", pr.CreatedAt.Local().Format("2006-01-02"))
		}

		fmt.Println()
	}

	return nil
}
