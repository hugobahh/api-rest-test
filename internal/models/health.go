package models

import "time"

// https://inadarei.github.io/rfc-healthcheck/

// HealthCheck :
type HealthCheck struct {
	Status      string                    `json:"status" example:"pass"`
	Description string                    `json:"description,omitempty" example:"Kafka and SQS health check"`
	Version     string                    `json:"version" example:"1.0"`
	Checks      map[string]ComponentCheck `json:"checks"`
}

// ComponentCheck :
type ComponentCheck struct {
	Name          string    `json:"componentName"`
	Type          string    `json:"componentType"`
	ObservedValue float64   `json:"observedValue,omitempty"`
	ObservedUnit  string    `json:"observedUnit,omitempty"`
	Time          time.Time `json:"time"`
	Status        string    `json:"status"`
	Output        string    `json:"output,omitempty"`
}

func NewHealthCheck(version, status string) HealthCheck {
	return HealthCheck{
		Version: version,
		Status:  status,
		Checks:  make(map[string]ComponentCheck)}
}

func NewComponentCheck(name, componentType string, now time.Time) ComponentCheck {
	return ComponentCheck{
		Name: name,
		Type: componentType,
		Time: now}
}
