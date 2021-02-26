package main

import (
	"os"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	log "github.com/sirupsen/logrus"
	"github.com/trivago/tgo/tcontainer"
)

type TransitionResponse struct {
	Transitions []jira.Transition `json:"transitions" structs:"transitions"`
}

type PostReqTransitionObject struct {
	Transition PostReqTransition `json:"transition"`
	Fields     PostReqFields     `json:"fields"`
}

type PostReqTransition struct {
	ID string `json:"id"`
}

type PostReqFields struct {
	Resolution PostReqResolution `json:"resolution"`
}

type PostReqResolution struct {
	Name string `json:"name"`
}

// ******************************************************************************
// Name				: createJiraIssue
// Description: Function to create JIRA issue ticket
// ******************************************************************************
func createJiraIssue(summary string, user ...User) (*jira.Issue, error) {
	jiraClient := getJIRAClient()
	customFields := tcontainer.NewMarshalMap()
	customFields["customfield_15201"] = time.Now().Format(time.RFC3339)
	i := jira.Issue{
		Fields: &jira.IssueFields{
			Type: jira.IssueType{
				ID: constants.JIRA.IssueTypeID, // To set issue type as incident
			},
			Project: jira.Project{
				ID: constants.JIRA.ProjectID, // To set project
			},
			Summary:  summary,
			Unknowns: customFields,
		},
	}
	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		log.Error("createJiraIssue IssueCreation Error: ", err)
		return issue, err
	}
	return issue, err
}

// ******************************************************************************
// Name				: addComment
// Description: Function to add comment to JIRA Ticket
// ******************************************************************************
func addComment(url string, user string, text string, status string) (*jira.Comment, *jira.Response, error) {
	jiraClient := getJIRAClient()
	var comment *jira.Comment
	var resp *jira.Response
	c := jira.Comment{
		Body: text,
	}
	userEmail := user + "@olx.com"
	jiraUser, _, err := jiraClient.User.Find(userEmail)
	if err != nil {
		log.Error("JIRA user not found", err)
		return comment, resp, err
	} else {
		c.Author = jira.User{
			AccountID: jiraUser[0].AccountID,
		}
	}
	c.Body = user + ": " + text
	url = strings.Trim(url, "<>")
	urlSplit := strings.Split(url, "/")
	issueID := urlSplit[len(urlSplit)-1]
	comment, resp, err = jiraClient.Issue.AddComment(issueID, &c)
	if err != nil {
		log.Error("Error in commenting on JIRA issue: ", err)
		return comment, resp, err
	}
	if status == "close" {
		err = changeJIRATicketStatus(issueID)
	}
	return comment, resp, err
}

// ******************************************************************************
// Name				: getJIRAClient
// Description: Function to get JIRA Client Object
// ******************************************************************************
func getJIRAClient() *jira.Client {
	base := constants.JIRA.Endpoint
	tp := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		log.Error("getJIRAClient Error: ", err)
	}
	return jiraClient
}

// ******************************************************************************
// Name				: getJIRAClient
// Description: Function to change JIRA Ticket Status
// ******************************************************************************
func changeJIRATicketStatus(issueId string) error {
	jiraClient := getJIRAClient()

	// Get Transition Id for Close Transition
	transitionReq, _ := jiraClient.NewRequest("GET", "rest/api/latest/issue/"+issueId+"/transitions?expand=transitions.fields", nil)
	transitions := new(TransitionResponse)
	_, err := jiraClient.Do(transitionReq, transitions)
	if err != nil {
		log.Error("JIRA TransitionRequest Error: ", err)
		return err
	}
	var transitionID string
	for _, transition := range *&transitions.Transitions {
		if strings.ToLower(transition.Name) == "close" {
			log.Debug(transition.ID, " : ", transition.Name)
			transitionID = transition.ID
		}
	}

	// Mark Incident as Done and close Ticket
	postData := PostReqTransitionObject{
		Transition: PostReqTransition{
			ID: transitionID,
		},
		Fields: PostReqFields{
			Resolution: PostReqResolution{
				Name: "Done",
			},
		},
	}
	statusUpdateReq, _ := jiraClient.NewRequest("POST", "rest/api/2/issue/"+issueId+"/transitions", postData)
	_, err = jiraClient.Do(statusUpdateReq, nil)
	if err != nil {
		log.Error("Error occurred while closing JIRA Ticket(" + issueId + ")")
		return err
	}
	return err
}
