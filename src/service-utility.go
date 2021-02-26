package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

var mutex sync.Mutex

var serviceMappings *ServiceMappings

var statusPageMappings *StatusPageMappings

type Constants struct {
	Slack              SlackConstants              `json:"slack"`
	ApplicationPort    string                      `json:"application_port"`
	PagerDuty          PagerDutyConstants          `json:"pagerduty"`
	JIRA               JIRAConstants               `json:"jira"`
	StatusPage         StatusPageConstants         `json:"statuspage"`
	ValidationMessages ValidationMessagesConstants `json:"validation_messages"`
}

type ValidationMessagesConstants struct {
	UseHelp                        string `json:"use_help"`
	IncorrectCommandFormat         string `json:"incorrect_command_format"`
	CommandParseError              string `json:"command_parse_error"`
	InvalidNumberOfArguments       string `json:"invalid_no_of_arguments"`
	InvalidStatus                  string `json:"invalid_status"`
	CommentCommandFormat           string `json:"comment_command_format"`
	CommentStatusPageCommandFormat string `json:"comment_statuspage_command_format"`
	CommentJiraCommandFormat       string `json:"comment_jira_command_format"`
	IssueCommandFormat             string `json:"issue_command_format"`
	AllowedStatusPageStatus        string `json:"allowed_statuspage_status"`
	AllowedJiraStatus              string `json:"allowed_jira_status"`
	TryAgain                       string `json:"try_again"`
}

type StatusPageConstants struct {
	PageID               string `json:"page_id"`
	DeliverNotifications bool   `json:"deliver_notifications"`
}

type SlackConstants struct {
	NotificationChannelIDs string `json:"notification_channel_ids"`
}

type PagerDutyConstants struct {
	Endpoint string `json:"endpoint"`
}

type JIRAConstants struct {
	Endpoint    string `json:"endpoint"`
	IssueTypeID string `json:"issue_type_id"`
	ProjectID   string `json:"project_id"`
}

var helpMessage string

var constants *Constants

type PDService struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type SPComponent struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type JiraComponent struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type ServiceMap struct {
	PDService      PDService       `json:"pdservice"`
	SPComponents   []SPComponent   `json:"spcomponents"`
	JiraComponents []JiraComponent `json:"jiracomponents"`
}

type ServiceMappings struct {
	ServiceMappings []ServiceMap `json:"service_mappings,omitempty"`
}

type StatusPageMap struct {
	Service     string      `json:"service,omitempty"`
	SPComponent SPComponent `json:"spcomponent,omitempty"`
}

type StatusPageMappings struct {
	StatusPageMappings []StatusPageMap `json:"statuspage_mappings,omitempty"`
}

// ******************************************************************************
// Name				: helpMessageInitializer
// Description: Function to load custom help message from config file
// ******************************************************************************
func helpMessageInitializer() {
	data, err := ioutil.ReadFile("./config/helpmessage.txt")
	if err != nil {
		log.Error("helpMessageInitializer Error: ", err)
	}
	helpMessage = string(data)
}

// ******************************************************************************
// Name				: constantsInitializer
// Description: Function to load constants from config file
// ******************************************************************************
func constantsInitializer() {
	data, err := ioutil.ReadFile("./config/constants.json")
	if err != nil {
		log.Error("constantsInitializer Error: ", err)
	}
	json.Unmarshal(data, &constants)
}

// ******************************************************************************
// Name				: statusPageMappingsInitializer
// Description: Function to load status page components mapping from config file
// ******************************************************************************
func statusPageMappingsInitializer() {
	if statusPageMappings == nil {
		data, err := ioutil.ReadFile("./config/statuspageMappings.json")
		if err != nil {
			log.Error("statusPageMappingsInitializer Error: ", err)
		}
		json.Unmarshal(data, &statusPageMappings)
	}
	for _, j := range statusPageMappings.StatusPageMappings {
		log.Debug(j.Service, " - ", j.SPComponent.Name)
	}
}

