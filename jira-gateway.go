// jira-gateway.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andygrunwald/go-jira"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PrometheusAlert struct {
	Status string `json:"status"`
	Alerts []struct {
		Labels struct {
			Alertname string `json:"alertname"`
			Severity  string `json:"severity"`
		} `json:"labels"`
		Annotations struct {
			Summary     string `json:"summary"`
			Description string `json:"description"`
		} `json:"annotations"`
	} `json:"alerts"`
}

func main() {
	// Kubernetes Client konfigurieren
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Fehler beim Erstellen der In-Cluster-Konfiguration: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Kubernetes-Clients: %v", err)
	}

	// Jira Secret abrufen
	secretName := "jira-secret"
	secretNamespace := "default"
	secret, err := clientset.CoreV1().Secrets(secretNamespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Fehler beim Abrufen des Jira-Secrets: %v", err)
	}

	jiraURL := string(secret.Data["JIRA_URL"])
	jiraEmail := string(secret.Data["JIRA_EMAIL"])
	jiraAPIToken := string(secret.Data["JIRA_API_TOKEN"])

	// Jira Client konfigurieren
	jiraClient, err := jira.NewClient(nil, jiraURL)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Jira-Clients: %v", err)
	}

	jiraClient.Authentication.SetBasicAuth(jiraEmail, jiraAPIToken)

	// Prometheus Alerts abrufen
	prometheusURL := "http://prometheus-server/api/v1/alerts"
	resp, err := http.Get(prometheusURL)
	if err != nil {
		log.Fatalf("Fehler beim Abrufen der Prometheus-Alerts: %v", err)
	}
	defer resp.Body.Close()

	var alerts PrometheusAlert
	if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
		log.Fatalf("Fehler beim Dekodieren der Prometheus-Alerts: %v", err)
	}

	// Alerts als Jira-Tickets senden
	for _, alert := range alerts.Alerts {
		issue := jira.Issue{
			Fields: &jira.IssueFields{
				Assignee: &jira.User{
					Name: "your-jira-username",
				},
				Reporter: &jira.User{
					Name: "your-jira-username",
				},
				Description: alert.Annotations.Description,
				Type: jira.IssueType{
					Name: "Task",
				},
				Project: jira.Project{
					Key: "PROJECTKEY",
				},
				Summary: fmt.Sprintf("Prometheus Alert: %s", alert.Labels.Alertname),
			},
		}

		_, _, err := jiraClient.Issue.Create(&issue)
		if err != nil {
			log.Printf("Fehler beim Erstellen des Jira-Tickets: %v", err)
		} else {
			log.Printf("Jira-Ticket für Alert %s erstellt", alert.Labels.Alertname)
		}
	}
}
