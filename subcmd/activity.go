package subcmd

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/go-github/v48/github"
	"github.com/winebarrel/frk"
	"github.com/winebarrel/frk/utils"
)

type ActivityCmd struct {
	Since int    `short:"d" default:"0"`
	From  string `short:"f"`
	To    string `short:"t"`
}

func (cmd *ActivityCmd) Run(binds *frk.Binds) error {
	from := time.Now().Truncate(24 * time.Hour)
	to := time.Now().Truncate(24 * time.Hour)

	if cmd.From != "" {
		var err error
		from, err = dateparse.ParseAny(cmd.From)

		if err != nil {
			return fmt.Errorf("failed to parse '--from': %s", cmd.From)
		}
	}

	if cmd.To != "" {
		var err error
		to, err = dateparse.ParseAny(cmd.To)

		if err != nil {
			return fmt.Errorf("failed to parse '--to': %s", cmd.To)
		}
	}

	if from.After(to) {
		return errors.New("'--from' is greater than '--to'")
	}

	diff := int(to.Sub(from).Hours() / 24)

	if diff == 0 {
		from = from.AddDate(0, 0, -cmd.Since)
	} else {
		cmd.Since = diff
	}

	eventsByRepo, err := binds.Github.EventsWithGrouping(from, to)

	if err != nil {
		return fmt.Errorf("failed to get events: %w", err)
	}

	var period string

	if cmd.Since == 0 {
		period = "Today's"
	} else {
		period = fmt.Sprintf("%d days", cmd.Since+1)
	}

	fmt.Printf("%s Activities\n-\n\n", period)

	for _, i := range eventsByRepo.Itor() {
		fmt.Printf("### %s\n\n", i.Repo)

		events := i.Events
		sort.Slice(events, func(i, j int) bool { return *events[i].Type < *events[j].Type })
		alreadyOutputs := map[string]struct{}{}
		var title string
		var htmlURL string

		for _, e := range events {
			payloadType := utils.PayloadType(*e.Type)
			payload, err := e.ParsePayload()

			if err != nil {
				return err
			}

			switch v := payload.(type) {
			case *github.IssueCommentEvent:
				htmlURL = *v.Comment.HTMLURL
				title = fmt.Sprintf("%s (%s)",
					utils.Cut(utils.Plain(*v.Comment.Body)),
					utils.CutN(utils.Plain(*v.Issue.Title), 30))
			case *github.CommitCommentEvent:
				htmlURL = *v.Comment.HTMLURL
				title = utils.Cut(utils.Plain(*v.Comment.Body))
			case *github.IssuesEvent:
				payloadType = fmt.Sprintf("%s_%s", *v.Action, payloadType)
				htmlURL = *v.Issue.HTMLURL
				title = utils.Cut(utils.Plain(*v.Issue.Title))
			case *github.PullRequestReviewCommentEvent:
				payloadType = "comment"
				htmlURL = *v.Comment.HTMLURL

				if v.PullRequest.Title != nil {
					title = fmt.Sprintf("%s (%s)",
						utils.Cut(utils.Plain(*v.Comment.Body)),
						utils.CutN(utils.Plain(*v.PullRequest.Title), 30))
				} else {
					title = utils.Cut(utils.Plain(*v.Comment.Body))
				}

			case *github.PullRequestEvent:
				htmlURL = *v.PullRequest.HTMLURL
				title = utils.Cut(utils.Plain(*v.PullRequest.Title))
			default:
				panic(fmt.Sprintf("unimplemented event type: %T", v))
			}

			key := fmt.Sprintf("%s-%s", payloadType, htmlURL)

			if _, ok := alreadyOutputs[key]; ok {
				continue
			}

			alreadyOutputs[key] = struct{}{}
			fmt.Printf("- [%s](%s): %s\n", payloadType, htmlURL, title)
		}

		fmt.Println()
	}

	return nil
}
