package pkg

import (
	"time"

	"github.com/shivanshs9/eb-worker-queue/pkg/sqs"
	"github.com/sirupsen/logrus"
)

func StartApp() {
	logger := logrus.New()
	client := sqs.NewSqsClient(logger)

	stop := make(chan struct{})
	options := sqs.ReceiveMessageOptions{
		QueueUrl:            "",
		MaxBufferedMessages: 10,
	}
	stream := client.ReceiveMessageStream(options, stop)
	for {
		select {
		case <-stream:
			logger.Info("Point!")
		case <-time.After(time.Second * 65):
			close(stop)
			return
		}
	}
}
