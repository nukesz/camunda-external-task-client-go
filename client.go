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
	BaseURL  string
	Username string
	Password string
}

type Task struct {
	ID                  string `json:"id"`
	ActivityID          string
	TopicName           string
	ProcessDefinitionID string
}

type TaskService struct {
	client *Client
}

func (ts TaskService) Complete(t Task) {
	fmt.Println("Sending task to camunda")
	ts.client.Complete(t.ID)
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
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
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
func (c Client) ExternalTasks(topic string) []Task {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/external-task?topicName=%s&notLocked=true", c.BaseURL, topic), nil)
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("could not get external tasks %v", err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	var tasks []Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	for _, task := range tasks {
		fmt.Printf("Task is %+v\n", task)
	}
	return tasks
}

func (c Client) FetchAndLock() []Task {
	var jsonStr = []byte(`{
		"workerId": "aWorkerId8",
		"maxTasks": 2,
		"usePriority": true,
		"topics": [{"topicName": "goTopic",
            "lockDuration": 10000
		}]
	}`)

	req, err := http.NewRequest("POST", c.BaseURL+"/external-task/fetchAndLock", bytes.NewBuffer(jsonStr))
	if c.Username != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}
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
	fmt.Printf("FetchAndLock Tasks: %v", tasks)
	return tasks
}

type Subscription struct {
	client   *Client
	topic    string
	handlers []func(t Task, ts TaskService)
}

// Handler is attaching handlers to the Subscription
func (s *Subscription) Handler(handler func(t Task, ts TaskService)) {
	s.handlers = append(s.handlers, handler)
}

// Open connects to camunda and start polling the external tasks
// It will call each handler if there is a new task on the topic
func (s *Subscription) Open() {
	go func() {
		tasks := s.client.FetchAndLock()
		fmt.Println(tasks)
		for _, task := range tasks {
			for _, handler := range s.handlers {
				handler(task, TaskService{s.client})
			}
		}

	}()
}

// Subscribe will create a new Subscription on the given topic
func (c Client) Subscribe(topicName string) *Subscription {
	fmt.Printf("Subscribing to: %v\n", topicName)
	return &Subscription{
		client: &c,
		topic:  topicName,
	}
}
