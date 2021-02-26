package main

import (
	"encoding/json"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// ******************************************************************************
// Name				: callPagerDuty
// Description: Helper function to prepare call to pagerduty api
// ******************************************************************************
func callPagerDuty(url string) (*http.Response, error) {
	var Authorization = "Token token=" + os.Getenv("PAGERDUTY_ACCESS_TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", Authorization)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("callPagerDuty Error: ", err)
	}
	return resp, err
}

// ******************************************************************************
// Name				: getPDUser
// Description: Function to get pagerduty user details
// ******************************************************************************
func getPDUser(url string) User {
	resp, err := callPagerDuty(url)
	if err == nil {
		user := new(UserWrapper)
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			log.Error("pdUser parsing Error: ", err)
		}
		return *user.User
	}
	return User{}
}

// ******************************************************************************
// Name				: loadPDTeamMembers
// Description: Function to get team members from pagerduty team
// ******************************************************************************
func loadPDTeamMembers(url string) TeamMembers {
	resp, err := callPagerDuty(url)
	if err == nil {
		var memberList TeamMembers
		err = json.NewDecoder(resp.Body).Decode(&memberList)
		if err != nil {
			log.Error("loadPDTeamMembers parsing Error: ", err)
		}
		return memberList
	}
	return TeamMembers{}
}

// ******************************************************************************
// Name				: getOnCall
// Description: Function to get on call user
// ******************************************************************************
func getOnCall(EscalationPolicy string) User {
	url := "https://api.pagerduty.com/oncalls?escalation_policy_ids[]=" + EscalationPolicy
	resp, err := callPagerDuty(url)
	if err == nil {
		var oncalls OncallWrapper
		json.NewDecoder(resp.Body).Decode(&oncalls)
		oncallUser := getPDUser("https://api.pagerduty.com/users/" + oncalls.Oncalls[0].User.ID)
		return oncallUser
	}
	return User{}
}
