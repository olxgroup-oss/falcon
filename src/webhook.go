package main

import (
	"time"
)

// Payload is the top level of the HTTP request body
// containing multiple messages inside.
type Payload struct {
	Messages []Message `json:"messages"`
}

// Message describes each incident
// ref: https://v2.developer.pagerduty.com/docs/webhooks-v2-overview#webhook-payloada
type Message struct {
	ID         string     `json:"id"`
	Event      string     `json:"event"`
	CreatedOn  string     `json:"created_on"`
	Incident   Incident   `json:"incident"`
	Webhook    Webhook    `json:"webhook"`
	LogEntries []LogEntry `json:"log_entries"`
}

// GetCreatedOn ...
func (msg *Message) GetCreatedOn() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, msg.CreatedOn)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

//OncallWrapper is a list of oncalls
type OncallWrapper struct {
	Oncalls []Oncall `json:"oncalls"`
}

// Oncall descries the oncall user
type Oncall struct {
	EscalationPolicy EscalationPolicy `json:"escalation_policy"`
	EscalationLevel  int              `json:"escalation_level"`
	Schedule         Schedule         `json:"schedule"`
	User             User             `json:"user"`
	Start            string           `json:"start"`
	End              string           `json:"end"`
}

//Schedule describes the schedule
type Schedule struct {
}

// Incident describes each incident
// ref: https://v2.developer.pagerduty.com/docs/webhooks-v2-overview#incident-details
type Incident struct {
	ID                   string               `json:"id"`
	Alerts               []Alert              `json:"alerts"`
	IncidentNumber       int                  `json:"incident_number"`
	Title                string               `json:"title"`
	CreatedAt            string               `json:"created_at"`
	Status               string               `json:"status"`
	IncidentKey          string               `json:"incident_key"`
	HTMLURL              string               `json:"html_url"`
	PendingActions       []PendingAction      `json:"pending_actions"`
	Service              Service              `json:"service"`
	Assignments          []Assignment         `json:"assignments"`
	Acknowledgements     []Acknowledgement    `json:"acknowledgements"`
	LastStatusChangeAt   string               `json:"last_status_change_at"`
	LastStatusChangeBy   LastStatusChangeBy   `json:"last_status_change_by"`
	FirstTriggerLogEntry FirstTriggerLogEntry `json:"first_trigger_log_entry"`
	EscalationPolicy     EscalationPolicy     `json:"escalation_policy"`
	Privilege            string               `json:"privilege"`
	Teams                []Team               `json:"teams"`
	Priority             Priority             `json:"priority"`
	Urgency              string               `json:"urgency"`
	ResolveReason        string               `json:"resolve_reason"`
	AlertCounts          AlertCounts          `json:"alert_counts"`
	Metadata             Metadata             `json:"metadata"`
	Type                 string               `json:"type"`
	Summary              string               `json:"summary"`
	Self                 string               `json:"self"`
	Description          string               `json:"description"`
	ImpactedServices     []ImpactedService    `json:"impacted_services"`
	IsMergeable          bool                 `json:"is_mergeable"`
	ExternalReferences   []ExternalReference  `json:"external_references"`
	Importance           string               `json:"importance"`
	BasicAlertGrouping   string               `json:"basic_alert_grouping"`
	IncidentsResponders  []IncidentsResponder `json:"incidents_responders"`
	ResponderRequests    []ResponderRequest   `json:"responder_requests"`
	SubscriberRequests   []SubscriberRequest  `json:"subscriber_requests"`
}

// GetCreatedAt ...
func (i *Incident) GetCreatedAt() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, i.CreatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// IncidentsResponder ...
type IncidentsResponder struct{}

// ResponderRequest ...
type ResponderRequest struct{}

// SubscriberRequest ...
type SubscriberRequest struct{}

// ExternalReference ...
type ExternalReference struct{}

