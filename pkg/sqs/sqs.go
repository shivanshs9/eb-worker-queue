package sqs

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

type Client struct {
	sqs *sqs.SQS
	log *logrus.Logger
}

type ReceiveMessageOptions struct {
	QueueUrl            string
	MaxBufferedMessages int
	RetryCount          int
}

type SQSMessage struct {
}

func NewSqsClient(log *logrus.Logger) *Client {
	mySession := session.Must(session.NewSession())
	return &Client{
		sqs: sqs.New(mySession),
		log: log,
	}
}

func (client *Client) ReceiveMessageStream(options ReceiveMessageOptions, stop chan struct{}) chan SQSMessage {
	// errors will be logged and ignored
	stream := make(chan SQSMessage, options.MaxBufferedMessages)
	go func() {
		for {
			select {
			case <-stop: // triggered when the stop channel is closed
				break // exit
			default:
				msgs, err := client.receiveMessage(options)
				if err != nil {
					client.log.WithError(err).Warn("Recieved error from SQS Client")
				}
				for _, msg := range msgs {
					stream <- msg
				}
			}
		}
	}()
	return stream
}

func (client *Client) receiveMessage(options ReceiveMessageOptions) (messages []SQSMessage, err error) {
	maxMsgCount := int64(options.MaxBufferedMessages)
	waitTime := int64(20)
	input := &sqs.ReceiveMessageInput{
		QueueUrl:            &options.QueueUrl,
		MaxNumberOfMessages: &maxMsgCount,
		WaitTimeSeconds:     &waitTime,
	}
	client.sqs.ReceiveMessage(input)
}
