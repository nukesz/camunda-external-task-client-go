package client

import "testing"

func TestSubscribe(t *testing.T) {
	topicName := "externalTopic"
	client := Client{"url", "demo", "demo"}
	client.Subscribe(topicName)
}
