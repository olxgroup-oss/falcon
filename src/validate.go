package main

import (
	"github.com/slack-go/slack"
)

// ******************************************************************************
// Name			  : isSlackDescriptionInvalid
// Description: Function to validate Slack Channel description before performing
// 							updates
// ******************************************************************************
func isSlackDescriptionInvalid(description []string) bool {
	if len(description) != 3 {
		return true
	}
	return false
}

// ******************************************************************************
//	Name			 : isArgumentFormatValid
// 	Description: Function to validate double quote symbols in slack command
// 							 arguments. To check that all arguments are starting and ending
// 							 with double quotes
// ******************************************************************************
func isArgumentFormatValid(runeText []rune, s slack.SlashCommand) bool {
	quotes := 0
	for i, j := range runeText {
		if int(j) == 8220 || int(j) == 8221 || j == '"' {
			runeText[i] = '"'
			quotes++
		}
	}
	if quotes%2 != 0 {
		return false
	}
	return true
}
