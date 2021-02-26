package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	helpMessageInitializer()
	// statusPageMappingsInitializer()
	// serviceMappingsInitializer()
	constantsInitializer()

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	router.HandleFunc("/pagerduty/webhook", pagerdutyController).Methods("POST")
	router.HandleFunc("/updateConfig", updateConfigController).Methods("GET")
	router.HandleFunc("/slack/comment", slackController)
	log.Info("Falcon Started on port : ", constants.ApplicationPort)
	log.Fatal(http.ListenAndServe(":8000", router))
}
