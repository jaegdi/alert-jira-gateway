package jira

import (
    "github.com/andygrunwald/go-jira"
)

func NewClient(jiraURL, jiraEmail, jiraAPIToken string) (*jira.Client, error) {
    tp := jira.BasicAuthTransport{
        Username: jiraEmail,
        Password: jiraAPIToken,
    }

    return jira.NewClient(tp.Client(), jiraURL)
}
