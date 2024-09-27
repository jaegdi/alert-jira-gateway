package main

import (
	"log"

	"alert-jira-gateway/jira"
	"alert-jira-gateway/kubernetes"
	"alert-jira-gateway/prometheus"
)

func main() {
	// Kubernetes Client konfigurieren
	clientset, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Kubernetes-Clients: %v", err)
	}

	// ConfigMap abrufen
	configMapName := "jira-gateway-config"
	configMapNamespace := "default"
	configMap, err := kubernetes.GetConfigMap(clientset, configMapNamespace, configMapName)
	if err != nil {
		log.Fatalf("Fehler beim Abrufen der ConfigMap: %v", err)
	}

	prometheusURL := configMap.Data["PROMETHEUS_URL"]
	jiraURL := configMap.Data["JIRA_URL"]
	jiraProjectKey := configMap.Data["JIRA_PROJECT_KEY"]
	jiraAssignee := configMap.Data["JIRA_ASSIGNEE"]

	// Jira Secret abrufen
	secretName := "jira-secret"
	secretNamespace := "default"
	secret, err := kubernetes.GetSecret(clientset, secretNamespace, secretName)
	if err != nil {
		log.Fatalf("Fehler beim Abrufen des Jira-Secrets: %v", err)
	}

	jiraEmail := string(secret.Data["JIRA_EMAIL"])
	jiraAPIToken := string(secret.Data["JIRA_API_TOKEN"])

	// Jira Client konfigurieren
	jiraClient, err := jira.NewClient(jiraURL, jiraEmail, jiraAPIToken)
	if err != nil {
		log.Fatalf("Fehler beim Erstellen des Jira-Clients: %v", err)
	}

	// Prometheus Alerts abrufen
	alerts, err := prometheus.GetAlerts(prometheusURL)
	if err != nil {
		log.Fatalf("Fehler beim Abrufen der Prometheus-Alerts: %v", err)
	}

	// Alerts als Jira-Tickets senden
	for _, alert := range alerts.Alerts {
		err := jira.CreateIssue(jiraClient, jiraProjectKey, jiraAssignee, alert)
		if err != nil {
			log.Printf("Fehler beim Erstellen des Jira-Tickets: %v", err)
		} else {
			log.Printf("Jira-Ticket für Alert %s erstellt", alert.Labels.Alertname)
		}
	}
}
