package sqs

import (
	"context"
	"encoding/json"
	"sync"
	"tech-challenge-payment/internal/canonical"
	"tech-challenge-payment/internal/config"
	"tech-challenge-payment/internal/service"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

var (
	once     sync.Once
	instance QueueInterface
)

const (
	PAYMENT = "payment"
	ORDER   = "order"
)

type QueueInterface interface {
	ReceiveMessage()
}

type queueSQS struct {
	sqsService    *sqs.SQS
	service       service.PaymentService
	queuesAddress string
}

func NewSQS() QueueInterface {
	once.Do(func() {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region:     aws.String(config.Get().SQS.Region),
				DisableSSL: aws.Bool(true),
			},
		}))

		sqs := &queueSQS{
			sqsService:    sqs.New(sess),
			service:       service.NewPaymentService(),
			queuesAddress: config.Get().SQS.PaymentPendingQueue,
		}

		instance = sqs
	})

	return instance
}

func (q *queueSQS) ReceiveMessage() {
	for {
		paramsOrder := &sqs.ReceiveMessageInput{
			QueueUrl:            &q.queuesAddress,
			MaxNumberOfMessages: aws.Int64(1),
		}

		resp, err := q.sqsService.ReceiveMessage(paramsOrder)
		if err != nil {
			log.Err(err).Msg("an error occurred when receive message from the queue")
			continue
		}

		if len(resp.Messages) > 0 {
			for _, msg := range resp.Messages {
				log.Info().Any("msg_id", msg.MessageId).Msg("msg received from payment queue")

				err := q.processPaymentMessage([]byte(*msg.Body))
				if err != nil {
					continue
				}

				_, err = q.sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &q.queuesAddress,
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					continue
				}
			}
		} else {
			logrus.Info("there aren't new messages")
			time.Sleep(time.Second * 10)
		}
	}
}

func (q *queueSQS) processPaymentMessage(msg []byte) error {
	var orderId string

	err := json.Unmarshal(msg, &orderId)
	if err != nil {
		log.Err(err).Msg("an error occurred when unmarshal order")
		return err
	}

	_, err = q.service.Create(context.Background(), canonical.Payment{
		OrderID: orderId,
	})
	if err != nil {
		log.Err(err).Any("order_id", orderId).Msg("an error occurred when create payment")
		return err
	}

	return nil
}
