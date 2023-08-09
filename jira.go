package main

import (
	"fmt"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
)

func get_jira_issue_estimate_as_string(estimate int) string {
	est, _ := time.ParseDuration(fmt.Sprintf("%ds", estimate))
	str_estimate := fmt.Sprintf(
		"%sh",
		strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", est.Hours()), "0"), "."),
	)
	return str_estimate
}

func get_jira_client(config *Config) *jira.Client {
	base := config.Jira_base_url
	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}
	client, _ := jira.NewClient(tp.Client(), base)

	return client

}

func get_all_jira_issues_for_assignee(client *jira.Client, config *Config) []Issue {
	opt := &jira.SearchOptions{
		MaxResults: 1000,
	}
	var jql string
	if config.Issues != "" {
		jql = fmt.Sprintf("assignee = currentuser() OR key in (%s)", config.Issues)
	} else {
		jql = "assignee = currentuser()"
	}
	issues, _, err := client.Issue.Search(jql, opt)
	if err != nil {
		return nil
	}

	var iss []Issue

	for _, issue := range issues {
		if issue.Fields.Status.Name != "Done" {
			iss = append(
				iss,
				Issue{
					title:             issue.Key,
					status:            issue.Fields.Status.Name,
					short_description: issue.Fields.Summary,
					original_estimate: get_jira_issue_estimate_as_string(
						issue.Fields.TimeOriginalEstimate,
					),
					logged_time: get_jira_issue_estimate_as_string(issue.Fields.TimeSpent),
				},
			)
		}
	}

	return iss
}

func log_hours_for_issue(client *jira.Client, id, time string) {
	client.Issue.AddWorklogRecord(id, &jira.WorklogRecord{TimeSpent: time})
}

func logHoursForIssuesScrumMeetings(client *jira.Client, issueId, timeToLog string) {
	issueCustomFields, _, _ := client.Issue.GetCustomFields(issueId)
	issuesEpic, _, _ := client.Issue.Get(issueCustomFields["customfield_12790"], nil)
	var scrumIssue string
	for _, issueLink := range issuesEpic.Fields.IssueLinks {
		outward_issue := issueLink.OutwardIssue
		if outward_issue != nil && strings.Contains(issueLink.OutwardIssue.Fields.Summary, "Scrum meetings") {
			scrumIssue = issueLink.OutwardIssue.Key
		}
	}
	if scrumIssue != "" {
		log_hours_for_issue(client, scrumIssue, timeToLog)
	}
}

func incrementIssueStatus(client *jira.Client, id, status string) {
	switch status {
	case "Open":
		doTransitionToStatus(client, id, "In Progress")
	case "In Progress":
		doTransitionToStatus(client, id, "In Review")
	case "In Review":
		doTransitionToStatus(client, id, "Done")
	}
}
func decrementIssueStatus(client *jira.Client, id, status string) {
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
