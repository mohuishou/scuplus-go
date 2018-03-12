package main

import (
	"log"

	"github.com/google/go-github/github"

	"github.com/mohuishou/scuplus-go/config"

	"github.com/mohuishou/scuplus-go/model"
	g "github.com/mohuishou/scuplus-go/util/github"
)

func main() {
	// 获取所有的issue
	c, ctx := g.Task()
	opt := &github.IssueListByRepoOptions{State: "all"}
	opt.PerPage = 100
	issues, _, err := c.Issues.ListByRepo(ctx, config.Get().Github.OwnerUser, config.Get().Github.Repo, opt)
	if err != nil {
		panic(err)
	}
	log.Println(len(issues))
	for _, issue := range issues {
		tags := ""
		for i, l := range issue.Labels {
			tags = tags + "," + l.GetName()
			if i == 0 {
				tags = l.GetName()
			}
		}
		model.UpdateFeedBack(issue.GetNumber(), issue.GetState(), tags)
	}
}
