package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-github/v33/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

var client *github.Client
var eventInfo map[string]interface{}

func init() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client = github.NewClient(tc)

	eventFile, err := os.Open(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer eventFile.Close()

	eventInfoBytes, err := ioutil.ReadAll(eventFile)
	if err != nil {
		log.Fatal(err)
	}
	eventInfo = make(map[string]interface{})

	err = json.Unmarshal(eventInfoBytes, &eventInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	currentPullRequest, err := getCurrentPullRequest()
	if err != nil {
		log.Warn(err)
		os.Exit(0)
	}

	changedFileNames, err := getChangedFileNames(*currentPullRequest.Number)
	if err != nil {
		log.Fatal(err)
	}

	err = selectChangedFiles("count.out", changedFileNames)
	if err != nil {
		log.Fatal(err)
	}

	outputs, err := exec.Command("go", "tool", "cover", "-func=count.out").Output()
	if err != nil {
		log.Fatal(err)
	}

	msg := buildCommentMsg(strings.Split(string(outputs), "\n"))

	err = createCoverageComment(*currentPullRequest.Number, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getRepoInfo() (string, string, error) {
	repoInfo, ok := eventInfo["repository"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("failed to get repository info")
	}

	name, ok := repoInfo["name"].(string)
	if !ok {
		return "", "", errors.New("failed to get repository name")
	}

	repoOwner, ok := repoInfo["owner"].(map[string]interface{})
	if !ok {
		return "", "", errors.New("failed to get repository owner info")
	}

	repoOwnerName, ok := repoOwner["name"].(string)
	if !ok {
		return "", "", errors.New("failed to get repository owner name")
	}

	return repoOwnerName, name, nil
}

func getPushSha() (string, error) {
	after, ok := eventInfo["after"].(string)
	if !ok {
		return "", errors.New("failed to get after commit sha")
	}

	return after, nil
}

func getCurrentPullRequest() (*github.PullRequest, error) {
	owner, name, err := getRepoInfo()
	if err != nil {
		log.Fatal(err)
	}

	pulls, _, err := client.PullRequests.List(context.Background(), owner, name, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		log.Fatal(err)
	}

	pushSha, err := getPushSha()
	if err != nil {
		log.Fatal(err)
	}

	for _, pull := range pulls {
		if pull.Head != nil && *pull.Head.SHA == pushSha {
			return pull, nil
		}
	}

	return nil, errors.New("current pull request not found")
}

func buildCommentMsg(outputs []string) string {
	msg := `
|file|row|func|coverage|
|--|--|--|--|
`
	if len(outputs) < 3 {
		return strings.Join(outputs, "\n")
	}

	msg = outputs[len(outputs)-2] + "  \n\n" + msg
	for _, o := range outputs[:len(outputs)-2] {
		arr1 := strings.Split(o, ":")
		if len(arr1) != 3 {
			log.Warnf("invalid output %s", o)
			continue
		}

		arr2 := strings.Fields(arr1[2])
		if len(arr2) != 2 {
			log.Warnf("invalid output %s", o)
			continue
		}

		msg = msg + fmt.Sprintf("|%s|%s|%s|%s|\n", arr1[0], arr1[1], arr2[0], arr2[1])
	}
	return msg
}

func createCoverageComment(pullRequestNumber int, message string) error {
	owner, name, err := getRepoInfo()
	if err != nil {
		return nil
	}

	comment := &github.IssueComment{
		Body: &message,
	}
	comment, _, err = client.Issues.CreateComment(context.Background(), owner, name, pullRequestNumber, comment)
	if err != nil {
		return err
	}

	log.Info(*comment.User.ID, *comment.ID)
	return deleteOldComments(pullRequestNumber, *comment.User.ID, *comment.ID)
}

func deleteOldComments(pullRequestNumber int, userID, keepCommentID int64) error {
	owner, name, err := getRepoInfo()
	if err != nil {
		return nil
	}

	oldComments, _, err := client.Issues.ListComments(context.Background(), owner, name, pullRequestNumber, &github.IssueListCommentsOptions{
		ListOptions: github.ListOptions{Page: 1, PerPage: 100},
	})
	if err != nil {
		return nil
	}

	for _, comment := range oldComments {
		log.Info(*comment.User.ID, *comment.ID)
		if *comment.ID == keepCommentID {
			continue
		}

		if *comment.User.ID == userID {
			_, err = client.Issues.DeleteComment(context.Background(), owner, name, *comment.ID)
			if err != nil {
				log.Warnf("failed to delete comment %d", *comment.ID)
			}
		}
	}

	return nil
}

func getChangedFileNames(pullRequestNumber int) ([]string, error) {
	owner, name, err := getRepoInfo()
	if err != nil {
		log.Fatal(err)
	}

	files, resp, err := client.PullRequests.ListFiles(context.Background(), owner, name, pullRequestNumber, &github.ListOptions{Page: 1, PerPage: 100})
	if err != nil {
		return nil, err
	}

	if resp.LastPage > 1 {
		for i := 2; i < resp.LastPage; i++ {
			tmpFiles, _, err := client.PullRequests.ListFiles(context.Background(), owner, name, pullRequestNumber, &github.ListOptions{Page: 1, PerPage: 100})
			if err != nil {
				return nil, err
			}

			files = append(files, tmpFiles...)
		}
	}

	filenames := make([]string, 0)
	for _, file := range files {
		if *file.Changes > 0 {
			filenames = append(filenames, file.GetFilename())
		}
	}
	return filenames, nil
}

func selectChangedFiles(filename string, changeFileNames []string) error {
	coverageFile, err := os.Open(os.Getenv("COVERAGE_OUTPUT_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer coverageFile.Close()

	coverageOutput, err := ioutil.ReadAll(coverageFile)
	if err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(string(coverageOutput), "\n")
	filteredLines := make([]string, 0)
	for i, line := range lines {
		if i == 0 {
			filteredLines = append(filteredLines, line)
			continue
		}

		for _, fn := range changeFileNames {
			if strings.Contains(line, fn) {
				filteredLines = append(filteredLines, line)
				break
			}
		}
	}

	countFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer countFile.Close()

	_, err = countFile.WriteString(strings.Join(filteredLines, "\n"))
	return err
}
