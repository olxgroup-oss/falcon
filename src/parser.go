package main

import (
	"errors"
	"strings"

	"github.com/nlopes/slack"
)

// ******************************************************************************
// Name			  : parseCommandArguments
// Description: Function to parse command arguments
// ******************************************************************************
func parseCommandArguments(s slack.SlashCommand, arguments []string) (string, error) {
	// Format check for `comments` command
	if arguments[0] == "comment" {
		if len(arguments) != 3 {
			response := constants.ValidationMessages.InvalidNumberOfArguments + ". " + constants.ValidationMessages.CommentCommandFormat + "\n" + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
		if !(arguments[1] == "current" || arguments[1] == "investigating" || arguments[1] == "identified" || arguments[1] == "monitoring" || arguments[1] == "resolved") {
			response := constants.ValidationMessages.InvalidStatus + ". " + constants.ValidationMessages.AllowedStatusPageStatus + "\n " + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
	}

	// Format check for `comment-statuspage` command
	if arguments[0] == "comment-statuspage" {
		if len(arguments) != 3 {
			response := constants.ValidationMessages.InvalidNumberOfArguments + ". " + constants.ValidationMessages.CommentStatusPageCommandFormat + "\n" + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
		if !(arguments[1] == "current" || arguments[1] == "investigating" || arguments[1] == "identified" || arguments[1] == "monitoring" || arguments[1] == "resolved") {
			response := constants.ValidationMessages.InvalidStatus + ". " + constants.ValidationMessages.AllowedStatusPageStatus + "\n " + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
	}

	// Format check for `comment-jira` command
	if arguments[0] == "comment-jira" {
		if len(arguments) > 3 {
			response := constants.ValidationMessages.InvalidNumberOfArguments + ". " + constants.ValidationMessages.CommentJiraCommandFormat + "\n" + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
		if len(arguments) == 3 && arguments[1] != "resolved" {
			response := constants.ValidationMessages.InvalidStatus + ". " + constants.ValidationMessages.AllowedJiraStatus + "\n" + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
	}

	// Format check for `issue` command
	if arguments[0] == "issue" || arguments[0] == "statuspage-incident" {
		if len(arguments) < 2 || len(arguments) > 4 {
			response := constants.ValidationMessages.InvalidNumberOfArguments + ". " + constants.ValidationMessages.IssueCommandFormat + "\n" + constants.ValidationMessages.UseHelp
			return response, errors.New("Invalid Arguments")
		}
		if len(arguments) == 3 {
			if !(arguments[2] == "minor" || arguments[2] == "major" || arguments[2] == "critical" || strings.Contains(arguments[2], "components = [") || strings.Contains(arguments[2], "components=[")) {
				response := constants.ValidationMessages.IncorrectCommandFormat + "\n" + constants.ValidationMessages.UseHelp
				return response, errors.New("Invalid Arguments")
			}
		}
		if len(arguments) == 4 {
			if !(arguments[2] == "minor" || arguments[2] == "major" || arguments[2] == "critical") {
				response := constants.ValidationMessages.IncorrectCommandFormat + "\n" + constants.ValidationMessages.UseHelp
				return response, errors.New("Invalid Arguments")
			}
		}
	}

	response := "Valid command"
	return response, nil
}

// ******************************************************************************
// Name: parseSubCommandArguments
// Description: Function to parse sub-command arguments Eg: `issue` sub-command
// ******************************************************************************
func parseSubCommandArguments(arguments []string) (string, []string) {
	var severity string
	var componentIDList []string

	// To set severity or components list if specified in the command
	if len(arguments) == 3 {
		if arguments[2] == "minor" || arguments[2] == "major" || arguments[2] == "critical" {
			severity = arguments[2]
		} else {
			componentIDList = parseAffectedComponents(arguments[2])
		}
	}

	if len(arguments) == 4 {
		severity = arguments[2]
		componentIDList = parseAffectedComponents(arguments[3])
	}

	return severity, componentIDList
}
