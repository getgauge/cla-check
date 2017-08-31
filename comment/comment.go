package comment

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/getgauge/cla-check/configuration"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	recheckComment = "@cla-bot check if all the contributors of this pull request has signed our CLA."
)

//PrInfo holds information about a Pull Request
type PrInfo struct {
	Owner    string
	Repo     string
	PrNumber int
}

func isGithubPrURL(url string) bool {
	return strings.HasPrefix(url, "https://github.com") && strings.Contains(url, "pull")
}

func getPrInfo(url string) (PrInfo, error) {
	pi := PrInfo{}
	if !isGithubPrURL(url) {
		return pi, fmt.Errorf("")
	}
	str := strings.TrimPrefix(url, "https://github.com/")
	info := strings.Split(str, "/")
	pn, err := strconv.Atoi(info[3])
	if err != nil {
		return pi, err
	}
	pi.Owner = info[0]
	pi.Repo = info[1]
	pi.PrNumber = pn
	return pi, nil
}

// CreateRecheckComment create a comment in given org, repository, and pr number
// cla-bot gets triggered by this comment
func CreateRecheckComment(url string) error {
	var pi PrInfo
	var err error
	if pi, err = getPrInfo(url); err != nil {
		return err
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: configuration.AccessToken()},
	)
	client := github.NewClient(oauth2.NewClient(ctx, ts))
	body := recheckComment
	comment := github.IssueComment{
		Body: &body,
	}
	_, _, err = client.Issues.CreateComment(ctx, pi.Owner, pi.Repo, pi.PrNumber, &comment)
	if err != nil {
		return err
	}
	return nil
}
