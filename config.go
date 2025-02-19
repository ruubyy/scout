package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the internal representation of the yaml that determines what
// the app listens to an enqueues
type Config struct {
	Redis RedisConfig `yaml:"redis"`
	AWS   AWSConfig   `yaml:"aws"`
	Queue QueueConfig `yaml:"queue"`
	SQS   SQSConfig
}

// RedisConfig is a nested config that contains the necessary parameters to
// connect to a redis instance and enqueue workers.
type RedisConfig struct {
	Host      string `yaml:"host"`
	Queue     string `yaml:"queue"`
	Namespace string `yaml:"namespace"` // optional
	Password  string `yaml:"password"`  // optional
}

// AWSConfig is a nested config that contains the necessary parameters to
// connect to AWS and read from SQS
type AWSConfig struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Region    string `yaml:"region"`
}

// QueueConfig is a nested config that gives the SQS queue to listen on
// and a mapping of topics to workeers
type QueueConfig struct {
	Name   string `yaml:"name"`
	Worker string `yaml:"worker"`
}

// SQSConfig is a nested config meant to be passed directly to the SQS client
type SQSConfig struct {
	maxNumberOfMessages int64 `yaml:"max_number_of_messages"`
	waitTimeSeconds     int64 `yaml:"wait_time_seconds"`
	visibilityTimeout   int64 `yaml:"visibility_timeout"`
}

// ReadConfig reads from a file with the given name and returns a config or
// an error if the file was unable to be parsed. It does no error checking
// as far as required fields.
func ReadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := new(Config)

	err = yaml.Unmarshal(data, config)
	return config, err
}
