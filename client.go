package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	BaseURL string
}

type Task struct {
	ID         string
	ActivityID string
	TopicName  string
}

func (c Client) Complete(id string) {
	var jsonStr = []byte(`{
		"workerId": "aWorkerId8",
		"variables": {
			"aVariable": {"value": "aStringValue"},
      "anotherVariable": {"value": 42},
			"aThirdVariable": {"value": true}
		},
		"localVariables": {
			"aLocalVariable": {"value": "aStringValue"}
		}
	}`)

	url := fmt.Sprintf("%s/external-task/%s/complete", c.BaseURL, id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("could not complete task: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var tasks []Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	fmt.Printf("Tasks: %v", tasks)
}

// ExternalTasks queries for all the external tasks
func (c Client) ExternalTasks() {
	resp, err := http.Get(c.BaseURL + "/external-task")
	if err != nil {
		log.Fatalf("could not get external tasks %v", err)
	}
	defer resp.Body.Close()

	var tasks []Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	for _, task := range tasks {
		fmt.Printf("Task is %+v\n", task)
	}
}

func (c Client) FetchAndLock() {
	var jsonStr = []byte(`{
		"workerId": "aWorkerId8",
		"maxTasks": 2,
		"usePriority": true,
		"topics": [{"topicName": "goTopic",
            "lockDuration": 50000,
            "processDefinitionId": "Process_0axyt3i:1:d70c9c8b-64ef-11e9-87bd-0242ac110005"
		}]
	}`)

	req, err := http.NewRequest("POST", c.BaseURL+"/external-task/fetchAndLock", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("could not fetch and lock task: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	var tasks []Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	fmt.Printf("Tasks: %v", tasks)
}

func (Client) Subscribe(topicName string) {
	fmt.Printf("Subscribing to: %v\n", topicName)
}
