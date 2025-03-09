package frk

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

type Github struct {
	Client *github.Client
	Login  string
}

func NewGithub(token string) (*Github, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	user, _, err := client.Users.Get(context.Background(), "")

	if err != nil {
		return nil, err
	}

	return &Github{
		Client: client,
		Login:  *user.Login,
	}, nil
}

type EventsByRepo map[string][]*github.Event

type RepoEvents struct {
	Repo   string
	Events []*github.Event
}

func (m EventsByRepo) Itor() []*RepoEvents {
	repos := make([]string, 0, len(m))

	for r := range m {
		repos = append(repos, r)
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i] < repos[j] })
	items := make([]*RepoEvents, 0, len(m))

	for _, r := range repos {
		items = append(items, &RepoEvents{
			Repo:   r,
			Events: m[r],
		})
	}

	return items
}

func (gh *Github) EventsWithGrouping(from time.Time, to time.Time) (EventsByRepo, error) {
	allEvents := []*github.Event{}
	opt := &github.ListOptions{}

	for {
		events, resp, err := gh.Client.Activity.ListEventsPerformedByUser(context.Background(), gh.Login, false, opt)

		if err != nil {
			return nil, err
		}

		allEvents = append(allEvents, events...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	eventsByRepo := EventsByRepo{}

	for _, e := range allEvents {
		if e == nil {
			continue
		}

		eventType := *e.Type

		if eventType != "IssuesEvent" &&
			eventType != "PullRequestEvent" &&
			eventType != "PullRequestReviewCommentEvent" &&
			eventType != "IssueCommentEvent" &&
			eventType != "CommitCommentEvent" {
			continue
		}

		createdAt := e.CreatedAt.Truncate(24 * time.Hour)

		if createdAt.Before(from) || createdAt.After(to) {
			continue
		}

		repoName := *e.Repo.Name
		repoEvents, ok := eventsByRepo[repoName]

		if !ok {
			repoEvents = []*github.Event{}
		}

		repoEvents = append(repoEvents, e)
		eventsByRepo[repoName] = repoEvents
	}

	return eventsByRepo, nil
}

type PullRequestByRepo map[string][]*github.PullRequest

type RepoPullRequests struct {
	Repo         string
	PullRequests []*github.PullRequest
}

func (m PullRequestByRepo) Itor() []*RepoPullRequests {
	repos := make([]string, 0, len(m))

	for r := range m {
		repos = append(repos, r)
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i] < repos[j] })
	items := make([]*RepoPullRequests, 0, len(m))

	for _, r := range repos {
		items = append(items, &RepoPullRequests{
			Repo:         r,
			PullRequests: m[r],
		})
	}

	return items
}

func (gh *Github) PullRequests() (PullRequestByRepo, error) {
	repoNames, err := gh.AllRepoNames()
	bar := pb.StartNew(len(repoNames))

	if err != nil {
		return nil, err
	}

	pullRequestsByRepo := PullRequestByRepo{}

	for _, r := range repoNames {
		openPRs, err := gh.PullRequest(r, "open")

		if err != nil {
			return nil, err
		}

		closedPRs, err := gh.PullRequest(r, "closed")

		if err != nil {
			return nil, err
		}

		pullRequestsByRepo[r] = append(openPRs, closedPRs...)
		bar.Increment()
	}

	bar.Finish()
	return pullRequestsByRepo, nil
}

func (gh *Github) PullRequest(fullRepoName string, state string) ([]*github.PullRequest, error) {
	ownerRepo := strings.SplitN(fullRepoName, "/", 2)
	allPRs := []*github.PullRequest{}
	opt := &github.PullRequestListOptions{State: state}

	for {
		prs, resp, err := gh.Client.PullRequests.List(context.Background(), ownerRepo[0], ownerRepo[1], opt)

		if err != nil {
			return nil, err
		}

		allPRs = append(allPRs, prs...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allPRs, nil
}

func (gh *Github) AllRepoNames() ([]string, error) {
	repoNames := []string{}
	userRepos, err := gh.RepoNames(gh.Login)

	if err != nil {
		return nil, err
	}

	repoNames = append(repoNames, userRepos...)
	orgs, err := gh.Orgs()

	if err != nil {
		return nil, err
	}

	for _, o := range orgs {
		repos, err := gh.RepoNames(*o.Login)

		if err != nil {
			return nil, err
		}

		repoNames = append(repoNames, repos...)
	}

	return repoNames, nil
}

func (gh *Github) RepoNames(login string) ([]string, error) {
	allRepos := []*github.Repository{}
	opt := &github.RepositoryListOptions{}

	for {
		repos, resp, err := gh.Client.Repositories.List(context.Background(), login, opt)

		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	repoNames := make([]string, 0, len(allRepos))

	for _, r := range allRepos {
		repoNames = append(repoNames, *r.FullName)
	}

	return repoNames, nil
}

func (gh *Github) Orgs() ([]*github.Organization, error) {
	allOrgs := []*github.Organization{}
	opt := &github.ListOptions{}

	for {
		orgs, resp, err := gh.Client.Organizations.List(context.Background(), "", opt)

		if err != nil {
			return nil, err
		}

		allOrgs = append(allOrgs, orgs...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allOrgs, nil
}
