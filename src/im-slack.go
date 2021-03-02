package main

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// ******************************************************************************
// Name				: createNewChannel
// Description: Function to create slack channel and add members to it
// ******************************************************************************
func createNewChannel(channelName string, users []User) (*slack.Channel, error) {
	slackAPI := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	channel, err := slackAPI.CreateConversation(channelName, false)
	if err != nil {
		log.Error("Slack channel creation Error:", err)
		return channel, err
	}

	if len(users) > 0 {
		userIDList := []string{}

		// Get users in Slack by email
		for _, j := range users {
			user, err := slackAPI.GetUserByEmail(j.Email)
			if err != nil {
				log.Error("Slack get user by email Error: ", err)
				return channel, err
			}
			userIDList = append(userIDList, user.ID)
		}

		// Invite users to Incident channnel
		channel, err = slackAPI.InviteUsersToConversation(channel.ID, userIDList...)
		if err != nil {
			log.Error("Slack add user to incident channel Error: ", err)
			return channel, err
		}
	}
	return channel, err
}

// ******************************************************************************
// Name				: setChannelPurpose
// Description: Function to set slack channel purpose
// ******************************************************************************
func setChannelPurpose(channelID string, purpose string) (*slack.Channel, error) {
	slackAPI := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	// response, err := slackAPI.SetChannelPurpose(channelID, purpose)
	response, err := slackAPI.SetPurposeOfConversation(channelID, purpose)
	if err != nil {
		log.Error("setChannelPurpose Error: ", err)
	}
	return response, err
}

// ******************************************************************************
// Name				: getChannelPurpose
// Description: Function to get slack channel purpose
// ******************************************************************************
func getChannelPurpose(channelID string) (string, error) {
	slackAPI := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
	// channel, err := slackAPI.GetChannelInfo(channelID)
	channel, err := slackAPI.GetConversationInfo(channelID, false)
	if err != nil {
		log.Error("getChannelPurpose Error: ", err)
	}
	// log.Info(channel.GroupConversation.Purpose.Value)
	return channel.GroupConversation.Purpose.Value, err
}

// ******************************************************************************
// Name				: postMessageToSlackChannel
// Description: Function to post custom message about incident to other slack
// 							channels
// ******************************************************************************
func postMessageToSlackChannel(channelID string, title string) {
	if constants.Slack.NotificationChannelIDs != "" {
		slackAPI := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
		attachment := slack.Attachment{
			Text: "All relevant members are requested to join the group <#" + channelID + ">",
		}
		messageText := "Incident Alert: " + title
		channelsIDs := strings.Split(constants.Slack.NotificationChannelIDs, ",")
		for i := 0; i < len(channelsIDs); i++ {
			channelID := channelsIDs[i]
			channel, timestamp, err := slackAPI.PostMessage(channelID, slack.MsgOptionText(messageText, false), slack.MsgOptionAttachments(attachment))
			if err != nil {
				log.Error("postMessage Error: ", err)
				return
			}
			log.Info("Message successfully sent to channel ", channel, " at ", timestamp)
		}
	}
	log.Info("No Channels configured for posting alerts")
}
