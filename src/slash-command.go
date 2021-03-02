package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type SlashResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

// ******************************************************************************
// Name				: slackCommandResponse
// Description: Function to make a POST request to send response back to Slack
// ******************************************************************************
func slackCommandResponse(response SlashResponse, s slack.SlashCommand) {
	json, _ := json.Marshal(response)
	reqBody := bytes.NewBuffer(json)
	endpoint := s.ResponseURL
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		log.Error("slackCommandResponse Build Request Error: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_ACCESS_TOKEN")))
	client := &http.Client{}
	req = req.WithContext(context.Background())
	resp, err := client.Do(req)
	if err != nil {
		log.Error("slackCommandResponse POST Request Error: ", err)
		return
	}
	defer resp.Body.Close()
}

// ******************************************************************************
// Name				: slashHelpResponse
// Description: Function to send help response
// ******************************************************************************
func slashHelpResponse(s slack.SlashCommand) {
	response := SlashResponse{"ephemeral", helpMessage}
	slackCommandResponse(response, s)
}
