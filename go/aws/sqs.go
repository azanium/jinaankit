package aws

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	AWSConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/azanium/jinaankit/go/aws/entity"
)

type SQSQueueAPI interface {
	GetQueueUrl(ctx context.Context,
		params *sqs.GetQueueUrlInput,
		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)

	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

func GetQueueURL(c context.Context, api SQSQueueAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	return api.GetQueueUrl(c, input)
}

func SendMessage(c context.Context, api SQSQueueAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return api.SendMessage(c, input)
}
func ReceiveMessage(ctx context.Context, api SQSQueueAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(ctx, input)
}

type SQS struct {
	client *sqs.Client
}

func NewSQS() (*SQS, error) {
	// load the default aws config along with custom resolver.
	cfg, err := AWSConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("configuration error: %v", err)
		return nil, err
	}

	return &SQS{
		client: sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		}),
	}, nil
}

func (s SQS) GetClient() sqs.Client {
	return *s.client
}

func getQueueURL(ctx context.Context, api SQSQueueAPI, queue string) (*string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: &queue,
	}
	result, err := GetQueueURL(ctx, api, input)
	if err != nil {
		log.Fatalf("Error getting the queue URL: %v", err)
		return nil, err
	}
	return result.QueueUrl, nil
}

func (s SQS) SendMessage(ctx context.Context, queue *string) (*string, error) {
	queueURL, err := getQueueURL(ctx, s.client, *queue)
	if err != nil {
		return nil, err
	}

	messageInput := &sqs.SendMessageInput{
		DelaySeconds: 10,
		MessageAttributes: map[string]types.MessageAttributeValue{
			"Blog": {
				DataType:    aws.String("String"),
				StringValue: aws.String("The Code Library"),
			},
			"Article": {
				DataType:    aws.String("Number"),
				StringValue: aws.String("10"),
			},
		},
		MessageBody: aws.String("article about sending a message to AWS SQS"),
		QueueUrl:    queueURL,
	}

	resp, err := SendMessage(ctx, s.client, messageInput)
	if err != nil {
		return nil, err
	}

	return resp.MessageId, nil
}

func (s SQS) ReceiveMessage(ctx context.Context, queue *string) (*entity.Message, error) {
	queueURL, err := getQueueURL(ctx, s.client, *queue)
	if err != nil {
		return nil, err
	}

	recvInput := &sqs.ReceiveMessageInput{
		QueueUrl:              queueURL,
		MessageAttributeNames: []string{"All"},
		MaxNumberOfMessages:   1,
		VisibilityTimeout:     *aws.Int32(10),
	}

	msg, err := ReceiveMessage(ctx, s.client, recvInput)
	if err != nil {
		return nil, err
	}

	if msg.Messages == nil {
		return nil, nil
	}
	attrs := make(map[string]string)
	for key, attr := range msg.Messages[0].MessageAttributes {
		attrs[key] = *attr.StringValue
	}

	return &entity.Message{
		ID:            *msg.Messages[0].MessageId,
		ReceiptHandle: *msg.Messages[0].ReceiptHandle,
		Body:          *msg.Messages[0].Body,
		Attributes:    attrs,
	}, nil
}

func (s SQS) DeleteMessage(ctx context.Context, queueURL string, receiptHandle string) error {
	timeout := time.Second * (10)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: &receiptHandle,
	})
	return err
}

func (s SQS) DecodeS3Message(body string) (*entity.S3Record, error) {
	var records entity.S3Records
	err := json.Unmarshal([]byte(body), &records)
	if err != nil {
		return nil, err
	}
	return &records.Records[0], nil
}
