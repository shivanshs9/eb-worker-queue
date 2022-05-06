package sqs

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

const MAX_ERROR_IGNORE = 5

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

func (msg SQSMessage) String() string {
	return "MSG"
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
	client.log.Info("Starting the SQS Messages Stream")
	stream := make(chan SQSMessage, options.MaxBufferedMessages)
	errorCnt := 0
	go func() {
		for {
			select {
			case <-stop: // triggered when the stop channel is closed
				client.log.Info("Stopping the Stream")
				return // exit
			default:
				msgs, err := client.receiveMessage(options)
				if err != nil {
					errorCnt++
					client.log.WithError(err).Warn("Received error from SQS Client")
					if errorCnt >= MAX_ERROR_IGNORE {
						client.log.Fatal("Max attempts reached, exiting...")
					}
				} else {
					errorCnt = 0
					client.log.WithField("NumMessages", len(msgs)).Info("Received messages")
					for _, msg := range msgs {
						stream <- msg
					}
				}
			}
		}
	}()
	return stream
}

func (client *Client) receiveMessage(options ReceiveMessageOptions) (messages []SQSMessage, err error) {
	maxMsgCount := int64(options.MaxBufferedMessages)
	waitTime := int64(20)
	attributeName := "*"
	input := &sqs.ReceiveMessageInput{
		QueueUrl:              &options.QueueUrl,
		MaxNumberOfMessages:   &maxMsgCount,
		WaitTimeSeconds:       &waitTime,
		MessageAttributeNames: []*string{&attributeName},
	}
	output, err := client.sqs.ReceiveMessage(input)
	if err != nil {
		return
	}
	for _, msg := range output.Messages {
		client.log.Info(msg.String())
		sqsMsg := SQSMessage{}
		messages = append(messages, sqsMsg)
	}
	return
}
