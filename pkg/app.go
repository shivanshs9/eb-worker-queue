package pkg

import (
	"errors"
	"fmt"
	"time"

	"github.com/shivanshs9/eb-worker-queue/pkg/http"
	"github.com/shivanshs9/eb-worker-queue/pkg/sqs"
	"github.com/sirupsen/logrus"
)

type AppOptions struct {
	sqs.ReceiveMessageOptions
	ApiHost string
}

type AppCls struct {
	sqsClient  *sqs.Client
	httpClient *http.APIClient
	options    *AppOptions
	log        *logrus.Logger
}

func (app *AppCls) processJob(job *sqs.SQSJobRequest) error {
	// to track execution time
	defer func(start time.Time) {
		app.log.Infof("[%v] Took %v", job.SqsMsgId, time.Since(start))
	}(time.Now())

	app.log.Infof("[%v] Sending POST to %v", job.SqsMsgId, job.AttrJobPath)
	resp, err := app.httpClient.PostRequest(*job)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Received %v from the API call", resp.Status))
	}
	return nil
}

func (app *AppCls) start() {
	stop := make(chan struct{})
	stream := app.sqsClient.ReceiveMessageStream(app.options.ReceiveMessageOptions, stop)
	for job := range stream {
		if err := app.processJob(job); err != nil {
			app.log.Warnf("[%v] Encountered error in processing job: %v", job.SqsMsgId, err)
		} else {
			app.log.Infof("[%v] Finished execution successfully.", job.SqsMsgId)
			app.sqsClient.AcknowledgeMessage(job)
		}
	}
}

func StartApp(options *AppOptions, log *logrus.Logger) {
	app := &AppCls{
		sqsClient:  sqs.NewSqsClient(log),
		httpClient: http.NewAPIClient(options.ApiHost, log),
		log:        log,
		options:    options,
	}
	app.start()
}