// ******************************************************************************
// Name				: statusPageMappingsInitializer
// Description: Function to load status page components mapping from config file
// ******************************************************************************
func serviceMappingsInitializer() {
	if serviceMappings == nil {
		data, err := ioutil.ReadFile("./config/config.json")
		if err != nil {
			log.Error("serviceMappingsInitializer Error: ", err)
		}
		json.Unmarshal(data, &serviceMappings)
	}
}

func readConfig(serviceID string) *ServiceMap {
	serviceMappingsInitializer()
	for _, j := range serviceMappings.ServiceMappings {
		if j.PDService.ID == serviceID {
			return &j
		}
	}
	return nil
}

func getSPMapping(serviceID string) []SPComponent {
	mapping := readConfig(serviceID)
	if mapping != nil {
		return mapping.SPComponents
	}
	return nil
}

func getAffectedSPComponents(serviceID string) []string {
	var componentIDs []string
	components := getSPMapping(serviceID)
	for _, j := range components {
		componentIDs = append(componentIDs, j.ID)
	}
	return componentIDs

}

func getAffectedJiraComponents(serviceID string) []string {
	mapping := readConfig(serviceID)
	var componentIDs []string
	for _, j := range mapping.JiraComponents {
		componentIDs = append(componentIDs, j.ID)
	}
	return componentIDs
}

func parseAffectedComponents(affectedComponents string) []string {
	componentSlice := strings.SplitN(affectedComponents, "[", 2)
	components := strings.Trim(componentSlice[1], "[], ")
	componentList := strings.Split(components, ",")
	log.Debug(componentList)
	var componentIDList []string

	for _, comp := range componentList {
		for _, j := range statusPageMappings.StatusPageMappings {
			if j.Service == strings.Trim(comp, " ") {
				componentIDList = append(componentIDList, j.SPComponent.ID)
			}
		}
	}
	log.Debug("componentIDList: ", componentIDList)
	return componentIDList
}

func prepareSlackChannelPurpose(incidentID string, statusPageLink string, jiraURL string) string {
	incidentID = "StatusPage Incident ID := " + incidentID
	statusPageLink = "StatusPage Link : " + statusPageLink
	jiraLink := "Jira Link : " + jiraURL
	purpose := incidentID + "\n\n" + statusPageLink + "\n\n" + jiraLink
	return purpose
}

func processSlackPurpose(w http.ResponseWriter, s slack.SlashCommand) (statusPageLInk string, jiraURL string, err error) {
	var statusPageLink, jiraLink string
	purpose, err := getChannelPurpose(s.ChannelID)
	if err != nil {
		msg := "ERROR!! Error reading the purpose of the channel: " + err.Error() + "\n" + "Please make sure the about info of the channel is unchanged and you are using the command from incident channel."
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		err = errors.New("InternalError")
		return statusPageLInk, jiraLink, err
	}
	description := strings.Split(purpose, "\n\n")
	if isSlackDescriptionValid(description) {
		msg := "ERROR!! Error parsing the purpose of the channel. \n" + "Please make sure the about info of the channel is unchanged and you are using the command from incident channel."
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		err = errors.New("InternalError")
		return statusPageLInk, jiraLink, err
	}
	statusPageID := description[0]
	statusPageSplit := strings.Split(statusPageID, " ")
	statusPageLink = "https://api.statuspage.io/v1/pages/" + constants.StatusPage.PageID + "/incidents/" + statusPageSplit[len(statusPageSplit)-1]
	jiraLink = description[len(description)-1]
	jiraLinkSplit := strings.Split(jiraLink, " ")
	jiraURL = jiraLinkSplit[len(jiraLinkSplit)-1]
	return statusPageLink, jiraURL, err
}

func getChannelName(issueKey string) string {
	return ("gl-" + strings.ToLower(issueKey))
}

func setJiraStatusForGenericComment(arguments []string) string {
	var jiraStatus string = ""
	if arguments[1] == "resolved" {
		jiraStatus = "close"
	}
	return jiraStatus
}

func setJiraStatusForJiraComment(arguments []string) string {
	var jiraStatus string = ""
	if len(arguments) == 3 {
		if arguments[1] == "resolved" {
			jiraStatus = "close"
		}
	}
	return jiraStatus
}
