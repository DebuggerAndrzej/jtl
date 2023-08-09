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

func get_all_jira_issues_for_assignee(client *jira.Client) []Issue {
	opt := &jira.SearchOptions{
		MaxResults: 1000,
	}
	jql := "assignee = currentuser()"
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
