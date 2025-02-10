package bill

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var TaskQueueName = "bill-task-queue"

//encore:service
type Service struct {
	client client.Client
	worker worker.Worker
}

func initService() (*Service, error) {
	c, err := client.Dial(client.Options{})
	if err != nil {
		return nil, fmt.Errorf("create temporal client: %v", err)
	}

	w := worker.New(c, TaskQueueName, worker.Options{})

	w.RegisterWorkflow(BillWorkflow)

	err = w.Start()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("start temporal worker: %v", err)
	}
	return &Service{client: c, worker: w}, nil
}

func (s *Service) Shutdown(_ context.Context) {
	s.client.Close()
	s.worker.Stop()
}
