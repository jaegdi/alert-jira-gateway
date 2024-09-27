# Verwende ein offizielles Golang-Image als Basis
FROM golang:1.22-alpine AS builder

# Setze das Arbeitsverzeichnis im Container
WORKDIR /app

# Kopiere die Go-Module-Dateien und installiere die Abhängigkeiten
COPY go.mod go.sum ./
RUN go mod download

# Kopiere den Rest des Quellcodes
COPY . .

# Baue die Applikation
RUN go build -o alert-jira-gateway main.go

# Verwende ein schlankes Image für die Ausführung
FROM alpine:latest

# Setze das Arbeitsverzeichnis im Container
WORKDIR /root/

# Kopiere das gebaute Binary aus dem vorherigen Schritt
COPY --from=builder /app/alert-jira-gateway .

# Exponiere den Port, falls notwendig (optional)
# EXPOSE 8080

# Setze den Befehl zum Ausführen des Binaries
CMD ["./aler-jira-gateway"]