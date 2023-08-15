package backend

import (
	"fmt"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"

	"jtl/backend/entities"
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
	client, _ := jira.NewClient(tp.Client(), base)
	return client

}

func GetAllJiraIssuesForAssignee(client *jira.Client, config *entities.Config) []entities.Issue {
	var jql string
	if config.AdditionalIssues != "" {
		jql = fmt.Sprintf("assignee = currentuser() OR key in (%s)", config.AdditionalIssues)
	} else {
		jql = "assignee = currentuser()"
	}
	issues, _, err := client.Issue.Search(jql, nil)
	if err != nil {
		return nil
	}

	var iss []entities.Issue

	for _, issue := range issues {
		if issue.Fields.Status.Name != "Done" {
			iss = append(
				iss,
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

	return iss
}

func LogHoursForIssue(client *jira.Client, id, time string) {
	client.Issue.AddWorklogRecord(id, &jira.WorklogRecord{TimeSpent: time})
}

func LogHoursForIssuesScrumMeetings(client *jira.Client, issueId, timeToLog string) {
	issueCustomFields, _, _ := client.Issue.GetCustomFields(issueId)
	issuesEpic, _, _ := client.Issue.Get(issueCustomFields["customfield_12790"], nil)
	var scrumIssue string
	for _, issueLink := range issuesEpic.Fields.IssueLinks {
		outwardIssue := issueLink.OutwardIssue
		if outwardIssue != nil && strings.Contains(issueLink.OutwardIssue.Fields.Summary, "Scrum meetings") {
			scrumIssue = issueLink.OutwardIssue.Key
		}
	}
	if scrumIssue != "" {
		LogHoursForIssue(client, scrumIssue, timeToLog)
	}
}

func IncrementIssueStatus(client *jira.Client, id, status string) {
	switch status {
	case "Open":
		doTransitionToStatus(client, id, "In Progress")
	case "In Progress":
		doTransitionToStatus(client, id, "In Review")
	case "In Review":
		doTransitionToStatus(client, id, "Done")
	}
}
func DecrementIssueStatus(client *jira.Client, id, status string) {
	switch status {
	case "In Progress":
		doTransitionToStatus(client, id, "Re-open")
	case "In Review":
		doTransitionToStatus(client, id, "In Progress")
	}
}

func doTransitionToStatus(client *jira.Client, id, status string) {
	var transitionID string
	possibleTransitions, _, _ := client.Issue.GetTransitions(id)
	for _, v := range possibleTransitions {
		if v.Name == status {
			transitionID = v.ID
			break
		}
	}

	client.Issue.DoTransition(id, transitionID)
}
