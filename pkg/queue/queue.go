package queue

import (
	"github.com/go-redis/redis"
)

type Queue struct {
	redisClient *redis.Client // Клиент Redis
}

func NewQueue() *Queue {
	// Инициализация клиента Redis
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	return &Queue{
		redisClient: client,
	}
}

func (q *Queue) AddTaskToQueue(taskID string) error {
	// Добавление задачи в очередь
	err := q.redisClient.LPush("task_queue", taskID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (q *Queue) GetTaskFromQueue() (string, error) {
	// Получение задачи из очереди
	taskID, err := q.redisClient.RPop("task_queue").Result()
	if err != nil {
		return "", err
	}
	return taskID, nil
}
