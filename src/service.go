package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// ******************************************************************************
// Name				: pagerDutyService
// Description: Function to handle incident if triggered from pagerduty
// ******************************************************************************
func pagerDutyService(payload Payload) {
	mutex.Lock()
	memberList := loadPDTeamMembers(payload.Messages[0].Incident.Teams[0].Self + "/members")
	users := []User{}

	for _, j := range memberList.Members {
		var user User
		user.ID = j.User.ID
		user = getPDUser(constants.PagerDuty.Endpoint + user.ID)
		users = append(users, user)
	}
	incidentSummary := payload.Messages[0].Incident.Summary
	issue, err := createJiraIssue(incidentSummary)
	if err != nil {
		log.Error("pagerDutyService JIRA Creation Error: ", err)
	}
	log.Info("JIRA issue created: ", issue.Key)

	channelName := getChannelName(issue.Key)
	channel, err := createNewChannel(channelName, users)
	if err != nil {
		log.Error("pagerDutyService Slack Channel Creation Error: ", err)
	}
	log.Info("Slack Channel Created: ", channel.Name)

	componentList := getAffectedSPComponents(payload.Messages[0].Incident.Service.ID)
	var severity string
	statusPageIncident, err := createStatusPageIncident(payload.Messages[0].Incident.Title, payload.Messages[0].Incident.Description, severity, componentList)
	if err != nil {
		log.Error("pagerDutyService StatusPage Incident Creation Error: ", err)
	}
	incidentID := "incidentID : " + statusPageIncident.ID
	statusPageLink := "StatusPage link : " + statusPageIncident.Shortlink
	pagerDutyLink := "PagerDuty link : " + payload.Messages[0].Incident.HTMLURL
	jiraLink := "Jira link : " + (constants.JIRA.Endpoint + "/browse/") + issue.Key
	purpose := incidentID + "\n\n" + statusPageLink + "\n\n" + pagerDutyLink + "\n\n" + jiraLink
	setChannelPurpose(channel.ID, purpose)
	postMessageToSlackChannel(channel.ID, payload.Messages[0].Incident.Title)
	mutex.Unlock()
}

// ******************************************************************************
// Name				: slashCommandService
// Description: Function to perform required actions based on Slack command
// ******************************************************************************
func slashCommandService(w http.ResponseWriter, s slack.SlashCommand, arguments []string) {
	switch arguments[0] {
	case "comment":
		statusPageLink, jiraURL, err := processSlackPurpose(w, s)
		if err != nil {
			return
		}
		err = updateStatePage(arguments[1:], statusPageLink, s)
		if err != nil {
			return
		}
		jiraStatus := setJiraStatusForGenericComment(arguments)
		err = addJiraComment(jiraURL, s.UserName, arguments, jiraStatus, s)
		if err != nil {
			return
		}
		response := SlashResponse{"in_channel", "Comment added to StatusPage and JIRA"}
		slackCommandResponse(response, s)
	case "comment-jira":
		_, jiraURL, err := processSlackPurpose(w, s)
		if err != nil {
			return
		}
		jiraStatus := setJiraStatusForJiraComment(arguments)
		err = addJiraComment(jiraURL, s.UserName, arguments, jiraStatus, s)
		if err != nil {
			return
		}
		response := SlashResponse{"in_channel", "Comment added to JIRA"}
		slackCommandResponse(response, s)
	case "comment-statuspage":
		statusPageLink, _, err := processSlackPurpose(w, s)
		if err != nil {
			return
		}
		err = updateStatePage(arguments[1:], statusPageLink, s)
		if err != nil {
			return
		}
		response := SlashResponse{"in_channel", "Comment added to StatusPage"}
		slackCommandResponse(response, s)
	case "issue":
		go issueCommandService(s, arguments)
	case "statuspage-incident":
		_, jiraURL, err := processSlackPurpose(w, s)
		if err != nil {
			return
		}
		go statuspageCommandService(s, arguments, jiraURL)
	case "help":
		slashHelpResponse(s)
	default:
		slashHelpResponse(s)
		return
	}
}

// ******************************************************************************
// Name				: issueCommandService
// Description: Function to create new Slack channel, StatusPage and JIRA ticket
//              for the incident
// ******************************************************************************
func issueCommandService(s slack.SlashCommand, arguments []string) {
	mutex.Lock()
	issueTitle := arguments[1]
	severity, componentIDList := parseSubCommandArguments(arguments)

	issueKey, err := createJIRAIssue(issueTitle, s)
	if err != nil {
		mutex.Unlock()
		return
	}

	channelID, err := createSlackChannel(issueKey, s)
	if err != nil {
		mutex.Unlock()
		return
	}

	statusPageIncident, err := createStatusPage(s, issueTitle, severity, componentIDList)
	if err != nil {
		mutex.Unlock()
		return
	}

	jiraURL := constants.JIRA.Endpoint + "/browse/" + issueKey
	err = setSlackChannelPurpose(s, statusPageIncident, jiraURL, channelID)
	if err != nil {
		mutex.Unlock()
		return
	}

	// Post message to other relevant channels
	postMessageToSlackChannel(channelID, s.Text)

	// Respond back to the user who executed the command
	responseText := "Success! All relevant members are requested to join the group <#" + channelID + ">"
	response := SlashResponse{"in_channel", responseText}
	slackCommandResponse(response, s)
	mutex.Unlock()
}

// ******************************************************************************
// Name				: statuspageCommandService
// Description: Function to create just StatusPage for the incident
// ******************************************************************************
func statuspageCommandService(s slack.SlashCommand, arguments []string, jiraURL string) {
	mutex.Lock()
	issueTitle := arguments[1]
	severity, componentIDList := parseSubCommandArguments(arguments)

	// Status Page Creation
	statusPageIncident, err := createStatusPage(s, issueTitle, severity, componentIDList)
	if err != nil {
		mutex.Unlock()
		return
	}

	// To set Slack Channel Description
	err = setSlackChannelPurpose(s, statusPageIncident, jiraURL, "")
	if err != nil {
		mutex.Unlock()
		return
	}

	responseText := "Success! Statuspage Incident created!!"
	response := SlashResponse{"in_channel", responseText}
	slackCommandResponse(response, s)
	mutex.Unlock()
}
