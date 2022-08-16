package http

import (
	"bytes"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/shivanshs9/eb-worker-queue/pkg/sqs"
	"github.com/sirupsen/logrus"
)

type APIClient struct {
	host        string
	contentType string
	log         *logrus.Logger
}

// type JobRequest struct {
// }

func NewAPIClient(host string, log *logrus.Logger) *APIClient {
	return &APIClient{
		host:        host,
		contentType: "application/json",
		log:         log,
	}
}

func (client *APIClient) WaitServerStartup() {
	attempt := 0
	for {
		client.log.Infof("[Attempt %v] Verify API by GET %v", attempt, client.host)
		req, err := http.NewRequest("GET", client.host, nil)
		if err != nil {
			return
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == 200 {
			return
		}
		time.Sleep(5 * time.Second)
	}
}

func (client *APIClient) PostRequest(request sqs.SQSJobRequest) (resp *http.Response, err error) {
	url, err := url.Parse(client.host)
	if err != nil {
		return
	}
	url.Path = path.Join(url.Path, request.AttrJobPath)

	httpReq, err := http.NewRequest("POST", url.String(), bytes.NewBufferString("")) // request.Body unused for now
	if err != nil {
		return
	}
	httpReq.Header.Set("User-Agent", "go-sqsd/1")
	httpReq.Header.Set("X-Aws-Sqsd-Msgid", request.SqsMsgId)
	httpReq.Header.Set("X-Aws-Sqsd-Queue", request.SqsQueueUrl)
	httpReq.Header.Set("X-Aws-Sqsd-Attr-beanstalk.sqsd.path", request.AttrJobPath)
	httpReq.Header.Set("Content-Type", client.contentType)
	resp, err = http.DefaultClient.Do(httpReq)
	if err != nil {
		return
	}
	return
}
