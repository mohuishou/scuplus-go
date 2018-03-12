package github

import (
	"context"
	"log"
	"net/http"

	"github.com/mohuishou/scuplus-go/job"

	"github.com/RichardKnop/machinery/v1/tasks"

	"github.com/mohuishou/scuplus-go/model"

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

func WebHook(r *http.Request) error {
	payload, err := github.ValidatePayload(r, []byte(config.Get().Github.WebhookSecret))
	if err != nil {
		log.Println("github webhook err:", err)
		return err
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Println("github webhook err:", err)
		return err
	}

	switch event := event.(type) {
	case *github.IssueCommentEvent:
		// 判断是否为用户的评论，仅通知管理员的评论

		if event.GetSender().GetName() != "mohuishou" {
			return nil
		}
		// 给用户发送通知
		sign := &tasks.Signature{
			Name: "notify_feedback",
			Args: []tasks.Arg{
				{
					Type:  "int",
					Value: event.GetIssue().GetNumber(),
				},
				{
					Type:  "string",
					Value: event.GetComment().GetBody(),
				},
			},
		}
		_, err := job.Server.SendTask(sign)
		if err != nil {
			log.Println("notify user feedback", err)
		}
	case *github.IssuesEvent:
		// 更新tags & stat
		issue := event.GetIssue()
		tags := ""
		for i, l := range issue.Labels {
			tags = tags + "," + l.GetName()
			if i == 0 {
				tags = l.GetName()
			}
		}
		model.UpdateFeedBack(issue.GetNumber(), issue.GetState(), tags)
	}
	return nil
}

// Comment 评论
func Comment(id int, content string) error {
	comment := github.IssueComment{Body: &content}
	_, _, err := client.Issues.CreateComment(ctx, config.Get().Github.OwnerUser, config.Get().Github.Repo, id, &comment)
	return err
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
