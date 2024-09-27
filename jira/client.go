package jira

import (
	"fmt"

	"alert-jira-gateway/prometheus"

	"github.com/andygrunwald/go-jira"
)

func CreateIssue(client *jira.Client, projectKey, assignee string, alert prometheus.Alert) error {
	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				Name: assignee,
			},
			Reporter: &jira.User{
				Name: assignee,
			},
			Description: alert.Annotations.Description,
			Type: jira.IssueType{
				Name: "Task",
			},
			Project: jira.Project{
				Key: projectKey,
			},
			Summary: fmt.Sprintf("Prometheus Alert: %s", alert.Labels.Alertname),
			Unknowns: map[string]interface{}{
				"customfield_12345": "Plattform Services", // Ersetzen Sie "customfield_12345" durch die ID Ihres benutzerdefinierten Feldes
			},
		},
	}

	_, _, err := client.Issue.Create(&issue)
	return err
}
