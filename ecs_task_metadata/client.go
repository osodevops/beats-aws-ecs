package ecs_task_metadata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

type TaskMetadataClient struct {
	client *http.Client
	config TaskMetadataClientConfig
}

type TaskMetadataClientConfig struct {
	TaskMetadataEndpoint   string
	MaxRetries             int
	DurationBetweenRetries time.Duration
}

func NewTaskMetadataClient(c *http.Client, cfg TaskMetadataClientConfig) *TaskMetadataClient {
	return &TaskMetadataClient{
		client: c,
		config: cfg,
	}
}

func GetDefaultConfig() TaskMetadataClientConfig {
	taskMetadataEndpoint := strings.TrimRight(os.Getenv("ECS_CONTAINER_METADATA_URI"), "/") + "/task"

	return TaskMetadataClientConfig{
		TaskMetadataEndpoint:   taskMetadataEndpoint,
		MaxRetries:             3,
		DurationBetweenRetries: 1 * time.Second,
	}
}

func (c *TaskMetadataClient) GetTaskMetadata() (*TaskMetadata, error) {
	data, err := c.request(c.config.TaskMetadataEndpoint)
	if err != nil {
		return nil, err
	}

	metadata, err := ParseTaskMetadata(data)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

func (c *TaskMetadataClient) request(endpoint string) ([]byte, error) {
	var resp []byte
	var err error
	for i := 0; i < c.config.MaxRetries; i++ {
		resp, err = c.requestOnce(endpoint)
		if err == nil {
			return resp, nil
		}
		logp.Warn("Attempt [%d/%d]: unable to get metadata response for from '%s': %v", i, c.config.MaxRetries, endpoint, err)
		time.Sleep(c.config.DurationBetweenRetries)
	}

	return nil, err
}

func (c *TaskMetadataClient) requestOnce(endpoint string) ([]byte, error) {
	resp, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Unable to get response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Incorrect status code  %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to read response body: %v", err)
	}

	return body, nil
}
