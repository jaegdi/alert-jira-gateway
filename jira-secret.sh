#!/usr/bin/env bash

oc create secret generic jira-secret \
  --from-literal=JIRA_URL=https://jira.sf-bk.de \
  --from-literal=JIRA_EMAIL="$(read -r -p "Your email-adress")" \
  --from-literal=JIRA_API_TOKEN="$(read -r -p "Your api-token")" \
  --dry-run=client \
  -o yaml | \
kubeseal \
  --controller-namespace sealed-secrets \
  > sealedsecret-jira.yaml

