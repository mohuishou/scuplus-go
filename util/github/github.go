package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/mohuishou/scuplus-go/config"
	"golang.org/x/oauth2"
)

var client *github.Client

var ctx context.Context

// 初始化，建立一个客户端
func init() {
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Get().Github.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}

// GetIssue 获取指定id的issue
func GetIssue(id int) (*github.Issue, []*github.IssueComment, error) {
	issue, _, err := client.Issues.Get(ctx, config.Get().Github.OwnerUser, config.Get().Github.Repo, id)
	if err != nil {
		return nil, nil, err
	}
	issueComments, _, err := client.Issues.ListComments(ctx, config.Get().Github.OwnerUser, config.Get().Github.Repo, id, &github.IssueListCommentsOptions{})
	return issue, issueComments, err
}

// NewIssue 创建一个新的issue
func NewIssue(title, body string, labels []string) (*github.Issue, error) {
	issue, _, err := client.Issues.Create(ctx, config.Get().Github.OwnerUser, config.Get().Github.Repo, &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	})
	if err != nil {
		return nil, err
	}
	return issue, nil
}
