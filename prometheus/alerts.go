package prometheus

import (
    "encoding/json"
    "net/http"
)

type Alert struct {
    Labels struct {
        Alertname string `json:"alertname"`
        Severity  string `json:"severity"`
    } `json:"labels"`
    Annotations struct {
        Summary     string `json:"summary"`
        Description string `json:"description"`
    } `json:"annotations"`
}

type PrometheusAlert struct {
    Status string  `json:"status"`
    Alerts []Alert `json:"alerts"`
}

func GetAlerts(prometheusURL string) (*PrometheusAlert, error) {
    resp, err := http.Get(prometheusURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var alerts PrometheusAlert
    if err := json.NewDecoder(resp.Body).Decode(&alerts); err != nil {
        return nil, err
    }

    return &alerts, nil
}
