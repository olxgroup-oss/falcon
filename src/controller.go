package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// ******************************************************************************
// Name				: healthcheck
// Description: Healthcheck endpoint
// ******************************************************************************
func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Falcon is up and running at port " + constants.ApplicationPort + "\n"))
}

// ******************************************************************************
// Name				: pagerdutyController
// Description: Function to process PagerDuty webhook
// ******************************************************************************
func pagerdutyController(w http.ResponseWriter, r *http.Request) {
	log.Info("Webhook received from PagerDuty")
	w.WriteHeader(http.StatusOK)

	var payload Payload
	json.NewDecoder(r.Body).Decode(&payload)
	priority := payload.Messages[0].Incident.Priority.Summary
	if payload.Messages[0].Event == "incident.trigger" && (priority == "P1" || priority == "P2") {
		go pagerDutyService(payload)
	}
}

// ******************************************************************************
// Name				: slackController
// Description: Entrypoint function for handling Slack commands.
// ******************************************************************************
func slackController(w http.ResponseWriter, r *http.Request) {
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := SlashResponse{"ephemeral", (constants.ValidationMessages.CommandParseError + "\n" + constants.ValidationMessages.UseHelp)}
		slackCommandResponse(response, s)
		log.Error("slackComment Parse Error: ", err)
		return
	}

	if respondWithHelpMessage(s) == true {
		return
	}
	runeText := []rune(s.Text) // To identify double quotes in slack command
	if !isArgumentFormatValid(runeText, s) {
		response := SlashResponse{"ephemeral", (constants.ValidationMessages.IncorrectCommandFormat + "\n" + constants.ValidationMessages.UseHelp)}
		slackCommandResponse(response, s)
		return
	}

	text := string(runeText)
	st := strings.Split(text, "\"")
	var arguments []string
	for _, j := range st {
		j = strings.Trim(j, "\" ")
		if j != "" {
			arguments = append(arguments, j)
		}
	}
	resp, err := parseCommandArguments(s, arguments)
	if err != nil {
		response := SlashResponse{"ephemeral", resp}
		slackCommandResponse(response, s)
		return
	}

	response := SlashResponse{"in_channel", "Processing request"}
	encode, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(encode))
	go slashCommandService(w, s, arguments)
}

// ******************************************************************************
// Name				: updateConfigController
// Description: Function to update config
// ******************************************************************************
func updateConfigController(w http.ResponseWriter, r *http.Request) {
	s := readServices()
	updateConfig(s)
	w.Write([]byte("config updated in config/config.json file"))
}
