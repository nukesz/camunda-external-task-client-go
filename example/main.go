package main

import (
	"fmt"

	client "github.com/nukesz/camunda-external-task-client-go"
)

func main() {
	// Create a new client to connect to camunda
	externalClient := client.Client{
		BaseURL: "http://localhost:8080/engine-rest",
	}
	// externalClient.ExternalTasks()
	// externalClient.FetchAndLock()
	//externalClient.Complete("e14186d0-64ef-11e9-87bd-0242ac110005")
	s := externalClient.Subscribe("myTopic")
	s.Handler(handle)
	s.Open()
}

func handle(t client.Task, ts client.TaskService) {
	fmt.Printf("Working on the External Task %s\n", t.ID)

	ts.Complete(t)
	fmt.Printf("The External Task %s has been completed!\n", t.ID)
}
