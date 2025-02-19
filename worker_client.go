package main

import (
	"encoding/json"
	"errors"

	"github.com/jrallison/go-workers"
	log "github.com/sirupsen/logrus"
)

// WorkerClient is an interface for enqueueing workers
type WorkerClient interface {
	// Push pushes a worker onto the queue
	Push(class string, args map[string]interface{}) (string, error)
}

type redisWorkerClient struct {
	queue string
}

// NewRedisWorkerClient creates a worker client that pushes the worker to redis
func NewRedisWorkerClient(redis RedisConfig) (WorkerClient, error) {
	if redis.Host == "" {
		return nil, errors.New("Redis host required")
	}

	if redis.Queue == "" {
		return nil, errors.New("Sidekiq queue required")
	}

	workerConfig := map[string]string{
		"server":   redis.Host,
		"database": "0",
		"pool":     "20",
		"process":  "1",
	}

	if redis.Namespace != "" {
		workerConfig["namespace"] = redis.Namespace
	}

	if redis.Password != "" {
		workerConfig["password"] = redis.Password
	}

	workers.Configure(workerConfig)

	return &redisWorkerClient{queue: redis.Queue}, nil
}

func (r *redisWorkerClient) Push(class string, args map[string]interface{}) (string, error) {
	// This will hopefully deserialize on the ruby end as a hash
	jsonBytes, err := json.Marshal(args)
    if err != nil {
		log.Error("Error converting to JSON: ", err.Error())
    }

    // Convert JSON bytes to string
    jsonString := string(jsonBytes)
	return workers.EnqueueIn(
		r.queue,
		class,
		5.0,
		[]string{jsonString},
	)
}