// Webhook describes configuration which resulted in this message
type Webhook struct {
	EndpointURL         string              `json:"endpoint_url"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	WebhookObject       WebhookObject       `json:"webhook_object"`
	Config              Config              `json:"config"`
	OutboundIntegration OutboundIntegration `json:"outbound_integration"`
	AccountsAddon       string              `json:"account_addon"`
	ID                  string              `json:"id"`
	Type                string              `json:"type"`
	Summary             string              `json:"summary"`
	Self                string              `json:"self"`
	HTMLURL             string              `json:"html_url"`
}

// LogEntry describes all events of the incidents.
// ref: https://v2.developer.pagerduty.com/docs/webhooks-v2-overview#log-entries
type LogEntry struct {
	ID           string       `json:"id"`
	Type         string       `json:"type"`
	Summary      string       `json:"summary"`
	Self         string       `json:"self"`
	HTMLURL      string       `json:"html_url"`
	CreatedAt    string       `json:"created_at"`
	Agent        Agent        `json:"agent"`
	Channel      Channel      `json:"channel"`
	Note         string       `json:"note"`
	Contexts     []Context    `json:"contexts"`
	Incident     Incident     `json:"incident"`
	Service      Service      `json:"service"`
	Teams        []Team       `json:"teams"`
	EventDetails EventDetails `json:"event_details"`
}

// GetCreatedAt ...
func (i *LogEntry) GetCreatedAt() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, i.CreatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// Agent ...
type Agent struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// Channel ...
type Channel struct {
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Subject string `json:"subject"`
	Details string `json:"details"`
}

// Alert describes the alert of the incident
type Alert struct {
	AlertKey string `json:"alert_key"`
}

// PendingAction describes the list of pending_actions on the incident.
// A pending_action object contains a type of action which can be escalate, unacknowledge, resolve or urgency_change.
// A pending_action object contains at, the time at which the action will take place.
// An urgency_change pending_action will contain to, the urgency that the incident will change to.
type PendingAction struct {
	Type string `json:"type"`
	At   string `json:"at"`
}

// Service describes the reference to the service associated with the incident.
type Service struct {
	ID                     string              `json:"id"`
	Name                   string              `json:"name"`
	Description            string              `json:"description"`
	AutoResolveTimeout     int                 `json:"auto_resolve_timeout"`
	AcknowledgementTimeout int                 `json:"acknowledgement_timeout"`
	CreatedAt              string              `json:"created_at"`
	Status                 string              `json:"status"`
	SupportHours           string              `json:"support_hours"`
	Addons                 []Addon             `json:"addons"`
	Privilege              string              `json:"privilege"`
	AlertCreation          string              `json:"alert_creation"`
	Integrations           []Integration       `json:"integrations"`
	ScheduledActions       []ScheduledAction   `json:"scheduled_actions"`
	LastIncidentTimestamp  string              `json:"last_incident_timestamp"`
	IncidentUrgencyRule    IncidentUrgencyRule `json:"incident_urgency_rule"`
	EscalationPolicy       EscalationPolicy    `json:"escalation_policy"`
	Teams                  []Team              `json:"teams"`
	Type                   string              `json:"type"`
	Summary                string              `json:"summary"`
	Self                   string              `json:"self"`
	HTMLURL                string              `json:"html_url"`
	Metadata               Metadata            `json:"metadata"`
}

// GetCreatedAt ...
func (i *Service) GetCreatedAt() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, i.CreatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// ImpactedService ...
type ImpactedService struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// Assignment a list of users assigned to the incident at the time of the webhook action.
// Each entry in the array indicates when the user was assigned.
type Assignment struct {
	At       string   `json:"at"`
	Assignee Assignee `json:"assignee"`
}

// Assignee describes the reference of the assignees.
type Assignee struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// Acknowledgement describes the list of all acknowledgements for this incident.
type Acknowledgement struct {
	At           string       `json:"at"`
	Acknowledger Acknowledger `json:"acknowledger"`
}

// Acknowledger describes who acknowledge the incident
type Acknowledger struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// LastStatusChangeBy describes the user or service which is responsible for the incident's last status change.
// If the incident is in the acknowledged or resolved status, this will be the user that took the first acknowledged or resolved action.
// If the incident was automatically resolved (say through the Events API), or if the incident is in the triggered state, this will be the incident's service.
type LastStatusChangeBy struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// FirstTriggerLogEntry the first trigger log entry for the incident.
type FirstTriggerLogEntry struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// EscalationPolicy the escalation policy that the incident is currently following.
type EscalationPolicy struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// Team describes the teams involved in the incident’s lifecycle.
type Team struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

type UserWrapper struct {
	User *User `json:"user"`
}

// User describes the user who is part of the team
type User struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
	Email   string `json:"email"`
}

// Member describes the member of a team user details and role in the team
type Member struct {
	User User   `json:"user"`
	Role string `json:"role"`
}

// TeamMembers describe the list of members who are part of the team, to be used while fetching the data from the api
type TeamMembers struct {
	Members []Member `json:"members"`
}

// Priority describes the priority of the incident.
type Priority struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// AlertCounts describes the summary of the number of alerts by status.
type AlertCounts struct {
	All       int `json:"all"`
	Triggered int `json:"triggered"`
	Resolved  int `json:"resolved"`
}

// Metadata describes the metadata saved on the service.
type Metadata struct{}

// WebhookObject ...
type WebhookObject struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// Config ...
type Config struct{}

// OutboundIntegration ...
type OutboundIntegration struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// Context ...
type Context struct{}

// EventDetails ...
type EventDetails struct {
	Description string `json:"description"`
}

// Addon ...
type Addon struct {
}

// Integration ...
type Integration struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
	Self    string `json:"self"`
	HTMLURL string `json:"html_url"`
}

// IncidentUrgencyRule ...
type IncidentUrgencyRule struct {
	Type    string `json:"type"`
	Urgency string `json:"urgency"`
}

// ScheduledAction ...
type ScheduledAction struct {
}
