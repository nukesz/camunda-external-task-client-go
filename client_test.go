package client

import "testing"

func TestSubscribe(t *testing.T) {
	topicName := "externalTopic"
	Subscribe(topicName)
}
