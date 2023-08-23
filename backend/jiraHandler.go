package backend

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"

	"github.com/DebuggerAndrzej/jtl/backend/entities"
)

func getJiraIssueEstimateAsString(estimate int) string {
	est, _ := time.ParseDuration(fmt.Sprintf("%ds", estimate))
	strEstimate := fmt.Sprintf(
		"%sh",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", est.Hours()), "0"), "."),
	)
	return strEstimate
}

func GetJiraClient(config *entities.Config) *jira.Client {
	base := config.JiraBaseUrl
	tp := jira.BasicAuthTransport{Username: config.Username, Password: config.Password}
	client, err := jira.NewClient(tp.Client(), base)

	if err != nil {
		panic("Couldn't get jira client. Please check provided config, internet connection and vpn if applicable")
	}

	return client

}

func GetAllJiraIssuesForAssignee(client *jira.Client, config *entities.Config) ([]entities.Issue, error) {
	var jql string
	if config.AdditionalIssues != "" {
		jql = fmt.Sprintf("assignee = currentuser() OR key in (%s)", config.AdditionalIssues)
	} else {
		jql = "assignee = currentuser()"
	}
	issues, _, err := client.Issue.Search(jql, nil)
	if err != nil {
		return nil, errors.New("Couldn't get issues from Jira API. Check internet connection and vpn if applicable.")
	}

	var mappedIssues []entities.Issue
	activeIssueStatuses := []string{"Open", "In Progress", "In Review"}

	for _, issue := range issues {
		if slices.Contains(activeIssueStatuses, issue.Fields.Status.Name) {
			mappedIssues = append(
				mappedIssues,
				entities.Issue{
					Key:              issue.Key,
					Status:           issue.Fields.Status.Name,
					ShortDescription: issue.Fields.Summary,
					Description:      issue.Fields.Description,
					OriginalEstimate: getJiraIssueEstimateAsString(
						issue.Fields.TimeOriginalEstimate,
					),
					LoggedTime: getJiraIssueEstimateAsString(issue.Fields.TimeSpent),
				},
			)
		}
	}

	return mappedIssues, nil
}

func LogHoursForIssue(client *jira.Client, id, time string) error {
	_, _, err := client.Issue.AddWorklogRecord(id, &jira.WorklogRecord{TimeSpent: time})
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't add worklog to %s issue", id))
	}
	return nil
}

func LogHoursForIssuesScrumMeetings(client *jira.Client, issueId, timeToLog string) error {
	issueCustomFields, _, err := client.Issue.GetCustomFields(issueId)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't get %s issue's custom fields", issueId))
	}
	issuesEpic, _, err := client.Issue.Get(issueCustomFields["customfield_12790"], nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Couldn't find epic for %s issue", issueId))
	}
	var scrumIssue string
	for _, issueLink := range issuesEpic.Fields.IssueLinks {
		outwardIssue := issueLink.OutwardIssue
		if outwardIssue != nil && strings.Contains(issueLink.OutwardIssue.Fields.Summary, "Scrum meetings") {
			scrumIssue = issueLink.OutwardIssue.Key
		}
	}
	if scrumIssue != "" {
		LogHoursForIssue(client, scrumIssue, timeToLog)
	} else {
		return errors.New(fmt.Sprintf("Couldn't find scrum issue for %s issue under %s epic", issueId, issuesEpic))
	}
	return nil
}

func IncrementIssueStatus(client *jira.Client, id, status string) error {
	var err error
	switch status {
	case "Open":
		err = doTransitionToStatus(client, id, "In Progress")
	case "In Progress":
		err = doTransitionToStatus(client, id, "In Review")
	case "In Review":
		err = doTransitionToStatus(client, id, "Done")
	}
	return err
}
func DecrementIssueStatus(client *jira.Client, id, status string) error {
	var err error
	switch status {
	case "In Progress":
		err = doTransitionToStatus(client, id, "Re-open")
	case "In Review":
		err = doTransitionToStatus(client, id, "In Progress")
	}
	return err
}

func doTransitionToStatus(client *jira.Client, id, status string) error {
	var transitionID string
	possibleTransitions, _, _ := client.Issue.GetTransitions(id)
	for _, v := range possibleTransitions {
		if v.Name == status {
			transitionID = v.ID
			break
		}
	}

	if transitionID != "" {
		_, err := client.Issue.DoTransition(id, transitionID)
		if err != nil {
			return errors.New(fmt.Sprintf("Couldn't change status for %s issue", id))
		}
		return nil
	}
	return errors.New(fmt.Sprintf("Couldn't find transitionID required to change issues status for %s issue", id))
}
