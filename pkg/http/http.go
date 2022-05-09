package http

import "github.com/sirupsen/logrus"

type APIClient struct {
	host string
	log  *logrus.Logger
}

type JobRequest struct {
}

func NewAPIClient(host string, log *logrus.Logger) *APIClient {
	return &APIClient{
		host: host,
		log:  log,
	}
}

func (client *APIClient) PostRequest(request JobRequest) {

}
