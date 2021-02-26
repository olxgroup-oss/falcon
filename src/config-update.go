package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	statuspage "github.com/nagelflorian/statuspage-go"
	log "github.com/sirupsen/logrus"
)

func readServices() []string {
	bs, err := ioutil.ReadFile(".config/services.txt")
	if err != nil {
		log.Error("readServices Error: ", err)
		os.Exit(1)
	}
	s := string(bs)
	fmt.Println(s)
	return strings.Split(s, "\n")
}

type PDservices struct {
	Pdservices []Service `json:"services"`
}

func updateConfig(serviceArray []string) {
	client := statuspage.NewClient(os.Getenv("STATUSPAGE_ACCESS_TOKEN"), nil)
	components, _ := client.Component.ListComponents(context.TODO(), constants.StatusPage.PageID)
	var serviceMappings ServiceMappings
	for _, j := range serviceArray {
		url := "https://api.pagerduty.com/services?query=" + j
		resp, err := callPagerDuty(url)
		if err != nil {
			log.Error("updateConfig Error: ", err)
		}
		var pdservices PDservices
		json.NewDecoder(resp.Body).Decode(&pdservices)
		for _, service := range pdservices.Pdservices {
			var serviceMap ServiceMap
			serviceMap.PDService.ID = service.ID
			serviceMap.PDService.Name = service.Name
			for _, comp := range *components {
				componentName := strings.ToLower(*comp.Name)
				var spcomponent SPComponent
				if strings.Contains(componentName, j) {
					spcomponent.ID = *comp.ID
					spcomponent.Name = *comp.Name
					serviceMap.SPComponents = append(serviceMap.SPComponents, spcomponent)
				}
			}
			serviceMappings.ServiceMappings = append(serviceMappings.ServiceMappings, serviceMap)
		}
	}
	log.Debug(serviceMappings)
	jsonString, _ := json.MarshalIndent(serviceMappings, "", "\t")
	ioutil.WriteFile("config/config.json", jsonString, os.ModePerm)
}
