package main

import (
	client "github.com/nukesz/camunda-external-task-client-go"
)

func main() {
	// Create a new client to connect to camunda
	externalClient := client.Client{
		BaseURL: "http://localhost:8080/engine-rest",
	}
	// externalClient.ExternalTasks()
	// externalClient.FetchAndLock()
	externalClient.Complete("e14186d0-64ef-11e9-87bd-0242ac110005")
	// externalClient.Subscribe("myTopic")
}
