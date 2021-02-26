package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// ******************************************************************************
// Name				: createStatusPageIncident
// Description: Function to create status page incident
// ******************************************************************************
func createStatusPageIncident(title string, description string, severity string, components []string) (*StatusPageIncident, error) {
	pageID := constants.StatusPage.PageID
	i := StatusPageIncident{
		PageID:               pageID,
		Name:                 title,
		ImpactOverride:       "minor",
		DeliverNotifications: constants.StatusPage.DeliverNotifications,
		ComponentIDs:         components,
		Body:                 description,
	}
	if severity == "major" || severity == "critical" {
		i.ImpactOverride = severity
		i.DeliverNotifications = true
	}
	incident, _, err := CreateIncident(context.TODO(), pageID, &i)
	if err != nil {
		log.Error("createStatusPageIncident Error: ", err)
		return incident, err
	}
	return incident, err
}

// ******************************************************************************
// Name				: updateStatusPageIncident
// Description: Function to update status page incident
// ******************************************************************************
func updateStatusPageIncident(processedMessage []string, incidentLink string) (*StatusPageIncident, error) {
	if processedMessage[0] == "current" {
		processedMessage[0] = ""
	}
	i := StatusPageIncident{
		Status:               processedMessage[0],
		DeliverNotifications: constants.StatusPage.DeliverNotifications,
		Body:                 processedMessage[1],
	}
	incident, _, err := UpdateIncident(context.TODO(), &i, incidentLink)
	if err != nil {
		log.Error("updateStatusPageIncident Error: ", err)
		return incident, err
	}
	return incident, err
}
