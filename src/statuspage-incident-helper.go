package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	statuspage "github.com/nagelflorian/statuspage-go"
)

//Incident is the struct type representing a statuspage.io incident
type StatusPageIncident struct {
	Name                          string                 `json:"name,omitempty" structs:"name,omitempty"`
	Status                        string                 `json:"status,omitempty"`
	CreatedAt                     *statuspage.Timestamp  `json:"created_at,omitempty"`
	UpdatedAt                     *statuspage.Timestamp  `json:"updated_at,omitempty"`
	MonitoringAt                  *statuspage.Timestamp  `json:"monitoring_at,omitempty"`
	ResolvedAt                    *statuspage.Timestamp  `json:"resolved_at,omitempty"`
	Impact                        string                 `json:"impact,omitempty"`
	Shortlink                     string                 `json:"shortlink,omitempty"`
	ScheduledFor                  *statuspage.Timestamp  `json:"scheduled_for,omitempty"`
	ScheduledUntil                *statuspage.Timestamp  `json:"scheduled_until,omitempty"`
	ScheduledRemindPrior          bool                   `json:"scheduled_remind_prior,omitempty"`
	SheduledRemindedAt            *statuspage.Timestamp  `json:"scheduled_reminded_at,omitempty"`
	ImpactOverride                string                 `json:"impact_override,omitempty"`
	ScheduledAutoInProgress       bool                   `json:"scheduled_auto_in_progress,omitempty"`
	ScheduledAutoCompleted        bool                   `json:"scheduled_auto_completed,omitempty"`
	Metadata                      StatusPageMetadata     `json:"metadata,omitempty"`
	StartedAt                     *statuspage.Timestamp  `json:"started_at,omitempty"`
	ID                            string                 `json:"id,omitempty"`
	PageID                        string                 `json:"page_id,omitempty"`
	DeliverNotifications          bool                   `json:"deliver_notifications"`
	IncidentUpdates               []IncidentUpdate       `json:"incident_updates,omitempty"`
	PostmortemBody                string                 `json:"postmortem_body,omitempty"`
	PostmortemBodyLastUpdatedAt   *statuspage.Timestamp  `json:"postmortem_body_last_updated_at,omitempty"`
	PostmortemIgnored             bool                   `json:"postmortem_ignored,omitempty"`
	PostmortemPublishedAt         *statuspage.Timestamp  `json:"postmortem_published_at,omitempty"`
	PostmortemNotifiedSubscribers bool                   `json:"postmortem_notified_subscribers,omitempty"`
	PostmortemNotifiedTwitter     bool                   `json:"postmortem_notified_twitter,omitempty"`
	ComponentIDs                  []string               `json:"component_ids,omitempty"`
	Components                    []statuspage.Component `json:"components,omitempty"`
	Body                          string                 `json:"body,omitempty"`
}

//CreateIncidentRequestBody assigns the incident struct to a json key named incident
type CreateIncidentRequestBody struct {
	Incident StatusPageIncident `json:"incident"`
}

//UpdateIncidentRequestBody assigns the incident struct to a json key named incident
type UpdateIncidentRequestBody struct {
	Incident StatusPageIncident `json:"incident"`
}

//Metadata stores metadata details of the incident
type StatusPageMetadata struct {
}

//IncidentUpdate is the struct for storing details about updates to the statuspage incidents
type IncidentUpdate struct {
	Status               string                `json:"status,omitempty"`
	Body                 string                `json:"body,omitempty"`
	CreatedAt            *statuspage.Timestamp `json:"created_at,omitempty"`
	WantsTwitterUpdate   bool                  `json:"wants_twitter_update,omitempty"`
	TwitterUpdatedAt     *statuspage.Timestamp `json:"twitter_updated_at,omitempty"`
	UpdatedAt            *statuspage.Timestamp `json:"updated_at,omitempty"`
	DisplayAt            *statuspage.Timestamp `json:"display_at,omitempty"`
	AffectedComponents   []AffectedComponent   `json:"affected_components,omitempty"`
	DeliverNotifications bool                  `json:"deliver_notifications,omitempty"`
	TweetID              string                `json:"tweet_id,omitempty"`
	ID                   string                `json:"id,omitempty"`
	IncidentID           string                `json:"incident_id,omitempty"`
	CustomTweet          string                `json:"custom_tweet,omitempty"`
}

func (i Incident) String() string {
	return statuspage.Stringify(i)
}

//AffectedComponent is struct to store details about the components affected by the incident
type AffectedComponent struct {
	Code      string `json:"code,omitempty"`
	Name      string `json:"name,omitempty"`
	OldStatus string `json:"old_status,omitempty"`
	NewStatus string `json:"new_status,omitempty"`
}

func prepareStatusPageRequest(method, path string, body interface{}) (*http.Request, error) {
	// rel := &url.URL{Path: path}
	rel := &url.URL{Host: "api.statuspage.io", Scheme: "https", Path: path}
	u := rel.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "OAuth "+os.Getenv("STATUSPAGE_ACCESS_TOKEN"))
	return req, nil
}

func callStatusPage(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		if v == nil {
			return resp, nil
		}

		err = json.NewDecoder(resp.Body).Decode(v)
		return resp, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	return nil, fmt.Errorf("response %s: %d – %s", resp.Status, resp.StatusCode, string(body))
}

//CreateIncident creates an incident for the pageID and incident parameters
func CreateIncident(ctx context.Context, pageID string, incident *StatusPageIncident) (*StatusPageIncident, *http.Response, error) {
	path := "v1/pages/" + pageID + "/incidents"
	payload := CreateIncidentRequestBody{Incident: *incident}
	req, err := prepareStatusPageRequest("POST", path, payload)
	if err != nil {
		return nil, nil, err
	}

	var inc StatusPageIncident
	resp, err := callStatusPage(ctx, req, &inc)
	if err != nil {
		return nil, resp, err
	}
	return &inc, resp, err
}

//UpdateIncident updates an incident for the pageID and incident parameters
func UpdateIncident(ctx context.Context, incident *StatusPageIncident, url string) (*StatusPageIncident, *http.Response, error) {
	index := strings.Index(url, "v1")
	path := url[index:]
	path = strings.Trim(path, "<>")
	fmt.Println("path :", path)
	payload := UpdateIncidentRequestBody{Incident: *incident}
	req, err := prepareStatusPageRequest("PATCH", path, payload)
	if err != nil {
		return nil, nil, err
	}
	var inc StatusPageIncident
	resp, err := callStatusPage(ctx, req, &inc)
	if err != nil {
		return nil, resp, err
	}
	return &inc, resp, err
}
