package main

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// ******************************************************************************
// Name				: createJIRAIssue
// Description: Helper function to create JIRA ticket
// ******************************************************************************
func createJIRAIssue(issueTitle string, s slack.SlashCommand) (string, error) {
	issue, err := createJiraIssue(issueTitle)
	if err != nil {
		msg := "ERROR!! Error creating JIRA Issue: " + err.Error() + "\n" + constants.ValidationMessages.TryAgain + ". " + constants.ValidationMessages.UseHelp
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		return "", errors.New("JiraTicketCreationError")
	}
	log.Info("JIRA Issue Created: ", issue.Key)
	return issue.Key, nil
}

// ******************************************************************************
// Name				: createSlackChannel
// Description: Helper function to create Slack Channel
// ******************************************************************************
func createSlackChannel(issueKey string, s slack.SlashCommand) (string, error) {
	channelName := getChannelName(issueKey)
	channel, err := createNewChannel(channelName, []User{})
	if err != nil {
		msg := "ERROR!! Error creating Slack Channel: " + err.Error() + "\n" + constants.ValidationMessages.TryAgain + ". " + constants.ValidationMessages.UseHelp
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		return "", errors.New("SlackChannelCreationError")
	}
	log.Info("Slack Channel Created: ", channel.Name)
	return channel.ID, nil
}

// ******************************************************************************
// Name				: createStatusPage
// Description: Helper function to create StatusPage Incident
// ******************************************************************************
func createStatusPage(s slack.SlashCommand, issueTitle string, severity string, componentIDList []string) (*StatusPageIncident, error) {
	var description string
	statusPageIncident, err := createStatusPageIncident(issueTitle, description, severity, componentIDList)
	if err != nil {
		msg := "ERROR!! Error creating StatusPage Incident: " + err.Error() + "\n" + constants.ValidationMessages.TryAgain + ". " + constants.ValidationMessages.UseHelp
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		return nil, err
	}
	log.Info("Status Page Created: ", statusPageIncident.ID)
	return statusPageIncident, nil
}

// ******************************************************************************
// Name				: updateStatePage
// Description: Helper function to update StatusPage Incident
// ******************************************************************************
func updateStatePage(arguments []string, statusPageLink string, s slack.SlashCommand) error {
	_, err := updateStatusPageIncident(arguments, statusPageLink)
	if err != nil {
		msg := "ERROR!! Error updating StatusPage Incident: " + err.Error() + "\n" + constants.ValidationMessages.TryAgain + ". " + constants.ValidationMessages.UseHelp
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		return errors.New("StatusPageUpdationError")
	}
	return nil
}

// ******************************************************************************
// Name				: addJiraComment
// Description: Helper function to add comment on JIRA ticket
// ******************************************************************************
func addJiraComment(jiraURL string, username string, arguments []string, jiraStatus string, s slack.SlashCommand) error {
	_, _, err := addComment(jiraURL, s.UserName, arguments[2], jiraStatus)
	if err != nil {
		msg := "ERROR!! Error updating JIRA Issue: " + err.Error() + "\n" + constants.ValidationMessages.TryAgain + ". " + constants.ValidationMessages.UseHelp
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		return errors.New("JiraAddCommentError")
	}
	return nil
}

// ******************************************************************************
// Name				: setSlackChannelPurpose
// Description: Helper function to add slack channel description
// ******************************************************************************
func setSlackChannelPurpose(s slack.SlashCommand, statusPageIncident *StatusPageIncident, jiraURL string, channelID string) error {
	purpose := prepareSlackChannelPurpose(statusPageIncident.ID, statusPageIncident.Shortlink, jiraURL)
	var err error
	if channelID == "" {
		_, err = setChannelPurpose(s.ChannelID, purpose)
	} else {
		_, err = setChannelPurpose(channelID, purpose)
	}
	if err != nil {
		msg := "ERROR!! Error in setting Slack Channel Purpose: " + err.Error() + "\n" + constants.ValidationMessages.TryAgain + ". " + constants.ValidationMessages.UseHelp
		response := SlashResponse{"ephemeral", msg}
		slackCommandResponse(response, s)
		return err
	}
	log.Info("Slack Channel Purpose Added")
	return nil
}

// ******************************************************************************
// Name				: respondWithHelpMessage
// Description: This check is needed if a user just types help after /falcon
// 							command or just types /falcon command
// ******************************************************************************
func respondWithHelpMessage(s slack.SlashCommand) bool {
	if len(s.Text) == 0 || s.Text == "help" {
		slashHelpResponse(s)
		return true
	}
	return false
}
