package pkg

import (
	"time"

	"github.com/shivanshs9/eb-worker-queue/pkg/http"
	"github.com/shivanshs9/eb-worker-queue/pkg/sqs"
	"github.com/sirupsen/logrus"
)

type AppOptions struct {
	sqs.ReceiveMessageOptions
	ApiHost string
}

func StartApp(options *AppOptions, log *logrus.Logger) {
	client := sqs.NewSqsClient(log)
	httpClient := http.NewAPIClient(options.ApiHost, log)

	stop := make(chan struct{})
	stream := client.ReceiveMessageStream(options.ReceiveMessageOptions, stop)
	for {
		select {
		case job := <-stream:
			log.Info(job.String())
		case <-time.After(time.Second * 65):
			close(stop)
			return
		}
	}
}
